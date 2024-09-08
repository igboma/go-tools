package qgit

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Qgit is a struct that represents a Git repository and provides methods to interact with it.
type Qgit struct {
	option *QRepoOptions
	repo   Repository
}

// QGitClient defines the methods to interact with a Git repository.
type QGitClient interface {
	SetRepo(repo Repository)
	Repo() Repository
	SetOption(option *QRepoOptions)
	Option() *QRepoOptions
	Head() (QReference, error)
	Fetch(ref string) error
	GetChangedFilesByPRNumber(pr int) ([]string, error)
	GetConfFileChangedByPRNumber(pr int) ([]string, error)
	Checkout(ref string) error
	GetChangedFilesByPRNumberFileExtMatch(prNumber int, fileExt string) ([]string, error)
	GetChangedFilesByPRNumberFilesMatching(prNumber int, fileName string) ([]string, error)
	GetChangedFilesByPRNumberFilesByRegex(prNumber int, regexFilter string) ([]string, error)
	GetChangedFilesByPRNumberFilesByFilter(prNumber int, filter func(string) bool) ([]string, error)
	GetFileContentFromBranch(branch, file string) (string, error)
	GetFileContentFromCommit(commitHash, file string) (string, error)
}

// Head retrieves the current HEAD reference of the repository and returns it as a QReference.
func (gr *Qgit) Head() (QReference, error) {
	ref, err := gr.Repo().Head()
	return ref, err
}

// Fetch fetches the given Git reference to ensure that it is available in the local repository.
// It synchronizes the local repository with the remote branch.
//
// Parameters:
//   - ref: The reference to fetch, typically a branch or tag.
//
// Returns:
//   - error: Returns an error if the fetch operation fails, or nil if successful.
func (gr *Qgit) Fetch(ref string) error {
	return gr.Repo().Fetch(ref)
}

// GetChangedFilesByPRNumber retrieves the list of files that have been changed in the specified pull request.
//
// Parameters:
//   - pr: The pull request number.
//
// Returns:
//   - []string: A list of file paths that were changed in the PR.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetChangedFilesByPRNumber(pr int) (changedFiles []string, err error) {
	changedFiles, err = gr.Repo().GetChangedFilesByPRNumber(pr)
	return
}

// GetChangedFilesByPRNumberFilesEndingWithYAML retrieves the list of changed files in the specified pull request
// that end with "conf.yaml".
//
// Parameters:
//   - pr: The pull request number.
//
// Returns:
//   - []string: A list of file paths that end with "conf.yaml" and were changed in the PR.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetConfFileChangedByPRNumber(pr int) (matchingFiles []string, err error) {
	matchingFiles, err = gr.GetChangedFilesByPRNumberFilesMatching(pr, "conf.yaml")
	return
}

// Checkout checks out the specified Git reference in the repository.
// The reference can be a branch, tag, or commit hash. The function checks
// the type of reference and performs the appropriate checkout operation.
//
// Parameters:
//   - ref: The Git reference to check out, which can be a branch name, tag name, or commit hash.
//
// Returns:
//   - error: Returns an error if the checkout operation fails, or if the reference type cannot be determined.
func (gr *Qgit) Checkout(ref string) error {
	// Determine if the reference is a branch, tag, or commit hash
	isBranch, isTag, isCommitHash, err := gr.Repo().CheckRemoteRef(ref)
	if err != nil {
		return fmt.Errorf("Checkout error: %w", err)
	}

	// Perform the appropriate checkout operation based on the reference type
	switch {
	case isBranch:
		fmt.Println("Checking out branch:", ref)
		return gr.Repo().CheckoutBranch(ref)
	case isTag:
		fmt.Println("Checking out tag:", ref)
		return gr.Repo().CheckoutTag(ref)
	case isCommitHash:
		fmt.Println("Checking out commit hash:", ref)
		return gr.Repo().CheckoutHash(ref)
	default:
		return fmt.Errorf("reference not found: %s", ref)
	}
}

// GetChangedFilesByPRNumberFileExtMatch retrieves the list of changed files from a pull request that match the specified file extension.
//
// Parameters:
//   - prNumber: The pull request number for which to retrieve changed files.
//   - fileExt: The file extension to filter the changed files (e.g., ".yaml").
//
// Returns:
//   - []string: A list of changed files that match the specified file extension.
//   - error: Returns an error if the changed files cannot be retrieved or filtered correctly.
func (gr *Qgit) GetChangedFilesByPRNumberFileExtMatch(prNumber int, fileExt string) (matchingFiles []string, err error) {
	// Call the existing function to get all changed files by PR number
	changedFiles, err := gr.Repo().GetChangedFilesByPRNumber(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files for PR %d: %w", prNumber, err)
	}

	// Filter the files to only include those that match the specified file extension
	for _, file := range changedFiles {
		// Compare file extension case-insensitively
		if strings.EqualFold(filepath.Ext(file), fileExt) {
			matchingFiles = append(matchingFiles, file)
		}
	}

	// Return the list of files with the matching extension
	return matchingFiles, nil
}

// GetChangedFilesByPRNumberFilesMatching retrieves the list of files that were changed in the specified pull request,
// and filters them to include only the files that exactly match the provided filename (case-insensitive).
//
// Parameters:
//   - prNumber: The pull request number for which to retrieve changed files.
//   - fileName: The filename to match against the changed files.
//
// Returns:
//   - []string: A list of file paths that were changed in the PR and match the provided filename.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetChangedFilesByPRNumberFilesMatching(prNumber int, fileName string) (matchingFiles []string, err error) {
	// Call the existing function to get all changed files by PR number
	changedFiles, err := gr.GetChangedFilesByPRNumber(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files for PR %d: %w", prNumber, err)
	}

	// Filter the files to only include those that match the provided filename
	for _, file := range changedFiles {
		// Extract the base name of the file and check if it matches the given filename
		if strings.EqualFold(filepath.Base(file), fileName) {
			matchingFiles = append(matchingFiles, file)
		}
	}

	return matchingFiles, nil
}

// GetChangedFilesByPRNumberFilesByRegex retrieves the list of files that were changed in the specified pull request,
// and filters them using the provided regular expression.
//
// Parameters:
//   - prNumber: The pull request number for which to retrieve changed files.
//   - regexFilter: A string representing the regular expression used to filter the files.
//
// Returns:
//   - []string: A list of file paths that were changed in the PR and match the regular expression.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetChangedFilesByPRNumberFilesByRegex(prNumber int, regexFilter string) (filteredFiles []string, err error) {
	changedFiles, err := gr.GetChangedFilesByPRNumber(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files for PR %d: %w", prNumber, err)
	}

	// Compile the regex filter
	filterRegex, err := regexp.Compile(regexFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex filter: %w", err)
	}

	// Filter the files based on the regex
	for _, file := range changedFiles {
		if filterRegex.MatchString(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

// GetChangedFilesByPRNumberFilesByFilter retrieves the list of files that were changed in the specified pull request,
// and filters them using the provided filter function.
//
// Parameters:
//   - prNumber: The pull request number for which to retrieve changed files.
//   - filter: A function that takes a file path as input and returns a boolean indicating whether the file matches the filter criteria.
//
// Returns:
//   - []string: A list of file paths that were changed in the PR and match the filter function.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetChangedFilesByPRNumberFilesByFilter(prNumber int, filter func(string) bool) (filteredFiles []string, err error) {
	changedFiles, err := gr.GetChangedFilesByPRNumber(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files for PR %d: %w", prNumber, err)
	}

	// Filter the files based on the provided function
	for _, file := range changedFiles {
		if filter(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

// GetFileContentFromBranch retrieves the content of a specified file from the latest commit on a given branch.
//
// Parameters:
//   - branch: The name of the branch from which to retrieve the file (e.g., "refs/heads/main").
//   - file: The path to the file within the repository whose content needs to be retrieved.
//
// Returns:
//   - content: A string containing the content of the specified file.
//   - error: An error if the operation fails, or nil if successful.
func (gr *Qgit) GetFileContentFromBranch(branch, file string) (content string, err error) {
	content, err = gr.Repo().GetFileContentFromBranch(branch, file)
	return
}

// GetFileContentFromCommit retrieves the content of a specified file from a given commit hash.
//
// Parameters:
//   - commitHash: The hash of the commit from which to retrieve the file.
//   - file: The path to the file within the repository whose content needs to be retrieved.
//
// Returns:
//   - content: A string containing the content of the specified file from the commit.
//   - error: An error if the operation fails, or nil if successful.
func (gr *Qgit) GetFileContentFromCommit(commitHash, file string) (content string, err error) {
	content, err = gr.Repo().GetFileContentFromCommit(commitHash, file)
	return
}

// SetRepo sets the Repository implementation for the Repository interface.
// This allows Qgit to interact with a specific Git repository.
//
// Parameters:
//   - repo: A Repository interface that provides methods for interacting with the Git repository.
func (gr *Qgit) SetRepo(repo Repository) {
	gr.repo = repo
}

// Repo retrieves the current Repository implementation being used by the Repository interface.
//
// Returns:
//   - Repository: The current Repository interface implementation.
func (gr *Qgit) Repo() Repository {
	return gr.repo
}

// SetOption sets the repository options (such as path, URL, and token) for the Qgit instance.
// This allows configuration of the repository parameters.
//
// Parameters:
//   - opt: A pointer to QRepoOptions struct containing the options for initializing or cloning the repository.
func (gr *Qgit) SetOption(opt *QRepoOptions) {
	gr.option = opt
}

// Option retrieves the current repository options (such as path, URL, and token) used by the Qgit instance.
//
// Returns:
//   - *QRepoOptions: A pointer to the QRepoOptions struct containing the repository options.
func (gr *Qgit) Option() *QRepoOptions {
	return gr.option
}

// NewQGit initializes a new Qgit instance with the provided repository options and repository implementation.
// It associates the QRepoOptions with the given Repository by calling SetOption, then returns a QGitClient.
//
// Parameters:
//   - o: A pointer to QRepoOptions containing the repository configuration (path, URL, token).
//   - repoInstance: A Repository interface implementation for interacting with the Git repository.
//
// Returns:
//   - QGitClient: The interface implemented by Qgit.
func NewQGit(o *QRepoOptions, repoInstance Repository) QGitClient {
	repoInstance.SetOption(o)
	qgit := &Qgit{option: o, repo: repoInstance}
	return qgit
}
