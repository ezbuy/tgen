package java

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ezbuy/tgen/langs"
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

	structpl *template.Template
	servtpl  *template.Template
}

type BaseJava struct{}

func (bj *BaseJava) PlainTypecast(t *parser.Type) string {
	return bj.typecast(t, true)
}

func (bj *BaseJava) ObjectTypecast(t *parser.Type) string {
	return bj.typecast(t, false)
}

func (bj *BaseJava) typecast(t *parser.Type, isplain bool) string {
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
		return fmt.Sprintf("ArrayList<%s>", bj.ObjectTypecast(t.ValueType))
	case langs.ThriftTypeMap:
		return fmt.Sprintf("Map<%s, %s>", bj.ObjectTypecast(t.KeyType), bj.ObjectTypecast(t.ValueType))
	default:
		return t.Name
	}
}

func (bj *BaseJava) AssembleParams(method *parser.Method) string {
	var buf bytes.Buffer

	for i, arg := range method.Arguments {
		if i != 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(fmt.Sprintf("final %s %s", bj.PlainTypecast(arg.Type), arg.Name))
	}

	if len(method.Arguments) == 0 {
		buf.WriteString("")
	} else {
		buf.WriteString(", ")
	}

	buf.WriteString(fmt.Sprintf("final Listener<%s> listener", bj.ObjectTypecast(method.ReturnType)))

	return buf.String()
}

func (bj *BaseJava) GetInnerType(t *parser.Type) string {
	if t == nil {
		return "Void"
	}

	// map is ignored
	if t.Name == langs.ThriftTypeList || t.Name == langs.ThriftTypeSet {
		return bj.GetInnerType(t.ValueType)
	}

	return bj.ObjectTypecast(t)
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

func (o *JavaGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	o.BaseGen.Init("java", parsedThrift)

	// key is the absoule path of thrift file
	for tf, t := range parsedThrift {
		// due to java's features,
		// we generate the struct and service in seperate template file

		namespace, ok := t.Namespaces["java"]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: namespace not found in file[%s] of language[java]\n", tf)
			return
		}

		for _, m := range t.Structs {
			if o.structpl == nil {
				o.structpl = langs.InitTemplate("tmpl/java/java_struct.gojava")
			}

			data := o.genStruct(namespace, m)

			// filename is the struct name
			name := m.Name + ".java"

			path := filepath.Join(output, name)

			// save to disk
			langs.Write(path, data)

			fmt.Printf("[%s] generated\n", path)
		}

		for _, s := range t.Services {
			if o.servtpl == nil {
				o.servtpl = langs.InitTemplate("tmpl/java/java_service.gojava")
			}

			data := o.genService(namespace, s)

			// filename is the service name plus 'Service'
			name := s.Name + "Service.java"

			path := filepath.Join(output, name)

			// save to disk
			langs.Write(path, data)

			fmt.Printf("[%s] generated\n", path)
		}
	}
}

func (o *JavaGen) genStruct(ns string, s *parser.Struct) []byte {
	js := &javaStruct{BaseJava: &BaseJava{}, Namespace: ns, Struct: s}

	data := langs.RenderTemplate(o.structpl, js)

	return data
}

func (o *JavaGen) genService(ns string, s *parser.Service) []byte {
	js := &javaService{BaseJava: &BaseJava{}, Namespace: ns, Service: s}

	data := langs.RenderTemplate(o.servtpl, js)

	return data
}

func init() {
	langs.Langs["java"] = &JavaGen{}
}
