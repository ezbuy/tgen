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
	TPL_NAME = "tgen/swift"
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

var tpl *template.Template

type SwiftGen struct {
	langs.BaseGen
}

type swiftThrift struct {
	*parser.Thrift
}

func (st *swiftThrift) PlainType(t *parser.Type) string {
	n := st.LastComponentOfDotStr(t.Name)

	if t, ok := typemapping[n]; ok {
		return t
	}

	switch n {
	case langs.ThriftTypeList, langs.ThriftTypeSet:
		return fmt.Sprintf("[%s]", st.PlainType(t.ValueType))
	case langs.ThriftTypeMap:
		return fmt.Sprintf("[%s: %s]", st.PlainType(t.KeyType), st.PlainType(t.ValueType))
	default:
		return n
	}
}

func (st *swiftThrift) LastComponentOfDotStr(str string) string {
	if strings.Contains(str, ".") == false {
		return str
	}

	strs := strings.Split(str, ".")
	return strs[len(strs)-1]
}

func (st *swiftThrift) ParamsJoinedByComma(args []*parser.Field) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer

	for i, arg := range args {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(arg.Name + ": " + st.Typecast(arg.Type, false))
	}

	return buf.String()
}

func (st *swiftThrift) AssignToDict(f *parser.Field) string {
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

func (st *swiftThrift) TypecastWithDefaultValue(t *parser.Type) string {
	return st.Typecast(t, true)
}

func (st *swiftThrift) TypecastWithoutDefaultValue(t *parser.Type) string {
	return st.Typecast(t, false)
}

func (st *swiftThrift) Typecast(t *parser.Type, flag bool) string {
	pt := st.PlainType(t)

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

func (o *SwiftGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	o.BaseGen.Init("swift", parsedThrift)

	// init template
	tpl = initTemplate("tmpl/swift/swift.goswift")

	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	// tp is the absoule path of thrift file
	for tp, t := range parsedThrift {
		// get file name
		name := o.filename(tp, t.Namespaces)

		// get output file path
		path := filepath.Join(output, name)

		data := &swiftThrift{Thrift: t}

		// sort structs & services

		if err := outputfile(path, TPL_NAME, data); err != nil {
			panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
		}
	}
}

// Filename returns the final name of the generated file
// generally, the name is parsed from the template's namespace of the specified language,
// if the name can't be parsed, it will be set as the template's name
func (o *SwiftGen) filename(tplfile string, ns map[string]string) string {
	name := strings.TrimRight(filepath.Base(tplfile), filepath.Ext(tplfile))

	if val, ok := ns["swift"]; ok {
		name = val
	}

	return name + ".swift"
}

func init() {
	langs.Langs["swift"] = &SwiftGen{}
}

func initTemplate(path string) *template.Template {
	data, err := tmpl.Asset(path)
	if err != nil {
		panic(err)
	}

	tpl, err := template.New(TPL_NAME).Parse(string(data))
	if err != nil {
		panic(err)
	}

	return tpl
}

func outputfile(fp string, tplname string, data interface{}) error {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return tpl.ExecuteTemplate(file, tplname, data)
}
