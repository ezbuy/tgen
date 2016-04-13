package swift

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/ezbuy/tgen/global"
	"github.com/ezbuy/tgen/langs"
	"github.com/ezbuy/tgen/tmpl"
	"github.com/samuel/go-thrift/parser"
)

const (
	SwiftTypeInt    = "Int"
	SwiftTypeInt64  = "Int64"
	SwiftTypeString = "String"
	SwiftTypeBool   = "Bool"
	SwiftTypeByte   = "Byte"
	SwiftTypeDouble = "Double"

	// other types (such as array, map, etc.) are implemented in the method 'PlainType'
)

const (
	TPL_STRUCT  = "tgen/swift/struct"
	TPL_SERVICE = "tgen/swift/serivce"
)

var typemapping = map[string]string{
	langs.ThriftTypeI16:    SwiftTypeInt,
	langs.ThriftTypeI32:    SwiftTypeInt,
	langs.ThriftTypeI64:    SwiftTypeInt64,
	langs.ThriftTypeString: SwiftTypeString,
	langs.ThriftTypeByte:   SwiftTypeByte,
	langs.ThriftTypeBool:   SwiftTypeBool,
	langs.ThriftTypeDouble: SwiftTypeDouble,
}

type SwiftGen struct {
	langs.BaseGen
}

type BaseSwift struct {
	Filepath string
	Thrift   *parser.Thrift
	Thrifts  *map[string]*parser.Thrift // use pointer to pass variable, so it doesn't need to copy
}

func (this *BaseSwift) PlainType(t *parser.Type) string {
	n := t.Name

	if t, ok := typemapping[n]; ok {
		return t
	}

	switch n {
	case langs.ThriftTypeList, langs.ThriftTypeSet:
		return fmt.Sprintf("[%s]", this.PlainType(t.ValueType))
	case langs.ThriftTypeMap:
		return fmt.Sprintf("[%s: %s]", this.PlainType(t.KeyType), this.PlainType(t.ValueType))
	default:
		return this.AssembleCustomizedTypeName(t)
	}
}

func (this *BaseSwift) GetInnerType(t *parser.Type) string {
	if t.Name == langs.ThriftTypeList || t.Name == langs.ThriftTypeSet {
		return this.GetInnerType(t.ValueType)
	}

	return this.PlainType(t)
}

func (this *BaseSwift) IsBasicType(t string) bool {
	switch t {
	case SwiftTypeBool, SwiftTypeByte, SwiftTypeDouble, SwiftTypeInt, SwiftTypeInt64, SwiftTypeString:
		return true
	default:
		return false
	}
}

func (this *BaseSwift) AssembleCustomizedTypeName(t *parser.Type) string {
	if t == nil {
		return "Void"
	}

	names := strings.Split(t.Name, ".")

	// if the type is in current thrift file
	// get namespace
	// else, iterator the included thrift files
	// found the very first of thrift file
	// get its namespace
	// strip the first letter, insert the namespace at the head of the left

	if len(names) == 1 {
		for n, _ := range this.Thrift.Structs {
			if n != t.Name {
				continue
			}

			// we have checked namespace earlier, so we assume it must have corresponding namespace
			ns, _ := this.Thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, t.Name[1:])
		}
	}

	for path, thrift := range *this.Thrifts {
		if thrift == this.Thrift {
			continue
		}

		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) // or use slice
		if filename != names[0] {
			continue
		}

		for n, _ := range thrift.Structs {
			if n != names[1] {
				continue
			}

			ns, _ := thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, names[1][1:])
		}
	}

	panic(fmt.Sprintf("thrift file '%s': namespace of customized type '%s' if not found\n", this.Filepath, t.Name))
}

func (this *BaseSwift) AssembleStructName(n string) string {
	ns, _ := this.Thrift.Namespaces["swift"]
	return fmt.Sprintf("%s%s", ns, n[1:])
}

// if the property name is the keyword of swift, rename it
// but encode/decode with its origin name
func (this *BaseSwift) FilterPropertory(n string) string {
	switch n {
	case "description":
		return "desc"
	default:
		return n
	}
}

func (this *BaseSwift) ParamsJoinedByComma(args []*parser.Field) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer

	for i, arg := range args {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(arg.Name + ": " + this.Typecast(arg.Type, false))
	}

	return buf.String()
}

func (this *BaseSwift) AssignToDict(f *parser.Field) string {
	if f.Type.Name == "list" {
		innertype := this.GetInnerType(f.Type)
		if this.IsBasicType(innertype) {
			if innertype == SwiftTypeInt64 {
				return fmt.Sprintf("%s?.map { value in NSNumber(longLong: value) }", this.FilterPropertory(f.Name))
			}
			return fmt.Sprintf("%s", this.FilterPropertory(f.Name))
		} else {
			return fmt.Sprintf("%s?.toJSON()", this.FilterPropertory(f.Name))
		}
	}

	switch f.Type.Name {
	case langs.ThriftTypeI16, langs.ThriftTypeI32, langs.ThriftTypeByte, langs.ThriftTypeString,
		langs.ThriftTypeBool, langs.ThriftTypeDouble,
		langs.ThriftTypeMap:
		return this.FilterPropertory(f.Name)
	case langs.ThriftTypeI64:
		return fmt.Sprintf("NSNumber(longLong: %s)", this.FilterPropertory(f.Name))
	default:
		return fmt.Sprintf("%s?.toJSON()", this.FilterPropertory(f.Name))
	}
}

func (this *BaseSwift) TypecastWithDefaultValue(t *parser.Type) string {
	return this.Typecast(t, true)
}

func (this *BaseSwift) TypecastWithoutDefaultValue(t *parser.Type) string {
	return this.Typecast(t, false)
}

func (this *BaseSwift) Typecast(t *parser.Type, flag bool) string {
	pt := this.PlainType(t)

	switch pt {
	case SwiftTypeInt, SwiftTypeInt64:
		if flag {
			return fmt.Sprintf("%s = 0", pt)
		}
		return pt
	case SwiftTypeByte:
		return pt
	case SwiftTypeBool:
		if flag {
			return fmt.Sprintf("%s = false", pt)
		}
		return pt
	case SwiftTypeDouble:
		if flag {
			return fmt.Sprintf("%s = 0.0", pt)
		}
		return pt
	default:
		return fmt.Sprintf("%s?", pt)
	}
}

type swiftStruct struct {
	*BaseSwift
	*parser.Struct
}

type swiftService struct {
	*BaseSwift
	*parser.Service
}

func generateWithModel(gen *SwiftGen, m string, output string, parsedThrift map[string]*parser.Thrift) {
	if m != global.MODE_REST && m != global.MODE_JSONRPC {
		log.Fatalf("mode '%s' is invalid", m)
	}

	gen.BaseGen.Init("swift", parsedThrift)

	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	// init templates
	structpl := initemplate(TPL_STRUCT, "tmpl/swift/struct.goswift")
	var servicetpl *template.Template
	if m == global.MODE_REST {
		servicetpl = initemplate(TPL_SERVICE, "tmpl/swift/rest_service.goswift")
	} else if m == global.MODE_JSONRPC {
		servicetpl = initemplate(TPL_SERVICE, "tmpl/swift/jsonrpc_service.goswift")
	}

	wg := sync.WaitGroup{}

	// key is the absoule path of thrift file
	for f, t := range parsedThrift {
		// check namespace
		if _, ok := t.Namespaces["swift"]; !ok {
			fmt.Printf("namespace of swift in file '%s' is not found\n", f)
			continue
		}

		wg.Add(1)

		go func(t *parser.Thrift) {
			defer wg.Done()

			for _, s := range t.Structs {
				baseSwift := &BaseSwift{Filepath: f, Thrift: t, Thrifts: &parsedThrift}

				name := fmt.Sprintf("%s.swift", baseSwift.AssembleStructName(s.Name))

				path := filepath.Join(output, name)

				data := &swiftStruct{BaseSwift: baseSwift, Struct: s}

				if err := outputfile(path, structpl, TPL_STRUCT, data); err != nil {
					panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
				}
			}
		}(t)

		wg.Add(1)

		go func(t *parser.Thrift) {
			defer wg.Done()

			for _, s := range t.Services {
				// filename is the service name plus 'Service'
				name := s.Name + "Service.swift"

				path := filepath.Join(output, name)

				data := &swiftService{BaseSwift: &BaseSwift{Filepath: f, Thrift: t, Thrifts: &parsedThrift}, Service: s}

				if err := outputfile(path, servicetpl, TPL_SERVICE, data); err != nil {
					panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
				}
			}
		}(t)
	}

	wg.Wait()
}

func generateAll(gen *SwiftGen, output string, parsedThrift map[string]*parser.Thrift) {
	generateWithModel(gen, global.MODE_REST, filepath.Join(output, global.MODE_REST), parsedThrift)
	generateWithModel(gen, global.MODE_JSONRPC, filepath.Join(output, global.MODE_JSONRPC), parsedThrift)
}

func (this *SwiftGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	if global.Mode != "" {
		generateWithModel(this, global.Mode, output, parsedThrift)
	} else {
		generateAll(this, output, parsedThrift)
	}
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

func outputfile(fp string, t *template.Template, tplname string, data interface{}) error {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return t.ExecuteTemplate(file, tplname, data)
}

func init() {
	langs.Langs["swift"] = &SwiftGen{}
}
