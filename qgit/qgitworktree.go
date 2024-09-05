package qgit

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitWt is a struct that implements the GitWorktree interface using the go-git library.
type GitWt struct {
	wt *git.Worktree
}

// GitWorktree defines an interface for performing worktree operations in a Git repository.
type GitWorktree interface {
	Checkout(opts *QgitCheckoutOptions) error
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
