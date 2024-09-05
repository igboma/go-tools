package qgit

import (
	"fmt"
	"os"
	"path/filepath"
)

// Qgit is a struct that represents a Git repository and provides methods to interact with it.
type Qgit struct {
	Option *QgitOptions
	Repo   GitRepository
}

// QReference holds information about a Git reference, such as the reference name and its hash.
type QReference struct {
	Hash          string
	ReferenceName string
}

// QgitOptions contains the options required for initializing or cloning a Git repository.
type QgitOptions struct {
	Path  string
	Url   string
	Token string
}

// QgitCheckoutOptions provides options for checking out a Git reference, including branches, tags, or commit hashes.
type QgitCheckoutOptions struct {
	branch string
	tag    string
	hash   string
}

// Head retrieves the current HEAD reference of the repository and returns it as a QReference.
func (gr *Qgit) Head() (QReference, error) {
	ref, err := gr.Repo.Head()
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
	return gr.Repo.Fetch(ref)
}

// GetChangedFilesByPRNumber retrieves the list of files that have been changed in the specified pull request.
//
// Parameters:
//   - pr: The pull request number.
//
// Returns:
//   - []string: A list of file paths that were changed in the PR.
//   - error: Returns an error if the operation fails, or nil if successful.
func (gr *Qgit) GetChangedFilesByPRNumber(pr int) ([]string, error) {
	return gr.Repo.GetChangedFilesByPRNumber(pr)
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
func (gr *Qgit) GetChangedFilesByPRNumberFilesEndingWithYAML(pr int) ([]string, error) {
	return gr.Repo.GetChangedFilesByPRNumberFilesMatching(pr, "conf.yaml")
}

// Checkout checks out the specified Git reference (branch, tag, or commit hash) in the repository.
func (gr *Qgit) Checkout(ref string) error {

	isBranch, isTag, isCommitHash := gr.Repo.CheckRemoteRef(ref)

	switch {
	case isBranch:
		fmt.Println("Checking out branch:", ref)
		return gr.Repo.CheckoutBranch(ref)
	case isTag:
		fmt.Println("Checking out tag:", ref)
		return gr.Repo.CheckoutTag(ref)
	case isCommitHash:
		fmt.Println("Checking out commit hash:", ref)
		return gr.Repo.CheckoutHash(ref)
	default:
		return fmt.Errorf("reference not found: %s", ref)
	}
}

// NewQGit initializes a new Qgit instance with the provided options and GitRepository implementation.
// It returns an error if the repository initialization fails.
func NewQGit(o *QgitOptions, gitRepo GitRepository) (*Qgit, error) {
	qgit := &Qgit{Option: o, Repo: gitRepo}
	err := qgit.Setup()
	if err != nil {
		return nil, err
	}
	return qgit, nil
}

// Stat checks if the directory exists locally.
func (gr *Qgit) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Init initializes the Git repository by either cloning it from a remote URL or opening it locally.
func (gr *Qgit) Setup() error {
	gitDir := filepath.Join(gr.Option.Path, ".git")

	// Call the Stat method from the interface
	if _, err := gr.Repo.Stat(gitDir); os.IsNotExist(err) {
		fmt.Println("Repository does not exist locally. Cloning...")
		err := gr.Repo.PlainClone(*gr.Option)
		if err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
		fmt.Println("Repository cloned successfully.")
	} else if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	} else {
		fmt.Println("Repository exists locally. Opening...")
		err := gr.Repo.PlainOpen(*gr.Option)
		if err != nil {
			return fmt.Errorf("error opening repository: %w", err)
		}
		fmt.Println("Repository opened successfully.")
	}
	return nil
}
