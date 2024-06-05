// Code generated by mockery. DO NOT EDIT.

package role_mock

import (
	role "github.com/omareloui/odinls/internal/application/core/role"
	mock "github.com/stretchr/testify/mock"
)

// MockRoleRepository is an autogenerated mock type for the RoleRepository type
type MockRoleRepository struct {
	mock.Mock
}

// CreateRole provides a mock function with given fields: roles
func (_m *MockRoleRepository) CreateRole(roles *role.Role) error {
	ret := _m.Called(roles)

	if len(ret) == 0 {
		panic("no return value specified for CreateRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*role.Role) error); ok {
		r0 = rf(roles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindRole provides a mock function with given fields: id
func (_m *MockRoleRepository) FindRole(id string) (*role.Role, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindRole")
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

// GetRoles provides a mock function with given fields:
func (_m *MockRoleRepository) GetRoles() ([]role.Role, error) {
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

// SeedRoles provides a mock function with given fields: roles
func (_m *MockRoleRepository) SeedRoles(roles []string) error {
	ret := _m.Called(roles)

	if len(ret) == 0 {
		panic("no return value specified for SeedRoles")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(roles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRoleRepository creates a new instance of MockRoleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRoleRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRoleRepository {
	mock := &MockRoleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
