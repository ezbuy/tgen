package grpc

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezbuy/tgen/langs"
	"github.com/ezbuy/tgen/tmpl"
	"github.com/samuel/go-thrift/parser"
)

const (
	langName = "grpc"
)

const TPL_SERVICE = "tgen/grpc/grpc"

type GrpcGen struct {
	langs.BaseGen
	thrift *parser.Thrift
}

func (g *GrpcGen) ServiceName() string {
	for key, _ := range g.thrift.Services {
		return key
	}
	return ""
}

func (g *GrpcGen) Includes() (includes []string) {
	for _, inc := range g.thrift.Includes {
		i := strings.LastIndex(inc, "/")
		if i > 0 {
			inc = inc[i+1:]
		}

		inc = strings.Replace(inc, ".thrift", ".proto", 1)
		includes = append(includes, inc)
	}

	return
}

type Struct struct {
	*parser.Struct
}

func (s *Struct) GetFields() (fields []*Field) {
	for _, inc := range s.Fields {
		fields = append(fields, &Field{inc})
	}
	return
}

type Field struct {
	*parser.Field
}

func getType(t *parser.Type) string {
	name := t.Name
	if name == "i32" {
		return "int32"
	}

	if name == "i64" {
		return "int64"
	}

	if name == "list" {
		return "repeated " + getType(t.ValueType)
	}

	return name
}

func (s *Field) GetType() string {
	return getType(s.Type)
}

func (g *GrpcGen) GetStructs() (structs []*Struct) {
	for _, inc := range g.thrift.Structs {
		structs = append(structs, &Struct{inc})
	}

	return
}

func initemplate(n string, path string) *template.Template {
	data, err := tmpl.Asset(path)
	if err != nil {
		panic(err)
	}

	tpl, err := template.New(n).Parse(string(data))
	if err != nil {
		panic(err)
	}

	return tpl
}

func genOutputPath(base string, fileName string) string {
	start := strings.LastIndex(fileName, "/")
	end := strings.LastIndex(fileName, ".")
	name := fileName[start+1 : end]
	return filepath.Join(base, name+".proto")
}

func outputfile(fp string, t *template.Template, tplname string, data interface{}) error {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return t.ExecuteTemplate(file, tplname, data)
}

func (this *GrpcGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	var servicetpl *template.Template
	servicetpl = initemplate(TPL_SERVICE, "tmpl/grpc/grpc.goproto")
	this.BaseGen.Init("grpc", parsedThrift)

	for fileName, t := range parsedThrift {
		outputPath := genOutputPath(output, fileName)
		this.thrift = t

		if err := outputfile(outputPath, servicetpl, TPL_SERVICE, this); err != nil {
			panic(fmt.Errorf("failed to write file %s. error: %v\n", outputPath, err))
		}

		log.Printf("%s", outputPath)
	}
}

func init() {
	langs.Langs[langName] = &GrpcGen{}
}
