package gogen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const (
	langName = "go"
)

type GoGen struct {
	langs.BaseGen
}

func (this *GoGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init(langName, parsedThrift)

	outputPath, err := filepath.Abs(output)
	if err != nil {
		exitWithError("fail to get absolute path for %q", output)
	}

	outputPackageDirs := make([]string, 0, len(parsedThrift))

	fmt.Println("##### Parsing:")

	packages := map[string]*Package{}

	// setup packages
	for filename, parsed := range parsedThrift {
		pkg := newPackage(parsed)
		packages[filename] = pkg
	}

	// setup includes
	for _, pkg := range packages {
		pkg.setupIncludes(packages)
	}

	for filename, pkg := range packages {
		fmt.Printf("##### Generating: %s\n", filename)

		// make output dir
		pkgDir := filepath.Join(outputPath, pkg.ImportPath)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			exitWithError("fail to make package directory %s\n", pkgDir)
		}

		outputPackageDirs = append(outputPackageDirs, pkgDir)

		if err := pkg.renderToFile(pkgDir, "defines", "defines_file"); err != nil {
			exitWithError("fail to write defines file: %s\n", err)
		}

		if err := pkg.renderToFile(pkgDir, "webapis", "echo_module"); err != nil {
			exitWithError("fail to write webapis file: %s\n", err)
		}
	}

	fmt.Println("##### gofmt")
	gofmt(outputPackageDirs...)
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
