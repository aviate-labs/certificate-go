package cert

import "crypto/sha256"

func domainSeperator(t string) []byte {
	return append(
		[]byte{uint8(len(t))},
		[]byte(t)...,
	)
}

type Empty struct{}

func (e Empty) Reconstruct() [32]byte {
	return sha256.Sum256(domainSeperator("ic-hashtree-empty"))
}

type Fork struct {
	LeftTree  HashTree
	RightTree HashTree
}

func (f Fork) Reconstruct() [32]byte {
	l := f.LeftTree.Reconstruct()
	r := f.RightTree.Reconstruct()
	return sha256.Sum256(append(
		domainSeperator("ic-hashtree-fork"),
		append(l[:], r[:]...)...,
	))
}

type HashTree interface {
	Reconstruct() [32]byte
}

type Label []byte

type Labeled struct {
	Label Label
	Tree  HashTree
}

func (l Labeled) Reconstruct() [32]byte {
	t := l.Tree.Reconstruct()
	return sha256.Sum256(append(
		domainSeperator("ic-hashtree-labeled"),
		append(l.Label, t[:]...)...,
	))
}

type Leaf []byte

func (l Leaf) Reconstruct() [32]byte {
	return sha256.Sum256(append(
		domainSeperator("ic-hashtree-leaf"),
		l...,
	))
}

type Pruned [32]byte

func (p Pruned) Reconstruct() [32]byte {
	return p
}
