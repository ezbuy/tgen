package gogen_test

import (
	"encoding/json"
	"reflect"
	"testing"
)

type testArgJsonUnmarshaler interface {
	json.Unmarshaler
	GetArg() interface{}
}

type testSimpleBinaryArgArguments struct {
	Arg []byte `thrift:"1,required" json:"arg"`
}

func (this *testSimpleBinaryArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg []byte `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleBinaryArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleBoolArgArguments struct {
	Arg bool `thrift:"1,required" json:"arg"`
}

func (this *testSimpleBoolArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg bool `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleBoolArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleByteArgArguments struct {
	Arg byte `thrift:"1,required" json:"arg"`
}

func (this *testSimpleByteArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg byte `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleByteArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleDoubleArgArguments struct {
	Arg float64 `thrift:"1,required" json:"arg"`
}

func (this *testSimpleDoubleArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg float64 `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleDoubleArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleI16ArgArguments struct {
	Arg int16 `thrift:"1,required" json:"arg"`
}

func (this *testSimpleI16ArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg int16 `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleI16ArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleI32ArgArguments struct {
	Arg int32 `thrift:"1,required" json:"arg"`
}

func (this *testSimpleI32ArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg int32 `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleI32ArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleI64ArgArguments struct {
	Arg int64 `thrift:"1,required" json:"arg"`
}

func (this *testSimpleI64ArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg int64 `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleI64ArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleListArgArguments struct {
	Arg []string `thrift:"1,required" json:"arg"`
}

func (this *testSimpleListArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg []string `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleListArgArguments) GetArg() interface{} {
	return this.Arg
}

type testSimpleStringArgArguments struct {
	Arg string `thrift:"1,required" json:"arg"`
}

func (this *testSimpleStringArgArguments) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &this.Arg); err == nil {
		return nil
	}

	result := struct {
		Arg string `json:"arg"`
	}{}

	err := json.Unmarshal(data, &result)
	if err == nil {
		this.Arg = result.Arg
	}

	return err
}

func (this *testSimpleStringArgArguments) GetArg() interface{} {
	return this.Arg
}

func TestSimpleArgumentsUnmarshal(t *testing.T) {
	oneByte := byte(47)
	encodedByte, err := json.Marshal(oneByte)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	binary := []byte{1, 3, 5, 7, 9}
	encodedBinary, err := json.Marshal(binary)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	testCases := []struct {
		data        string
		unmarshaler testArgJsonUnmarshaler
		arg         interface{}
	}{
		{
			string(encodedBinary),
			&testSimpleBinaryArgArguments{},
			binary,
		},
		{
			"{\"arg\": " + string(encodedBinary) + "}",
			&testSimpleBinaryArgArguments{},
			binary,
		},
		{
			string(encodedByte),
			&testSimpleByteArgArguments{},
			oneByte,
		},
		{
			"{\"arg\": " + string(encodedByte) + "}",
			&testSimpleByteArgArguments{},
			oneByte,
		},
		{
			"true",
			&testSimpleBoolArgArguments{},
			true,
		},
		{
			"{\"arg\": true}",
			&testSimpleBoolArgArguments{},
			true,
		},
		{
			"17",
			&testSimpleI16ArgArguments{},
			int16(17),
		},
		{
			"{\"arg\": 17}",
			&testSimpleI16ArgArguments{},
			int16(17),
		},
		{
			"33",
			&testSimpleI32ArgArguments{},
			int32(33),
		},
		{
			"{\"arg\": 33}",
			&testSimpleI32ArgArguments{},
			int32(33),
		},
		{
			"65",
			&testSimpleI64ArgArguments{},
			int64(65),
		},
		{
			"{\"arg\": 65}",
			&testSimpleI64ArgArguments{},
			int64(65),
		},
		{
			"12.34",
			&testSimpleDoubleArgArguments{},
			float64(12.34),
		},
		{
			"{\"arg\": 12.34}",
			&testSimpleDoubleArgArguments{},
			float64(12.34),
		},
		{
			"\"this is a string\"",
			&testSimpleStringArgArguments{},
			"this is a string",
		},
		{
			"{\"arg\": \"this is a string\"}",
			&testSimpleStringArgArguments{},
			"this is a string",
		},
		{
			"[\"a\", \"b\", \"c\"]",
			&testSimpleListArgArguments{},
			[]string{"a", "b", "c"},
		},
		{
			"{\"arg\": [\"a\", \"b\", \"c\"]}",
			&testSimpleListArgArguments{},
			[]string{"a", "b", "c"},
		},
	}

	for idx, one := range testCases {
		err := one.unmarshaler.UnmarshalJSON([]byte(one.data))
		if err != nil {
			t.Errorf("case %d: %#v : expected nil error, got %s", idx, one, err)
		}

		if got := one.unmarshaler.GetArg(); !reflect.DeepEqual(got, one.arg) {
			t.Errorf("case %d: expected %#v, got %#v", idx, one.arg, got)
		}
	}
}
