package javascript

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ezbuy/tgen/global"
	"github.com/ezbuy/tgen/langs"
	"github.com/ezbuy/tgen/tmpl"
	"github.com/samuel/go-thrift/parser"
)

const TPL_SERVICE = "tgen/javascript/service"

type JavaScriptGen struct {
	langs.BaseGen
}

func (this *JavaScriptGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	generateWithModel(this, global.MODE_REST, output, parsedThrift)
}

type BaseJavaScript struct {
	t  *parser.Thrift
	ts *map[string]*parser.Thrift
}

func (this *BaseJavaScript) AssembleParams(method *parser.Method) string {
	params := []string{}

	for _, arg := range method.Arguments {
		params = append(params, fmt.Sprintf("%s: %s", arg.Name, this.typecast(arg.Type)))
		// params = append(params, fmt.Sprintf("%s", arg.Name))
	}
	return strings.Join(params, " , ")
}

func (this *BaseJavaScript) MethodReturnType(method *parser.Method) string {
	return this.typecast(method.ReturnType)
}

func (this *BaseJavaScript) typecast(t *parser.Type) string {
	if t != nil {
		switch t.Name {
		case langs.ThriftTypeI16, langs.ThriftTypeI32, langs.ThriftTypeI64, langs.ThriftTypeDouble:
			return "number"
		case langs.ThriftTypeString:
			return "string"
		case langs.ThriftTypeBool:
			return "boolean"
		case langs.ThriftTypeList, langs.ThriftTypeSet:
			return "JSONArray"
		case langs.ThriftTypeMap:
			return "JSONObject"
		default:
			return "object"

		}
	}
	return "null"
}

type javaService struct {
	*BaseJavaScript
	*parser.Service
}

func generateWithModel(gen *JavaScriptGen, m string, output string, parsedThrift map[string]*parser.Thrift) {
	if m != global.MODE_REST && m != global.MODE_JSONRPC {
		log.Fatalf("mode '%s' is invalid", m)
	}

	// gen.Init("javascript", parsedThrift)
	gen.Lang = "javascript"
	gen.Thrifts = parsedThrift

	if err := os.MkdirAll(output, 0755); err != nil {

		panic(fmt.Errorf("failed to create output directory %s", output))
	}

	// init templates
	var servicetpl *template.Template
	if m == global.MODE_REST {
		servicetpl = initemplate(TPL_SERVICE, "tmpl/javascript/rest_service.gojavascript")
	} else if m == global.MODE_JSONRPC {
		servicetpl = initemplate(TPL_SERVICE, "tmpl/javascript/jsonrpc_service.gojavascript")
	}

	for _, t := range parsedThrift {

		log.Printf("## services")

		for _, s := range t.Services {
			// filename is the service name plus 'Service'
			name := s.Name + "Service.js"

			// fix java file path
			if err := os.MkdirAll(output, 0755); err != nil {
				panic(fmt.Errorf("failed to create output directory %s", output))
			}

			path := filepath.Join(output, name)

			base := BaseJavaScript{t: t, ts: &parsedThrift}
			data := &javaService{BaseJavaScript: &base, Service: s}

			if err := outputfile(path, servicetpl, TPL_SERVICE, data); err != nil {
				panic(fmt.Errorf("failed to write file %s. error: %v\n", path, err))
			}

			log.Printf("%s", path)
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
	langs.Langs["javascript"] = &JavaScriptGen{}
}
