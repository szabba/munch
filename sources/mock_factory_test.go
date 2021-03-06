// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/szabba/munch/sources (interfaces: InputFactory,ParserFactory)

// Package sources_test is a generated GoMock package.
package sources_test

import (
	json "encoding/json"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockInputFactory is a mock of InputFactory interface
type MockInputFactory struct {
	ctrl     *gomock.Controller
	recorder *MockInputFactoryMockRecorder
}

// MockInputFactoryMockRecorder is the mock recorder for MockInputFactory
type MockInputFactoryMockRecorder struct {
	mock *MockInputFactory
}

// NewMockInputFactory creates a new mock instance
func NewMockInputFactory(ctrl *gomock.Controller) *MockInputFactory {
	mock := &MockInputFactory{ctrl: ctrl}
	mock.recorder = &MockInputFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInputFactory) EXPECT() *MockInputFactoryMockRecorder {
	return m.recorder
}

// NewInput mocks base method
func (m *MockInputFactory) NewInput(arg0 json.RawMessage) (io.ReadCloser, error) {
	ret := m.ctrl.Call(m, "NewInput", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewInput indicates an expected call of NewInput
func (mr *MockInputFactoryMockRecorder) NewInput(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewInput", reflect.TypeOf((*MockInputFactory)(nil).NewInput), arg0)
}

// MockParserFactory is a mock of ParserFactory interface
type MockParserFactory struct {
	ctrl     *gomock.Controller
	recorder *MockParserFactoryMockRecorder
}

// MockParserFactoryMockRecorder is the mock recorder for MockParserFactory
type MockParserFactoryMockRecorder struct {
	mock *MockParserFactory
}

// NewMockParserFactory creates a new mock instance
func NewMockParserFactory(ctrl *gomock.Controller) *MockParserFactory {
	mock := &MockParserFactory{ctrl: ctrl}
	mock.recorder = &MockParserFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockParserFactory) EXPECT() *MockParserFactoryMockRecorder {
	return m.recorder
}

// NewParser mocks base method
func (m *MockParserFactory) NewParser(arg0 json.RawMessage) (io.WriteCloser, error) {
	ret := m.ctrl.Call(m, "NewParser", arg0)
	ret0, _ := ret[0].(io.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewParser indicates an expected call of NewParser
func (mr *MockParserFactoryMockRecorder) NewParser(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewParser", reflect.TypeOf((*MockParserFactory)(nil).NewParser), arg0)
}
