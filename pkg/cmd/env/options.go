/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package env

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/gardener/gardenctl-v2/internal/gardenclient"
	"github.com/gardener/gardenctl-v2/internal/util"
	"github.com/gardener/gardenctl-v2/pkg/ac"
	"github.com/gardener/gardenctl-v2/pkg/cmd/base"
	"github.com/gardener/gardenctl-v2/pkg/config"
	"github.com/gardener/gardenctl-v2/pkg/target"
)

type options struct {
	base.Options

	// Unset resets environment variables and configuration of the cloudprovider CLI for your shell.
	Unset bool
	// Shell to configure.
	Shell string
	// GardenDir is the configuration directory of gardenctl.
	GardenDir string
	// SessionDir is the session directory of gardenctl.
	SessionDir string
	// CmdPath is the path of the called command.
	CmdPath string
	// Target is the target used when executing the command
	Target target.Target
	// TargetFlags are the target override flags
	TargetFlags target.TargetFlags
	// ProviderType is the name of the cloud provider
	ProviderType string
	// Template is the script template
	Template Template
	// Symlink indicates if KUBECONFIG environment variable should point to the session stable symlink
	Symlink bool
	// Force generates the script even if there are access restrictions to be confirmed
	Force bool
}

// Complete adapts from the command line args to the data required.
func (o *options) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	o.Shell = cmd.Name()
	o.CmdPath = cmd.Parent().CommandPath()
	o.GardenDir = f.GardenHomeDir()
	o.Template = newTemplate("helpers")

	//nolint:gocritic // accept singleCaseSwitch to be consistent with rest of the file. Will be resolved once we refactor to have own options for each provider type
	switch o.ProviderType {
	case "kubernetes":
		filename := filepath.Join(o.GardenDir, "templates", "kubernetes.tmpl")
		if err := o.Template.ParseFiles(filename); err != nil {
			return err
		}
	}

	manager, err := f.Manager()
	if err != nil {
		return err
	}

	o.Symlink = manager.Configuration().SymlinkTargetKubeconfig()
	o.SessionDir = manager.SessionDir()
	o.TargetFlags = manager.TargetFlags()

	return nil
}

// Validate validates the provided command options.
func (o *options) Validate() error {
	if o.Shell == "" {
		return pflag.ErrHelp
	}

	s := Shell(o.Shell)
	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}

// AddFlags binds the command options to a given flagset.
func (o *options) AddFlags(flags *pflag.FlagSet) {
	var text string

	switch o.ProviderType {
	case "kubernetes":
		text = "the KUBECONFIG environment variable"
	default:
		text = "the cloud provider CLI environment variables and logout"

		flags.BoolVarP(&o.Force, "force", "f", false, "Generate the script even if there are access restrictions to be confirmed")
	}

	usage := fmt.Sprintf("Generate the script to unset %s for %s", text, o.Shell)
	flags.BoolVarP(&o.Unset, "unset", "u", o.Unset, usage)
}

// Run does the actual work of the command.
func (o *options) Run(f util.Factory) error {
	ctx := f.Context()

	manager, err := f.Manager()
	if err != nil {
		return err
	}

	o.Target, err = manager.CurrentTarget()
	if err != nil {
		return err
	}

	switch o.ProviderType {
	case "kubernetes":
		if !o.Symlink && o.Target.GardenName() == "" {
			return target.ErrNoGardenTargeted
		}

		return o.runKubernetes(ctx, manager)
	default:
		if o.Target.GardenName() == "" {
			return target.ErrNoGardenTargeted
		}

		return o.run(ctx, manager)
	}
}

func (o *options) runKubernetes(ctx context.Context, manager target.Manager) error {
	data := map[string]interface{}{
		"__meta": generateMetadata(o),
	}

	if !o.Unset {
		var filename string

		if o.Symlink {
			filename = filepath.Join(o.SessionDir, "kubeconfig.yaml")

			if !o.Target.IsEmpty() {
				_, err := os.Lstat(filename)
				if os.IsNotExist(err) {
					return fmt.Errorf("symlink to targeted cluster does not exist: %w", err)
				}
			}
		} else {
			config, err := manager.ClientConfig(ctx, o.Target)
			if err != nil {
				return err
			}

			filename, err = manager.WriteClientConfig(config)
			if err != nil {
				return err
			}
		}

		data["filename"] = filename
	}

	return o.Template.ExecuteTemplate(o.IOStreams.Out, o.Shell, data)
}

func (o *options) run(ctx context.Context, manager target.Manager) error {
	t := o.Target

	client, err := manager.GardenClient(t.GardenName())
	if err != nil {
		return fmt.Errorf("failed to create garden cluster client: %w", err)
	}

	if t.ShootName() == "" && t.SeedName() != "" {
		if shoot, err := client.GetShootOfManagedSeed(ctx, t.SeedName()); err != nil {
			if apierrors.IsNotFound(err) {
				return fmt.Errorf("cannot generate cloud provider CLI configuration script for non-managed seeds: %w", err)
			}

			return err
		} else if shoot != nil {
			o.Target = o.Target.WithProjectName("garden").WithShootName(shoot.Name)
		}
	}

	if o.Target.ShootName() == "" {
		return target.ErrNoShootTargeted
	}

	shoot, err := client.FindShoot(ctx, o.Target.AsListOption())
	if err != nil {
		return err
	}

	secretBinding, err := client.GetSecretBinding(ctx, shoot.Namespace, shoot.Spec.SecretBindingName)
	if err != nil {
		return err
	}

	secret, err := client.GetSecret(ctx, secretBinding.SecretRef.Namespace, secretBinding.SecretRef.Name)
	if err != nil {
		return err
	}

	cloudProfile, err := client.GetCloudProfile(ctx, shoot.Spec.CloudProfileName)
	if err != nil {
		return err
	}

	// check access restrictions
	messages, err := o.checkAccessRestrictions(manager.Configuration(), t.GardenName(), shoot)
	if err != nil {
		return err
	}

	return execTmpl(o, shoot, secret, cloudProfile, messages)
}

func execTmpl(o *options, shoot *gardencorev1beta1.Shoot, secret *corev1.Secret, cloudProfile *gardencorev1beta1.CloudProfile, messages ac.AccessRestrictionMessages) error {
	o.ProviderType = shoot.Spec.Provider.Type

	metadata := generateMetadata(o)

	if len(messages) > 0 {
		b := &bytes.Buffer{}
		messages.Render(b)

		if o.TargetFlags.ShootName() == "" || o.Force {
			metadata["notification"] = b.String()
		} else {
			s := Shell(o.Shell)
			return o.Template.ExecuteTemplate(o.IOStreams.Out, "printf", map[string]interface{}{
				"format": b.String() + "\n%s %s\n%s\n",
				"arguments": []string{
					"The cloud provider CLI configuration script can only be generated if you confirm the access despite the existing restrictions.",
					"Use the --force flag to confirm the access.",
					s.Prompt(runtime.GOOS) + s.EvalCommand(fmt.Sprintf("%s --force %s", o.CmdPath, o.Shell)),
				},
			})
		}
	}

	data := map[string]interface{}{
		"__meta": metadata,
		"region": shoot.Spec.Region,
	}

	for key, value := range secret.Data {
		data[key] = string(value)
	}

	switch o.ProviderType {
	case "azure":
		if !o.Unset {
			configDir, err := createProviderConfigDir(o.SessionDir, o.ProviderType)
			if err != nil {
				return err
			}

			data["configDir"] = configDir
		}
	case "gcp":
		credentials := make(map[string]interface{})

		serviceaccountJSON, err := parseGCPCredentials(secret, &credentials)
		if err != nil {
			return err
		}

		if !o.Unset {
			configDir, err := createProviderConfigDir(o.SessionDir, o.ProviderType)
			if err != nil {
				return err
			}

			data["configDir"] = configDir
		}

		data["credentials"] = credentials
		data["serviceaccount.json"] = string(serviceaccountJSON)
	case "openstack":
		authURL, err := getKeyStoneURL(cloudProfile, shoot.Spec.Region)
		if err != nil {
			return err
		}

		data["authURL"] = authURL
	}

	filename := filepath.Join(o.GardenDir, "templates", o.ProviderType+".tmpl")
	if err := o.Template.ParseFiles(filename); err != nil {
		return fmt.Errorf("failed to generate the cloud provider CLI configuration script: %w", err)
	}

	return o.Template.ExecuteTemplate(o.IOStreams.Out, o.Shell, data)
}

func generateMetadata(o *options) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["unset"] = o.Unset
	metadata["shell"] = o.Shell
	metadata["commandPath"] = o.CmdPath
	metadata["cli"] = getProviderCLI(o.ProviderType)
	metadata["prompt"] = Shell(o.Shell).Prompt(runtime.GOOS)
	metadata["targetFlags"] = getTargetFlags(o.Target)

	return metadata
}

func getProviderCLI(providerType string) string {
	switch providerType {
	case "alicloud":
		return "aliyun"
	case "gcp":
		return "gcloud"
	case "azure":
		return "az"
	case "kubernetes":
		return "kubectl"
	default:
		return providerType
	}
}

func getTargetFlags(t target.Target) string {
	if t.ProjectName() != "" {
		return fmt.Sprintf("--garden %s --project %s --shoot %s", t.GardenName(), t.ProjectName(), t.ShootName())
	}

	return fmt.Sprintf("--garden %s --seed %s --shoot %s", t.GardenName(), t.SeedName(), t.ShootName())
}

func getKeyStoneURL(cloudProfile *gardencorev1beta1.CloudProfile, region string) (string, error) {
	config, err := gardenclient.CloudProfile(*cloudProfile).GetOpenstackProviderConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get openstack provider config: %w", err)
	}

	for _, keyStoneURL := range config.KeyStoneURLs {
		if keyStoneURL.Region == region {
			return keyStoneURL.URL, nil
		}
	}

	if config.KeyStoneURL != "" {
		return config.KeyStoneURL, nil
	}

	return "", fmt.Errorf("cannot find keystone URL for region %q in cloudprofile %q", region, cloudProfile.Name)
}

func parseGCPCredentials(secret *corev1.Secret, credentials interface{}) ([]byte, error) {
	data := secret.Data["serviceaccount.json"]
	if data == nil {
		return nil, fmt.Errorf("no \"serviceaccount.json\" data in Secret %q", secret.Name)
	}

	if err := json.Unmarshal(data, credentials); err != nil {
		return nil, err
	}

	return json.Marshal(credentials)
}

func createProviderConfigDir(sessionDir string, providerType string) (string, error) {
	cli := getProviderCLI(providerType)
	configDir := filepath.Join(sessionDir, ".config", cli)

	err := os.MkdirAll(configDir, 0o700)
	if err != nil {
		return "", fmt.Errorf("failed to create %s configuration directory: %w", cli, err)
	}

	return configDir, nil
}

func (o *options) checkAccessRestrictions(cfg *config.Config, gardenName string, shoot *gardencorev1beta1.Shoot) (ac.AccessRestrictionMessages, error) {
	if cfg == nil {
		return nil, errors.New("garden configuration is required")
	}

	garden, err := cfg.Garden(gardenName)
	if err != nil {
		return nil, err
	}

	messages := ac.CheckAccessRestrictions(garden.AccessRestrictions, shoot)

	return messages, nil
}
