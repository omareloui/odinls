// Code generated by mockery. DO NOT EDIT.

package role_mock

import (
	role "github.com/omareloui/odinls/internal/application/core/role"
	mock "github.com/stretchr/testify/mock"
)

// MockRoleService is an autogenerated mock type for the RoleService type
type MockRoleService struct {
	mock.Mock
}

// CreateRole provides a mock function with given fields: _a0
func (_m *MockRoleService) CreateRole(_a0 *role.Role) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*role.Role) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRoleByID provides a mock function with given fields: id
func (_m *MockRoleService) GetRoleByID(id string) (*role.Role, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetRoleByID")
	}

	var r0 *role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*role.Role, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *role.Role); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*role.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoleByName provides a mock function with given fields: _a0
func (_m *MockRoleService) GetRoleByName(_a0 string) (*role.RoleEnum, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetRoleByName")
	}

	var r0 *role.RoleEnum
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*role.RoleEnum, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *role.RoleEnum); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*role.RoleEnum)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoles provides a mock function with given fields:
func (_m *MockRoleService) GetRoles() ([]role.Role, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRoles")
	}

	var r0 []role.Role
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]role.Role, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []role.Role); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]role.Role)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SeedRoles provides a mock function with given fields:
func (_m *MockRoleService) SeedRoles() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SeedRoles")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRoleService creates a new instance of MockRoleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRoleService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRoleService {
	mock := &MockRoleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}