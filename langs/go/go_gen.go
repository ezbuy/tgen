package gogen

import (
	"fmt"
	"os"

	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const langName = "go"

type GoGen struct {
	langs.BaseGen
}

func (this *GoGen) Generate(parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init(langName, parsedThrift)

	// for filename, parsed := range parsedThrift {
	// 	for nKey, nValue := range parsed.Namespaces {
	// 		if nKey == langName {
	// 			fmt.Printf("namespace: %s\n", nValue)
	// 		}
	// 	}

	// 	fmt.Printf("name: %s\n", filename)
	// 	fmt.Printf("include: %s\n", parsed.Includes)
	// 	for structName, pStruct := range parsed.Structs {
	// 		fmt.Printf("struct name %s\n", structName)
	// 		fmt.Printf("struct structname %s\n", pStruct.Name)
	// 		for _, field := range pStruct.Fields {
	// 			fmt.Printf("field name %s\n", field.Name)

	// 			typ := field.Type
	// 			fmt.Println("=======")
	// 			fmt.Printf("field type %s\n", typ.Name)
	// 			if typ.KeyType != nil {
	// 				fmt.Printf("field type key %s \n", typ.KeyType.Name)
	// 			}

	// 			if typ.ValueType != nil {
	// 				fmt.Printf("field type value %s\n", typ.ValueType.Name)
	// 			}

	// 			fmt.Println("=======")
	// 		}
	// 	}

	// 	fmt.Println(">>>>>>>>>>>>>>")

	// }

	for filename, parsed := range parsedThrift {
		fmt.Printf("Parsing: %s >>>>>>>>>>>>>>>\n", filename)

		importPath, pkgName := genNamespace(getNamespace(parsed.Namespaces))
		fmt.Printf("import path: %s\n", importPath)

		includes := getIncludes(parsedThrift, parsed.Includes)

		data := &structsFileData{
			Package:  pkgName,
			Includes: includes,
		}

		for structName, parsedStruct := range parsed.Structs {
			data.Structs = append(data.Structs, &structData{
				Name:   upperHead(structName),
				Fields: parsedStruct.Fields,
			})
		}

		tpl.ExecuteTemplate(os.Stderr, "structs_file", data)

		fmt.Printf("<<<<<<<<<<< Parsed %s\n", filename)
	}
}

func init() {
	langs.Langs[langName] = &GoGen{}
}
