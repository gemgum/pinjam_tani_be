// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	users "projectBE23/internal/features/users"

	mock "github.com/stretchr/testify/mock"
)

// ServiceUserInterface is an autogenerated mock type for the ServiceUserInterface type
type ServiceUserInterface struct {
	mock.Mock
}

// DeleteAccount provides a mock function with given fields: userid
func (_m *ServiceUserInterface) DeleteAccount(userid uint) error {
	ret := _m.Called(userid)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(userid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProfile provides a mock function with given fields: userid
func (_m *ServiceUserInterface) GetProfile(userid uint) (*users.User, error) {
	ret := _m.Called(userid)

	if len(ret) == 0 {
		panic("no return value specified for GetProfile")
	}

	var r0 *users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*users.User, error)); ok {
		return rf(userid)
	}
	if rf, ok := ret.Get(0).(func(uint) *users.User); ok {
		r0 = rf(userid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginAccount provides a mock function with given fields: email, password
func (_m *ServiceUserInterface) LoginAccount(email string, password string) (*users.User, string, error) {
	ret := _m.Called(email, password)

	if len(ret) == 0 {
		panic("no return value specified for LoginAccount")
	}

	var r0 *users.User
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string) (*users.User, string, error)); ok {
		return rf(email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) *users.User); ok {
		r0 = rf(email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) string); ok {
		r1 = rf(email, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(email, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// LogoutAccount provides a mock function with given fields: token
func (_m *ServiceUserInterface) LogoutAccount(token string) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for LogoutAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegistrasiAccount provides a mock function with given fields: accounts
func (_m *ServiceUserInterface) RegistrasiAccount(accounts users.User) (uint, error) {
	ret := _m.Called(accounts)

	if len(ret) == 0 {
		panic("no return value specified for RegistrasiAccount")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(users.User) (uint, error)); ok {
		return rf(accounts)
	}
	if rf, ok := ret.Get(0).(func(users.User) uint); ok {
		r0 = rf(accounts)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(users.User) error); ok {
		r1 = rf(accounts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfile provides a mock function with given fields: userid, accounts
func (_m *ServiceUserInterface) UpdateProfile(userid uint, accounts users.User) error {
	ret := _m.Called(userid, accounts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, users.User) error); ok {
		r0 = rf(userid, accounts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServiceUserInterface creates a new instance of ServiceUserInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceUserInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceUserInterface {
	mock := &ServiceUserInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}