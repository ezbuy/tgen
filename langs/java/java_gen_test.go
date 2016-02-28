package java

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ezbuy/tgen/utils"
	"github.com/samuel/go-thrift/parser"
)

func TestGenerate(t *testing.T) {
	// 1 read thrift files from folder 'cases'
	// 2 generate & output
	// 3 read generated files, compared with corresponding files in folder 'test'

	casedir, _ := filepath.Abs(filepath.Dir("./cases/"))

	// create output dir
	outdir := filepath.Dir("./output/")
	if !utils.PathExists(outdir) {
		os.MkdirAll(outdir, 0775)
	}

	outdir, _ = filepath.Abs(outdir)
	testdir, _ := filepath.Abs("./test/")

	gen := &JavaGen{}
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

		for _, thrift := range parsedThrift {
			for _, m := range thrift.Structs {
				name := m.Name + ".java"

				outfile := filepath.Join(outdir, name)
				testfile := filepath.Join(testdir, name)

				fileCompare(t, outfile, testfile)
			}

			for _, s := range thrift.Services {
				name := s.Name + "Service.java"

				outfile := filepath.Join(outdir, name)
				testfile := filepath.Join(testdir, name)

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
	if !utils.PathExists(src) {
		t.Error("geenerate error\n")
	} else if !utils.PathExists(dest) {
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
