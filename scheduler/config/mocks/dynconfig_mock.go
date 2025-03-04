// Code generated by MockGen. DO NOT EDIT.
// Source: dynconfig.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	types "d7y.io/dragonfly/v2/manager/types"
	config "d7y.io/dragonfly/v2/scheduler/config"
	gomock "github.com/golang/mock/gomock"
)

// MockDynconfigInterface is a mock of DynconfigInterface interface.
type MockDynconfigInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDynconfigInterfaceMockRecorder
}

// MockDynconfigInterfaceMockRecorder is the mock recorder for MockDynconfigInterface.
type MockDynconfigInterfaceMockRecorder struct {
	mock *MockDynconfigInterface
}

// NewMockDynconfigInterface creates a new mock instance.
func NewMockDynconfigInterface(ctrl *gomock.Controller) *MockDynconfigInterface {
	mock := &MockDynconfigInterface{ctrl: ctrl}
	mock.recorder = &MockDynconfigInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynconfigInterface) EXPECT() *MockDynconfigInterfaceMockRecorder {
	return m.recorder
}

// Deregister mocks base method.
func (m *MockDynconfigInterface) Deregister(arg0 config.Observer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Deregister", arg0)
}

// Deregister indicates an expected call of Deregister.
func (mr *MockDynconfigInterfaceMockRecorder) Deregister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deregister", reflect.TypeOf((*MockDynconfigInterface)(nil).Deregister), arg0)
}

// Get mocks base method.
func (m *MockDynconfigInterface) Get() (*config.DynconfigData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*config.DynconfigData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDynconfigInterfaceMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDynconfigInterface)(nil).Get))
}

// GetCDNClusterConfig mocks base method.
func (m *MockDynconfigInterface) GetCDNClusterConfig(arg0 uint) (types.CDNClusterConfig, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCDNClusterConfig", arg0)
	ret0, _ := ret[0].(types.CDNClusterConfig)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetCDNClusterConfig indicates an expected call of GetCDNClusterConfig.
func (mr *MockDynconfigInterfaceMockRecorder) GetCDNClusterConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCDNClusterConfig", reflect.TypeOf((*MockDynconfigInterface)(nil).GetCDNClusterConfig), arg0)
}

// GetSchedulerClusterClientConfig mocks base method.
func (m *MockDynconfigInterface) GetSchedulerClusterClientConfig() (types.SchedulerClusterClientConfig, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedulerClusterClientConfig")
	ret0, _ := ret[0].(types.SchedulerClusterClientConfig)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSchedulerClusterClientConfig indicates an expected call of GetSchedulerClusterClientConfig.
func (mr *MockDynconfigInterfaceMockRecorder) GetSchedulerClusterClientConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedulerClusterClientConfig", reflect.TypeOf((*MockDynconfigInterface)(nil).GetSchedulerClusterClientConfig))
}

// GetSchedulerClusterConfig mocks base method.
func (m *MockDynconfigInterface) GetSchedulerClusterConfig() (types.SchedulerClusterConfig, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedulerClusterConfig")
	ret0, _ := ret[0].(types.SchedulerClusterConfig)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSchedulerClusterConfig indicates an expected call of GetSchedulerClusterConfig.
func (mr *MockDynconfigInterfaceMockRecorder) GetSchedulerClusterConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedulerClusterConfig", reflect.TypeOf((*MockDynconfigInterface)(nil).GetSchedulerClusterConfig))
}

// Notify mocks base method.
func (m *MockDynconfigInterface) Notify() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Notify")
	ret0, _ := ret[0].(error)
	return ret0
}

// Notify indicates an expected call of Notify.
func (mr *MockDynconfigInterfaceMockRecorder) Notify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockDynconfigInterface)(nil).Notify))
}

// Register mocks base method.
func (m *MockDynconfigInterface) Register(arg0 config.Observer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0)
}

// Register indicates an expected call of Register.
func (mr *MockDynconfigInterfaceMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockDynconfigInterface)(nil).Register), arg0)
}

// Serve mocks base method.
func (m *MockDynconfigInterface) Serve() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serve")
	ret0, _ := ret[0].(error)
	return ret0
}

// Serve indicates an expected call of Serve.
func (mr *MockDynconfigInterfaceMockRecorder) Serve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serve", reflect.TypeOf((*MockDynconfigInterface)(nil).Serve))
}

// Stop mocks base method.
func (m *MockDynconfigInterface) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockDynconfigInterfaceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockDynconfigInterface)(nil).Stop))
}

// MockObserver is a mock of Observer interface.
type MockObserver struct {
	ctrl     *gomock.Controller
	recorder *MockObserverMockRecorder
}

// MockObserverMockRecorder is the mock recorder for MockObserver.
type MockObserverMockRecorder struct {
	mock *MockObserver
}

// NewMockObserver creates a new mock instance.
func NewMockObserver(ctrl *gomock.Controller) *MockObserver {
	mock := &MockObserver{ctrl: ctrl}
	mock.recorder = &MockObserverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObserver) EXPECT() *MockObserverMockRecorder {
	return m.recorder
}

// OnNotify mocks base method.
func (m *MockObserver) OnNotify(arg0 *config.DynconfigData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnNotify", arg0)
}

// OnNotify indicates an expected call of OnNotify.
func (mr *MockObserverMockRecorder) OnNotify(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnNotify", reflect.TypeOf((*MockObserver)(nil).OnNotify), arg0)
}
