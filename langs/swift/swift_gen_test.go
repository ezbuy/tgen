package swift

import (
	"strings"
	"testing"

	// "github.com/ezbuy/tgen/langs/swift"

	"github.com/samuel/go-thrift/parser"
)

type testcase struct {
	infile   string
	contains [][]string
}

type testcases []testcase

var cases = testcases{
	testcase{
		infile: "./../../sample.thrift",
		contains: [][]string{
			{"var id: Int64", "var id: Int"},
		},
	},
	testcase{
		infile: "./../../a.thrift",
		contains: [][]string{
			{"dict[\"pendingWithdrawAmount\"] = pendingWithdrawAmount"},
			{"params[\"key\"] = key"},
		},
	},
}

func TestGen(t *testing.T) {
	p := &parser.Parser{}
	gen := &SwiftGen{}

	for _, c := range cases {
		// generate
		parsedThrift, _, err := p.ParseFile(c.infile)
		if err != nil {
			t.Errorf("parse error: %s\n", err.Error())
		}

		results, err := gen.Generate(parsedThrift)
		if err != nil {
			t.Errorf("generate error: %s\n", err.Error())
		}

		for idx, result := range results {
			strs := c.contains[idx]
			data := string(result.Data)

			for _, str := range strs {
				if strings.Contains(data, str) {
					continue
				}

				t.Errorf("mismatch found! [infile: %s, contain: %s]\n", result.Infile, str)
			}
		}
	}
}
