package qgit_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gitpkg/qgit"
	"gitpkg/qgit/mocks"

	"github.com/stretchr/testify/assert"
)

func TestQgit_Head(t *testing.T) {
	t.Run("Head returns the current reference successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		expectedRef := qgit.QReference{
			ReferenceName: "refs/heads/main",
			Hash:          "abc123",
		}

		// Mock the Head method of the repository
		mockRepo.On("Head").Return(expectedRef, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

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
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		expectedErr := errors.New("failed to get HEAD")

		// Mock the Head method of the repository to return an error
		mockRepo.On("Head").Return(qgit.QReference{}, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		ref, err := qgitInstance.Head()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, qgit.QReference{}, ref)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_Fetch(t *testing.T) {
	t.Run("Fetch returns no error when the fetch succeeds", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}

		// Mock the Fetch method of the repository
		mockRepo.On("Fetch", "refs/heads/main").Return(nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Fetch("refs/heads/main")

		// Assert
		assert.NoError(t, err)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fetch returns an error when the fetch fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		expectedErr := errors.New("network error")

		// Mock the Fetch method of the repository to return an error
		mockRepo.On("Fetch", "refs/heads/main").Return(expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Fetch("refs/heads/main")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetChangedFilesByPRNumber(t *testing.T) {
	t.Run("GetChangedFilesByPRNumber returns files when the repository fetch succeeds", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		expectedFiles := []string{"file1.txt", "file2.txt"}

		// Mock the GetChangedFilesByPRNumber method of the repository
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(expectedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		files, err := qgitInstance.GetChangedFilesByPRNumber(prNumber)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedFiles, files)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumber returns an error when the repository fetch fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		expectedErr := errors.New("repository error")

		// Mock the GetChangedFilesByPRNumber method of the repository to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		files, err := qgitInstance.GetChangedFilesByPRNumber(prNumber)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, files)
		assert.Equal(t, expectedErr, err)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetChangedFilesByPRNumberFileExtMatch(t *testing.T) {
	t.Run("GetChangedFilesByPRNumberFileExtMatch returns matching files based on file extension", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "file2.json", "file3.yaml", "README.md"}
		expectedMatchingFiles := []string{"file1.yaml", "file3.yaml"}

		fileExt := ".yaml" // The file extension to match

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFileExtMatch(prNumber, fileExt)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedMatchingFiles, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFileExtMatch returns no files when no match is found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "file2.json", "file3.yaml"}
		fileExt := ".md" // File extension that does not match any files
		//expectedMatchingFiles := []string{} // No matching files

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFileExtMatch(prNumber, fileExt)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, matchingFiles) // Use assert.Empty to handle both nil and empty slice

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFileExtMatch returns an error when GetChangedFilesByPRNumber fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		fileExt := ".yaml"
		expectedErr := errors.New("failed to fetch changed files")

		// Mock the GetChangedFilesByPRNumber method to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFileExtMatch(prNumber, fileExt)

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		assert.Nil(t, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetChangedFilesByPRNumberFilesMatching(t *testing.T) {
	t.Run("GetChangedFilesByPRNumberFilesMatching returns matching files successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "config/file2.yaml", "docs/file3.yaml", "README.md"}
		expectedMatchingFiles := []string{"config/file2.yaml"}

		fileName := "file2.yaml" // The filename to match

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesMatching(prNumber, fileName)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedMatchingFiles, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFilesMatching returns no files when no match is found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "config/file2.yaml", "docs/file3.yaml"}

		fileName := "nonexistent.yaml" // The filename that does not exist in the list

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesMatching(prNumber, fileName)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, matchingFiles) // Use assert.Empty to handle both nil and empty slice

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFilesMatching returns an error when GetChangedFilesByPRNumber fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		fileName := "file2.yaml" // The filename to match
		expectedErr := errors.New("failed to fetch changed files")

		// Mock the GetChangedFilesByPRNumber method to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesMatching(prNumber, fileName)

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		assert.Nil(t, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetChangedFilesByPRNumberFilesByRegex(t *testing.T) {
	t.Run("GetChangedFilesByPRNumberFilesByRegex returns filtered files successfully based on regex", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "file2.json", "docs/file3.yaml", "README.md"}
		expectedFilteredFiles := []string{"docs/file3.yaml"}

		// Regex to match files that are inside the "docs" folder
		regexFilter := `^docs/.*\.yaml$`

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		filteredFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesByRegex(prNumber, regexFilter)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedFilteredFiles, filteredFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFilesByRegex returns an error when regex compilation fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "file2.json", "docs/file3.yaml"}
		invalidRegex := `[` // Invalid regex pattern
		//expectedErr := errors.New("failed to compile regex filter")

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		filteredFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesByRegex(prNumber, invalidRegex)

		// Assert
		assert.Error(t, err)
		//assert.Contains(t, err.Error(), "failed to compile regex filter")
		assert.Nil(t, filteredFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFilesByRegex returns an error when GetChangedFilesByPRNumber fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		expectedErr := errors.New("failed to fetch changed files")
		validRegex := `.*\.yaml$` // Valid regex

		// Mock the GetChangedFilesByPRNumber method to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		filteredFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesByRegex(prNumber, validRegex)

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		assert.Nil(t, filteredFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetChangedFilesByPRNumberFilesByFilter(t *testing.T) {
	t.Run("GetChangedFilesByPRNumberFilesByFilter returns filtered files successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"file1.yaml", "file2.json", "file3.yaml"}
		expectedFilteredFiles := []string{"file1.yaml", "file3.yaml"}

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		filteredFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesByFilter(prNumber, func(file string) bool {
			return strings.HasSuffix(file, ".yaml") // Only return files ending with .yaml
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedFilteredFiles, filteredFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetChangedFilesByPRNumberFilesByFilter returns an error when GetChangedFilesByPRNumber fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		expectedErr := errors.New("failed to fetch changed files")

		// Mock the GetChangedFilesByPRNumber method to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		filteredFiles, err := qgitInstance.GetChangedFilesByPRNumberFilesByFilter(prNumber, func(file string) bool {
			return strings.HasSuffix(file, ".yaml")
		})

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		assert.Nil(t, filteredFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_GetConfFileChangedByPRNumber(t *testing.T) {
	t.Run("GetConfFileChangedByPRNumber returns matching YAML files", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"config/conf.yaml", "docs/readme.md", "src/conf.yaml"}
		expectedMatchingFiles := []string{"config/conf.yaml", "src/conf.yaml"}

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetConfFileChangedByPRNumber(prNumber)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedMatchingFiles, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetConfFileChangedByPRNumber returns no files when no YAML files are found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		changedFiles := []string{"docs/readme.md", "src/app.json"}

		// Mock the GetChangedFilesByPRNumber method to return the list of changed files
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(changedFiles, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetConfFileChangedByPRNumber(prNumber)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, matchingFiles) // No YAML files found

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetConfFileChangedByPRNumber returns an error when GetChangedFilesByPRNumber fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		prNumber := 123
		expectedErr := errors.New("failed to fetch changed files")

		// Mock the GetChangedFilesByPRNumber method to return an error
		mockRepo.On("GetChangedFilesByPRNumber", prNumber).Return(nil, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		matchingFiles, err := qgitInstance.GetConfFileChangedByPRNumber(prNumber)

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		assert.Nil(t, matchingFiles)

		// Verify that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}

func TestQgit_Checkout(t *testing.T) {
	t.Run("Checkout a branch reference successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		ref := "refs/heads/main"

		// Mock the CheckRemoteRef to return true for branch
		mockRepo.On("CheckRemoteRef", ref).Return(true, false, false, nil)
		mockRepo.On("CheckoutBranch", ref).Return(nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Checkout(ref)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Checkout a tag reference successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		ref := "v1.0.0"

		// Mock the CheckRemoteRef to return true for tag
		mockRepo.On("CheckRemoteRef", ref).Return(false, true, false, nil)
		mockRepo.On("CheckoutTag", ref).Return(nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Checkout(ref)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Checkout a commit hash reference successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		ref := "abc123"

		// Mock the CheckRemoteRef to return true for commit hash
		mockRepo.On("CheckRemoteRef", ref).Return(false, false, true, nil)
		mockRepo.On("CheckoutHash", ref).Return(nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Checkout(ref)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Checkout returns an error when CheckRemoteRef fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		ref := "refs/heads/main"
		expectedErr := errors.New("failed to check remote ref")

		// Mock the CheckRemoteRef to return an error
		mockRepo.On("CheckRemoteRef", ref).Return(false, false, false, expectedErr)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Checkout(ref)

		// Assert
		assert.Error(t, err)
		//assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Checkout returns an error when reference is not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.Repository)
		options := qgit.QRepoOptions{
			Path: "/test/repo",
			Url:  "https://github.com/test/repo.git",
		}
		ref := "unknown_ref"

		// Mock the CheckRemoteRef to return all false
		mockRepo.On("CheckRemoteRef", ref).Return(false, false, false, nil)

		// Create a Qgit instance
		qgitInstance := qgit.NewQGit(&options, mockRepo)

		// Act
		err := qgitInstance.Checkout(ref)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("reference not found: %s", ref), err.Error())
		mockRepo.AssertExpectations(t)
	})
}
