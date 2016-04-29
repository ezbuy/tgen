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
	tplStruct  = "tgen/swift/struct"
	tplEnum    = "tgen/swift/enum"
	tplService = "tgen/swift/serivce"
)

var typeMapping = map[string]string{
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
	Thrifts  map[string]*parser.Thrift
}

func (this *BaseSwift) PlainType(t *parser.Type) string {
	n := t.Name

	if t, ok := typeMapping[n]; ok {
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

// IsEnum checks whether a type is enum.
// it first checks in its own definition, than check from included files
func (b *BaseSwift) IsEnum(t *parser.Type) bool {
	if t == nil {
		return false
	}

	names := strings.Split(t.Name, ".")

	if len(names) == 1 {
		for n := range b.Thrift.Enums {
			if n == t.Name {
				return true
			}
		}

		return false
	}

	for path, thrift := range b.Thrifts {
		if thrift == b.Thrift {
			continue
		}

		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if filename != names[0] {
			continue
		}

		for n := range thrift.Enums {
			if n == names[1] {
				return true
			}
		}
	}

	return false
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
		for n := range this.Thrift.Structs {
			if n != t.Name {
				continue
			}

			// we have checked namespace earlier, so we assume it must have corresponding namespace
			ns, _ := this.Thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, t.Name[1:])
		}

		for n := range this.Thrift.Enums {
			if n != t.Name {
				continue
			}

			ns, _ := this.Thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, t.Name[1:])
		}
	}

	for path, thrift := range this.Thrifts {
		if thrift == this.Thrift {
			continue
		}

		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) // or use slice
		if filename != names[0] {
			continue
		}

		for n := range thrift.Structs {
			if n != names[1] {
				continue
			}

			ns, _ := thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, names[1][1:])
		}

		for n := range thrift.Enums {
			if n != names[1] {
				continue
			}

			ns, _ := thrift.Namespaces["swift"]

			return fmt.Sprintf("%s%s", ns, names[1][1:])
		}
	}

	log.Fatalf("thrift file '%s': namespace of customized type '%s' if not found\n", this.Filepath, t.Name)

	return ""
}

// AssembleStructName returns object name of generated file name, assembled with namespace
func (b *BaseSwift) AssembleStructName(n string) string {
	ns, _ := b.Thrift.Namespaces["swift"]
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

func (b *BaseSwift) AssignToDict(f *parser.Field) string {
	if f.Type.Name == "list" {
		innertype := b.GetInnerType(f.Type)
		if b.IsBasicType(innertype) {
			if innertype == SwiftTypeInt64 {
				return fmt.Sprintf("%s?.map { value in NSNumber(longLong: value) }", b.FilterPropertory(f.Name))
			}
			return fmt.Sprintf("%s", b.FilterPropertory(f.Name))
		}

		return fmt.Sprintf("%s?.toJSON()", b.FilterPropertory(f.Name))
	}

	switch f.Type.Name {
	case langs.ThriftTypeI16, langs.ThriftTypeI32, langs.ThriftTypeByte, langs.ThriftTypeString,
		langs.ThriftTypeBool, langs.ThriftTypeDouble,
		langs.ThriftTypeMap:
		return b.FilterPropertory(f.Name)
	case langs.ThriftTypeI64:
		return fmt.Sprintf("NSNumber(longLong: %s)", b.FilterPropertory(f.Name))
	default:
		if b.IsEnum(f.Type) {
			return fmt.Sprintf("%s?.rawValue", b.FilterPropertory(f.Name))
		}
		return fmt.Sprintf("%s?.toJSON()", b.FilterPropertory(f.Name))
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

type swiftEnum struct {
	*BaseSwift
	*parser.Enum
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
	structpl := initemplate(tplStruct, "tmpl/swift/struct.goswift")
	enumtpl := initemplate(tplEnum, "tmpl/swift/enum.goswift")

	var servicetpl *template.Template
	if m == global.MODE_REST {
		servicetpl = initemplate(tplService, "tmpl/swift/rest_service.goswift")
	} else if m == global.MODE_JSONRPC {
		servicetpl = initemplate(tplService, "tmpl/swift/jsonrpc_service.goswift")
	}

	wg := sync.WaitGroup{}

	// key is the absoule path of thrift file
	for f, t := range parsedThrift {
		if f != global.InputFile {
			continue // ignore
		}

		// enum
		wg.Add(1)

		go func(t *parser.Thrift, f string) {
			defer wg.Done()

			for _, e := range t.Enums {
				baseSwift := &BaseSwift{Filepath: f, Thrift: t, Thrifts: parsedThrift}

				name := fmt.Sprintf("%s.swift", baseSwift.AssembleStructName(e.Name))
				path := filepath.Join(output, name)

				data := &swiftEnum{BaseSwift: baseSwift, Enum: e}

				if err := outputfile(path, enumtpl, tplEnum, data); err != nil {
					panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
				}
			}
		}(t, f)

		wg.Add(1)

		go func(t *parser.Thrift, f string) {
			defer wg.Done()

			for _, s := range t.Structs {
				b := &BaseSwift{Filepath: f, Thrift: t, Thrifts: parsedThrift}

				name := fmt.Sprintf("%s.swift", b.AssembleStructName(s.Name))
				path := filepath.Join(output, name)

				data := &swiftStruct{BaseSwift: b, Struct: s}

				if err := outputfile(path, structpl, tplStruct, data); err != nil {
					panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
				}
			}
		}(t, f)

		wg.Add(1)

		go func(t *parser.Thrift, f string) {
			defer wg.Done()

			for _, s := range t.Services {
				// filename is the service name plus 'Service'
				name := s.Name + "Service.swift"
				path := filepath.Join(output, name)

				data := &swiftService{BaseSwift: &BaseSwift{Filepath: f, Thrift: t, Thrifts: parsedThrift}, Service: s}

				if err := outputfile(path, servicetpl, tplService, data); err != nil {
					panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
				}
			}
		}(t, f)
	}

	wg.Wait()
}

func generateAll(gen *SwiftGen, output string, parsedThrift map[string]*parser.Thrift) {
	generateWithModel(gen, global.MODE_REST, filepath.Join(output, global.MODE_REST), parsedThrift)
	generateWithModel(gen, global.MODE_JSONRPC, filepath.Join(output, global.MODE_JSONRPC), parsedThrift)
}

// Generate is the entry of the logic
func (g *SwiftGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	if global.Mode != "" {
		generateWithModel(g, global.Mode, output, parsedThrift)
	} else {
		generateAll(g, output, parsedThrift)
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
