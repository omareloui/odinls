// Code generated by mockery. DO NOT EDIT.

package merchant_mock

import (
	merchant "github.com/omareloui/odinls/internal/application/core/merchant"
	mock "github.com/stretchr/testify/mock"
)

// MockMerchantRepository is an autogenerated mock type for the MerchantRepository type
type MockMerchantRepository struct {
	mock.Mock
}

// CreateMerchant provides a mock function with given fields: _a0
func (_m *MockMerchantRepository) CreateMerchant(_a0 *merchant.Merchant) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateMerchant")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*merchant.Merchant) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindMerchant provides a mock function with given fields: id
func (_m *MockMerchantRepository) FindMerchant(id string) (*merchant.Merchant, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindMerchant")
	}

	var r0 *merchant.Merchant
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*merchant.Merchant, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *merchant.Merchant); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*merchant.Merchant)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMerchants provides a mock function with given fields:
func (_m *MockMerchantRepository) GetMerchants() ([]merchant.Merchant, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMerchants")
	}

	var r0 []merchant.Merchant
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]merchant.Merchant, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []merchant.Merchant); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]merchant.Merchant)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMerchantByID provides a mock function with given fields: id, _a1
func (_m *MockMerchantRepository) UpdateMerchantByID(id string, _a1 *merchant.Merchant) error {
	ret := _m.Called(id, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMerchantByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *merchant.Merchant) error); ok {
		r0 = rf(id, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockMerchantRepository creates a new instance of MockMerchantRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMerchantRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMerchantRepository {
	mock := &MockMerchantRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}