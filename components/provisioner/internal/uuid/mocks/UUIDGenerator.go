// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UUIDGenerator is an autogenerated mock type for the UUIDGenerator type
type UUIDGenerator struct {
	mock.Mock
}

// New provides a mock function with given fields:
func (_m *UUIDGenerator) New() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewUUIDGenerator creates a new instance of UUIDGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUUIDGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *UUIDGenerator {
	mock := &UUIDGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
