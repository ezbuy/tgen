package gogen

import (
	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

const langName = "go"

type GoGen struct {
	langs.BaseGen
}

func (this *GoGen) Generate(parsedThrift map[string]*parser.Thrift) {
	this.BaseGen.Init(langName, parsedThrift)

}

func init() {
	langs.Langs[langName] = &GoGen{}
}
