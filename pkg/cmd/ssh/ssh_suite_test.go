/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package ssh_test

import (
	"os"
	"testing"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	operationsv1alpha1 "github.com/gardener/gardener/pkg/apis/operations/v1alpha1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	utilruntime.Must(gardencorev1beta1.AddToScheme(scheme.Scheme))
	utilruntime.Must(operationsv1alpha1.AddToScheme(scheme.Scheme))
}

func TestCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SSH Command Test Suite")
}

var _ = BeforeSuite(func() {
	sessionID := uuid.NewString()
	Expect(os.Setenv("GCTL_SESSION_ID", sessionID)).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(os.Unsetenv("GCTL_SESSION_ID"))
})
