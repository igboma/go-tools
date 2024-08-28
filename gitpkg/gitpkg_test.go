package gitpkg_test

import (
	"errors"
	"fmt"
	"os"
	"qgit/gitpkg"
	"qgit/gitpkg/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockFileInfo implements os.FileInfo interface for testing purposes.
type mockFileInfo struct{}

func (m mockFileInfo) Name() string       { return "mock" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return true }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestInit(t *testing.T) {
	t.Run("Repository does not exist locally; cloning succeeds", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository1)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Mock the Stat and PlainClone methods
		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})

		// Also mock the Init method
		//mockRepo.On("Init", options).Return(nil)

		// Initialize Qgit instance
		_, err := gitpkg.NewQGit(options, mockRepo)
		assert.NoError(t, err)

		// Verify that methods were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository exists locally; opening succeeds", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository1)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Mock Stat and PlainOpen methods
		mockRepo.On("Stat", "/test/repo/.git").Return(&mockFileInfo{}, nil)
		mockRepo.On("PlainOpen", options).Return(nil)

		// Initialize Qgit instance
		_, err := gitpkg.NewQGit(options, mockRepo)
		assert.NoError(t, err)

		// Verify that methods were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository exists locally; opening fails", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository1)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Mock Stat and PlainOpen methods
		mockRepo.On("Stat", "/test/repo/.git").Return(&mockFileInfo{}, nil)
		mockRepo.On("PlainOpen", options).Return(errors.New("open error"))

		// Initialize Qgit instance
		_, err := gitpkg.NewQGit(options, mockRepo)
		assert.Error(t, err)

		// Verify that methods were called
		mockRepo.AssertExpectations(t)
	})
}

// TestNewQGit tests the initialization logic.
func TestNewQGit(t *testing.T) {

	t.Run("Repository does not exist locally; cloning succeeds", func(t *testing.T) {
		// Create a new mock GitRepository using the NewGitRepository function
		mockRepo := new(mocks.GitRepository1)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Mock the Stat, PlainClone, and PlainOpen methods
		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})

		_, err := gitpkg.NewQGit(options, mockRepo)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

		// Expectations will be asserted automatically in the Cleanup function registered by NewGitRepository
	})
	// t.Run("Repository does not exist locally; cloning fails", func(t *testing.T) {
	// 	mockRepo := new(mocks.GitRepository1)
	// 	options := gitpkg.QgitOptions{
	// 		Path:   "/test/repo",
	// 		Url:    "https://github.com/test/repo.git",
	// 		IsBare: false,
	// 	}

	// 	// Simulate Init failure due to cloning error
	// 	mockRepo.On("Init", options).Return(errors.New("cloning failed"))

	// 	qgitInstance, err := gitpkg.NewQGit(options, mockRepo)
	// 	assert.Error(t, err)
	// 	assert.Nil(t, qgitInstance)
	// 	mockRepo.AssertExpectations(t)
	// })

	// t.Run("Repository exists locally; opening succeeds", func(t *testing.T) {
	// 	mockRepo := new(mocks.GitRepository1)
	// 	options := gitpkg.QgitOptions{
	// 		Path:   "/test/repo",
	// 		Url:    "https://github.com/test/repo.git",
	// 		IsBare: false,
	// 	}

	// 	// Simulate successful Init with existing repo
	// 	mockRepo.On("Init", options).Return(nil)

	// 	qgitInstance, err := gitpkg.NewQGit(options, mockRepo)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, qgitInstance)
	// 	mockRepo.AssertExpectations(t)
	// })

	// t.Run("Repository exists locally; opening fails", func(t *testing.T) {
	// 	mockRepo := new(mocks.GitRepository)
	// 	options := gitpkg.QgitOptions{
	// 		Path:   "/test/repo",
	// 		Url:    "https://github.com/test/repo.git",
	// 		IsBare: false,
	// 	}

	// 	// Simulate Init failure due to open repo error
	// 	mockRepo.On("Init", options).Return(errors.New("repository corrupted"))

	// 	qgitInstance, err := gitpkg.NewQGit(options, mockRepo)
	// 	assert.Error(t, err)
	// 	assert.Nil(t, qgitInstance)
	// 	mockRepo.AssertExpectations(t)
	// })

}

// TestCheckout tests the checkout logic for branches, tags, and commits.
func TestCheckout(t *testing.T) {
	// t.Run("Branch exists and is successfully checked out", func(t *testing.T) {
	// 	mockRepo := new(mocks.GitRepository1)
	// 	//mockWorktree := new(mocks.GitWorktree)
	// 	options := gitpkg.QgitOptions{
	// 		Path:   "/test/repo",
	// 		Url:    "https://github.com/test/repo.git",
	// 		IsBare: false,
	// 	}

	// 	// mockRepo.On("Checkout", "main").Return(nil).Run(func(args mock.Arguments) {
	// 	// 	fmt.Printf("Checkout called:\n")
	// 	// })
	// 	// Simulate Init and CheckRemoteRef finding a branch
	// 	mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
	// 		fmt.Println("Stat called: /test/repo/.git does not exist")
	// 	})
	// 	mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
	// 		fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
	// 	})

	// 	mockRepo.On("CheckRemoteRef", "main").Return(true, false, false)
	// 	mockRepo.On("CheckoutBranch", "main").Return(nil)

	// 	// Create the Qgit instance and perform the checkout operation

	// 	qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)
	// 	err := qgitInstance.Checkout("main")
	// 	assert.NoError(t, err)

	// 	// Verify that the methods were called with the correct arguments
	// 	mockRepo.AssertExpectations(t)
	// })

	t.Run("Tag exists and is successfully checked out", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository)
		//mockWorktree := new(mocks.GitWorktree)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Simulate Init and CheckRemoteRef finding a tag
		//mockRepo.On("Init", options).Return(nil)
		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		mockRepo.On("CheckRemoteRef", "v1.0.0").Return(false, true, false)
		mockRepo.On("CheckoutTag", "v1.0.0").Return(nil)

		// Create the Qgit instance and perform the checkout operation
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)
		err := qgitInstance.Checkout("v1.0.0")
		assert.NoError(t, err)

		// Verify that the methods were called with the correct arguments
		mockRepo.AssertExpectations(t)
	})

	t.Run("Commit hash exists and is successfully checked out", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository)
		//mockWorktree := new(mocks.GitWorktree)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Simulate Init and CheckRemoteRef finding a commit hash
		//mockRepo.On("Init", options).Return(nil)
		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		mockRepo.On("CheckRemoteRef", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef").Return(false, false, true)
		mockRepo.On("CheckoutHash", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef").Return(nil)

		// Create the Qgit instance and perform the checkout operation
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)
		err := qgitInstance.Checkout("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
		assert.NoError(t, err)

		// Verify that the methods were called with the correct arguments
		mockRepo.AssertExpectations(t)
	})

	t.Run("Reference does not exist", func(t *testing.T) {
		mockRepo := new(mocks.GitRepository)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		// Simulate Init and CheckRemoteRef failing to find the reference
		//mockRepo.On("Init", options).Return(nil)
		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		mockRepo.On("CheckRemoteRef", "unknown-ref").Return(false, false, false)

		// Create the Qgit instance and perform the checkout operation
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)
		err := qgitInstance.Checkout("unknown-ref")
		assert.Error(t, err)
		assert.Equal(t, "reference not found: unknown-ref", err.Error())

		// Verify that the methods were called with the correct arguments
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_Head(t *testing.T) {
	t.Run("Head returns the current reference successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.GitRepository)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}
		expectedRef := gitpkg.QReference{
			ReferenceName: "refs/heads/main",
			Hash:          "abc123",
		}

		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		// mockRepo.On("CheckRemoteRef", "main").Return(true, false, false)
		// mockRepo.On("CheckoutBranch", "main").Return(nil)

		// Mock the Head method of the repository
		mockRepo.On("Head").Return(expectedRef, nil)

		// Create a Qgit instance
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)

		// Act
		ref, err := qgitInstance.Head()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedRef, ref)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Head returns an error when the repository fails to get HEAD", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.GitRepository)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}
		expectedErr := errors.New("failed to get HEAD")

		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})

		// Mock the Head method of the repository to return an error
		mockRepo.On("Head").Return(gitpkg.QReference{}, expectedErr)

		// Create a Qgit instance
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)

		// Act
		ref, err := qgitInstance.Head()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, gitpkg.QReference{}, ref)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_Fetch(t *testing.T) {
	t.Run("Fetch returns no error when the fetch succeeds", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.GitRepository)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}

		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		// Mock the Fetch method of the repository
		mockRepo.On("Fetch", "refs/heads/main").Return(nil)

		// Create a Qgit instance
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)

		// Act
		err := qgitInstance.Fetch("refs/heads/main")

		// Assert
		assert.NoError(t, err)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fetch returns an error when the fetch fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.GitRepository)
		options := gitpkg.QgitOptions{
			Path:   "/test/repo",
			Url:    "https://github.com/test/repo.git",
			IsBare: false,
		}
		expectedErr := errors.New("network error")

		mockRepo.On("Stat", "/test/repo/.git").Return(nil, os.ErrNotExist).Run(func(args mock.Arguments) {
			fmt.Println("Stat called: /test/repo/.git does not exist")
		})
		mockRepo.On("PlainClone", options).Return(nil).Run(func(args mock.Arguments) {
			fmt.Printf("PlainClone called: Cloning repository with options: %+v\n", options)
		})
		// Mock the Fetch method of the repository to return an error
		mockRepo.On("Fetch", "refs/heads/main").Return(expectedErr)

		// Create a Qgit instance
		qgitInstance, _ := gitpkg.NewQGit(options, mockRepo)

		// Act
		err := qgitInstance.Fetch("refs/heads/main")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}
