// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sunsingerus/mservice/pkg/api/mservice (interfaces: MServiceControlPlaneClient,MServiceControlPlane_DataClient)

// Package mock_mservice is a generated GoMock package.
package mock_mservice

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	mservice "github.com/sunsingerus/mservice/pkg/api/mservice"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockMServiceControlPlaneClient is a mock of MServiceControlPlaneClient interface
type MockMServiceControlPlaneClient struct {
	ctrl     *gomock.Controller
	recorder *MockMServiceControlPlaneClientMockRecorder
}

// MockMServiceControlPlaneClientMockRecorder is the mock recorder for MockMServiceControlPlaneClient
type MockMServiceControlPlaneClientMockRecorder struct {
	mock *MockMServiceControlPlaneClient
}

// NewMockMServiceControlPlaneClient creates a new mock instance
func NewMockMServiceControlPlaneClient(ctrl *gomock.Controller) *MockMServiceControlPlaneClient {
	mock := &MockMServiceControlPlaneClient{ctrl: ctrl}
	mock.recorder = &MockMServiceControlPlaneClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMServiceControlPlaneClient) EXPECT() *MockMServiceControlPlaneClientMockRecorder {
	return m.recorder
}

// Commands mocks base method
func (m *MockMServiceControlPlaneClient) Commands(arg0 context.Context, arg1 ...grpc.CallOption) (mservice.MServiceControlPlane_CommandsClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Commands", varargs...)
	ret0, _ := ret[0].(mservice.MServiceControlPlane_CommandsClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Commands indicates an expected call of Commands
func (mr *MockMServiceControlPlaneClientMockRecorder) Commands(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commands", reflect.TypeOf((*MockMServiceControlPlaneClient)(nil).Commands), varargs...)
}

// Data mocks base method
func (m *MockMServiceControlPlaneClient) Data(arg0 context.Context, arg1 ...grpc.CallOption) (mservice.MServiceControlPlane_DataClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Data", varargs...)
	ret0, _ := ret[0].(mservice.MServiceControlPlane_DataClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Data indicates an expected call of Data
func (mr *MockMServiceControlPlaneClientMockRecorder) Data(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockMServiceControlPlaneClient)(nil).Data), varargs...)
}

// Metrics mocks base method
func (m *MockMServiceControlPlaneClient) Metrics(arg0 context.Context, arg1 ...grpc.CallOption) (mservice.MServiceControlPlane_MetricsClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Metrics", varargs...)
	ret0, _ := ret[0].(mservice.MServiceControlPlane_MetricsClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Metrics indicates an expected call of Metrics
func (mr *MockMServiceControlPlaneClientMockRecorder) Metrics(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metrics", reflect.TypeOf((*MockMServiceControlPlaneClient)(nil).Metrics), varargs...)
}

// MockMServiceControlPlane_DataClient is a mock of MServiceControlPlane_DataClient interface
type MockMServiceControlPlane_DataClient struct {
	ctrl     *gomock.Controller
	recorder *MockMServiceControlPlane_DataClientMockRecorder
}

// MockMServiceControlPlane_DataClientMockRecorder is the mock recorder for MockMServiceControlPlane_DataClient
type MockMServiceControlPlane_DataClientMockRecorder struct {
	mock *MockMServiceControlPlane_DataClient
}

// NewMockMServiceControlPlane_DataClient creates a new mock instance
func NewMockMServiceControlPlane_DataClient(ctrl *gomock.Controller) *MockMServiceControlPlane_DataClient {
	mock := &MockMServiceControlPlane_DataClient{ctrl: ctrl}
	mock.recorder = &MockMServiceControlPlane_DataClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMServiceControlPlane_DataClient) EXPECT() *MockMServiceControlPlane_DataClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockMServiceControlPlane_DataClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockMServiceControlPlane_DataClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockMServiceControlPlane_DataClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockMServiceControlPlane_DataClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).Context))
}

// Header mocks base method
func (m *MockMServiceControlPlane_DataClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockMServiceControlPlane_DataClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).Header))
}

// Recv mocks base method
func (m *MockMServiceControlPlane_DataClient) Recv() (*mservice.DataChunk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*mservice.DataChunk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockMServiceControlPlane_DataClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockMServiceControlPlane_DataClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockMServiceControlPlane_DataClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockMServiceControlPlane_DataClient) Send(arg0 *mservice.DataChunk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockMServiceControlPlane_DataClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).Send), arg0)
}

// SendMsg mocks base method
func (m *MockMServiceControlPlane_DataClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockMServiceControlPlane_DataClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockMServiceControlPlane_DataClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockMServiceControlPlane_DataClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockMServiceControlPlane_DataClient)(nil).Trailer))
}
