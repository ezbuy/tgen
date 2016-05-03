package gogen

import (
	"fmt"
	"testing"

	"github.com/samuel/go-thrift/parser"
)

func TestType(t *testing.T) {
	cases := []struct {
		expect string
		actual string
	}{
		{
			"bool",
			TypeBool,
		},
		{
			"byte",
			TypeByte,
		},
		{
			"i16",
			TypeI16,
		},
		{
			"i32",
			TypeI32,
		},
		{
			"i64",
			TypeI64,
		},
		{
			"double",
			TypeDouble,
		},
		{
			"binary",
			TypeBinary,
		},
		{
			"string",
			TypeString,
		},
		{
			"list",
			TypeList,
		},
		{
			"set",
			TypeSet,
		},
		{
			"map",
			TypeMap,
		},
	}

	for _, one := range cases {
		if one.expect != one.actual {
			t.Errorf("expect: %q; actual: %q", one.expect, one.actual)
		}
	}
}

func TestGenTypeString(t *testing.T) {
	fieldName := "testfield"

	cases := []struct {
		typ      *parser.Type
		parent   *parser.Type
		optional bool
		result   string
	}{
		// bool
		{
			&parser.Type{
				Name: TypeBool,
			},
			nil,
			false,
			"bool",
		},
		{
			&parser.Type{
				Name: TypeBool,
			},
			nil,
			true,
			"*bool",
		},

		// byte
		{
			&parser.Type{
				Name: TypeByte,
			},
			nil,
			false,
			"byte",
		},
		{
			&parser.Type{
				Name: TypeByte,
			},
			nil,
			true,
			"*byte",
		},

		// int16
		{
			&parser.Type{
				Name: TypeI16,
			},
			nil,
			false,
			"int16",
		},
		{
			&parser.Type{
				Name: TypeI16,
			},
			nil,
			true,
			"*int16",
		},

		// int32
		{
			&parser.Type{
				Name: TypeI32,
			},
			nil,
			false,
			"int32",
		},
		{
			&parser.Type{
				Name: TypeI32,
			},
			nil,
			true,
			"*int32",
		},

		// int64
		{
			&parser.Type{
				Name: TypeI64,
			},
			nil,
			false,
			"int64",
		},
		{
			&parser.Type{
				Name: TypeI64,
			},
			nil,
			true,
			"*int64",
		},

		// double
		{
			&parser.Type{
				Name: TypeDouble,
			},
			nil,
			false,
			"float64",
		},
		{
			&parser.Type{
				Name: TypeDouble,
			},
			nil,
			true,
			"*float64",
		},

		// binary
		{
			&parser.Type{
				Name: TypeBinary,
			},
			nil,
			false,
			"[]byte",
		},
		{
			&parser.Type{
				Name: TypeBinary,
			},
			nil,
			true,
			"[]byte",
		},

		// string
		{
			&parser.Type{
				Name: TypeString,
			},
			nil,
			false,
			"string",
		},
		{
			&parser.Type{
				Name: TypeString,
			},
			nil,
			true,
			"*string",
		},

		// custom name
		{
			&parser.Type{
				Name: "SomeStruct",
			},
			nil,
			false,
			"*SomeStruct",
		},
		{
			&parser.Type{
				Name: "SomeStruct",
			},
			nil,
			true,
			"*SomeStruct",
		},

		// included custom name
		{
			&parser.Type{
				Name: "SomeIncludes.UpperStruct",
			},
			nil,
			false,
			"*SomeIncludes.UpperStruct",
		},
		{
			&parser.Type{
				Name: "SomeIncludes.UpperStruct",
			},
			nil,
			true,
			"*SomeIncludes.UpperStruct",
		},
		{
			&parser.Type{
				Name: "SomeIncludes.lowerStruct",
			},
			nil,
			false,
			"*SomeIncludes.LowerStruct",
		},
		{
			&parser.Type{
				Name: "SomeIncludes.lowerStruct",
			},
			nil,
			true,
			"*SomeIncludes.LowerStruct",
		},

		// list<bool>
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeBool,
				},
			},
			nil,
			false,
			"[]bool",
		},
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeBool,
				},
			},
			nil,
			true,
			"[]bool",
		},

		// list<SomeStruct>
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: "SomeStruct",
				},
			},
			nil,
			false,
			"[]*SomeStruct",
		},
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: "SomeStruct",
				},
			},
			nil,
			true,
			"[]*SomeStruct",
		},

		// list<list<bool>>
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: TypeBool,
					},
				},
			},
			nil,
			false,
			"[][]bool",
		},
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: TypeBool,
					},
				},
			},
			nil,
			true,
			"[][]bool",
		},

		// list<list<SomeStruct>>
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: "SomeStruct",
					},
				},
			},
			nil,
			false,
			"[][]*SomeStruct",
		},
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: "SomeStruct",
					},
				},
			},
			nil,
			true,
			"[][]*SomeStruct",
		},

		// map<string, bool>
		{
			&parser.Type{
				Name: TypeMap,
				KeyType: &parser.Type{
					Name: TypeString,
				},
				ValueType: &parser.Type{
					Name: TypeBool,
				},
			},
			nil,
			false,
			"map[string]bool",
		},
		{
			&parser.Type{
				Name: TypeMap,
				KeyType: &parser.Type{
					Name: TypeString,
				},
				ValueType: &parser.Type{
					Name: TypeBool,
				},
			},
			nil,
			true,
			"map[string]bool",
		},

		// map<string, list<SomeStruct>>
		{
			&parser.Type{
				Name: TypeMap,
				KeyType: &parser.Type{
					Name: TypeString,
				},
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: "SomeStruct",
					},
				},
			},
			nil,
			false,
			"map[string][]*SomeStruct",
		},
		{
			&parser.Type{
				Name: TypeMap,
				KeyType: &parser.Type{
					Name: TypeString,
				},
				ValueType: &parser.Type{
					Name: TypeList,
					ValueType: &parser.Type{
						Name: "SomeStruct",
					},
				},
			},
			nil,
			true,
			"map[string][]*SomeStruct",
		},
	}

	thrift := &parser.Thrift{}
	thrift.Structs = map[string]*parser.Struct{
		"SomeStruct": &parser.Struct{},
	}

	pkg := newPackage(thrift)

	includeTrhift := &parser.Thrift{}
	includeTrhift.Structs = map[string]*parser.Struct{
		"lowerStruct": &parser.Struct{},
		"UpperStruct": &parser.Struct{},
	}

	includePkg := newPackage(includeTrhift)
	includePkg.PkgName = "SomeIncludes"

	pkg.includes = map[string]*Package{
		"SomeIncludes": includePkg,
	}

	for _, one := range cases {
		str := pkg.GenTypeString(fieldName, one.typ, one.parent, one.optional)
		if str != one.result {
			t.Errorf("expected: %q, got: %q", one.result, str)
		}
	}
}

type genTypeStringPanicTestCase struct {
	typ       *parser.Type
	parent    *parser.Type
	optional  bool
	err       string
	recovered interface{}
}

func TestGenTypeStringPanics(t *testing.T) {
	fieldName := "testfield"

	typeMapWithBinaryKey := &parser.Type{
		Name: TypeMap,
		KeyType: &parser.Type{
			Name: TypeBinary,
		},
	}

	typeMapWithListKey := &parser.Type{
		Name: TypeMap,
		KeyType: &parser.Type{
			Name: TypeList,
		},
	}

	typeListWithNilValue := &parser.Type{
		Name:      TypeList,
		ValueType: nil,
	}

	typeMapWithMapKey := &parser.Type{
		Name: TypeMap,
		KeyType: &parser.Type{
			Name: TypeMap,
		},
	}

	typeMapWithNilKey := &parser.Type{
		Name:    TypeMap,
		KeyType: nil,
	}

	typeMapWithNilValue := &parser.Type{
		Name: TypeMap,
		KeyType: &parser.Type{
			Name: TypeString,
		},
		ValueType: nil,
	}

	cases := []genTypeStringPanicTestCase{
		// nil type
		{
			nil,
			nil,
			false,
			fmt.Sprintf("field %s with nil type", fieldName),
			nil,
		},

		// map with binary key
		{
			typeMapWithBinaryKey.KeyType,
			typeMapWithBinaryKey,
			false,
			fmt.Sprintf("map field %s with binary key", fieldName),
			nil,
		},

		// map with list key
		{
			typeMapWithListKey.KeyType,
			typeMapWithListKey,
			false,
			fmt.Sprintf("map field %s with list key", fieldName),
			nil,
		},

		// list with nil value
		{
			typeListWithNilValue,
			nil,
			false,
			fmt.Sprintf("list field %s with nil value type", fieldName),
			nil,
		},

		// map with map key
		{
			typeMapWithMapKey.KeyType,
			typeMapWithMapKey,
			false,
			fmt.Sprintf("map field %s with map key", fieldName),
			nil,
		},

		// map with nil key
		{
			typeMapWithNilKey,
			nil,
			false,
			fmt.Sprintf("map field %s with nil key type", fieldName),
			nil,
		},

		// map with nil value
		{
			typeMapWithNilValue,
			nil,
			false,
			fmt.Sprintf("map field %s with nil value type", fieldName),
			nil,
		},

		// without type name
		{
			&parser.Type{
				Name: "",
			},
			nil,
			false,
			fmt.Sprintf("field %s without type name", fieldName),
			nil,
		},
	}

	for _, one := range cases {
		ptr := &one
		getGenTypeStringPanic(fieldName, ptr)

		err, ok := (ptr.recovered).(error)
		if !ok {
			t.Errorf("expected an error, got %s", ptr.recovered)
			continue
		}

		if str := err.Error(); str != ptr.err {
			t.Errorf("expected: %q, got: %q", ptr.err, str)
		}

	}

}

func TestGenConstant(t *testing.T) {
	cases := []struct {
		constant *parser.Constant
		out      string
	}{
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeI16,
				},
				Name:  "ConstantI16",
				Value: 16,
			},
			out: "ConstantI16 int16 = 16",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeI32,
				},
				Name:  "ConstantI32",
				Value: 32,
			},
			out: "ConstantI32 int32 = 32",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeI64,
				},
				Name:  "ConstantI64",
				Value: 64,
			},
			out: "ConstantI64 int64 = 64",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeDouble,
				},
				Name:  "ConstantDouble",
				Value: 128.128,
			},
			out: "ConstantDouble float64 = 128.128000",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeBool,
				},
				Name:  "ConstantBool",
				Value: parser.Identifier("false"),
			},
			out: "ConstantBool bool = false",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeByte,
				},
				Name:  "ConstantByte",
				Value: 65,
			},
			out: "ConstantByte byte = 65",
		},
		{
			constant: &parser.Constant{
				Type: &parser.Type{
					Name: TypeString,
				},
				Name:  "ConstantString",
				Value: "this is string",
			},
			out: `ConstantString string = "this is string"`,
		},
	}

	pkg := &Package{}

	for _, one := range cases {
		if out := pkg.GenConstants(one.constant); out != one.out {
			t.Errorf("expected %q, got %q", one.out, out)
		}
	}
}

func getGenTypeStringPanic(fieldName string, testCase *genTypeStringPanicTestCase) {
	defer func() {
		testCase.recovered = recover()
	}()

	pkg := &Package{}
	pkg.GenTypeString(fieldName, testCase.typ, testCase.parent, testCase.optional)
	return
}
