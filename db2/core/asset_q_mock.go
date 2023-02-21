package core

import "github.com/stretchr/testify/mock"

type AssetQMock struct {
	mock.Mock
}

func (q *AssetQMock) ByCode(assetCode string) (*Asset, error) {
	args := q.Called(assetCode)
	return args.Get(0).(*Asset), args.Error(1)
}

func (q *AssetQMock) ForOwner(owner string) AssetQI {
	args := q.Called(owner)
	return args.Get(0).(AssetQI)
}

func (q *AssetQMock) ForCodes(codes []string) AssetQI {
	args := q.Called(codes)
	return args.Get(0).(AssetQI)
}

func (q *AssetQMock) ForPolicy(policy uint32) AssetQI {
	args := q.Called(policy)
	return args.Get(0).(AssetQI)
}

func (q *AssetQMock) Select() ([]Asset, error) {
	args := q.Called()
	return args.Get(0).([]Asset), args.Error(1)
}
