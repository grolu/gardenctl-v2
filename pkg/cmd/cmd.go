/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gardener/gardenctl-v2/internal/util"
	cmdconfig "github.com/gardener/gardenctl-v2/pkg/cmd/config"
	cmdenv "github.com/gardener/gardenctl-v2/pkg/cmd/env"
	cmdssh "github.com/gardener/gardenctl-v2/pkg/cmd/ssh"
	cmdtarget "github.com/gardener/gardenctl-v2/pkg/cmd/target"
	cmdversion "github.com/gardener/gardenctl-v2/pkg/cmd/version"
	"github.com/gardener/gardenctl-v2/pkg/target"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

const (
	envPrefix        = "GCTL"
	envGardenHomeDir = envPrefix + "_HOME"
	envConfigName    = envPrefix + "_CONFIG_NAME"
	envSessionID     = envPrefix + "_SESSION_ID"
	envTermSessionID = "TERM_SESSION_ID"

	gardenHomeFolder = ".garden"
	configName       = "gardenctl-v2"
	targetFilename   = "target.yaml"
)

var (
	sidRegexp  = regexp.MustCompile(`^[\w-]{1,128}$`)
	uuidRegexp = regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12})`)
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the root cmd.
func Execute() {
	cmd := NewDefaultGardenctlCommand()
	// any error would already be printed, so avoid doing it again here
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}

// NewDefaultGardenctlCommand creates the `gardenctl` command with defaults
func NewDefaultGardenctlCommand() *cobra.Command {
	factory := util.FactoryImpl{
		TargetFlags: target.NewTargetFlags("", "", "", "", false),
	}
	ioStreams := util.NewIOStreams()

	return NewGardenctlCommand(&factory, ioStreams)
}

// NewGardenctlCommand creates the `gardenctl` command
func NewGardenctlCommand(f *util.FactoryImpl, ioStreams util.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gardenctl",
		Short: "gardenctl is a utility to interact with Gardener installations",
		Long: `Gardenctl is a utility to interact with Gardener installations.

The state of gardenctl is bound to a shell session and is not shared across windows, tabs or panes.
A shell session is defined by the environment variable GCTL_SESSION_ID. If this is not defined,
the value of the TERM_SESSION_ID environment variable is used instead. If both are not defined,
this leads to an error and gardenctl cannot be executed. The target.yaml and temporary
kubeconfig.*.yaml files are store in the following directory ${TMPDIR}/garden/${GCTL_SESSION_ID}.

You can make sure that GCTL_SESSION_ID or TERM_SESSION_ID is always present by adding
the following code to your terminal profile ~/.profile, ~/.bashrc or comparable file.
  bash and zsh: [ -n "$GCTL_SESSION_ID" ] || [ -n "$TERM_SESSION_ID" ] || export GCTL_SESSION_ID=$(uuidgen)
  fish:         [ -n "$GCTL_SESSION_ID" ] || [ -n "$TERM_SESSION_ID" ] || set -gx GCTL_SESSION_ID (uuidgen)
  powershell:   if ( !(Test-Path Env:GCTL_SESSION_ID) -and !(Test-Path Env:TERM_SESSION_ID) ) { $Env:GCTL_SESSION_ID = [guid]::NewGuid().ToString() }

Find more information at: https://github.com/gardener/gardenctl-v2/blob/master/README.md
`,
		SilenceUsage: true,
	}

	cmd.SetIn(ioStreams.In)
	cmd.SetOut(ioStreams.Out)
	cmd.SetErr(ioStreams.ErrOut)

	// register initializers
	cobra.OnInitialize(func() {
		initConfig(f)
	})

	flags := cmd.PersistentFlags()

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. klog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	addKlogFlags(flags)

	// Do not precalculate what $HOME is for the help text, because it prevents
	// usage where the current user has no home directory (which might _just_ be
	// the reason the user chose to specify an explicit config file).
	flags.StringVar(&f.ConfigFile, "config", "", fmt.Sprintf("config file (default is %s)", filepath.Join("~", gardenHomeFolder, configName+".yaml")))

	// allow to temporarily re-target a different cluster
	f.TargetFlags.AddFlags(flags)

	registerCompletionFuncForGlobalFlags(cmd, f, ioStreams)

	// add subcommands
	cmd.AddCommand(cmdssh.NewCmdSSH(f, cmdssh.NewSSHOptions(ioStreams)))
	cmd.AddCommand(cmdtarget.NewCmdTarget(f, cmdtarget.NewTargetOptions(ioStreams)))
	cmd.AddCommand(cmdversion.NewCmdVersion(f, cmdversion.NewVersionOptions(ioStreams)))
	cmd.AddCommand(cmdconfig.NewCmdConfig(f, ioStreams))
	cmd.AddCommand(cmdenv.NewCmdProviderEnv(f, ioStreams))
	cmd.AddCommand(cmdenv.NewCmdKubectlEnv(f, ioStreams))

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(f *util.FactoryImpl) {
	var err error

	if f.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(f.ConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		configPath := filepath.Join(home, gardenHomeFolder)

		// Search config in ~/.garden or in path provided with the env variable GCTL_HOME with name "gardenctl-v2" (without extension) or name from env variable GCTL_CONFIG_NAME.
		envHomeDir, err := homedir.Expand(os.Getenv(envGardenHomeDir))
		cobra.CheckErr(err)

		viper.AddConfigPath(envHomeDir)
		viper.AddConfigPath(configPath)
		if os.Getenv(envConfigName) != "" {
			viper.SetConfigName(os.Getenv(envConfigName))
		} else {
			viper.SetConfigName(configName)
		}
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		klog.Errorf("failed to read config file: %v", err)
	}

	// initialize the factory

	// prefer an explicit GCTL_HOME env,
	// but fallback to the system-defined home directory
	home := os.Getenv(envGardenHomeDir)
	if len(home) == 0 {
		home, err = homedir.Dir()
		cobra.CheckErr(err)

		home = filepath.Join(home, gardenHomeFolder)
	}

	f.ConfigFile = viper.ConfigFileUsed()
	f.GardenHomeDirectory = home

	sid, err := getSessionID()
	cobra.CheckErr(err)

	f.SessionDirectory = filepath.Join(os.TempDir(), "garden", sid)
	err = os.MkdirAll(f.SessionDirectory, 0700)
	cobra.CheckErr(err)

	f.TargetFile = filepath.Join(f.SessionDirectory, targetFilename)
}

func registerCompletionFuncForGlobalFlags(cmd *cobra.Command, f *util.FactoryImpl, ioStreams util.IOStreams) {
	utilruntime.Must(cmd.RegisterFlagCompletionFunc("garden", completionWrapper(f, ioStreams, gardenFlagCompletionFunc)))
	utilruntime.Must(cmd.RegisterFlagCompletionFunc("project", completionWrapper(f, ioStreams, projectFlagCompletionFunc)))
	utilruntime.Must(cmd.RegisterFlagCompletionFunc("seed", completionWrapper(f, ioStreams, seedFlagCompletionFunc)))
	utilruntime.Must(cmd.RegisterFlagCompletionFunc("shoot", completionWrapper(f, ioStreams, shootFlagCompletionFunc)))
}

type cobraCompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
type cobraCompletionFuncWithError func(ctx context.Context, manager target.Manager, tf target.TargetFlags) ([]string, error)

func completionWrapper(f *util.FactoryImpl, ioStreams util.IOStreams, completer cobraCompletionFuncWithError) cobraCompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		tf := f.TargetFlags

		// By default, the factory will provide a target manager that uses the dynamicTargetProvider (DTP)
		// implementation, i.e. is based on the file just as much as the CLI flags.
		// The DTP tries to allow users to "move up", i.e. when they already targeted a shoot, just adding
		// "--garden foo" should not just change the used garden cluster, but _target_ the garden (instead
		// of the shoot). This behaviour is not suitable for the CLI completion functions, because
		// when completing "gardenctl --garden foo --shoot [tab]", the DTP would consider this as
		// "user wants to target the garden" and will therefore throw away the project/seed information.
		// Project and seed information however are important for the completion functions.
		//
		// To work around this, all completion functions use a manager with a regular filesystem based
		// target provider without considering the given target flags.
		manager, err := f.WithoutTargetFlags().Manager()

		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "%v\n", err)
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		result, err := completer(f.Context(), manager, tf)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "%v\n", err)
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return util.FilterStringsByPrefix(toComplete, result), cobra.ShellCompDirectiveNoFileComp
	}
}

func getCurrentTarget(tf target.TargetFlags, manager target.Manager) (target.Target, error) {
	// any --garden flag has precedence over the config file
	if tf.GardenName() != "" {
		return target.NewTarget(tf.GardenName(), "", "", ""), nil
	}

	currentTarget, err := manager.CurrentTarget()
	if err != nil {
		return nil, fmt.Errorf("failed to read current target: %w", err)
	}

	if tf.ProjectName() != "" {
		currentTarget = currentTarget.WithProjectName(tf.ProjectName()).WithSeedName("")
	} else if tf.SeedName() != "" {
		currentTarget = currentTarget.WithSeedName(tf.SeedName()).WithProjectName("")
	}

	return currentTarget, nil
}

func gardenFlagCompletionFunc(ctx context.Context, manager target.Manager, tf target.TargetFlags) ([]string, error) {
	return util.GardenNames(manager)
}

func projectFlagCompletionFunc(ctx context.Context, manager target.Manager, tf target.TargetFlags) ([]string, error) {
	currentTarget, err := getCurrentTarget(tf, manager)
	if err != nil {
		return nil, err
	}

	return util.ProjectNamesForTarget(ctx, manager, currentTarget)
}

func seedFlagCompletionFunc(ctx context.Context, manager target.Manager, tf target.TargetFlags) ([]string, error) {
	currentTarget, err := getCurrentTarget(tf, manager)
	if err != nil {
		return nil, err
	}

	return util.SeedNamesForTarget(ctx, manager, currentTarget)
}

func shootFlagCompletionFunc(ctx context.Context, manager target.Manager, tf target.TargetFlags) ([]string, error) {
	currentTarget, err := getCurrentTarget(tf, manager)
	if err != nil {
		return nil, err
	}

	return util.ShootNamesForTarget(ctx, manager, currentTarget)
}

// addKlogFlags adds flags from k8s.io/klog
func addKlogFlags(fs *pflag.FlagSet) {
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)

	local.VisitAll(func(fl *flag.Flag) {
		fs.AddGoFlag(fl)
	})
}

func getSessionID() (string, error) {
	if value, ok := os.LookupEnv(envSessionID); ok {
		if sidRegexp.MatchString(value) {
			return value, nil
		}

		return "", fmt.Errorf("Environment variable %s must only contain alphanumeric characters, underscore and dash and have a minimum length of 1 and a maximum length of 128", envSessionID)
	}

	if value, ok := os.LookupEnv(envTermSessionID); ok {
		match := uuidRegexp.FindStringSubmatch(strings.ToLower(value))
		if len(match) > 1 {
			return match[1], nil
		}
	}

	return "", fmt.Errorf("Environment variable %s is required. Use \"gardenctl help\" for more information about the requirements of gardenctl", envSessionID)
}