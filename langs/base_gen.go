package langs

import (
	"log"

	"github.com/samuel/go-thrift/parser"
)

type BaseGen struct {
	Lang    string
	Thrifts map[string]*parser.Thrift
}

func (g *BaseGen) Init(lang string, parsedThrift map[string]*parser.Thrift) {
	g.Lang = lang
	g.Thrifts = parsedThrift
	g.CheckNamespace()
}

func (g *BaseGen) CheckNamespace() {
	for f, t := range g.Thrifts {
		if _, ok := t.Namespaces[g.Lang]; !ok {
			log.Fatalf("Namespace not found for language '%s' in file '%s'", g.Lang, f)
		}
	}
}
