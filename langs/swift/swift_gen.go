package swift

import (
	"github.com/ezbuy/tgen/langs"
	"github.com/samuel/go-thrift/parser"
)

type SwiftGen struct {
}

func (o *SwiftGen) Generate(parsedThrift map[string]*parser.Thrift) {

}

func init() {
	langs.Langs["swift"] = &SwiftGen{}
}
