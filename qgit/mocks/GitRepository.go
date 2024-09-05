// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	qgit "gitpkg/qgit"
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// GitRepository is an autogenerated mock type for the GitRepository type
type GitRepository struct {
	mock.Mock
}

// CheckRemoteRef provides a mock function with given fields: ref
func (_m *GitRepository) CheckRemoteRef(ref string) (bool, bool, bool) {
	ret := _m.Called(ref)

	if len(ret) == 0 {
		panic("no return value specified for CheckRemoteRef")
	}

	var r0 bool
	var r1 bool
	var r2 bool
	if rf, ok := ret.Get(0).(func(string) (bool, bool, bool)); ok {
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

	return r0, r1, r2
}

// Checkout provides a mock function with given fields: ref
func (_m *GitRepository) Checkout(ref string) error {
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
func (_m *GitRepository) CheckoutBranch(branch string) error {
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
func (_m *GitRepository) CheckoutHash(hash string) error {
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
func (_m *GitRepository) CheckoutTag(tag string) error {
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
func (_m *GitRepository) Fetch(refSpecStr string) error {
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
func (_m *GitRepository) GetChangedFilesByPRNumber(prNumber int) ([]string, error) {
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

// GetChangedFilesByPRNumberFileExtMatch provides a mock function with given fields: prNumber, fileExt
func (_m *GitRepository) GetChangedFilesByPRNumberFileExtMatch(prNumber int, fileExt string) ([]string, error) {
	ret := _m.Called(prNumber, fileExt)

	if len(ret) == 0 {
		panic("no return value specified for GetChangedFilesByPRNumberFileExtMatch")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string) ([]string, error)); ok {
		return rf(prNumber, fileExt)
	}
	if rf, ok := ret.Get(0).(func(int, string) []string); ok {
		r0 = rf(prNumber, fileExt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(prNumber, fileExt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChangedFilesByPRNumberFilesByFilter provides a mock function with given fields: prNumber, functionThatImplementsFilter
func (_m *GitRepository) GetChangedFilesByPRNumberFilesByFilter(prNumber int, functionThatImplementsFilter func(string) bool) ([]string, error) {
	ret := _m.Called(prNumber, functionThatImplementsFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetChangedFilesByPRNumberFilesByFilter")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, func(string) bool) ([]string, error)); ok {
		return rf(prNumber, functionThatImplementsFilter)
	}
	if rf, ok := ret.Get(0).(func(int, func(string) bool) []string); ok {
		r0 = rf(prNumber, functionThatImplementsFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int, func(string) bool) error); ok {
		r1 = rf(prNumber, functionThatImplementsFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChangedFilesByPRNumberFilesByRegex provides a mock function with given fields: prNumber, regexFilter
func (_m *GitRepository) GetChangedFilesByPRNumberFilesByRegex(prNumber int, regexFilter string) ([]string, error) {
	ret := _m.Called(prNumber, regexFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetChangedFilesByPRNumberFilesByRegex")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string) ([]string, error)); ok {
		return rf(prNumber, regexFilter)
	}
	if rf, ok := ret.Get(0).(func(int, string) []string); ok {
		r0 = rf(prNumber, regexFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(prNumber, regexFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChangedFilesByPRNumberFilesMatching provides a mock function with given fields: prNumber, fileName
func (_m *GitRepository) GetChangedFilesByPRNumberFilesMatching(prNumber int, fileName string) ([]string, error) {
	ret := _m.Called(prNumber, fileName)

	if len(ret) == 0 {
		panic("no return value specified for GetChangedFilesByPRNumberFilesMatching")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string) ([]string, error)); ok {
		return rf(prNumber, fileName)
	}
	if rf, ok := ret.Get(0).(func(int, string) []string); ok {
		r0 = rf(prNumber, fileName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(prNumber, fileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileContentFromBranch provides a mock function with given fields: branch, file
func (_m *GitRepository) GetFileContentFromBranch(branch string, file string) (string, error) {
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
func (_m *GitRepository) GetFileContentFromCommit(commitHash string, file string) (string, error) {
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
func (_m *GitRepository) Head() (qgit.QReference, error) {
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

// PlainClone provides a mock function with given fields: o
func (_m *GitRepository) PlainClone(o qgit.QgitOptions) error {
	ret := _m.Called(o)

	if len(ret) == 0 {
		panic("no return value specified for PlainClone")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(qgit.QgitOptions) error); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PlainOpen provides a mock function with given fields: o
func (_m *GitRepository) PlainOpen(o qgit.QgitOptions) error {
	ret := _m.Called(o)

	if len(ret) == 0 {
		panic("no return value specified for PlainOpen")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(qgit.QgitOptions) error); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stat provides a mock function with given fields: path
func (_m *GitRepository) Stat(path string) (fs.FileInfo, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for Stat")
	}

	var r0 fs.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (fs.FileInfo, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) fs.FileInfo); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fs.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Worktree provides a mock function with given fields:
func (_m *GitRepository) Worktree() (qgit.GitWorktree, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Worktree")
	}

	var r0 qgit.GitWorktree
	var r1 error
	if rf, ok := ret.Get(0).(func() (qgit.GitWorktree, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() qgit.GitWorktree); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(qgit.GitWorktree)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGitRepository creates a new instance of GitRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGitRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *GitRepository {
	mock := &GitRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
