package typescript

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"strings"

	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const TPL_SERVICE = "tgen/typescript/service"

const (
	langName = "ts"
)

type TypeScriptGen struct {
	langs.BaseGen
}

type Method struct {
	Name       string
	Arguments  string
	ReturnType string
}

type InterfaceField struct {
	Name string
	Type string
}

type Interface struct {
	Name   string
	Fields []*InterfaceField
}

type EnumVal struct {
	Name string
	Val  int
}

type Enum struct {
	Name   string
	Values []*parser.EnumValue
}

type Thrift struct {
	Methods    []*Method
	Interfaces []*Interface
	Includes   []string
	Enums      map[string]*parser.Enum
}

func initemplate(n string, path string) *template.Template {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}

	return tpl
}

func outputfile(fp string, t *template.Template, data interface{}) error {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return t.Execute(file, data)
}

func typeCast(t *parser.Type) string {
	if t != nil {
		switch t.Name {
		case langs.ThriftTypeI16, langs.ThriftTypeI32, langs.ThriftTypeI64, langs.ThriftTypeDouble:
			return "number"
		case langs.ThriftTypeString:
			return "string"
		case langs.ThriftTypeBool:
			return "boolean"
		case langs.ThriftTypeList, langs.ThriftTypeSet:
			valType := typeCast(t.ValueType)
			return valType + "[]"
		case langs.ThriftTypeMap:
			return "JSONObject"
		default:
			return t.Name
		}
	}
	return "null"
}

func genOutputPath(base string, fileName string) string {
	start := strings.LastIndex(fileName, "/")
	end := strings.LastIndex(fileName, ".")
	name := fileName[start+1 : end]
	return filepath.Join(base, name+"Service.ts")
}

func (this *TypeScriptGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	var servicetpl *template.Template
	servicetpl = initemplate(TPL_SERVICE, "./tmpl/typescript/rest_service.gots")

	for fileName, t := range parsedThrift {
		data := &Thrift{
			Methods:    make([]*Method, 0),
			Interfaces: make([]*Interface, 0),
			Includes:   make([]string, 0),
			Enums:      make(map[string]*parser.Enum),
		}
		outputPath := genOutputPath(output, fileName)

		// fill in Enums
		data.Enums = t.Enums

		// fill in Includes
		for name, _ := range t.Includes {
			data.Includes = append(data.Includes, name)
		}

		// fill in Methods
		for _, s := range t.Services {
			for mName, mVal := range s.Methods {
				m := &Method{}
				m.Name = mName

				args := make([]string, 0)
				for _, arg := range mVal.Arguments {
					argStr := arg.Name + ": " + typeCast(arg.Type)
					args = append(args, argStr)
				}
				m.Arguments = strings.Join(args, ", ")

				m.ReturnType = typeCast(mVal.ReturnType)

				data.Methods = append(data.Methods, m)
			}
		}

		// fill in Interfaces
		interfaces := make([]*Interface, 0)
		for _, s := range t.Structs {
			ife := &Interface{}
			ife.Name = s.Name

			fields := make([]*InterfaceField, 0)
			for _, rawFiled := range s.Fields {
				field := &InterfaceField{}
				field.Name = rawFiled.Name
				field.Type = typeCast(rawFiled.Type)
				fields = append(fields, field)
			}
			ife.Fields = fields
			interfaces = append(interfaces, ife)
		}
		data.Interfaces = interfaces

		if err := outputfile(outputPath, servicetpl, data); err != nil {
			panic(fmt.Errorf("failed to write file %s. error: %v\n", outputPath, err))
		}

		log.Printf("%s", outputPath)
	}

}

func init() {
	langs.Langs[langName] = &TypeScriptGen{}
}
