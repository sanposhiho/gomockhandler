// Code generated by MockGen. DO NOT EDIT.
// Source: input.go

// Package const_length is a generated GoMock package.
package const_length

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockI is a mock of I interface.
type MockI struct {
	ctrl     *gomock.Controller
	recorder *MockIMockRecorder
}

// MockIMockRecorder is the mock recorder for MockI.
type MockIMockRecorder struct {
	mock *MockI
}

// NewMockI creates a new mock instance.
func NewMockI(ctrl *gomock.Controller) *MockI {
	mock := &MockI{ctrl: ctrl}
	mock.recorder = &MockIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockI) EXPECT() *MockIMockRecorder {
	return m.recorder
}

// Bar mocks base method.
func (m *MockI) Bar() [2]int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bar")
	ret0, _ := ret[0].([2]int)
	return ret0
}

// Bar indicates an expected call of Bar.
func (mr *MockIMockRecorder) Bar() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bar", reflect.TypeOf((*MockI)(nil).Bar))
}

// Baz mocks base method.
func (m *MockI) Baz() [127]int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Baz")
	ret0, _ := ret[0].([127]int)
	return ret0
}

// Baz indicates an expected call of Baz.
func (mr *MockIMockRecorder) Baz() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Baz", reflect.TypeOf((*MockI)(nil).Baz))
}

// Foo mocks base method.
func (m *MockI) Foo() [2]int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Foo")
	ret0, _ := ret[0].([2]int)
	return ret0
}

// Foo indicates an expected call of Foo.
func (mr *MockIMockRecorder) Foo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Foo", reflect.TypeOf((*MockI)(nil).Foo))
}
