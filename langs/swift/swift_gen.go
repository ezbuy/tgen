package swift

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ezbuy/tgen/langs"
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

	tpl *template.Template
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
	o.tpl = langs.InitTemplate("tmpl/swift/swift.goswift")

	// tp is the absoule path of thrift file
	for tp, thrift := range parsedThrift {
		data := o.gen(thrift)

		name := o.filename(tp, thrift.Namespaces)

		path := filepath.Join(output, name)

		// save to disk
		langs.Write(path, data)

		fmt.Printf("[%s] generated\n", path)
	}
}

func (o *SwiftGen) gen(thrift *parser.Thrift) []byte {
	st := &swiftThrift{Thrift: thrift}

	data := langs.RenderTemplate(o.tpl, st)

	return data
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
