// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Manan007224/sigo/pkg/proto (interfaces: SchedulerClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	protobuf "github.com/Manan007224/sigo/pkg/proto"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

// MockSchedulerClient is a mock of SchedulerClient interface
type MockSchedulerClient struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerClientMockRecorder
}

// MockSchedulerClientMockRecorder is the mock recorder for MockSchedulerClient
type MockSchedulerClientMockRecorder struct {
	mock *MockSchedulerClient
}

// NewMockSchedulerClient creates a new mock instance
func NewMockSchedulerClient(ctrl *gomock.Controller) *MockSchedulerClient {
	mock := &MockSchedulerClient{ctrl: ctrl}
	mock.recorder = &MockSchedulerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSchedulerClient) EXPECT() *MockSchedulerClientMockRecorder {
	return m.recorder
}

// Acknowledge mocks base method
func (m *MockSchedulerClient) Acknowledge(arg0 context.Context, arg1 *protobuf.JobPayload, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Acknowledge", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Acknowledge indicates an expected call of Acknowledge
func (mr *MockSchedulerClientMockRecorder) Acknowledge(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acknowledge", reflect.TypeOf((*MockSchedulerClient)(nil).Acknowledge), varargs...)
}

// BroadCast mocks base method
func (m *MockSchedulerClient) BroadCast(arg0 context.Context, arg1 *protobuf.JobPayload, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BroadCast", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BroadCast indicates an expected call of BroadCast
func (mr *MockSchedulerClientMockRecorder) BroadCast(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadCast", reflect.TypeOf((*MockSchedulerClient)(nil).BroadCast), varargs...)
}

// Discover mocks base method
func (m *MockSchedulerClient) Discover(arg0 context.Context, arg1 *protobuf.ClientConfig, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Discover", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Discover indicates an expected call of Discover
func (mr *MockSchedulerClientMockRecorder) Discover(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discover", reflect.TypeOf((*MockSchedulerClient)(nil).Discover), varargs...)
}

// Fail mocks base method
func (m *MockSchedulerClient) Fail(arg0 context.Context, arg1 *protobuf.FailPayload, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Fail", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fail indicates an expected call of Fail
func (mr *MockSchedulerClientMockRecorder) Fail(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockSchedulerClient)(nil).Fail), varargs...)
}

// Fetch mocks base method
func (m *MockSchedulerClient) Fetch(arg0 context.Context, arg1 *protobuf.Queue, arg2 ...grpc.CallOption) (*protobuf.JobPayload, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Fetch", varargs...)
	ret0, _ := ret[0].(*protobuf.JobPayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockSchedulerClientMockRecorder) Fetch(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockSchedulerClient)(nil).Fetch), varargs...)
}

// HeartBeat mocks base method
func (m *MockSchedulerClient) HeartBeat(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HeartBeat", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeartBeat indicates an expected call of HeartBeat
func (mr *MockSchedulerClientMockRecorder) HeartBeat(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeartBeat", reflect.TypeOf((*MockSchedulerClient)(nil).HeartBeat), varargs...)
}
