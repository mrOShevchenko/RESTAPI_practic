// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import mock "github.com/stretchr/testify/mock"

// MockGenerator is an autogenerated mock type for the Generator type
type MockGenerator struct {
	mock.Mock
}

type MockGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGenerator) EXPECT() *MockGenerator_Expecter {
	return &MockGenerator_Expecter{mock: &_m.Mock}
}

// GeneratePasswordHash provides a mock function with given fields: password
func (_m *MockGenerator) GeneratePasswordHash(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(password)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerator_GeneratePasswordHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GeneratePasswordHash'
type MockGenerator_GeneratePasswordHash_Call struct {
	*mock.Call
}

// GeneratePasswordHash is a helper method to define mock.On call
//   - password string
func (_e *MockGenerator_Expecter) GeneratePasswordHash(password interface{}) *MockGenerator_GeneratePasswordHash_Call {
	return &MockGenerator_GeneratePasswordHash_Call{Call: _e.mock.On("GeneratePasswordHash", password)}
}

func (_c *MockGenerator_GeneratePasswordHash_Call) Run(run func(password string)) *MockGenerator_GeneratePasswordHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockGenerator_GeneratePasswordHash_Call) Return(_a0 string, _a1 error) *MockGenerator_GeneratePasswordHash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerator_GeneratePasswordHash_Call) RunAndReturn(run func(string) (string, error)) *MockGenerator_GeneratePasswordHash_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGenerator creates a new instance of MockGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGenerator {
	mock := &MockGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
