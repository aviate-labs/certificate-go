package cert

import "github.com/aviate-labs/principal-go"

type Certificate struct {
	Tree       HashTree
	Signature  []byte
	Delegation *Delegation
}

type Delegation struct {
	SubnetId principal.Principal
	// The nested certificate typically does not itself again contain a
	// delegation, although there is no reason why agents should enforce that
	// property.
	Certificate Certificate
}
