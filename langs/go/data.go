package gogen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/samuel/go-thrift/parser"
)

type Package struct {
	PkgName    string
	ImportPath string

	includes map[string]*Package
	thrift   *parser.Thrift

	TplUtils
}

func newPackage(thrift *parser.Thrift) *Package {
	pkg := &Package{}
	pkg.setup(thrift)
	return pkg
}

func (this *Package) setup(thrift *parser.Thrift) {
	namespace := thrift.Namespaces[langName]
	this.PkgName, this.ImportPath = this.GenNamespace(namespace)

	this.thrift = thrift
}

func (this *Package) setupIncludes(packages map[string]*Package) {
	pkgMap := map[string]*Package{}

	for includeName, filename := range this.thrift.Includes {
		pkg, ok := packages[filename]
		if !ok {
			exitWithError("include thrift %q ( %s ) not found", includeName, filename)
		}

		pkgMap[includeName] = pkg
	}

	this.includes = pkgMap
}

func (this *Package) Includes() map[string]*Package {
	return this.includes
}

func (this *Package) Enums() map[string]*parser.Enum {
	return this.thrift.Enums
}

func (this *Package) Structs() map[string]*parser.Struct {
	return this.thrift.Structs
}

func (this *Package) Services() map[string]*parser.Service {
	return this.thrift.Services
}

func (this *Package) Constants() map[string]*parser.Constant {
	return this.thrift.Constants
}

func (this *Package) Typedefs() map[string]*parser.Typedef {
	return this.thrift.Typedefs
}

func (this *Package) Exceptions() map[string]*parser.Struct {
	return this.thrift.Exceptions
}

func (this *Package) Unions() map[string]*parser.Struct {
	return this.thrift.Unions
}

func (this *Package) Namespaces() map[string]string {
	return this.thrift.Namespaces
}

func (this *Package) Namespace() string {
	return this.thrift.Namespaces[langName]
}

func (this *Package) WebApiPrefix() string {
	namespace := this.thrift.Namespaces["webapi"]
	if namespace != "" {
		namespace = slash + strings.Replace(namespace, dot, slash, -1)
	}

	return namespace
}

func (this *Package) GenTypeString(fieldName string, typ, parent *parser.Type, optional bool) string {
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
		if parent != nil && typ == parent.KeyType {
			panicWithErr("map field %s with binary key", fieldName)
		}
		str = typeStrs[TypeBinary]

	case TypeList:
		if parent != nil && typ == parent.KeyType {
			panicWithErr("map field %s with list key", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("list field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("[]%s", this.GenTypeString(fieldName, typ.ValueType, typ, false))

	case TypeMap:
		if parent != nil && typ == parent.KeyType {
			panicWithErr("map field %s with map key", fieldName)
		}

		if typ.KeyType == nil {
			panicWithErr("map field %s with nil key type", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("map field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("map[%s]%s",
			this.GenTypeString(fieldName, typ.KeyType, typ, false),
			this.GenTypeString(fieldName, typ.ValueType, typ, false),
		)

	case TypeSet:
		// TODO: support set

	default:
		if typ.Name == "" {
			panicWithErr("field %s without type name", fieldName)
		}

		// TODO check if is Enum, Const, TypeDef etc.
		name := typ.Name
		if dotIdx := strings.Index(name, "."); dotIdx != -1 {
			name = typ.Name[:dotIdx+1] + this.UpperHead(typ.Name[dotIdx+1:])
		}

		str = "*" + name
	}

	return str
}

func (this *Package) GenServiceMethodArguments(fields []*parser.Field) string {
	var str string

	maxIdx := len(fields) - 1
	for idx, field := range fields {
		str += fmt.Sprintf("%s %s", field.Name, this.GenTypeString(field.Name, field.Type, nil, field.Optional))
		if idx != maxIdx {
			str += ", "
		}
	}

	return str
}

func (this *Package) GenServiceMethodReturn(method *parser.Method) string {
	if method.ReturnType == nil {
		return "error"
	}

	return fmt.Sprintf("(%s, error)", this.GenTypeString("method return value", method.ReturnType, nil, false))
}

func (this *Package) GenWebApiServiceParams(fields []*parser.Field) string {
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

func (this *Package) GenConstants(constant *parser.Constant) string {
	format, ok := constantValueFormat[constant.Type.Name]
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s %s = "+format, constant.Name, typeStrs[constant.Type.Name], constant.Value)
}

func (this *Package) genOutputFilename(typ string) string {
	return fmt.Sprintf("gen_%s_%s.go", this.PkgName, typ)
}

func (this *Package) render(tplName string, wr io.Writer) error {
	return tpl.ExecuteTemplate(wr, tplName, this)
}

func (this *Package) renderToFile(dir, typ, tplName string) error {
	filename := this.genOutputFilename(typ)

	path := filepath.Join(dir, filename)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return this.render(tplName, file)
}
