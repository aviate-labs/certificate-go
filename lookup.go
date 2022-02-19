package cert

import "bytes"

func LookupLabel(tree Node, label Label) (Node, LabelResult) {
	switch n := tree.(type) {
	case Labeled:
		switch bytes.Compare(label, n.Label) {
		case 1: // a > b
			return nil, Continue
		case -1: // a < b
			return nil, Absent
		}
		return n.Tree, Found
	case Fork:
		switch v, lr := LookupLabel(n.LeftTree, label); lr {
		case Continue, Unknown:
			switch v, rr := LookupLabel(n.RightTree, label); rr {
			case Absent:
				if lr == Unknown {
					return nil, Unknown
				}
				return nil, Absent
			default:
				return v, rr
			}
		default:
			return v, lr
		}
	case Pruned:
		return nil, Unknown
	default:
		return nil, Absent
	}
}

type LabelResult string

const (
	Absent   LabelResult = "absent"
	Continue LabelResult = "continue"
	Found    LabelResult = "found"
	Unknown  LabelResult = "unknown"
)
