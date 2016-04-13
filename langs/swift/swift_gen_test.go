package swift

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ezbuy/tgen/global"
	"github.com/samuel/go-thrift/parser"
)

func TestGenerate(t *testing.T) {
	// 1 read thrift files from folder 'cases'
	// 2 generate & output
	// 3 read generated files, compared with corresponding files in folder 'test'

	casedir, _ := filepath.Abs("./../../example/swift")

	// create output dir
	outdir, _ := filepath.Abs("./output")
	if err := os.MkdirAll(outdir, 0755); err != nil {
		t.Errorf("failed to create output directory %s", outdir)
	}

	testdir, _ := filepath.Abs("./../../example/swift/ref")

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

		for f, thrift := range parsedThrift {
			for _, s := range thrift.Structs {
				baseSwift := &BaseSwift{Filepath: f, Thrift: thrift}
				name := fmt.Sprintf("%s.swift", baseSwift.AssembleStructName(s.Name))

				// jsonrpc
				outfile := filepath.Join(outdir, global.MODE_JSONRPC, name)
				testfile := filepath.Join(testdir, global.MODE_JSONRPC, name)

				fileCompare(t, outfile, testfile)

				// rest
				outfile = filepath.Join(outdir, global.MODE_REST, name)
				testfile = filepath.Join(testdir, global.MODE_REST, name)

				fileCompare(t, outfile, testfile)
			}

			for _, s := range thrift.Services {
				name := s.Name + "Service.swift"

				// jsonrpc
				outfile := filepath.Join(outdir, global.MODE_JSONRPC, name)
				testfile := filepath.Join(testdir, global.MODE_JSONRPC, name)

				fileCompare(t, outfile, testfile)

				// rest
				outfile = filepath.Join(outdir, global.MODE_REST, name)
				testfile = filepath.Join(testdir, global.MODE_REST, name)

				fileCompare(t, outfile, testfile)
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

func fileCompare(t *testing.T, src string, dest string) {
	if !pathexists(src) {
		t.Error("generate error\n")
	} else if !pathexists(dest) {
		t.Errorf("no test file found [%s]\n", dest)
	} else {
		// compare the output file with the case
		srcdata, srcerr := ioutil.ReadFile(src)
		destdata, desterr := ioutil.ReadFile(dest)

		if srcerr != nil || desterr != nil {
			t.Error("compare error [reading]")
		} else if string(srcdata) != string(destdata) {
			t.Errorf("mismatch: [%s, %s]", src, dest)
		} else {
			t.Log("PASS")
		}
	}
}

func pathexists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
