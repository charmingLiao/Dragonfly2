// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/cdn/supervisor (interfaces: CDNService)

// Package progress is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	task "d7y.io/dragonfly/v2/cdn/supervisor/task"
	gomock "github.com/golang/mock/gomock"
)

// MockCDNService is a mock of CDNService interface.
type MockCDNService struct {
	ctrl     *gomock.Controller
	recorder *MockCDNServiceMockRecorder
}

// MockCDNServiceMockRecorder is the mock recorder for MockCDNService.
type MockCDNServiceMockRecorder struct {
	mock *MockCDNService
}

// NewMockCDNService creates a new mock instance.
func NewMockCDNService(ctrl *gomock.Controller) *MockCDNService {
	mock := &MockCDNService{ctrl: ctrl}
	mock.recorder = &MockCDNServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCDNService) EXPECT() *MockCDNServiceMockRecorder {
	return m.recorder
}

// GetSeedPieces mocks base method.
func (m *MockCDNService) GetSeedPieces(arg0 string) ([]*task.PieceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeedPieces", arg0)
	ret0, _ := ret[0].([]*task.PieceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSeedPieces indicates an expected call of GetSeedPieces.
func (mr *MockCDNServiceMockRecorder) GetSeedPieces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeedPieces", reflect.TypeOf((*MockCDNService)(nil).GetSeedPieces), arg0)
}

// GetSeedTask mocks base method.
func (m *MockCDNService) GetSeedTask(arg0 string) (*task.SeedTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeedTask", arg0)
	ret0, _ := ret[0].(*task.SeedTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSeedTask indicates an expected call of GetSeedTask.
func (mr *MockCDNServiceMockRecorder) GetSeedTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeedTask", reflect.TypeOf((*MockCDNService)(nil).GetSeedTask), arg0)
}

// RegisterSeedTask mocks base method.
func (m *MockCDNService) RegisterSeedTask(arg0 context.Context, arg1 string, arg2 *task.SeedTask) (*task.SeedTask, <-chan *task.PieceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterSeedTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(*task.SeedTask)
	ret1, _ := ret[1].(<-chan *task.PieceInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RegisterSeedTask indicates an expected call of RegisterSeedTask.
func (mr *MockCDNServiceMockRecorder) RegisterSeedTask(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSeedTask", reflect.TypeOf((*MockCDNService)(nil).RegisterSeedTask), arg0, arg1, arg2)
}
