package gogen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const (
	doubleDot = ".."
	dot       = "."
	slash     = "/"
)

const (
	TypeBool   = langs.ThriftTypeBool
	TypeByte   = langs.ThriftTypeByte
	TypeI16    = langs.ThriftTypeI16
	TypeI32    = langs.ThriftTypeI32
	TypeI64    = langs.ThriftTypeI64
	TypeDouble = langs.ThriftTypeDouble
	TypeBinary = langs.ThriftTypeBinary
	TypeString = langs.ThriftTypeString
	TypeList   = langs.ThriftTypeList
	TypeMap    = langs.ThriftTypeMap
	TypeSet    = langs.ThriftTypeSet
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

var constantValueFormat = map[string]string{
	TypeBool:   "%s", // parser.Identifier("true") and parser.Identifier("false")
	TypeByte:   "%d",
	TypeI16:    "%d",
	TypeI32:    "%d",
	TypeI64:    "%d",
	TypeDouble: "%f",
	TypeString: "%q",
}

var invalidWords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,

	"rune":    true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"float32": true,
	"float64": true,

	"close":   true,
	"len":     true,
	"cap":     true,
	"new":     true,
	"make":    true,
	"append":  true,
	"copy":    true,
	"delete":  true,
	"complex": true,
	"real":    true,
	"imag":    true,
	"panic":   true,
	"recover": true,
	"print":   true,
	"println": true,

	"init": true,
	"main": true,
}

type TplUtils struct {
}

func (this *TplUtils) IsInvalidTypeName(str string) bool {
	invalid, _ := invalidWords[str]
	return invalid
}

func (this *TplUtils) GenNamespace(namespace string) (pkgName string, importPath string) {
	if strings.Contains(namespace, doubleDot) {
		importPath = strings.Replace(namespace, doubleDot, slash, -1)
	} else {
		importPath = strings.Replace(namespace, dot, slash, -1)
	}

	pkgName = filepath.Base(importPath)

	return pkgName, importPath
}

func (this *TplUtils) UpperHead(name string) string {
	if name == "" {
		return name
	}

	head := name[0:1]
	return strings.ToUpper(head) + name[1:]
}

func (this *TplUtils) IsSimpleArguments(args []*parser.Field) bool {
	if len(args) != 1 {
		return false
	}

	arg := args[0]

	if arg == nil || arg.Type == nil {
		return false
	}

	switch arg.Type.Name {
	case TypeBool, TypeByte, TypeI16, TypeI32, TypeI64, TypeDouble, TypeBinary, TypeString, TypeList:
		return true

	default:
		return false
	}
}

func (this *TplUtils) IsNilType(typ *parser.Type) bool {
	return typ == nil
}

func (this *TplUtils) FieldTagThrift(field *parser.Field) string {
	str := fmt.Sprintf("%d", field.ID)

	if field.Optional {
		return str
	}

	return str + ",required"
}

func (this *TplUtils) FieldTagJson(field *parser.Field) string {
	str := fmt.Sprintf("%s", field.Name)

	if !field.Optional {
		return str
	}

	return str + ",omitempty"
}

func gofmt(paths ...string) {
	args := []string{"-l", "-w"}
	args = append(args, paths...)

	buf := new(bytes.Buffer)

	cmd := exec.Command("gofmt", args...)
	cmd.Stdout = buf
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fail to gofmt %s", err)
		fmt.Fprintln(os.Stderr, "##### gofmt trace info #####")

		if _, err := io.Copy(os.Stderr, buf); err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stderr, "##### gofmt trace info end #####")
		os.Exit(1)
	}
}

func panicWithErr(format string, msg ...interface{}) {
	panic(fmt.Errorf(format, msg...))
}

func exitWithError(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, msg...)
	os.Exit(1)
}
