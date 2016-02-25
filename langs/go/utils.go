package gogen

import (
	"fmt"

	"github.com/samuel/go-thrift/parser"
)

const (
	TypeBool   = "bool"
	TypeByte   = "byte"
	TypeI16    = "i16"
	TypeI32    = "i32"
	TypeI64    = "164"
	TypeDouble = "double"
	TypeBinary = "binary"
	TypeString = "string"
	TypeList   = "list"
	TypeMap    = "map"
	TypeSet    = "set"
)

var typeStrs = map[string]string{
	TypeBool:   "bool",
	TypeByte:   "byte",
	TypeI16:    "int16",
	TypeI32:    "int32",
	TypeI64:    "int64",
	TypeDouble: "float64",
	TypeBinary: "[]byte",
	TypeString: "string",
}

func genTypeString(fieldName string, typ *parser.Type, optional bool, isMapKey bool) string {
	if typ == nil {
		panicWithErr("field %s with nil type", fieldName)
	}

	var str string

	switch typ.Name {
	case TypeBool, TypeByte, TypeI16, TypeI32, TypeI64, TypeDouble, TypeString:
		if optional {
			str = "*"
		}
		str += typeStrs[typ.Name]

	case TypeBinary:
		if isMapKey {
			panicWithErr("map field %s with binary key", fieldName)
		}
		str = typeStrs[TypeBinary]

	case TypeList:
		if isMapKey {
			panicWithErr("map field %s with list key", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("list field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("[]%s", genTypeString(fieldName, typ.ValueType, false, false))

	case TypeMap:
		if isMapKey {
			panicWithErr("map field %s with map key", fieldName)
		}

		if typ.KeyType == nil {
			panicWithErr("map field %s with nil key type", fieldName)
		}

		if typ.ValueType == nil {
			panicWithErr("map field %s with nil value type", fieldName)
		}

		str = fmt.Sprintf("map[%s]%s",
			genTypeString(fieldName, typ.KeyType, false, true),
			genTypeString(fieldName, typ.ValueType, false, false),
		)

	case TypeSet:
		// TODO: support set

	default:
		if typ.Name == "" {
			panicWithErr("field %s without type name", fieldName)
		}

		if optional {
			str = "*"
		}
		str += typ.Name
	}

	return str
}

func panicWithErr(format string, msg ...interface{}) {
	panic(fmt.Errorf(format, msg...))
}
