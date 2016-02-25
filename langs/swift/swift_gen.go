package swift

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/ezbuy/tgen/langs"
	"github.com/ezbuy/tgen/tmpl"
	"github.com/samuel/go-thrift/parser"
)

const (
	EXT = ".swift"
)

type SwiftGen struct {
	langs.BaseGen
}

func (o *SwiftGen) Generate(parsedThrift map[string]*parser.Thrift) ([]langs.GenResult, error) {
	o.BaseGen.Init("swift", "tmpl/swift.goswift", parsedThrift)

	funcMap := template.FuncMap{
		"plainType":                   plainType,
		"typecastWithDefaultValue":    typecastWithDefaultValue,
		"typecastWithoutDefaultValue": typecastWithoutDefaultValue,
		"paramsJoinedByComma":         paramsJoinedByComma,
		"lastComponentOfDotStr":       lastComponentOfDotStr,
	}

	// read template
	tpldata, err := tmpl.Asset(o.BaseGen.Tplpath)
	if err != nil {
		return nil, err
	}

	t, err := template.New("swift").Funcs(funcMap).Parse(string(tpldata))
	if err != nil {
		return nil, err
	}

	results := []langs.GenResult{}

	// key is the absoule path of thrift file
	// we may not need it
	for k, v := range parsedThrift {
		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return nil, err
		}

		name := k

		if val, ok := v.Namespaces["swift"]; ok {
			name = val
		}

		results = append(results, langs.GenResult{Infile: k, Filename: name + EXT, Data: buf.Bytes()})
	}

	return results, nil
}

func init() {
	langs.Langs["swift"] = &SwiftGen{}
}

func paramsJoinedByComma(args []*parser.Field) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer

	for i, arg := range args {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(arg.Name + ": " + typecast(arg.Type, false))
	}

	return buf.String()
}

func typecastWithDefaultValue(t *parser.Type) string {
	return typecast(t, true)
}

func typecastWithoutDefaultValue(t *parser.Type) string {
	return typecast(t, false)
}

func typecast(t *parser.Type, flag bool) string {
	pt := plainType(t)

	switch pt {
	case "Int", "Int64":
		if flag {
			return fmt.Sprintf("%s = 0", pt)
		}
		return pt
	case "Byte":
		return pt
	case "Bool":
		if flag {
			return fmt.Sprintf("%s = false", pt)
		}
		return pt
	case "Double":
		if flag {
			return fmt.Sprintf("%s = 0.0", pt)
		}
		return pt
	default:
		return fmt.Sprintf("%s?", pt)
	}
}

func plainType(t *parser.Type) string {
	tn := lastComponentOfDotStr(t.Name)

	switch tn {
	case "i16", "i32":
		return "Int"
	case "i64":
		return "Int64"
	case "string", "byte", "bool", "double":
		return fmt.Sprintf("%s%s", strings.ToUpper(tn[:1]), strings.ToLower(tn[1:]))
	case "list", "set":
		return fmt.Sprintf("[%s]", plainType(t.ValueType))
	case "map":
		return fmt.Sprintf("[%s: %s]", plainType(t.KeyType), plainType(t.ValueType))
	default:
		return tn
	}
}

func lastComponentOfDotStr(str string) string {
	if strings.Contains(str, ".") == false {
		return str
	}

	strs := strings.Split(str, ".")
	return strs[len(strs)-1]
}
