package qgit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitRepo is a struct that implements the GitRepository interface using the go-git library.
type GitRepo struct {
	Repo   *git.Repository
	Option *QgitOptions
}

// Stat checks if the directory exists locally.
func (gr *GitRepo) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// GitRepository defines an interface for performing Git repository operations.
type GitRepository interface {
	Head() (QReference, error)
	Worktree() (GitWorktree, error)
	PlainClone(o QgitOptions) error
	PlainOpen(o QgitOptions) error
	Fetch(refSpecStr string) error
	Checkout(ref string) error
	CheckoutBranch(branch string) error
	CheckoutTag(tag string) error
	CheckoutHash(hash string) error
	CheckRemoteRef(ref string) (bool, bool, bool)
	Stat(path string) (os.FileInfo, error)
	GetFileContentFromBranch(branch, file string) (string, error)
	GetFileContentFromCommit(commitHash, file string) (string, error)
	GetChangedFilesByPRNumber(prNumber int) ([]string, error)
	GetChangedFilesByPRNumberFilesEndingWithYAML(prNumber int) ([]string, error)
}

type PullRequest struct {
	Number int
	Title  string
	URL    string
	State  string
}

// Checkout checks out the specified Git reference (branch, tag, or commit hash) in the repository.
func (gr *GitRepo) Checkout(ref string) error {

	isBranch, isTag, isCommitHash := gr.CheckRemoteRef(ref)

	switch {
	case isBranch:
		fmt.Println("Checking out branch:", ref)
		return gr.CheckoutBranch(ref)
	case isTag:
		fmt.Println("Checking out tag:", ref)
		return gr.CheckoutTag(ref)
	case isCommitHash:
		fmt.Println("Checking out commit hash:", ref)
		return gr.CheckoutHash(ref)
	default:
		return fmt.Errorf("reference not found: %s", ref)
	}
}

func (gr *GitRepo) PlainClone(o QgitOptions) error {
	fmt.Println("Repository does not exist locally. Cloning here...")
	var err error = nil
	repo, err := git.PlainClone(o.Path, o.IsBare, &git.CloneOptions{
		URL: o.Url,
		Auth: &http.BasicAuth{
			Username: "git",   // Can be anything, GitHub ignores it
			Password: o.Token, // Use the token as the password
		},
	})
	gr.Repo = repo
	return err
}

func (gr *GitRepo) PlainOpen(o QgitOptions) error {
	fmt.Println("Repository exists locally. Opening...")
	repo, err := git.PlainOpen(o.Path)
	if err != nil {
		return fmt.Errorf("error opening repository: %w", err)
	}
	gr.Repo = repo
	fmt.Println("Repository opened successfully.")

	return nil
}

// Head retrieves the current HEAD reference of the repository.
func (gr *GitRepo) Head() (QReference, error) {
	ref, err := gr.Repo.Head()
	if err != nil {
		return QReference{}, err
	}
	return QReference{
		ReferenceName: ref.Name().String(),
		Hash:          ref.Hash().String(),
	}, nil
}

// Worktree retrieves the worktree of the Git repository and returns it as a GitWorktree.
func (gr *GitRepo) Worktree() (GitWorktree, error) {
	wt, err := gr.Repo.Worktree()
	if err != nil {
		return nil, err
	}
	return &GitWt{wt}, nil
}

// Fetch fetches changes from the remote repository, based on the specified refSpec.
func (gr *GitRepo) Fetch(refSpecStr string) error {
	remote, err := gr.Repo.Remote("origin")
	if err != nil {
		return err
	}

	var refSpecs []config.RefSpec
	if refSpecStr == "" {
		refSpecs = []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
			config.RefSpec("+refs/tags/*:refs/tags/*"),
		}
	} else {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err := remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("fetch origin failed: %w", err)
	}

	fmt.Println("Fetched remote references successfully.")
	return nil
}

// CheckoutBranch checks out the specified branch in the repository's worktree.
func (gr *GitRepo) CheckoutBranch(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	checkoutOption := QgitCheckoutOptions{branch: ref}

	if err := wt.Checkout(&checkoutOption); err != nil {
		if err := gr.Fetch(fmt.Sprintf("refs/heads/%s:refs/heads/%s", ref, ref)); err != nil {
			return err
		}
		return wt.Checkout(&checkoutOption)
	}
	return nil
}

// CheckoutTag checks out the specified tag in the repository's worktree.
func (gr *GitRepo) CheckoutTag(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	err = gr.Fetch(fmt.Sprintf("+refs/tags/%s:refs/tags/%s", ref, ref))
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching tag %s: %w", ref, err)
	}

	checkoutOption := QgitCheckoutOptions{tag: ref}

	err = wt.Checkout(&checkoutOption)
	if err != nil {
		return fmt.Errorf("error checking out tag %s: %w", ref, err)
	}
	return nil
}

// CheckoutHash checks out the specified commit hash in the repository's worktree.
func (gr *GitRepo) CheckoutHash(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&QgitCheckoutOptions{hash: ref})
	if err != nil {
		return fmt.Errorf("error checking out commit %s: %w", ref, err)
	}

	return nil
}

// checkRemoteRef checks if the specified reference is a branch, tag, or commit hash by querying the remote and local repository.
func (gr *GitRepo) CheckRemoteRef(ref string) (bool, bool, bool) {
	remote, err := gr.Repo.Remote("origin")
	if err != nil {
		fmt.Printf("Error getting remote: %v\n", err)
		return false, false, false
	}

	refs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: "git",           // GitHub ignores the username but requires it
			Password: gr.Option.Token, // Use the token for authentication
		},
	})
	if err != nil {
		fmt.Printf("Error listing remote references: %v\n", err)
		return false, false, false
	}
	return gr.classifyRef(ref, refs)
}

// classifyRef checks the given reference against the list of remote references and identifies it as a branch, tag, or commit hash.
func (gr *GitRepo) classifyRef(ref string, refs []*plumbing.Reference) (bool, bool, bool) {
	isBranch := false
	isTag := false
	isCommitHash := false

	for _, r := range refs {
		if r.Name().IsBranch() && r.Name().Short() == ref {
			isBranch = true
		} else if r.Name().IsTag() && r.Name().Short() == ref {
			isTag = true
		}
	}

	if len(ref) == 40 {
		_, err := gr.Repo.CommitObject(plumbing.NewHash(ref))
		if err == nil {
			isCommitHash = true
		} else {
			fmt.Printf("Error fetching commit object: %v\n", err)
		}
	}

	return isBranch, isTag, isCommitHash
}

func (gr *GitRepo) GetFileContentFromBranch(branch, file string) (string, error) {
	ref, err := gr.Repo.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		return "", err
	}

	commit, err := gr.Repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	content, err := entry.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (gr *GitRepo) GetFileContentFromCommit(commitHash, file string) (string, error) {
	commit, err := gr.Repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	content, err := entry.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

// GetChangedFilesByPRNumber fetches the changed files between the main branch and the PR branch.
func (gr *GitRepo) GetChangedFilesByPRNumber(prNumber int) ([]string, error) {
	// Convert the PR number into a reference that exists in the Git repository
	// Usually PR references are in the form: refs/pull/{prNumber}/head
	prRef := fmt.Sprintf("refs/pull/%d/head", prNumber)

	fmt.Printf("pref %v==>", prRef)

	// Fetch the remote branch (PR branch) to ensure the reference exists locally
	err := gr.Repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", prRef, prRef))},
		Auth: &http.BasicAuth{
			Username: "git", // GitHub ignores the username but requires it
			Password: gr.Option.Token,
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("failed to fetch remote branch %s: %w", prRef, err)
	}

	// Get the current HEAD reference (main branch)
	currentRef, err := gr.Repo.Head()

	fmt.Printf("currentRef %v==>", currentRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Resolve the comparison reference from the PR (origin)
	compareRef, err := gr.Repo.Reference(plumbing.ReferenceName(prRef), true)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve reference %s: %w", prRef, err)
	}

	fmt.Printf("compareRef %v==>", compareRef)

	// Get the commit for the comparison reference (PR branch)
	compareCommit, err := gr.Repo.CommitObject(compareRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for ref %s: %w", prRef, err)
	}

	fmt.Printf("compareCommit %v==>", compareCommit)
	// Get the commit for the current HEAD reference (main branch)
	currentCommit, err := gr.Repo.CommitObject(currentRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get current commit: %w", err)
	}

	fmt.Printf("currentCommit %v==>", currentCommit)

	// Get the file changes between the two commits
	patch, err := currentCommit.Patch(compareCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate patch: %w", err)
	}

	fmt.Printf("patch %v==>", patch)

	// Collect the list of changed files
	var changedFiles []string
	for _, fileStat := range patch.Stats() {
		changedFiles = append(changedFiles, fileStat.Name)
	}

	fmt.Printf("changedFiles %v==>", changedFiles)

	return changedFiles, nil
}

// GetChangedFilesByPRNumberFilesEndingWithYAML fetches the changed files in a given PR number and filters for files that end with `.yaml`.
func (gr *GitRepo) GetChangedFilesByPRNumberFilesEndingWithYAML(prNumber int) ([]string, error) {
	// Call the existing function to get all changed files by PR number
	changedFiles, err := gr.GetChangedFilesByPRNumber(prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files for PR %d: %w", prNumber, err)
	}

	// Filter the files to only include those that end with `.yaml`
	var yamlFiles []string
	for _, file := range changedFiles {
		if strings.EqualFold(filepath.Ext(file), ".yaml") {
			yamlFiles = append(yamlFiles, file)
		}
	}

	return yamlFiles, nil
}
