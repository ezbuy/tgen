package gogen

import (
	"os"
	"text/template"

	"github.com/ezbuy/tgen/tmpl/golang"
	"github.com/samuel/go-thrift/parser"
)

var tpl *template.Template

func Tpl() *template.Template {
	return tpl
}

func init() {
	funcMap := template.FuncMap{
		"upperHead":     upperHead,
		"genTypeString": genTypeString,
		"isLast":        tplIsLast,
		"isNilType":     tplIsNilType,
	}

	tpl = template.New("tgen/golang").Funcs(funcMap)

	files := []string{
		"tmpl/golang/include.gogo",
		"tmpl/golang/struct.gogo",
		"tmpl/golang/structs_file.gogo",
		"tmpl/golang/service.gogo",
		"tmpl/golang/services_file.gogo",
		"tmpl/golang/echo_module.gogo",
	}

	for _, filename := range files {
		data, err := gotpl.Asset(filename)
		if err != nil {
			panic(err)
		}

		if _, err = tpl.Parse(string(data)); err != nil {
			panic(err)
		}
	}
}

func outputFile(path string, tplName string, data interface{}) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return tpl.ExecuteTemplate(file, tplName, data)
}

func tplIsLast(idx, size int) bool {
	return idx == size-1
}

func tplIsNilType(typ *parser.Type) bool {
	return typ == nil
}
