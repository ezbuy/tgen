package langs

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samuel/go-thrift/parser"
)

type BaseGen struct {
	Lang      string
	Namespace string
	Thrifts   map[string]*parser.Thrift
}

func (g *BaseGen) Init(lang string, parsedThrift map[string]*parser.Thrift) {
	g.Lang = lang
	g.Thrifts = parsedThrift
	g.CheckNamespace()
}

func (g *BaseGen) CheckNamespace() {
	for _, thrift := range g.Thrifts {
		for lang, namepace := range thrift.Namespaces {
			if lang == g.Lang {
				g.Namespace = namepace
				return
			}
		}
	}

	w := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(w, "Namespace not found for: %s\n", g.Lang)
	os.Exit(2)
}
