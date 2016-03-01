package java

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ezbuy/tgen/langs"
	"github.com/ezbuy/tgen/tmpl"
	"github.com/samuel/go-thrift/parser"
)

const (
	JavaTypeshort  = "short"
	JavaTypeint    = "int"
	JavaTypelong   = "long"
	JavaTypebool   = "boolean"
	JavaTypebyte   = "byte"
	JavaTypedouble = "double"

	JavaTypeString = "String"

	JavaTypeShort  = "Short"
	JavaTypeInt    = "Integer"
	JavaTypeLong   = "Long"
	JavaTypeBool   = "Boolean"
	JavaTypeByte   = "Byte"
	JavaTypeDouble = "Double"

	// other types (such as array, map, etc.) are implemented in the method 'Typecast'
)

const (
	TPL_STRUCT  = "tgen/java/struct"
	TPL_SERVICE = "tgen/java/service"
)

var plaintypemapping = map[string]string{
	langs.ThriftTypeI16:    JavaTypeshort,
	langs.ThriftTypeI32:    JavaTypeint,
	langs.ThriftTypeI64:    JavaTypelong,
	langs.ThriftTypeString: JavaTypeString,
	langs.ThriftTypeByte:   JavaTypebyte,
	langs.ThriftTypeBool:   JavaTypebool,
	langs.ThriftTypeDouble: JavaTypedouble,
}

var objecttypemapping = map[string]string{
	langs.ThriftTypeI16:    JavaTypeShort,
	langs.ThriftTypeI32:    JavaTypeInt,
	langs.ThriftTypeI64:    JavaTypeLong,
	langs.ThriftTypeString: JavaTypeString,
	langs.ThriftTypeByte:   JavaTypeByte,
	langs.ThriftTypeBool:   JavaTypeBool,
	langs.ThriftTypeDouble: JavaTypeDouble,
}

type JavaGen struct {
	langs.BaseGen
}

type BaseJava struct{}

func (this *BaseJava) PlainTypecast(t *parser.Type) string {
	return this.typecast(t, true)
}

func (this *BaseJava) ObjectTypecast(t *parser.Type) string {
	return this.typecast(t, false)
}

func (this *BaseJava) typecast(t *parser.Type, isplain bool) string {
	if t == nil {
		if isplain {
			return "void"
		} else {
			return "Void"
		}
	}

	var typemapping map[string]string

	if isplain {
		typemapping = plaintypemapping
	} else {
		typemapping = objecttypemapping
	}

	if t, ok := typemapping[t.Name]; ok {
		return t
	}

	switch t.Name {
	case langs.ThriftTypeList, langs.ThriftTypeSet:
		return fmt.Sprintf("ArrayList<%s>", this.ObjectTypecast(t.ValueType))
	case langs.ThriftTypeMap:
		return fmt.Sprintf("Map<%s, %s>", this.ObjectTypecast(t.KeyType), this.ObjectTypecast(t.ValueType))
	default:
		return t.Name
	}
}

func (this *BaseJava) AssembleParams(method *parser.Method) string {
	var buf bytes.Buffer

	for i, arg := range method.Arguments {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(fmt.Sprintf("final %s %s", this.PlainTypecast(arg.Type), arg.Name))
	}

	if len(method.Arguments) == 0 {
		buf.WriteString("")
	} else {
		buf.WriteString(", ")
	}

	buf.WriteString(fmt.Sprintf("final Listener<%s> listener", this.ObjectTypecast(method.ReturnType)))

	return buf.String()
}

func (this *BaseJava) GetInnerType(t *parser.Type) string {
	if t == nil {
		return "Void"
	}

	// map is ignored
	if t.Name == langs.ThriftTypeList || t.Name == langs.ThriftTypeSet {
		return this.GetInnerType(t.ValueType)
	}

	return this.ObjectTypecast(t)
}

type javaStruct struct {
	*BaseJava
	Namespace string
	*parser.Struct
}

type javaService struct {
	*BaseJava
	Namespace string
	*parser.Service
}

func (this *JavaGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init("java", parsedThrift)

	generatejsonrpc(filepath.Join(output, "jsonrpc"), parsedThrift)
	genraterest(filepath.Join(output, "rest"), parsedThrift)
}

func generatejsonrpc(output string, parsedThrift map[string]*parser.Thrift) {
	dogenerate(output, 0, parsedThrift)
}

func genraterest(output string, parsedThrift map[string]*parser.Thrift) {
	dogenerate(output, 1, parsedThrift)
}

// flag: 0-jsonrpc, 1-rest
func dogenerate(output string, flag int16, parsedThrift map[string]*parser.Thrift) {
	if err := os.MkdirAll(output, 0755); err != nil {
		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	var structpl *template.Template
	var servicetpl *template.Template

	// key is the absoule path of thrift file
	for tf, t := range parsedThrift {
		// due to java's features,
		// we generate the struct and service in seperate template file

		namespace, ok := t.Namespaces["java"]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: namespace not found in file[%s] of language[java]\n", tf)
			return
		}

		for _, s := range t.Structs {
			if structpl == nil {
				if flag == 0 {
					structpl = initemplate(TPL_STRUCT, "tmpl/java/jsonrpc_struct.gojava")
				} else if flag == 1 {
					structpl = initemplate(TPL_STRUCT, "tmpl/java/rest_struct.gojava")
				}
			}

			// filename is the struct name
			name := s.Name + ".java"

			path := filepath.Join(output, name)

			data := &javaStruct{BaseJava: &BaseJava{}, Namespace: namespace, Struct: s}

			if err := outputfile(path, structpl, TPL_STRUCT, data); err != nil {
				panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
			}
		}

		for _, s := range t.Services {
			if servicetpl == nil {
				if flag == 0 {
					servicetpl = initemplate(TPL_SERVICE, "tmpl/java/jsonrpc_service.gojava")
				} else if flag == 1 {
					servicetpl = initemplate(TPL_SERVICE, "tmpl/java/rest_service.gojava")
				}
			}

			// filename is the service name plus 'Service'
			name := s.Name + "Service.java"

			path := filepath.Join(output, name)

			data := &javaService{BaseJava: &BaseJava{}, Namespace: namespace, Service: s}

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
	langs.Langs["java"] = &JavaGen{}
}
