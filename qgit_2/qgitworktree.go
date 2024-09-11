package qgit_2

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitWorkTree is a struct that provides an implementation of the GitWorktree interface using the go-git library.
// It wraps a go-git Worktree, allowing operations such as branch, tag, or commit checkout.
type GitWorkTree struct {
	wt *git.Worktree // wt is the underlying go-git Worktree that performs Git operations.
}

// GitWorktree defines an interface for performing operations in a Git repository's worktree.
type Worktree interface {
	// Checkout checks out a branch, tag, or commit hash in the worktree based on the provided options.
	//
	Checkout(opts *QRepoCheckoutOptions) error
}

// Checkout checks out a branch, tag, or commit hash in the Git worktree based on the options provided.
//
// It uses the following logic:
//   - If a branch name is provided in opts, it checks out the branch.
//   - If a tag name is provided in opts, it checks out the tag.
//   - If a commit hash is provided in opts, it checks out the specific commit.
//   - If none of these are provided, it returns an error indicating an unknown checkout operation.
//
// Parameters:
//   - opts: A pointer to QRepoCheckoutOptions containing the branch, tag, or commit hash to check out.
//
// Returns:
//   - error: Returns an error if the checkout operation fails, or nil if successful.
func (gwt *GitWorkTree) Checkout(opts *QRepoCheckoutOptions) error {
	checkoutOpt := git.CheckoutOptions{}

	// Check out a branch if the branch name is provided
	if opts.branch != "" {
		branchRefName := plumbing.NewBranchReferenceName(opts.branch)
		checkoutOpt = git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchRefName),
			Force:  true,
		}
		return gwt.wt.Checkout(&checkoutOpt)

		// Check out a tag if the tag name is provided
	} else if opts.tag != "" {
		tagRefName := plumbing.NewTagReferenceName(opts.tag)
		checkoutOpt = git.CheckoutOptions{
			Branch: tagRefName,
			Force:  true,
		}
		return gwt.wt.Checkout(&checkoutOpt)

		// Check out a specific commit if the hash is provided
	} else if opts.hash != "" {
		checkoutOpt = git.CheckoutOptions{
			Hash:  plumbing.NewHash(opts.hash),
			Force: true,
		}
		return gwt.wt.Checkout(&checkoutOpt)
	}

	// Return an error if none of branch, tag, or hash is provided
	return fmt.Errorf("unknown checkout operation")
}
