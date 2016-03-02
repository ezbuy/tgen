package swift

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

type BaseSwift struct{}

func (this *BaseSwift) PlainType(t *parser.Type) string {
	n := this.LastComponentOfDotStr(t.Name)

	if t, ok := typemapping[n]; ok {
		return t
	}

	switch n {
	case langs.ThriftTypeList, langs.ThriftTypeSet:
		return fmt.Sprintf("[%s]", this.PlainType(t.ValueType))
	case langs.ThriftTypeMap:
		return fmt.Sprintf("[%s: %s]", this.PlainType(t.KeyType), this.PlainType(t.ValueType))
	default:
		return n
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

func (this *BaseSwift) LastComponentOfDotStr(str string) string {
	if strings.Contains(str, ".") == false {
		return str
	}

	strs := strings.Split(str, ".")
	return strs[len(strs)-1]
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
		if this.IsBasicType(this.GetInnerType(f.Type)) {
			return fmt.Sprintf("%s", f.Name)
		} else {
			return fmt.Sprintf("%s?.toJSON()", f.Name)
		}
	}

	switch f.Type.Name {
	case langs.ThriftTypeI16, langs.ThriftTypeI32, langs.ThriftTypeByte, langs.ThriftTypeString,
		langs.ThriftTypeBool, langs.ThriftTypeDouble,
		langs.ThriftTypeMap:
		return f.Name
	case langs.ThriftTypeI64:
		return fmt.Sprintf("NSNumber(longLong: %s)", f.Name)
	default:
		return fmt.Sprintf("%s?.toJSON()", f.Name)
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

func (o *SwiftGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	o.BaseGen.Init("swift", parsedThrift)

	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	// templates
	var structpl *template.Template
	var servicetpl *template.Template

	// tp is the absoule path of thrift file
	for _, t := range parsedThrift {
		for _, s := range t.Structs {
			if structpl == nil {
				structpl = initemplate(TPL_STRUCT, "tmpl/swift/struct.goswift")
			}

			// filename is the struct name
			name := s.Name + ".swift"

			path := filepath.Join(output, name)

			data := &swiftStruct{BaseSwift: &BaseSwift{}, Struct: s}

			if err := outputfile(path, structpl, TPL_STRUCT, data); err != nil {
				panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
			}
		}

		for _, s := range t.Services {
			if servicetpl == nil {
				servicetpl = initemplate(TPL_SERVICE, "tmpl/swift/service.goswift")
			}

			// filename is the service name plus 'Service'
			name := s.Name + "Service.swift"

			path := filepath.Join(output, name)

			data := &swiftService{BaseSwift: &BaseSwift{}, Service: s}

			if err := outputfile(path, servicetpl, TPL_SERVICE, data); err != nil {
				panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
			}
		}
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
