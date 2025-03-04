// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/rpc/manager/client/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	manager "d7y.io/dragonfly/v2/pkg/rpc/manager"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// GetScheduler mocks base method.
func (m *MockClient) GetScheduler(arg0 *manager.GetSchedulerRequest) (*manager.Scheduler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScheduler", arg0)
	ret0, _ := ret[0].(*manager.Scheduler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScheduler indicates an expected call of GetScheduler.
func (mr *MockClientMockRecorder) GetScheduler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScheduler", reflect.TypeOf((*MockClient)(nil).GetScheduler), arg0)
}

// KeepAlive mocks base method.
func (m *MockClient) KeepAlive(arg0 time.Duration, arg1 *manager.KeepAliveRequest) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "KeepAlive", arg0, arg1)
}

// KeepAlive indicates an expected call of KeepAlive.
func (mr *MockClientMockRecorder) KeepAlive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepAlive", reflect.TypeOf((*MockClient)(nil).KeepAlive), arg0, arg1)
}

// ListSchedulers mocks base method.
func (m *MockClient) ListSchedulers(arg0 *manager.ListSchedulersRequest) (*manager.ListSchedulersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSchedulers", arg0)
	ret0, _ := ret[0].(*manager.ListSchedulersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSchedulers indicates an expected call of ListSchedulers.
func (mr *MockClientMockRecorder) ListSchedulers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSchedulers", reflect.TypeOf((*MockClient)(nil).ListSchedulers), arg0)
}

// UpdateCDN mocks base method.
func (m *MockClient) UpdateCDN(arg0 *manager.UpdateCDNRequest) (*manager.CDN, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCDN", arg0)
	ret0, _ := ret[0].(*manager.CDN)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCDN indicates an expected call of UpdateCDN.
func (mr *MockClientMockRecorder) UpdateCDN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCDN", reflect.TypeOf((*MockClient)(nil).UpdateCDN), arg0)
}

// UpdateScheduler mocks base method.
func (m *MockClient) UpdateScheduler(arg0 *manager.UpdateSchedulerRequest) (*manager.Scheduler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScheduler", arg0)
	ret0, _ := ret[0].(*manager.Scheduler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateScheduler indicates an expected call of UpdateScheduler.
func (mr *MockClientMockRecorder) UpdateScheduler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScheduler", reflect.TypeOf((*MockClient)(nil).UpdateScheduler), arg0)
}
