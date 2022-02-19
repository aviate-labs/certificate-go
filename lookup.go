package cert

import "bytes"

func LookupLabel(tree HashTree, label Label) (HashTree, LabelResult) {
	switch t := tree.(type) {
	case Labeled:
		switch bytes.Compare(label, t.Label) {
		case 1: // a > b
			return nil, Continue
		case -1: // a < b
			return nil, Absent
		}
		return t.Tree, Found
	case Fork:
		switch v, lr := LookupLabel(t.LeftTree, label); lr {
		case Continue, Unknown:
			switch v, rr := LookupLabel(t.RightTree, label); rr {
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
