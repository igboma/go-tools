package gitpkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// Qgit is a struct that represents a Git repository and provides methods to interact with it.
type Qgit struct {
	option QgitOptions
	repo   GitRepository
}

// QReference holds information about a Git reference, such as the reference name and its hash.
type QReference struct {
	Hash          string
	ReferenceName string
}

// GitWt is a struct that implements the GitWorktree interface using the go-git library.
type GitWt struct {
	wt *git.Worktree
}

// QgitOptions contains the options required for initializing or cloning a Git repository.
type QgitOptions struct {
	Path   string
	Url    string
	IsBare bool
}

// QgitCheckoutOptions provides options for checking out a Git reference, including branches, tags, or commit hashes.
type QgitCheckoutOptions struct {
	branch string
	tag    string
	hash   string
}

// GitRepo is a struct that implements the GitRepository interface using the go-git library.
type GitRepo struct {
	repo *git.Repository
}

// Stat checks if the directory exists locally.
func (gr *GitRepo) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// GitRepository defines an interface for performing Git repository operations.
type GitRepository interface {
	//Init(o QgitOptions) error
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
}

// GitWorktree defines an interface for performing worktree operations in a Git repository.
type GitWorktree interface {
	Checkout(opts *QgitCheckoutOptions) error
}

// Head retrieves the current HEAD reference of the repository and returns it as a QReference.
func (gr *Qgit) Head() (QReference, error) {
	ref, err := gr.repo.Head()
	return ref, err
}

func (gr *Qgit) Fetch(ref string) error {
	return gr.repo.Fetch(ref)
}

func (gr *Qgit) PR() error {
	//gr.repo.PR();
	return nil
}

// Checkout checks out the specified Git reference (branch, tag, or commit hash) in the repository.
func (gr *Qgit) Checkout(ref string) error {

	//return gr.repo.Checkout(ref)
	isBranch, isTag, isCommitHash := gr.repo.CheckRemoteRef(ref)

	switch {
	case isBranch:
		fmt.Println("Checking out branch:", ref)
		return gr.repo.CheckoutBranch(ref)
	case isTag:
		fmt.Println("Checking out tag:", ref)
		return gr.repo.CheckoutTag(ref)
	case isCommitHash:
		fmt.Println("Checking out commit hash:", ref)
		return gr.repo.CheckoutHash(ref)
	default:
		return fmt.Errorf("reference not found: %s", ref)
	}
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

// NewQGit initializes a new Qgit instance with the provided options and GitRepository implementation.
// It returns an error if the repository initialization fails.
func NewQGit(o QgitOptions, gitRepo GitRepository) (*Qgit, error) {
	qgit := &Qgit{option: o, repo: gitRepo}
	err := qgit.Init()
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
func (gr *Qgit) Init() error {
	gitDir := filepath.Join(gr.option.Path, ".git")

	// Call the Stat method from the interface
	if _, err := gr.repo.Stat(gitDir); os.IsNotExist(err) {
		fmt.Println("Repository does not exist locally. Cloning...")
		err := gr.repo.PlainClone(gr.option)
		if err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
		fmt.Println("Repository cloned successfully.")
		//fmt.Println("Repository cloned successfully.")
	} else if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	} else {
		fmt.Println("Repository exists locally. Opening...")
		err := gr.repo.PlainOpen(gr.option)
		if err != nil {
			return fmt.Errorf("error opening repository: %w", err)
		}
		//fmt.Println("Repository opened successfully.")
	}
	return nil
}

// Init initializes the Git repository by either cloning it from a remote URL or opening it locally.
// func (gr *GitRepo) Init(o QgitOptions) error {
// 	gitDir := filepath.Join(o.Path, ".git")

// 	// Call the Stat method from the interface
// 	if _, err := gr.Stat(gitDir); os.IsNotExist(err) {
// 		fmt.Println("Repository does not exist locally. Cloning...")
// 		err := gr.PlainClone(o)
// 		if err != nil {
// 			return fmt.Errorf("error cloning repository: %w", err)
// 		}
// 		//fmt.Println("Repository cloned successfully.")
// 	} else if err != nil {
// 		return fmt.Errorf("error checking repository: %w", err)
// 	} else {
// 		fmt.Println("Repository exists locally. Opening...")
// 		err := gr.PlainOpen(o)
// 		if err != nil {
// 			return fmt.Errorf("error opening repository: %w", err)
// 		}
// 		//fmt.Println("Repository opened successfully.")
// 	}
// 	return nil
// }

func (gr *GitRepo) PlainClone(o QgitOptions) error {
	fmt.Println("Repository does not exist locally. Cloning...")
	var err error = nil
	repo, err := git.PlainClone(o.Path, o.IsBare, &git.CloneOptions{
		URL: o.Url,
	})
	gr.repo = repo
	return err
}

func (gr *GitRepo) PlainOpen(o QgitOptions) error {
	fmt.Println("Repository exists locally. Opening...")
	repo, err := git.PlainOpen(o.Path)
	if err != nil {
		return fmt.Errorf("error opening repository: %w", err)
	}
	gr.repo = repo
	fmt.Println("Repository opened successfully.")

	return nil
}

// Head retrieves the current HEAD reference of the repository.
func (gr *GitRepo) Head() (QReference, error) {
	ref, err := gr.repo.Head()
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
	wt, err := gr.repo.Worktree()
	if err != nil {
		return nil, err
	}
	return &GitWt{wt}, nil
}

// Fetch fetches changes from the remote repository, based on the specified refSpec.
func (gr *GitRepo) Fetch(refSpecStr string) error {
	remote, err := gr.repo.Remote("origin")
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
	remote, err := gr.repo.Remote("origin")
	if err != nil {
		fmt.Printf("Error getting remote: %v\n", err)
		return false, false, false
	}

	refs, err := remote.List(&git.ListOptions{})
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
		_, err := gr.repo.CommitObject(plumbing.NewHash(ref))
		if err == nil {
			isCommitHash = true
		} else {
			fmt.Printf("Error fetching commit object: %v\n", err)
		}
	}

	return isBranch, isTag, isCommitHash
}

// Checkout checks out a branch, tag, or commit hash in the worktree based on the options provided.
func (gwt *GitWt) Checkout(opts *QgitCheckoutOptions) error {
	checkoutOpt := git.CheckoutOptions{}

	if opts.branch != "" {
		branchRefName := plumbing.NewBranchReferenceName(opts.branch)
		checkoutOpt = git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchRefName),
			Force:  true,
		}
		return gwt.wt.Checkout(&checkoutOpt)

	} else if opts.tag != "" {
		tagRefName := plumbing.NewTagReferenceName(opts.tag)
		checkoutOpt = git.CheckoutOptions{
			Branch: tagRefName,
			Force:  true,
		}
		return gwt.wt.Checkout(&checkoutOpt)
	} else if opts.hash != "" {
		checkoutOpt = git.CheckoutOptions{
			Hash:  plumbing.NewHash(opts.hash),
			Force: true,
		}
		return gwt.wt.Checkout(&checkoutOpt)
	}
	return fmt.Errorf("unknown checkout operation")
}

// Fetch fetches changes from the remote repository, based on the specified refSpec.
func Fetch(refSpecStr string, gr *GitRepo) error {
	remote, err := gr.repo.Remote("origin")
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
