package qgit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Client struct {
	repo *git.Repository
	opts *Options
}

// NewClient creates a new instance of GitClient with the provided options.
//
// The client is not ready to use until InitRepo(), Clone() or Open() is used.
func NewClient(opts ...Option) (*Client, error) {
	options, err := compileOptions(opts...)
	if err != nil {
		return nil, err
	}
	return &Client{opts: options}, nil
}

// InitRepo clones the repo from remote if it does not exist locally
func (c *Client) InitRepo() (err error) {
	gitDir := filepath.Join(c.opts.RepoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		fmt.Println("Repository does not exist locally. Cloning...")
		err = c.Clone()
		if err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
		fmt.Println("Repository cloned successfully.")
	}
	if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	}

	fmt.Println("Repository exists locally. Opening...")
	if err = c.Open(); err != nil {
		return fmt.Errorf("error opening repository: %w", err)
	}
	fmt.Println("Repository opened successfully.")
	return nil
}

// Clone clones a Git repository from a remote URL to the specified local path using basic authentication.
func (c *Client) Clone() (err error) {
	c.repo, err = git.PlainClone(c.opts.RepoPath, false, &git.CloneOptions{
		URL: c.opts.RepoUrl,
		Auth: &http.BasicAuth{
			Username: c.opts.Username, // GitHub ignores the username, can be anything
			Password: c.opts.Token,
		},
	})
	return err
}

// Open opens an existing Git repository from the specified RepoPath.
func (c *Client) Open() (err error) {
	c.repo, err = git.PlainOpen(c.opts.RepoPath)
	if err != nil {
		return fmt.Errorf("failed to open repo at %s: %w", c.opts.RepoPath, err)
	}
	return
}

// Checkout checks out the specified Git reference (branch, tag, or commit hash) in the repository.
//
// Use the shortname of the branch and tag
func (c *Client) Checkout(ref string) error {
	isBranch, isTag, isCommitHash, err := c.CheckLocalRef(ref)
	if err != nil {
		fmt.Printf("%v", fmt.Errorf("failed to resolve ref locally: %w", err))
		fmt.Printf("Checking remote refs")
		isBranch, isTag, isCommitHash, err = c.CheckRemoteRef(ref)
		if err != nil {
			return fmt.Errorf("failed to resolved ref on remote: %w", err)
		}
		c.Fetch(ref) // fetch the ref from the remote
	}

	checkoutOpts := git.CheckoutOptions{}
	switch {
	case isBranch:
		fmt.Println("Checking out branch:", ref)
		checkoutOpts = git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(ref),
		}
	case isTag:
		fmt.Println("Checking out tag:", ref)
		ref, err := c.repo.Tag(ref)
		if err != nil {
			return fmt.Errorf("failed to get tag reference: %w", err)
		}
		checkoutOpts = git.CheckoutOptions{
			Hash: ref.Hash(),
		}
	case isCommitHash:
		fmt.Println("Checking out commit hash:", ref)
		checkoutOpts = git.CheckoutOptions{
			Hash: plumbing.NewHash(ref),
		}
	default:
		return fmt.Errorf("reference not found: %s", ref)
	}
	return c.checkout(&checkoutOpts)
}

func (c *Client) checkout(opts *git.CheckoutOptions) error {
	wt, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get repo work tree: %w", err)
	}

	return wt.Checkout(opts)
}

// Head retrieves the hash of the current HEAD reference of the repository.
func (c *Client) Head() (string, error) {
	ref, err := c.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve head: %w", err)
	}
	return ref.Hash().String(), nil
}

// Fetch fetches updates from the remote repository, ensuring that the specified references are up to date.
//
// Parameters:
//   - refSpecStr: A string specifying the shortname of the ref to fetch. If empty, default refspecs for branches and tags are used.
func (c *Client) Fetch(ref string) error {
	// Get the remote repository (assumed to be named "origin")
	remote, err := c.repo.Remote("origin")
	if err != nil {
		return fmt.Errorf(`failed to get remote "origin": %w`, err)
	}

	var refSpecs []config.RefSpec
	if ref == "" {
		// If no specific refspec is provided, fetch all branches and tags
		refSpecs = []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"), // Fetch all branches
			config.RefSpec("+refs/tags/*:refs/tags/*"),            // Fetch all tags
		}
	} else {
		refSpecs = []config.RefSpec{config.RefSpec(ref)}
	}

	if err := remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("fetch origin failed: %w", err)
	}

	return nil
}

// CheckRemoteRef checks if the specified reference is a branch, tag, or commit hash by querying the remote repository.
// It returns three boolean flags indicating whether the ref is a branch, tag, or commit hash.
func (c *Client) CheckRemoteRef(ref string) (isBranch, isTag, isCommitHash bool, err error) {
	remote, err := c.repo.Remote("origin")
	if err != nil {
		err = fmt.Errorf("error getting remote: %w", err)
		return
	}

	// List the remote references
	refs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: c.opts.Username, // GitHub ignores the username but requires it
			Password: c.opts.Token,    // Use the token for authentication
		},
	})
	if err != nil {
		err = fmt.Errorf("error listing remote references: %w", err)
		return
	}
	for _, r := range refs {
		if r.Name().Short() == ref {
			_, isBranch, isTag, isCommitHash, err = c.resolveRef(ref)
			break
		}
	}
	return
}

// CheckRemoteRef checks if the specified reference is a branch, tag, or commit hash by querying the local repository.
// It returns three boolean flags indicating whether the ref is a branch, tag, or commit hash.
func (c *Client) CheckLocalRef(ref string) (isBranch, isTag, isCommitHash bool, err error) {
	_, isBranch, isTag, isCommitHash, err = c.resolveRef(ref)
	return
}

func (c *Client) resolveRef(ref string) (hash string, isBranch, isTag, isCommitHash bool, err error) {
	// Try to resolve the reference as a branch/tag
	if gitRef, err := c.repo.Reference(plumbing.ReferenceName("refs/heads/"+ref), true); err == nil {
		isBranch = true
		hash = gitRef.Hash().String()
	} else if _, err = c.repo.Reference(plumbing.ReferenceName("refs/tags/"+ref), true); err == nil {
		isTag = true
		hash = gitRef.Hash().String()
	}
	// If the ref has a length of 40, check if it is a commit hash
	if len(ref) == 40 {
		if commit, err := c.repo.CommitObject(plumbing.NewHash(ref)); err == nil {
			isCommitHash = true
			hash = commit.Hash.String()
		}
	}
	if err != nil {
		err = fmt.Errorf("failed to resolve ref %s: %w", ref, err)
	}
	return hash, isBranch, isTag, isCommitHash, err
}

// FileContentFromCommit retrieves the content of a specified file from a given commit hash in the repository.
//
// Parameters:
//   - commitHash: The hash of the commit from which to retrieve the file.
//   - file: The path of the file whose content needs to be retrieved.
func (c *Client) FileContentFromCommit(commitHash, file string) (content string, err error) {
	// Get the commit object using the commit hash
	commit, err := c.repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return "", fmt.Errorf("failed to get commit object: %w", err)
	}

	// Get the tree associated with the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get commit tree: %w", err)
	}

	// Find the file entry in the tree
	entry, err := tree.File(file)
	if err != nil {
		return "", fmt.Errorf("failed to get file from tree: %w", err)
	}

	// Retrieve the file contents
	content, err = entry.Contents()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file content: %w", err)
	}

	return content, nil
}

// FileContentFromBranch retrieves the content of a specified file from the latest commit on a given branch.
//
// Parameters:
//   - branch: The branch name from which to retrieve the file (e.g., "main").
//   - file: The path of the file whose content needs to be retrieved.
func (c *Client) FileContentFromBranch(branch, file string) (content string, err error) {
	// Get the reference to the specified branch
	ref, err := c.repo.Reference(plumbing.NewBranchReferenceName(branch), true)
	if err != nil {
		return "", fmt.Errorf("failed to resolved branch ref: %w", err)
	}
	return c.FileContentFromCommit(ref.Hash().String(), file)
}

func (c *Client) changedFiles(base, current string) (*object.Changes, error) {
	baseHashStr, _, _, _, err := c.resolveRef(base)
	if err != nil {
		return nil, err
	}
	currentHashStr, _, _, _, err := c.resolveRef(current)
	if err != nil {
		return nil, err
	}
	baseCommit, err := c.repo.CommitObject(plumbing.NewHash(baseHashStr))
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for base ref: %s err:%v", base, err)
	}
	currentCommit, err := c.repo.CommitObject(plumbing.NewHash(currentHashStr))
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for current ref: %s err:%v", current, err)
	}

	// Get the trees for the commits
	currentTree, err := currentCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get current commit tree: %w", err)
	}
	baseTree, err := baseCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get base commit tree: %w", err)
	}

	changes, err := object.DiffTree(baseTree, currentTree)
	if err != nil {
		return nil, fmt.Errorf("failed to diff commits: %w", err)
	}
	return &changes, err
}

// ChangedFiles returns the changed files between the base ref and the current ref.
func (c *Client) ChangedFiles(base, current string) (changedFiles []string, err error) {
	changes, err := c.changedFiles(base, current)
	if err != nil {
		return nil, err
	}
	patch, err := changes.Patch()
	if err != nil {
		return nil, fmt.Errorf("failed to get patch from changes: %w", err)
	}
	for _, s := range patch.Stats() {
		changedFiles = append(changedFiles, s.Name)
	}
	return
}

// ChangedFilesByFilter returns the changed filepaths between the base ref and the current ref, matching the given filter.
func (c *Client) ChangedFilesByFilter(base, current string, filter func(string) bool) ([]string, error) {
	changedFiles, err := c.ChangedFiles(base, current)
	if err != nil {
		return nil, err
	}
	var filteredFiles []string
	for _, file := range changedFiles {
		if filter(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

// ChangedFilesByExt returns the changed filepaths between the base ref and the current ref, matching the given extension.
func (c *Client) ChangedFilesByExt(base, current string, fileExt string) ([]string, error) {
	return c.ChangedFilesByFilter(base, current,
		func(file string) bool {
			return strings.EqualFold(filepath.Ext(file), fileExt)
		},
	)
}

// ChangedFilesByName returns the changed filepaths between the base ref and the current ref, matching the given basename.
func (c *Client) ChangedFilesByName(base, current, baseName string) ([]string, error) {
	return c.ChangedFilesByFilter(base, current,
		func(file string) bool {
			return strings.EqualFold(filepath.Base(file), baseName)
		},
	)
}

// ChangedFilesByRegex returns the changed filepaths between the base ref and the current ref, matching the given regex.
func (c *Client) ChangedFilesByRegex(base, current, regexFilter string) ([]string, error) {
	filterRegex, err := regexp.Compile(regexFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex filter: %w", err)
	}
	return c.ChangedFilesByFilter(base, current, filterRegex.MatchString)
}
