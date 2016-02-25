package gogen

import (
	"testing"

	"github.com/samuel/go-thrift/parser"
)

func TestGenTypeString(t *testing.T) {
	fieldName := "testfield"

	cases := []struct {
		typ      *parser.Type
		optional bool
		isMapKey bool
		result   string
	}{
		// bool
		{
			&parser.Type{
				Name: TypeBool,
			},
			false,
			false,
			"bool",
		},
		{
			&parser.Type{
				Name: TypeBool,
			},
			true,
			false,
			"*bool",
		},

		// byte
		{
			&parser.Type{
				Name: TypeByte,
			},
			false,
			false,
			"byte",
		},
		{
			&parser.Type{
				Name: TypeByte,
			},
			true,
			false,
			"*byte",
		},

		// int16
		{
			&parser.Type{
				Name: TypeI16,
			},
			false,
			false,
			"int16",
		},
		{
			&parser.Type{
				Name: TypeI16,
			},
			true,
			false,
			"*int16",
		},

		// int32
		{
			&parser.Type{
				Name: TypeI32,
			},
			false,
			false,
			"int32",
		},
		{
			&parser.Type{
				Name: TypeI32,
			},
			true,
			false,
			"*int32",
		},

		// int64
		{
			&parser.Type{
				Name: TypeI64,
			},
			false,
			false,
			"int64",
		},
		{
			&parser.Type{
				Name: TypeI64,
			},
			true,
			false,
			"*int64",
		},

		// double
		{
			&parser.Type{
				Name: TypeDouble,
			},
			false,
			false,
			"float64",
		},
		{
			&parser.Type{
				Name: TypeDouble,
			},
			true,
			false,
			"*float64",
		},

		// binary
		{
			&parser.Type{
				Name: TypeBinary,
			},
			false,
			false,
			"[]byte",
		},
		{
			&parser.Type{
				Name: TypeBinary,
			},
			true,
			false,
			"[]byte",
		},

		// string
		{
			&parser.Type{
				Name: TypeString,
			},
			false,
			false,
			"string",
		},
		{
			&parser.Type{
				Name: TypeString,
			},
			true,
			false,
			"*string",
		},

		// custom name
		{
			&parser.Type{
				Name: "SomeStruct",
			},
			false,
			false,
			"SomeStruct",
		},
		{
			&parser.Type{
				Name: "SomeStruct",
			},
			true,
			false,
			"*SomeStruct",
		},

		// list<bool>
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: TypeBool,
				},
			},
			false,
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
			true,
			false,
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
			false,
			false,
			"[]SomeStruct",
		},
		{
			&parser.Type{
				Name: TypeList,
				ValueType: &parser.Type{
					Name: "SomeStruct",
				},
			},
			true,
			false,
			"[]SomeStruct",
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
			false,
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
			true,
			false,
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
			false,
			false,
			"[][]SomeStruct",
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
			true,
			false,
			"[][]SomeStruct",
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
			false,
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
			true,
			false,
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
			false,
			false,
			"map[string][]SomeStruct",
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
			true,
			false,
			"map[string][]SomeStruct",
		},
	}

	for _, one := range cases {
		str := genTypeString(fieldName, one.typ, one.optional, one.isMapKey)
		if str != one.result {
			t.Errorf("expected: %q, got: %q", one.result, str)
		}
	}
}
