package grpc

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bradfitz/slice"
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
	Thrift      *parser.Thrift
	reqStructs  []*Struct
	respStructs []*Struct
}

func (g *GrpcGen) ServiceName() string {
	for key, _ := range g.Thrift.Services {
		return key
	}

	for k, v := range g.Thrift.Namespaces {
		if k == "*" {
			return v
		}
	}

	return ""
}

func (g *GrpcGen) SetThrift(t *parser.Thrift) {
	g.Thrift = t
	g.reqStructs = nil
	g.respStructs = nil

	for _, svr := range g.Thrift.Services {
		for _, method := range svr.Methods {
			args := &Struct{&parser.Struct{}}
			args.Name = method.Name + "Request"
			args.Fields = method.Arguments
			g.reqStructs = append(g.reqStructs, args)

			args = &Struct{&parser.Struct{}}
			args.Name = method.Name + "Response"
			if method.ReturnType != nil {
				f := &parser.Field{}
				f.ID = 1
				f.Name = "Result"
				f.Type = method.ReturnType
				args.Fields = append(args.Fields, f)
			}
			g.respStructs = append(g.respStructs, args)
		}
	}

	slice.Sort(g.respStructs, func(i, j int) bool {
		return g.respStructs[i].Name < g.respStructs[j].Name
	})

	slice.Sort(g.reqStructs, func(i, j int) bool {
		return g.reqStructs[i].Name < g.reqStructs[j].Name
	})

	return
}

func (g *GrpcGen) Includes() (includes []string) {
	for _, inc := range g.Thrift.Includes {
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
	if t == nil {
		return ""
	}

	name := t.Name
	if name == "i32" || name == "i16" {
		return "int32"
	}

	if name == "i64" {
		return "int64"
	}

	if name == "list" {
		return "repeated " + getType(t.ValueType)
	}

	if name == "map" {
		return getMapType(t)
	}

	return name
}

func genListMessageName(typeName string) string {
	return "ListOf" + strings.Title(typeName)
}

func genMapDefine(keyType, valueType string) string {
	return "map<" + keyType + "," + valueType + ">"
}

func getMapType(t *parser.Type) string {

	if t.ValueType.Name == "list" {
		return genMapDefine(getType(t.KeyType), genListMessageName(getType(t.ValueType.ValueType)))
	}

	return genMapDefine(getType(t.KeyType), getType(t.ValueType))

}

func (s *Field) GetType() template.HTML {
	return template.HTML(getType(s.Type))
}

func (g *GrpcGen) GetPackages() (result map[string]string) {
	result = make(map[string]string)
	for k, v := range g.Thrift.Namespaces {
		if k != "webapi" && k != "objc" && k != "javascript" &&
			k != "csharp" && k != "swift" && k != "*" {
			result[k+"_package"] = v
		} else if k == "csharp" {
			result["csharp_namespace"] = v
		}
	}
	return
}

func (g *GrpcGen) GetStructs() (structs []*Struct) {
	for _, inc := range g.Thrift.Structs {
		structs = append(structs, &Struct{inc})
	}

	return
}

func (g *GrpcGen) GetReqStructs() (structs []*Struct) {
	return g.reqStructs
}

func (g *GrpcGen) GetRespStructs() (structs []*Struct) {
	return g.respStructs
}

func initemplate(n string, path string) *template.Template {
	data, err := tmpl.Asset(path)
	if err != nil {
		panic(err)
	}

	tpl, err := template.New(n).Funcs(template.FuncMap{
		"listEnumValue": listEnumValue,
		"getType":       getType,
	}).Parse(string(data))
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

func listEnumValue(enums map[string]*parser.EnumValue) (result []*parser.EnumValue) {
	zeroKey := "Unknown"
	for _, v := range enums {
		if v.Value == 0 {
			zeroKey = v.Name
		}
	}

	result = append(result, &parser.EnumValue{
		Name:  zeroKey,
		Value: 0,
	})

	for _, v := range enums {
		if v.Name != zeroKey {
			result = append(result, v)
		}
	}

	slice.Sort(result, func(i, j int) bool {
		return result[i].Value < result[j].Value
	})
	return
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
		this.preProcessMapType(t)
		outputPath := genOutputPath(output, fileName)
		this.SetThrift(t)

		if err := outputfile(outputPath, servicetpl, TPL_SERVICE, this); err != nil {
			panic(fmt.Errorf("failed to write file %s. error: %v\n", outputPath, err))
		}

		log.Printf("%s", outputPath)
	}
}

func (this *GrpcGen) preProcessMapType(thrift *parser.Thrift) {

	expandMap := make(map[string]*parser.Struct)

	appendExpandMap := func(tType *parser.Type) {
		if tType.Name == "map" && tType.ValueType.Name == "list" {
			structName := genListMessageName(getType(tType.ValueType.ValueType))
			fmt.Printf("warning: convert thrift map list value [map<_,list<%s>] to [map<_, %s] \n", getType(tType.ValueType.ValueType), genListMessageName(getType(tType.ValueType.ValueType)))
			expandMap[structName] = &parser.Struct{
				Name: structName,
				Fields: []*parser.Field{
					{
						ID:   1,
						Name: "Data",
						Type: &parser.Type{
							Name:      "list",
							ValueType: tType.ValueType.ValueType,
						},
					},
				},
			}
		}

	}

	eraseEnumMapkey := func(tType *parser.Type) {
		if tType.Name == "map" {
			if _, ok := thrift.Enums[tType.KeyType.Name]; ok {
				fmt.Printf("warning: convert thrift map enum key [%s] to map<int32,_>\n", getMapType(tType))
				tType.KeyType = &parser.Type{
					Name: "i32",
				}
			}
		}
	}

	for _, tstruct := range thrift.Structs {
		for i, field := range tstruct.Fields {
			appendExpandMap(field.Type)
			eraseEnumMapkey(field.Type)
			field.ID = i + 1 // reorder
		}
	}

	for _, tservice := range thrift.Services {
		for _, method := range tservice.Methods {
			for _, arg := range method.Arguments {
				appendExpandMap(arg.Type)
				eraseEnumMapkey(arg.Type)
			}
			appendExpandMap(method.ReturnType)
			eraseEnumMapkey(method.ReturnType)
		}
	}
	for name, expand := range expandMap {
		thrift.Structs[name] = expand
	}

}
func init() {
	langs.Langs[langName] = &GrpcGen{}
}
