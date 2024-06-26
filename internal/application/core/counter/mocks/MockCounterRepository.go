// Code generated by mockery. DO NOT EDIT.

package counter_mock

import (
	counter "github.com/omareloui/odinls/internal/application/core/counter"
	mock "github.com/stretchr/testify/mock"
)

// MockCounterRepository is an autogenerated mock type for the CounterRepository type
type MockCounterRepository struct {
	mock.Mock
}

// AddOneToOrder provides a mock function with given fields: merchantId
func (_m *MockCounterRepository) AddOneToOrder(merchantId string) (uint, error) {
	ret := _m.Called(merchantId)

	if len(ret) == 0 {
		panic("no return value specified for AddOneToOrder")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (uint, error)); ok {
		return rf(merchantId)
	}
	if rf, ok := ret.Get(0).(func(string) uint); ok {
		r0 = rf(merchantId)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(merchantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddOneToProduct provides a mock function with given fields: merchantId, category
func (_m *MockCounterRepository) AddOneToProduct(merchantId string, category string) (uint8, error) {
	ret := _m.Called(merchantId, category)

	if len(ret) == 0 {
		panic("no return value specified for AddOneToProduct")
	}

	var r0 uint8
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (uint8, error)); ok {
		return rf(merchantId, category)
	}
	if rf, ok := ret.Get(0).(func(string, string) uint8); ok {
		r0 = rf(merchantId, category)
	} else {
		r0 = ret.Get(0).(uint8)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(merchantId, category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCounter provides a mock function with given fields: _a0
func (_m *MockCounterRepository) CreateCounter(_a0 *counter.Counter) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateCounter")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*counter.Counter) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCounterByID provides a mock function with given fields: id
func (_m *MockCounterRepository) GetCounterByID(id string) (*counter.Counter, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetCounterByID")
	}

	var r0 *counter.Counter
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*counter.Counter, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *counter.Counter); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*counter.Counter)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCounterByMerchantID provides a mock function with given fields: merchantId
func (_m *MockCounterRepository) GetCounterByMerchantID(merchantId string) (*counter.Counter, error) {
	ret := _m.Called(merchantId)

	if len(ret) == 0 {
		panic("no return value specified for GetCounterByMerchantID")
	}

	var r0 *counter.Counter
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*counter.Counter, error)); ok {
		return rf(merchantId)
	}
	if rf, ok := ret.Get(0).(func(string) *counter.Counter); ok {
		r0 = rf(merchantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*counter.Counter)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(merchantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockCounterRepository creates a new instance of MockCounterRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCounterRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCounterRepository {
	mock := &MockCounterRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
