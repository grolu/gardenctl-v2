// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/gardenctl-v2/pkg/target (interfaces: Manager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	gardenclient "github.com/gardener/gardenctl-v2/internal/gardenclient"
	config "github.com/gardener/gardenctl-v2/pkg/config"
	target "github.com/gardener/gardenctl-v2/pkg/target"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// ClientConfig mocks base method.
func (m *MockManager) ClientConfig(arg0 context.Context, arg1 target.Target) (clientcmd.ClientConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientConfig", arg0, arg1)
	ret0, _ := ret[0].(clientcmd.ClientConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientConfig indicates an expected call of ClientConfig.
func (mr *MockManagerMockRecorder) ClientConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientConfig", reflect.TypeOf((*MockManager)(nil).ClientConfig), arg0, arg1)
}

// Configuration mocks base method.
func (m *MockManager) Configuration() *config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configuration")
	ret0, _ := ret[0].(*config.Config)
	return ret0
}

// Configuration indicates an expected call of Configuration.
func (mr *MockManagerMockRecorder) Configuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configuration", reflect.TypeOf((*MockManager)(nil).Configuration))
}

// CurrentTarget mocks base method.
func (m *MockManager) CurrentTarget() (target.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentTarget")
	ret0, _ := ret[0].(target.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentTarget indicates an expected call of CurrentTarget.
func (mr *MockManagerMockRecorder) CurrentTarget() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentTarget", reflect.TypeOf((*MockManager)(nil).CurrentTarget))
}

// GardenClient mocks base method.
func (m *MockManager) GardenClient(arg0 string) (gardenclient.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GardenClient", arg0)
	ret0, _ := ret[0].(gardenclient.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GardenClient indicates an expected call of GardenClient.
func (mr *MockManagerMockRecorder) GardenClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GardenClient", reflect.TypeOf((*MockManager)(nil).GardenClient), arg0)
}

// GardenNames mocks base method.
func (m *MockManager) GardenNames() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GardenNames")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GardenNames indicates an expected call of GardenNames.
func (mr *MockManagerMockRecorder) GardenNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GardenNames", reflect.TypeOf((*MockManager)(nil).GardenNames))
}

// ProjectNames mocks base method.
func (m *MockManager) ProjectNames(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProjectNames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProjectNames indicates an expected call of ProjectNames.
func (mr *MockManagerMockRecorder) ProjectNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProjectNames", reflect.TypeOf((*MockManager)(nil).ProjectNames), arg0)
}

// SeedClient mocks base method.
func (m *MockManager) SeedClient(arg0 context.Context, arg1 target.Target) (client.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedClient", arg0, arg1)
	ret0, _ := ret[0].(client.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedClient indicates an expected call of SeedClient.
func (mr *MockManagerMockRecorder) SeedClient(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedClient", reflect.TypeOf((*MockManager)(nil).SeedClient), arg0, arg1)
}

// SeedNames mocks base method.
func (m *MockManager) SeedNames(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeedNames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SeedNames indicates an expected call of SeedNames.
func (mr *MockManagerMockRecorder) SeedNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeedNames", reflect.TypeOf((*MockManager)(nil).SeedNames), arg0)
}

// SessionDir mocks base method.
func (m *MockManager) SessionDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SessionDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// SessionDir indicates an expected call of SessionDir.
func (mr *MockManagerMockRecorder) SessionDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SessionDir", reflect.TypeOf((*MockManager)(nil).SessionDir))
}

// ShootClient mocks base method.
func (m *MockManager) ShootClient(arg0 context.Context, arg1 target.Target) (client.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShootClient", arg0, arg1)
	ret0, _ := ret[0].(client.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShootClient indicates an expected call of ShootClient.
func (mr *MockManagerMockRecorder) ShootClient(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShootClient", reflect.TypeOf((*MockManager)(nil).ShootClient), arg0, arg1)
}

// ShootNames mocks base method.
func (m *MockManager) ShootNames(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShootNames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShootNames indicates an expected call of ShootNames.
func (mr *MockManagerMockRecorder) ShootNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShootNames", reflect.TypeOf((*MockManager)(nil).ShootNames), arg0)
}

// TargetControlPlane mocks base method.
func (m *MockManager) TargetControlPlane(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetControlPlane", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetControlPlane indicates an expected call of TargetControlPlane.
func (mr *MockManagerMockRecorder) TargetControlPlane(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetControlPlane", reflect.TypeOf((*MockManager)(nil).TargetControlPlane), arg0)
}

// TargetFlags mocks base method.
func (m *MockManager) TargetFlags() target.TargetFlags {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetFlags")
	ret0, _ := ret[0].(target.TargetFlags)
	return ret0
}

// TargetFlags indicates an expected call of TargetFlags.
func (mr *MockManagerMockRecorder) TargetFlags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetFlags", reflect.TypeOf((*MockManager)(nil).TargetFlags))
}

// TargetGarden mocks base method.
func (m *MockManager) TargetGarden(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetGarden", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetGarden indicates an expected call of TargetGarden.
func (mr *MockManagerMockRecorder) TargetGarden(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetGarden", reflect.TypeOf((*MockManager)(nil).TargetGarden), arg0, arg1)
}

// TargetMatchPattern mocks base method.
func (m *MockManager) TargetMatchPattern(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetMatchPattern", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetMatchPattern indicates an expected call of TargetMatchPattern.
func (mr *MockManagerMockRecorder) TargetMatchPattern(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetMatchPattern", reflect.TypeOf((*MockManager)(nil).TargetMatchPattern), arg0, arg1)
}

// TargetProject mocks base method.
func (m *MockManager) TargetProject(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetProject", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetProject indicates an expected call of TargetProject.
func (mr *MockManagerMockRecorder) TargetProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetProject", reflect.TypeOf((*MockManager)(nil).TargetProject), arg0, arg1)
}

// TargetSeed mocks base method.
func (m *MockManager) TargetSeed(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetSeed", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetSeed indicates an expected call of TargetSeed.
func (mr *MockManagerMockRecorder) TargetSeed(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetSeed", reflect.TypeOf((*MockManager)(nil).TargetSeed), arg0, arg1)
}

// TargetShoot mocks base method.
func (m *MockManager) TargetShoot(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetShoot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TargetShoot indicates an expected call of TargetShoot.
func (mr *MockManagerMockRecorder) TargetShoot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetShoot", reflect.TypeOf((*MockManager)(nil).TargetShoot), arg0, arg1)
}

// UnsetTargetControlPlane mocks base method.
func (m *MockManager) UnsetTargetControlPlane(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsetTargetControlPlane", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsetTargetControlPlane indicates an expected call of UnsetTargetControlPlane.
func (mr *MockManagerMockRecorder) UnsetTargetControlPlane(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetTargetControlPlane", reflect.TypeOf((*MockManager)(nil).UnsetTargetControlPlane), arg0)
}

// UnsetTargetGarden mocks base method.
func (m *MockManager) UnsetTargetGarden(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsetTargetGarden", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnsetTargetGarden indicates an expected call of UnsetTargetGarden.
func (mr *MockManagerMockRecorder) UnsetTargetGarden(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetTargetGarden", reflect.TypeOf((*MockManager)(nil).UnsetTargetGarden), arg0)
}

// UnsetTargetProject mocks base method.
func (m *MockManager) UnsetTargetProject(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsetTargetProject", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnsetTargetProject indicates an expected call of UnsetTargetProject.
func (mr *MockManagerMockRecorder) UnsetTargetProject(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetTargetProject", reflect.TypeOf((*MockManager)(nil).UnsetTargetProject), arg0)
}

// UnsetTargetSeed mocks base method.
func (m *MockManager) UnsetTargetSeed(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsetTargetSeed", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnsetTargetSeed indicates an expected call of UnsetTargetSeed.
func (mr *MockManagerMockRecorder) UnsetTargetSeed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetTargetSeed", reflect.TypeOf((*MockManager)(nil).UnsetTargetSeed), arg0)
}

// UnsetTargetShoot mocks base method.
func (m *MockManager) UnsetTargetShoot(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsetTargetShoot", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnsetTargetShoot indicates an expected call of UnsetTargetShoot.
func (mr *MockManagerMockRecorder) UnsetTargetShoot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetTargetShoot", reflect.TypeOf((*MockManager)(nil).UnsetTargetShoot), arg0)
}

// WriteClientConfig mocks base method.
func (m *MockManager) WriteClientConfig(arg0 clientcmd.ClientConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteClientConfig", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteClientConfig indicates an expected call of WriteClientConfig.
func (mr *MockManagerMockRecorder) WriteClientConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteClientConfig", reflect.TypeOf((*MockManager)(nil).WriteClientConfig), arg0)
}
