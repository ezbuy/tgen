package gogen

import (
	"github.com/samuel/go-thrift/parser"
)

type structsFileData struct {
	Package  string
	Includes [][2]string
	Structs  []*structData
}

type structData struct {
	Name   string
	Fields []*parser.Field
}
