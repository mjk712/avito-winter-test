// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source=internal/storage/storage.go -destination=internal/storage/mock_storage.go -package=mocks
//

// Package mocks is a generated GoMock package.
package storage

import (
	dao "avito-winter-test/internal/models/dao"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
	isgomock struct{}
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// BuyItem mocks base method.
func (m *MockStorage) BuyItem(ctx context.Context, userID, itemID, price int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyItem", ctx, userID, itemID, price)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyItem indicates an expected call of BuyItem.
func (mr *MockStorageMockRecorder) BuyItem(ctx, userID, itemID, price any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyItem", reflect.TypeOf((*MockStorage)(nil).BuyItem), ctx, userID, itemID, price)
}

// CheckUserAuth mocks base method.
func (m *MockStorage) CheckUserAuth(ctx context.Context, username string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserAuth", ctx, username)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserAuth indicates an expected call of CheckUserAuth.
func (mr *MockStorageMockRecorder) CheckUserAuth(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserAuth", reflect.TypeOf((*MockStorage)(nil).CheckUserAuth), ctx, username)
}

// CreateNewUser mocks base method.
func (m *MockStorage) CreateNewUser(ctx context.Context, username, password string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewUser", ctx, username, password)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewUser indicates an expected call of CreateNewUser.
func (mr *MockStorageMockRecorder) CreateNewUser(ctx, username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewUser", reflect.TypeOf((*MockStorage)(nil).CreateNewUser), ctx, username, password)
}

// GetMerchByName mocks base method.
func (m *MockStorage) GetMerchByName(ctx context.Context, name string) (dao.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchByName", ctx, name)
	ret0, _ := ret[0].(dao.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchByName indicates an expected call of GetMerchByName.
func (mr *MockStorageMockRecorder) GetMerchByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchByName", reflect.TypeOf((*MockStorage)(nil).GetMerchByName), ctx, name)
}

// GetUserByID mocks base method.
func (m *MockStorage) GetUserByID(ctx context.Context, userID int) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStorageMockRecorder) GetUserByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStorage)(nil).GetUserByID), ctx, userID)
}

// GetUserCoinHistory mocks base method.
func (m *MockStorage) GetUserCoinHistory(ctx context.Context, userID int) ([]dao.TransactionHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCoinHistory", ctx, userID)
	ret0, _ := ret[0].([]dao.TransactionHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCoinHistory indicates an expected call of GetUserCoinHistory.
func (mr *MockStorageMockRecorder) GetUserCoinHistory(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCoinHistory", reflect.TypeOf((*MockStorage)(nil).GetUserCoinHistory), ctx, userID)
}

// GetUserIDByUsername mocks base method.
func (m *MockStorage) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByUsername", ctx, username)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByUsername indicates an expected call of GetUserIDByUsername.
func (mr *MockStorageMockRecorder) GetUserIDByUsername(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByUsername", reflect.TypeOf((*MockStorage)(nil).GetUserIDByUsername), ctx, username)
}

// GetUserInventory mocks base method.
func (m *MockStorage) GetUserInventory(ctx context.Context, userID int) ([]dao.Inventory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInventory", ctx, userID)
	ret0, _ := ret[0].([]dao.Inventory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInventory indicates an expected call of GetUserInventory.
func (mr *MockStorageMockRecorder) GetUserInventory(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInventory", reflect.TypeOf((*MockStorage)(nil).GetUserInventory), ctx, userID)
}

// TransferCoins mocks base method.
func (m *MockStorage) TransferCoins(ctx context.Context, fromUserID, toUserID, amount int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferCoins", ctx, fromUserID, toUserID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferCoins indicates an expected call of TransferCoins.
func (mr *MockStorageMockRecorder) TransferCoins(ctx, fromUserID, toUserID, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferCoins", reflect.TypeOf((*MockStorage)(nil).TransferCoins), ctx, fromUserID, toUserID, amount)
}
