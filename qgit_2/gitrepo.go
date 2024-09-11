package qgit_2

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// plainClone clones a Git repository from a remote URL to the specified local path using basic authentication.
func plainClone(o QRepoOptions) (*git.Repository, error) {
	repo, err := git.PlainClone(o.Path, false, &git.CloneOptions{
		URL: o.Url,
		Auth: &http.BasicAuth{
			Username: "git",   // Can be anything, GitHub ignores it
			Password: o.Token, // Use the token as the password
		},
	})
	return repo, err
}

// plainOpen opens an existing Git repository from the specified local path.
func plainOpen(o QRepoOptions) (*git.Repository, error) {
	repo, err := git.PlainOpen(o.Path)
	if err != nil {
		return nil, fmt.Errorf("error opening repository: %w", err)
	}

	return repo, err
}

// getRepo checks if a Git repository exists at the specified local path.
// If it doesn't exist, it clones the repository. Otherwise, it opens the existing repository.
func getRepo(options QRepoOptions) (git *git.Repository, err error) {
	gitDir := filepath.Join(options.Path, ".git")

	// Call the Stat method from the interface to check if the repository exists
	if _, err := Stat(gitDir); os.IsNotExist(err) {
		fmt.Println("Repository does not exist locally. Cloning...")
		git, err = plainClone(options)
		if err != nil {
			return nil, fmt.Errorf("error cloning repository: %w", err)
		}
		fmt.Println("Repository cloned successfully.")
	} else if err != nil {
		return nil, fmt.Errorf("error checking repository: %w", err)
	} else {
		fmt.Println("Repository exists locally. Opening...")
		git, err = plainOpen(options)
		if err != nil {
			return nil, fmt.Errorf("error opening repository: %w", err)
		}
		fmt.Println("Repository opened successfully.")
	}
	return git, nil
}

// QReference holds information about a Git reference, such as the reference name and its hash.
type QReference struct {
	Hash          string
	ReferenceName string
	IsBranch      bool
	IsTag         bool
	Name          string
}

// QRepoOptions contains the options required for initializing or cloning a Git repository.
type QRepoOptions struct {
	Path  string
	Url   string
	Token string
}

// QRepoCheckoutOptions provides options for checking out a Git reference, including branches, tags, or commit hashes.
type QRepoCheckoutOptions struct {
	branch string
	tag    string
	hash   string
}

// GitRepo is a struct that implements the Repository interface using the go-git library.
type QGitRepo struct {
	option *QRepoOptions
}

// Stat checks if the directory exists locally.
func (gr *QGitRepo) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Repository defines an interface for performing Git repository operations.
type Repository interface {
	Head() (QReference, error)
	SetOption(option *QRepoOptions)
	Option() *QRepoOptions
	Worktree() (Worktree, error)
	PlainClone(o QRepoOptions) error
	PlainOpen(o QRepoOptions) error
	Fetch(refSpecStr string) error
	Checkout(ref string) error
	CheckoutBranch(branch string) error
	CheckoutTag(tag string) error
	CheckoutHash(hash string) error
	CheckRemoteRef(ref string) (isBranch, isTag, isCommitHash bool, err error)
	GetFileContentFromBranch(branch, file string) (string, error)
	GetFileContentFromCommit(commitHash, file string) (string, error)
	GetChangedFilesByPRNumber(prNumber int) ([]string, error)
}

// Checkout checks out the specified Git reference (branch, tag, or commit hash) in the repository.
func (gr *QGitRepo) Checkout(ref string) error {

	isBranch, isTag, isCommitHash, err := gr.CheckRemoteRef(ref)
	if err != nil {
		return fmt.Errorf("error conneting to repo %v", err)
	}
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

// PlainClone clones a Git repository from a remote URL to the specified local path using basic authentication.
//
// Parameters:
//   - o: QRepoOptions struct containing the repository URL, path, and authentication token.
//
// Returns:
//   - error: Returns an error if the cloning process fails, otherwise nil.
func (gr *QGitRepo) PlainClone(o QRepoOptions) error {
	// Clone the Git repository to the specified path
	_, err := git.PlainClone(o.Path, false, &git.CloneOptions{
		URL: o.Url, // URL of the remote Git repository
		Auth: &http.BasicAuth{
			Username: "git",   // GitHub ignores the username, can be anything
			Password: o.Token, // Use the provided token as the password for authentication
		},
	})

	// Return the error, if any
	return err
}

// PlainOpen opens an existing Git repository from the specified local path.
//
// Parameters:
//   - o: QRepoOptions struct containing the local repository path.
//
// Returns:
//   - error: Returns an error if the repository cannot be opened, otherwise nil.
func (gr *QGitRepo) PlainOpen(o QRepoOptions) error {
	// Attempt to open the Git repository from the specified local path
	_, err := git.PlainOpen(o.Path)
	if err != nil {
		// If an error occurs, wrap it with additional context and return it
		return fmt.Errorf("error opening repository: %w", err)
	}

	// Return nil if the repository opens successfully
	return nil
}

// SetOption sets the repository options (such as path, URL, and token) for the Git repository.
// This allows configuration of the repository parameters.
//
// Parameters:
//   - opt: A pointer to QRepoOptions struct containing the options for initializing or cloning the repository.
func (gr *QGitRepo) SetOption(opt *QRepoOptions) {
	gr.option = opt
}

// Option retrieves the current repository options (such as path, URL, and token) used for the Git repository.
//
// Returns:
//   - *QRepoOptions: A pointer to the QRepoOptions struct containing the repository options.
func (gr *QGitRepo) Option() *QRepoOptions {
	return gr.option
}

// Head retrieves the current HEAD reference of the repository.
func (gr *QGitRepo) Head() (QReference, error) {
	repo, err := getRepo(*gr.Option())
	if err != nil {
		return QReference{}, fmt.Errorf("error conneting to repo %w", err)
	}
	ref, err := repo.Head()
	if err != nil {
		return QReference{}, err
	}
	return QReference{
		ReferenceName: ref.Name().String(),
		Hash:          ref.Hash().String(),
	}, nil
}

// Worktree retrieves the working tree (the directory containing the checked-out files) of the repository.
//
// Returns:
//   - GitWorktree: A reference to the Git working tree, allowing for file system operations such as checking out branches.
//   - error: Returns an error if the repository connection fails or if the working tree cannot be accessed.
func (gr *QGitRepo) Worktree() (Worktree, error) {
	// Get the repository instance using the provided options
	repo, err := getRepo(*gr.Option())
	if err != nil {
		// Return an error if unable to connect to the repository
		return nil, fmt.Errorf("error connecting to repo: %v", err)
	}

	// Retrieve the working tree (the directory where the repository is checked out)
	wt, err := repo.Worktree()
	if err != nil {
		// Return an error if the working tree cannot be accessed
		return nil, err
	}

	// Return the working tree wrapped in a GitWt struct
	return &GitWorkTree{wt}, nil
}

// Fetch fetches updates from the remote repository, ensuring that the specified references are up to date.
// if refSpec is empty all refSpecs are fetched by default
//
// Parameters:
//   - refSpecStr: A string specifying the refspec to fetch. If empty, default refspecs for branches and tags are used.
//
// Returns:
//   - error: Returns an error if fetching from the remote fails, otherwise nil.
func (gr *QGitRepo) Fetch(refSpecStr string) error {
	// Get the repository instance using the provided options
	repo, err := getRepo(*gr.Option())
	if err != nil {
		// Return an error if unable to connect to the repository
		return fmt.Errorf("error connecting to repo: %w", err)
	}

	// Get the remote repository (assumed to be named "origin")
	remote, err := repo.Remote("origin")
	if err != nil {
		// Return an error if the remote "origin" is not found
		return err
	}

	// Define refspecs (specifying which refs to fetch)
	var refSpecs []config.RefSpec
	if refSpecStr == "" {
		// If no specific refspec is provided, fetch all branches and tags
		refSpecs = []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"), // Fetch all branches
			config.RefSpec("+refs/tags/*:refs/tags/*"),            // Fetch all tags
		}
	} else {
		// Use the provided refspec
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	// Perform the fetch operation with the specified refspecs
	if err := remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		// If an error occurs and it's not the "already up-to-date" error, return the error
		return fmt.Errorf("fetch origin failed: %w", err)
	}

	// Return nil if fetching was successful
	return nil
}

// CheckoutBranch checks out the specified branch in the repository's worktree.
func (gr *QGitRepo) CheckoutBranch(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	checkoutOption := QRepoCheckoutOptions{branch: ref}

	if err := wt.Checkout(&checkoutOption); err != nil {
		if err := gr.Fetch(fmt.Sprintf("refs/heads/%s:refs/heads/%s", ref, ref)); err != nil {
			return err
		}
		return wt.Checkout(&checkoutOption)
	}
	return nil
}

// CheckoutTag checks out the specified tag in the repository's worktree.
func (gr *QGitRepo) CheckoutTag(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	err = gr.Fetch(fmt.Sprintf("+refs/tags/%s:refs/tags/%s", ref, ref))
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching tag %s: %w", ref, err)
	}

	checkoutOption := QRepoCheckoutOptions{tag: ref}

	err = wt.Checkout(&checkoutOption)
	if err != nil {
		return fmt.Errorf("error checking out tag %s: %w", ref, err)
	}
	return nil
}

// CheckoutHash checks out the specified commit hash in the repository's worktree.
func (gr *QGitRepo) CheckoutHash(ref string) error {
	wt, err := gr.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&QRepoCheckoutOptions{hash: ref})
	if err != nil {
		return fmt.Errorf("error checking out commit %s: %w", ref, err)
	}

	return nil
}

// checkRemoteRef checks if the specified reference is a branch, tag, or commit hash by querying the remote and local repository.
// It returns three boolean flags indicating whether the ref is a branch, tag, or commit hash.
func (gr *QGitRepo) CheckRemoteRef(ref string) (isBranch, isTag, isCommitHash bool, err error) {
	repo, err := getRepo(*gr.Option())
	if err != nil {
		err = fmt.Errorf("error connecting to repo: %w", err)
		return false, false, false, err
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		err = fmt.Errorf("error getting remote: %w", err)
		return false, false, false, err
	}

	// List the remote references
	refs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: "git",             // GitHub ignores the username but requires it
			Password: gr.Option().Token, // Use the token for authentication
		},
	})
	if err != nil {
		err = fmt.Errorf("error listing remote references: %w", err)
		return false, false, false, err
	}

	// Convert go-git references to QReference
	qReferences := gr.convertToQReferences(refs)

	// Classify the ref
	return gr.classifyRef(ref, qReferences)
}

// convertToQReferences converts go-git references to a slice of QReference.
func (gr *QGitRepo) convertToQReferences(refs []*plumbing.Reference) []*QReference {
	var qRefs []*QReference

	for _, r := range refs {
		qRef := &QReference{
			Name:     r.Name().Short(),
			IsBranch: r.Name().IsBranch(),
			IsTag:    r.Name().IsTag(),
		}
		qRefs = append(qRefs, qRef)
	}

	return qRefs
}

// classifyRef checks the given reference against the list of QReferences and identifies it as a branch, tag, or commit hash.
func (gr *QGitRepo) classifyRef(ref string, refs []*QReference) (isBranch, isTag, isCommitHash bool, err error) {
	repo, err := getRepo(*gr.Option())
	if err != nil {
		//fmt.Printf("Error conneting to repo: %w\n", err)
		return
	}

	// Check if the reference matches any branches or tags
	for _, r := range refs {
		if r.IsBranch && r.Name == ref {
			isBranch = true
		} else if r.IsTag && r.Name == ref {
			isTag = true
		}
	}

	// If the ref has a length of 40, check if it is a commit hash
	if len(ref) == 40 {
		_, err := repo.CommitObject(plumbing.NewHash(ref))
		if err == nil {
			isCommitHash = true
		} else {
			fmt.Printf("error fetching commit object: %v\n", err)
		}
	}
	return isBranch, isTag, isCommitHash, err
}

// GetFileContentFromBranch retrieves the content of a specified file from the latest commit on a given branch.
//
// Parameters:
//   - branch: The branch name from which to retrieve the file (e.g., "refs/heads/main").
//   - file: The path of the file whose content needs to be retrieved.
//
// Returns:
//   - content: The content of the specified file from the latest commit on the branch.
//   - error: Returns an error if the branch or file cannot be found, or if reading the file content fails.
func (gr *QGitRepo) GetFileContentFromBranch(branch, file string) (content string, err error) {
	// Get the repository instance
	repo, err := getRepo(*gr.Option())
	if err != nil {
		return
	}

	// Get the reference to the specified branch
	ref, err := repo.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		return "", err
	}

	// Get the commit object for the latest commit on the branch
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	// Get the tree (file structure) associated with the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	// Find the file entry in the tree
	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	// Retrieve the file contents
	content, err = entry.Contents()
	if err != nil {
		return "", err
	}

	// Return the file content
	return content, nil
}

// GetFileContentFromCommit retrieves the content of a specified file from a given commit hash in the repository.
//
// Parameters:
//   - commitHash: The hash of the commit from which to retrieve the file.
//   - file: The path of the file whose content needs to be retrieved.
//
// Returns:
//   - content: The content of the specified file from the given commit.
//   - error: Returns an error if the commit or file cannot be found, or if reading the file content fails.
func (gr *QGitRepo) GetFileContentFromCommit(commitHash, file string) (content string, err error) {
	// Get the repository instance
	repo, err := getRepo(*gr.Option())
	if err != nil {
		return
	}

	// Get the commit object using the commit hash
	commit, err := repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return "", err
	}

	// Get the tree associated with the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	// Find the file entry in the tree
	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	// Retrieve the file contents
	content, err = entry.Contents()
	if err != nil {
		return "", err
	}

	// Return the file content
	return content, nil
}

// GetChangedFilesByPRNumber fetches the changed files between the main branch and the PR branch.
func (gr *QGitRepo) GetChangedFilesByPRNumber(prNumber int) (changedFiles []string, err error) {
	repo, err := getRepo(*gr.Option())
	if err != nil {
		return nil, fmt.Errorf("error conneting to repo: %w", err)
	}
	// Convert the PR number into a reference that exists in the Git repository
	// Usually PR references are in the form: refs/pull/{prNumber}/head
	prRef := fmt.Sprintf("refs/pull/%d/head", prNumber)

	// Fetch the remote branch (PR branch) to ensure the reference exists locally
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", prRef, prRef))},
		Auth: &http.BasicAuth{
			Username: "git", // GitHub ignores the username but requires it
			Password: gr.Option().Token,
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("failed to fetch remote branch %s: %w", prRef, err)
	}

	// Get the current HEAD reference (main branch)
	currentRef, err := repo.Head()

	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Resolve the comparison reference from the PR (origin)
	compareRef, err := repo.Reference(plumbing.ReferenceName(prRef), true)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve reference %s: %w", prRef, err)
	}
	// Get the commit for the comparison reference (PR branch)
	compareCommit, err := repo.CommitObject(compareRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for ref %s: %w", prRef, err)
	}

	// Get the commit for the current HEAD reference (main branch)
	currentCommit, err := repo.CommitObject(currentRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get current commit: %w", err)
	}

	// Get the file changes between the two commits
	patch, err := currentCommit.Patch(compareCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate patch: %w", err)
	}

	// Collect the list of changed files
	for _, fileStat := range patch.Stats() {
		changedFiles = append(changedFiles, fileStat.Name)
	}

	return changedFiles, nil
}

// NewGitRepo creates a new instance of GitRepo with the provided options.
//
// Parameters:
//   - options: A QRepoOptions struct that contains the repository configuration such as path, URL, and token.
//
// Returns:
//   - *QGitRepo: A pointer to a new instance of GitRepo initialized with the given options.
func NewGitRepo(options *QRepoOptions) Repository {
	return &QGitRepo{
		option: options,
	}
}
