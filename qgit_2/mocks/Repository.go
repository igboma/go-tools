// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	qgit "gitpkg/qgit_2"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CheckRemoteRef provides a mock function with given fields: ref
func (_m *Repository) CheckRemoteRef(ref string) (bool, bool, bool, error) {
	ret := _m.Called(ref)

	if len(ret) == 0 {
		panic("no return value specified for CheckRemoteRef")
	}

	var r0 bool
	var r1 bool
	var r2 bool
	var r3 error
	if rf, ok := ret.Get(0).(func(string) (bool, bool, bool, error)); ok {
		return rf(ref)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(ref)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(string) bool); ok {
		r2 = rf(ref)
	} else {
		r2 = ret.Get(2).(bool)
	}

	if rf, ok := ret.Get(3).(func(string) error); ok {
		r3 = rf(ref)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Checkout provides a mock function with given fields: ref
func (_m *Repository) Checkout(ref string) error {
	ret := _m.Called(ref)

	if len(ret) == 0 {
		panic("no return value specified for Checkout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckoutBranch provides a mock function with given fields: branch
func (_m *Repository) CheckoutBranch(branch string) error {
	ret := _m.Called(branch)

	if len(ret) == 0 {
		panic("no return value specified for CheckoutBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(branch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckoutHash provides a mock function with given fields: hash
func (_m *Repository) CheckoutHash(hash string) error {
	ret := _m.Called(hash)

	if len(ret) == 0 {
		panic("no return value specified for CheckoutHash")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(hash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckoutTag provides a mock function with given fields: tag
func (_m *Repository) CheckoutTag(tag string) error {
	ret := _m.Called(tag)

	if len(ret) == 0 {
		panic("no return value specified for CheckoutTag")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tag)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: refSpecStr
func (_m *Repository) Fetch(refSpecStr string) error {
	ret := _m.Called(refSpecStr)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(refSpecStr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetChangedFilesByPRNumber provides a mock function with given fields: prNumber
func (_m *Repository) GetChangedFilesByPRNumber(prNumber int) ([]string, error) {
	ret := _m.Called(prNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetChangedFilesByPRNumber")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]string, error)); ok {
		return rf(prNumber)
	}
	if rf, ok := ret.Get(0).(func(int) []string); ok {
		r0 = rf(prNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(prNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileContentFromBranch provides a mock function with given fields: branch, file
func (_m *Repository) GetFileContentFromBranch(branch string, file string) (string, error) {
	ret := _m.Called(branch, file)

	if len(ret) == 0 {
		panic("no return value specified for GetFileContentFromBranch")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(branch, file)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(branch, file)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(branch, file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileContentFromCommit provides a mock function with given fields: commitHash, file
func (_m *Repository) GetFileContentFromCommit(commitHash string, file string) (string, error) {
	ret := _m.Called(commitHash, file)

	if len(ret) == 0 {
		panic("no return value specified for GetFileContentFromCommit")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(commitHash, file)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(commitHash, file)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(commitHash, file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Head provides a mock function with given fields:
func (_m *Repository) Head() (qgit.QReference, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Head")
	}

	var r0 qgit.QReference
	var r1 error
	if rf, ok := ret.Get(0).(func() (qgit.QReference, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() qgit.QReference); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(qgit.QReference)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Option provides a mock function with given fields:
func (_m *Repository) Option() *qgit.QRepoOptions {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Option")
	}

	var r0 *qgit.QRepoOptions
	if rf, ok := ret.Get(0).(func() *qgit.QRepoOptions); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*qgit.QRepoOptions)
		}
	}

	return r0
}

// PlainClone provides a mock function with given fields: o
func (_m *Repository) PlainClone(o qgit.QRepoOptions) error {
	ret := _m.Called(o)

	if len(ret) == 0 {
		panic("no return value specified for PlainClone")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(qgit.QRepoOptions) error); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PlainOpen provides a mock function with given fields: o
func (_m *Repository) PlainOpen(o qgit.QRepoOptions) error {
	ret := _m.Called(o)

	if len(ret) == 0 {
		panic("no return value specified for PlainOpen")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(qgit.QRepoOptions) error); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetOption provides a mock function with given fields: option
func (_m *Repository) SetOption(option *qgit.QRepoOptions) {
	_m.Called(option)
}

// Worktree provides a mock function with given fields:
func (_m *Repository) Worktree() (qgit.Worktree, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Worktree")
	}

	var r0 qgit.Worktree
	var r1 error
	if rf, ok := ret.Get(0).(func() (qgit.Worktree, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() qgit.Worktree); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(qgit.Worktree)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}