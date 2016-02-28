package langs

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/ezbuy/tgen/tmpl"
)

func InitTemplate(tplpath string) *template.Template {
	tpldata, err := tmpl.Asset(tplpath)
	if err != nil {
		panic(err)
	}

	tpl, err := template.New("").Parse(string(tpldata))
	if err != nil {
		panic(err)
	}

	return tpl
}

func RenderTemplate(tpl *template.Template, data interface{}) []byte {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Write(path string, data []byte) {
	// save to disk
	// if the file exist, we intend to overwrite it
	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		panic(err)
	}
}
