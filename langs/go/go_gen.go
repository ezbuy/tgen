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

		// output struct file
		if len(parsed.Structs) != 0 {
			dataForStructsFile := getStructsFileData(pkgName, pkgDir, includes, parsed.Structs)

			if err := outputFile(dataForStructsFile.FilePath, "structs_file", dataForStructsFile); err != nil {
				panicWithErr("fail to write structs file %q : %s", dataForStructsFile.FilePath, err)
			}
		}

		// output service file
		if len(parsed.Services) != 0 {
			dataForServicesFile := getServicesFileData(pkgName, pkgDir, includes, parsed.Services)

			if err := outputFile(dataForServicesFile.FilePath, "services_file", dataForServicesFile); err != nil {
				panicWithErr("fail to write services file %q : %s", dataForServicesFile.FilePath, err)
			}

			for _, sData := range dataForServicesFile.Services {
				dataForEchoModule := getEchoFileData(pkgName, pkgDir, sData)
				if err := outputFile(dataForEchoModule.FilePath, "echo_module", dataForEchoModule); err != nil {
					panicWithErr("fail to write web apis file %q : %s", dataForEchoModule.FilePath, err)
				}
			}
		}
	}

	fmt.Println("##### gofmt")
	gofmt(outputPackageDirs...)
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
