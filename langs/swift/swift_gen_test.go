package swift

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/samuel/go-thrift/parser"
)

func TestGenerate(t *testing.T) {
	// 1 read thrift files from folder 'cases'
	// 2 generate & output
	// 3 read generated files, compared with corresponding files in folder 'test'

	casedir, _ := filepath.Abs(filepath.Dir("./cases/"))

	// create output dir
	outdir := filepath.Dir("./output/")
	if err := os.MkdirAll(outdir, 0755); err != nil {
		t.Errorf("failed to create output directory %s", outdir)
	}

	outdir, _ = filepath.Abs(outdir)
	testdir, _ := filepath.Abs("./test/")

	gen := &SwiftGen{}
	p := &parser.Parser{}

	visitfunc := func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), ".") || filepath.Ext(path) != ".thrift" {
			return nil
		}

		parsedThrift, _, err := p.ParseFile(path)
		if err != nil {
			t.Errorf("parse error: %s\n", err.Error())
		}

		gen.Generate(outdir, parsedThrift)

		for tp, thrift := range parsedThrift {
			name := gen.filename(tp, thrift.Namespaces)

			outfile := filepath.Join(outdir, name)
			testfile := filepath.Join(testdir, name)

			if !pathExists(outfile) {
				t.Errorf("geenerate error: thrift [%s]\n", tp)
			} else if !pathExists(testfile) {
				t.Errorf("no test file found [%s]\n", testfile)
			} else {
				// compare the output file with the case
				outdata, outerr := ioutil.ReadFile(outfile)
				testdata, testerr := ioutil.ReadFile(testfile)

				if outerr != nil || testerr != nil {
					t.Error("compare error [reading]")
				} else if string(outdata) != string(testdata) {
					t.Errorf("mismatch: [%s, %s]", testfile, outfile)
				} else {
					t.Log("PASS")
				}
			}
		}

		return nil
	}

	if err := filepath.Walk(casedir, visitfunc); err != nil {
		t.Errorf("walk error: %s\n", err.Error())
	}

	// do some clean
	os.RemoveAll(outdir)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
