// Code generated by mockery v1.0.0. DO NOT EDIT.

package operations

import mock "github.com/stretchr/testify/mock"

// mockOperationIDProvider is an autogenerated mock type for the operationIDProvider type
type mockOperationIDProvider struct {
	mock.Mock
}

// GetOperationID provides a mock function with given fields:
func (_m *mockOperationIDProvider) GetOperationID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}
