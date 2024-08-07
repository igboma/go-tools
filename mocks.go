package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/mock"
)

// MockGitRepository is a mock implementation of GitRepository
type MockGitRepository struct {
	mock.Mock
}

func (m *MockGitRepository) FetchPRs() ([]PR, error) {
	args := m.Called()
	return args.Get(0).([]PR), args.Error(1)
}

func (m *MockGitRepository) FetchPRBranch(pr PR) (*git.Repository, error) {
	args := m.Called(pr)
	return args.Get(0).(*git.Repository), args.Error(1)
}

func (m *MockGitRepository) GetFileContent(repo *git.Repository, filePath string) ([]byte, error) {
	args := m.Called(repo, filePath)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockGitRepository) GetChangedFiles(pr PR) ([]string, error) {
	args := m.Called(pr)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGitRepository) UpdateCountLabel(pr PR, count int) error {
	args := m.Called(pr, count)
	return args.Error(0)
}

func (m *MockGitRepository) MergePR(pr PR) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockGitRepository) UpdateBranch(pr PR) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockGitRepository) ListLabels(prID int) ([]string, error) {
	args := m.Called(prID)
	return args.Get(0).([]string), args.Error(1)
}
