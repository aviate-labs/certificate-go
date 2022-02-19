package cert_test

import (
	"testing"

	cert "github.com/aviate-labs/certificate-go"
)

func TestLookupLabel(t *testing.T) {
	for _, l := range []string{
		"a", "b", "c", "d",
	} {
		_, r := cert.LookupLabel(tree, []byte(l))
		if r != cert.Found {
			t.Error(r, l)
		}
	}

	for _, l := range []string{
		"x", "y", "e",
	} {
		_, r := cert.LookupLabel(tree, []byte(l))
		if r != cert.Continue {
			t.Error(r, l)
		}
	}

	for _, l := range []string{
		"absent",
	} {
		_, r := cert.LookupLabel(tree, []byte(l))
		if r != cert.Absent {
			t.Error(r, l)
		}
	}
}
