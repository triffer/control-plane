// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

// KubeconfigProvider is an autogenerated mock type for the KubeconfigProvider type
type KubeconfigProvider struct {
	mock.Mock
}

// FetchRaw provides a mock function with given fields: _a0, _a1
func (_m *KubeconfigProvider) FetchRaw(_a0 context.Context, _a1 v1beta1.Shoot) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1beta1.Shoot) ([]byte, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1beta1.Shoot) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1beta1.Shoot) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewKubeconfigProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewKubeconfigProvider creates a new instance of KubeconfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKubeconfigProvider(t mockConstructorTestingTNewKubeconfigProvider) *KubeconfigProvider {
	mock := &KubeconfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
