package gogen

import (
	"text/template"

	"github.com/ezbuy/tgen/tmpl/golang"
)

var tpl *template.Template

func Tpl() *template.Template {
	return tpl
}

func init() {
	funcMap := template.FuncMap{
		"upperHead":     upperHead,
		"genTypeString": genTypeString,
	}

	tpl = template.New("tgen/golang").Funcs(funcMap)

	files := []string{
		"tmpl/golang/include.gogo",
		"tmpl/golang/struct.gogo",
		"tmpl/golang/structs_file.gogo",
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
