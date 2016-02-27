package swift

import (
	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

type SwiftGen struct {
	langs.BaseGen
}

func (o *SwiftGen) Generate(output string, parsedThrift map[string]*parser.Thrift) {
	o.BaseGen.Init("swift", parsedThrift)
}

func init() {
	langs.Langs["swift"] = &SwiftGen{}
}
