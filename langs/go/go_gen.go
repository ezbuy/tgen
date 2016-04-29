package gogen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const langName = "go"

type GoGen struct {
	langs.BaseGen
}

func (this *GoGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init(langName, parsedThrift)

	outputPath, err := filepath.Abs(output)
	if err != nil {
		panicWithErr("fail to get absolute path for %q", output)
	}

	outputPackageDirs := make([]string, 0, len(parsedThrift))

	fmt.Println("##### Parsing:")
	for filename, parsed := range parsedThrift {
		fmt.Printf("%s\n", filename)

		importPath, pkgName := genNamespace(getNamespace(parsed.Namespaces))

		includes := getIncludes(parsedThrift, parsed.Includes)

		// make output dir
		pkgDir := filepath.Join(outputPath, importPath)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			panicWithErr("fail to make package directory %s", pkgDir)
		}

		outputPackageDirs = append(outputPackageDirs, pkgDir)

		// output defines file
		dataForDefinesFile := getDefinesFileData(pkgName, pkgDir, includes, parsed)
		if err := outputFile(dataForDefinesFile.FilePath, "defines_file", dataForDefinesFile); err != nil {
			panicWithErr("fail to write defines file %q : %s", dataForDefinesFile.FilePath, err)
		}

		// output webapi file
		for _, sData := range dataForDefinesFile.Services {
			dataForEchoModule := getEchoFileData(pkgName, pkgDir, dataForDefinesFile.Includes, sData)
			if err := outputFile(dataForEchoModule.FilePath, "echo_module", dataForEchoModule); err != nil {
				panicWithErr("fail to write web apis file %q : %s", dataForEchoModule.FilePath, err)
			}
		}
	}

	fmt.Println("##### gofmt")
	gofmt(outputPackageDirs...)
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
