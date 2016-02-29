package gogen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/samuel/go-thrift/parser"
)

const (
	TypeBool   = "bool"
	TypeByte   = "byte"
	TypeI16    = "i16"
	TypeI32    = "i32"
	TypeI64    = "i64"
	TypeDouble = "double"
	TypeBinary = "binary"
	TypeString = "string"
	TypeList   = "list"
	TypeMap    = "map"
	TypeSet    = "set"
)

var typeStrs = map[string]string{
	TypeBool:   "bool",
	TypeByte:   "byte",
	TypeI16:    "int16",
	TypeI32:    "int32",
	TypeI64:    "int64",
	TypeDouble: "float64",
	TypeBinary: "[]byte",
	TypeString: "string",
}

func getNamespace(namespaces map[string]string) string {
	if namespace, ok := namespaces[langName]; ok {
		return namespace
	}

	return ""
}

func getIncludes(parsedThrift map[string]*parser.Thrift, includes map[string]string) [][2]string {
	results := make([][2]string, 0, len(includes))

	for includeName, filename := range includes {
		parsed, ok := parsedThrift[filename]
		if !ok {
			panicWithErr("include thrift %q not found %s", includeName, parsedThrift)
		}

		importPath, _ := genNamespace(getNamespace(parsed.Namespaces))

		results = append(results, [2]string{includeName, importPath})
	}

	return results
}

func genNamespace(namespace string) (string, string) {
	path := strings.Replace(namespace, ".", "/", -1)
	pkgName := filepath.Base(path)
	return path, pkgName
}

func panicWithErr(format string, msg ...interface{}) {
	panic(fmt.Errorf(format, msg...))
}

func gofmt(paths ...string) {
	args := []string{"-l", "-w"}
	args = append(args, paths...)

	cmd := exec.Command("gofmt", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fail to gofmt %s", err)
	}
}

type TplUtils struct {
}

func (this *TplUtils) UpperHead(name string) string {
	if name == "" {
		return name
	}

	head := name[0:1]
	return strings.ToUpper(head) + name[1:]
}

func (this *TplUtils) GenTypeString(fieldName string, typ *parser.Type, optional bool, isMapKey bool) string {
	if typ == nil {
		panicWithErr("field %s with nil type", fieldName)
	}

	var str string

	switch typ.Name {
	case TypeBool, TypeByte, TypeI16, TypeI32, TypeI64, TypeDouble, TypeString:
		if optional {
			str = "*"
		}
		str += typeStrs[typ.Name]

	case TypeBinary:
		if isMapKey {
			panicWithErr("map field %s with binary key", fieldName)
		}
		str = typeStrs[TypeBinary]

	case TypeList:
		if isMapKey {
			panicWithErr("map field %s with list key", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("list field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("[]%s", this.GenTypeString(fieldName, typ.ValueType, false, false))

	case TypeMap:
		if isMapKey {
			panicWithErr("map field %s with map key", fieldName)
		}

		if typ.KeyType == nil {
			panicWithErr("map field %s with nil key type", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("map field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("map[%s]%s",
			this.GenTypeString(fieldName, typ.KeyType, false, true),
			this.GenTypeString(fieldName, typ.ValueType, false, false),
		)

	case TypeSet:
		// TODO: support set

	default:
		if typ.Name == "" {
			panicWithErr("field %s without type name", fieldName)
		}

		// TODO check if is Enum, Const, TypeDef etc.
		str = "*" + typ.Name
	}

	return str
}

func (this *TplUtils) IsNilType(typ *parser.Type) bool {
	return typ == nil
}

func (this *TplUtils) GenServiceMethodArguments(fields []*parser.Field) string {
	var str string

	maxIdx := len(fields) - 1
	for idx, field := range fields {
		str += fmt.Sprintf("%s %s", field.Name, this.GenTypeString(field.Name, field.Type, field.Optional, false))
		if idx != maxIdx {
			str += ", "
		}
	}

	return str
}

func (this *TplUtils) GenWebApiServiceParams(fields []*parser.Field) string {
	var str string

	maxIdx := len(fields) - 1
	for idx, field := range fields {
		str += fmt.Sprintf("params.%s", this.UpperHead(field.Name))
		if idx != maxIdx {
			str += ", "
		}
	}

	return str
}
