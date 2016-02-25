package langs

import "github.com/samuel/go-thrift/parser"

type GenResult struct {
	Filename string
	Data     []byte
}

type ApiGen interface {
	Generate(tplpath string, parsedThrift map[string]*parser.Thrift) ([]GenResult, error)
}

// the key of Langs is language
var Langs = make(map[string]ApiGen)
