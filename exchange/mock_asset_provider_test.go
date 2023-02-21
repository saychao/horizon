// Code generated by mockery v1.0.1. DO NOT EDIT.

package exchange

import (
	mock "github.com/stretchr/testify/mock"
	core2 "github.com/SafeRE-IT/horizon/db2/core2"
)

// mockAssetProvider is an autogenerated mock type for the assetProvider type
type mockAssetProvider struct {
	mock.Mock
}

// GetByCode provides a mock function with given fields: code
func (_m *mockAssetProvider) GetByCode(code string) (*core2.Asset, error) {
	ret := _m.Called(code)

	var r0 *core2.Asset
	if rf, ok := ret.Get(0).(func(string) *core2.Asset); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core2.Asset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectByAssets provides a mock function with given fields: baseAssets, quoteAssets
func (_m *mockAssetProvider) SelectByAssets(baseAssets []string, quoteAssets []string) ([]core2.AssetPair, error) {
	ret := _m.Called(baseAssets, quoteAssets)

	var r0 []core2.AssetPair
	if rf, ok := ret.Get(0).(func([]string, []string) []core2.AssetPair); ok {
		r0 = rf(baseAssets, quoteAssets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core2.AssetPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, []string) error); ok {
		r1 = rf(baseAssets, quoteAssets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectByPolicy provides a mock function with given fields: policy
func (_m *mockAssetProvider) SelectByPolicy(policy uint64) ([]core2.Asset, error) {
	ret := _m.Called(policy)

	var r0 []core2.Asset
	if rf, ok := ret.Get(0).(func(uint64) []core2.Asset); ok {
		r0 = rf(policy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core2.Asset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(policy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}