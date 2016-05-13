package gogen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezbuy/tgen/global"
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

	targetPkg, ok := packages[global.InputFile]
	if !ok {
		exitWithError("target package for %q not found in %#v", global.InputFile, packages)
	}

	fmt.Printf("##### Generating: %s\n", global.InputFile)

	// make output dir
	pkgDir := filepath.Join(outputPath, targetPkg.ImportPath)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		exitWithError("fail to make package directory %s\n", pkgDir)
	}

	if err := targetPkg.renderToFile(pkgDir, "defines", "defines_file"); err != nil {
		exitWithError("fail to write defines file: %s\n", err)
	}

	if global.GenRpcClient {
		fmt.Printf("##### Generating Rpc Client File")
		if err := targetPkg.renderToFile(pkgDir, "rpc_client", "rpc_client"); err != nil {
			exitWithError("fail to write rpcclient file: %s\n", err)
		}
	}

	if global.GenWebApi {
		fmt.Printf("##### Generating Rpc WebApis File")
		if err := targetPkg.renderToFile(pkgDir, "webapis", "echo_module"); err != nil {
			exitWithError("fail to write webapis file: %s\n", err)
		}
	}
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
