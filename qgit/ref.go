package qgit

import "github.com/go-git/go-git/v5/plumbing"

// Ref holds information about a Git reference, such as the reference name and its hash.
type Ref struct {
	Hash          string
	ReferenceName string
	IsBranch      bool
	IsTag         bool
	Name          string
}

// Refs converts go-git references to a slice of Ref.
func Refs(refs []*plumbing.Reference) []*Ref {
	var qRefs []*Ref

	for _, r := range refs {
		qRef := &Ref{
			Name:     r.Name().Short(),
			IsBranch: r.Name().IsBranch(),
			IsTag:    r.Name().IsTag(),
		}
		qRefs = append(qRefs, qRef)
	}

	return qRefs
}
