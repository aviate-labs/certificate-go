package cert_test

import (
	"fmt"

	cert "github.com/aviate-labs/certificate-go"
)

var tree = cert.Fork{
	LeftTree: cert.Fork{
		LeftTree: cert.Labeled{
			Label: []byte("a"),
			Tree: cert.Fork{
				LeftTree: cert.Fork{
					LeftTree: cert.Labeled{
						Label: []byte("x"),
						Tree:  cert.Leaf("hello"),
					},
					RightTree: cert.Empty{},
				},
				RightTree: cert.Labeled{
					Label: []byte("y"),
					Tree:  cert.Leaf("world"),
				},
			},
		},
		RightTree: cert.Labeled{
			Label: []byte("b"),
			Tree:  cert.Leaf("good"),
		},
	},
	RightTree: cert.Fork{
		LeftTree: cert.Labeled{
			Label: []byte("c"),
			Tree:  cert.Empty{},
		},
		RightTree: cert.Labeled{
			Label: []byte("d"),
			Tree:  cert.Leaf("morning"),
		},
	},
}

func ExampleB() {
	fmt.Printf("%X", cert.Leaf("good").Reconstruct())
	// Output:
	// 7B32AC0C6BA8CE35AC82C255FC7906F7FC130DAB2A090F80FE12F9C2CAE83BA6
}

func ExampleC() {
	fmt.Printf("%X", cert.Labeled{
		Label: []byte("c"),
		Tree:  cert.Empty{},
	}.Reconstruct())
	// Output:
	// EC8324B8A1F1AC16BD2E806EDBA78006479C9877FED4EB464A25485465AF601D
}

func ExampleRoot() {
	fmt.Printf("%X", tree.Reconstruct())
	// Output:
	// EB5C5B2195E62D996B84C9BCC8259D19A83786A2F59E0878CEC84C811F669AA0
}

// Source: https://sdk.dfinity.org/docs/interface-spec/index.html#_example
// ─┬─┬╴"a" ─┬─┬╴"x" ─╴"hello"
//  │ │      │ └╴Empty
//  │ │      └╴  "y" ─╴"world"
//  │ └╴"b" ──╴"good"
//  └─┬╴"c" ──╴Empty
//    └╴"d" ──╴"morning"
func ExampleX() {
	fmt.Printf("%X", cert.Fork{
		LeftTree: cert.Labeled{
			Label: []byte("x"),
			Tree:  cert.Leaf("hello"),
		},
		RightTree: cert.Empty{},
	}.Reconstruct())
	// Output:
	// 1B4FEFF9BEF8131788B0C9DC6DBAD6E81E524249C879E9F10F71CE3749F5A638
}
