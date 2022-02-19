package cert

import (
	"crypto/sha256"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

func Serialize(node Node) ([]byte, error) {
	return cbor.Marshal(serialize(node))
}

func domainSeperator(t string) []byte {
	return append(
		[]byte{uint8(len(t))},
		[]byte(t)...,
	)
}

func serialize(node Node) []interface{} {
	switch n := node.(type) {
	case Empty:
		return []interface{}{0x00}
	case Fork:
		return []interface{}{
			0x01,
			serialize(n.LeftTree),
			serialize(n.RightTree),
		}
	case Labeled:
		return []interface{}{
			0x02,
			[]byte(n.Label),
			serialize(n.Tree),
		}
	case Leaf:
		return []interface{}{
			0x03,
			[]byte(n),
		}
	case Pruned:
		return []interface{}{
			0x04,
			n,
		}
	}
	return nil
}

type Empty struct{}

func (e Empty) Reconstruct() [32]byte {
	return sha256.Sum256(domainSeperator("ic-hashtree-empty"))
}

type Fork struct {
	LeftTree  Node
	RightTree Node
}

func (f Fork) Reconstruct() [32]byte {
	l := f.LeftTree.Reconstruct()
	r := f.RightTree.Reconstruct()
	return sha256.Sum256(append(
		domainSeperator("ic-hashtree-fork"),
		append(l[:], r[:]...)...,
	))
}

type Label []byte

func (l Label) String() string {
	return string(l)
}

type Labeled struct {
	Label Label
	Tree  Node
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

type Node interface {
	Reconstruct() [32]byte
}

func Deserialize(data []byte) (Node, error) {
	var s []interface{}
	if err := cbor.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return deserialize(s)
}

func deserialize(s []interface{}) (Node, error) {
	tag, ok := s[0].(uint64)
	if !ok {
		return nil, fmt.Errorf("unknown tag: %v", s[0])
	}
	switch tag {
	case 0:
		if l := len(s); l != 1 {
			return nil, fmt.Errorf("invalid len: %d", l)
		}
		return Empty{}, nil
	case 1:
		if l := len(s); l != 3 {
			return nil, fmt.Errorf("invalid len: %d", l)
		}
		lt, ok := s[1].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[1])
		}
		l, err := deserialize(lt)
		if err != nil {
			return nil, err
		}
		rt, ok := s[2].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[2])
		}
		r, err := deserialize(rt)
		if err != nil {
			return nil, err
		}
		return Fork{
			LeftTree:  l,
			RightTree: r,
		}, nil
	case 2:
		if l := len(s); l != 3 {
			return nil, fmt.Errorf("invalid len: %d", l)
		}
		l, ok := s[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[1])
		}
		rt, ok := s[2].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[2])
		}
		t, err := deserialize(rt)
		if err != nil {
			return nil, err
		}
		return Labeled{
			Label: l,
			Tree:  t,
		}, nil
	case 3:
		if l := len(s); l != 2 {
			return nil, fmt.Errorf("invalid len: %d", l)
		}
		l, ok := s[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[1])
		}
		return Leaf(l), nil
	case 4:
		if l := len(s); l != 2 {
			return nil, fmt.Errorf("invalid len: %d", l)
		}
		p, ok := s[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("unknown value: %v", s[1])
		}
		if l := len(p); l != 32 {
			return nil, fmt.Errorf("invalid hash len: %d", l)
		}
		var p32 [32]byte
		copy(p32[:], p)
		return Pruned(p32), nil
	default:
		return nil, fmt.Errorf("invalid tag: %d", tag)
	}
}

type Pruned [32]byte

func (p Pruned) Reconstruct() [32]byte {
	return p
}
