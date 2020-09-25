package json

import (
	"github.com/beckjiang/go-json/pool"
	"io"
)

type ValueType uint8

const (
	None ValueType = iota
	Null
	Bool
	Int
	Float
	String
	Object
	Array
	Source
)

type Value interface {
	pool.RefCountable

	Type() ValueType

	// getter
	Bool() bool
	Int() int64
	Float() float64
	String(dst []byte) []byte

	// ObjectVal 获取Key对应的value，但Value的 ownership 仍归属原先的 object
	ObjectVal(key []byte) (found bool, val Value)

	// GrabObjectVal 获取Key对应的value，并将 ownership 转移至 caller,
	// 用于解决可能多次 Release() 的问题
	GrabObjectVal(key []byte) (found bool, val Value)

	// ArrayElem 获取 index 对应的value，但Value的 ownership 仍归属原先 array
	ArrayElem(index int) Value
	SetArrayElem(index int, val Value)

	ArraySize() int

	// setter
	SetBool(val bool)
	SetInt(val int64)
	SetFloat(val float64)
	SetString(val []byte)
	SetNull()
	SetObjectKV(key []byte, val Value)
	AddArrayElem(val Value)
	SetType(t ValueType)

	// DelArrayElem 删除 index 对应的value，并将 ownership 转移至 caller,
	// 用于解决可能多次 Release() 的问题
	DelArrayElem(index int) Value

	Decode(str []byte) error
	DecodeWithSchema(str []byte, schema *DecodeSchema) error
	Encode(dst []byte) []byte
	EncodeTo(w io.Writer)

	LookUp(path []byte) Value
	AppendSource(b []byte)
	HasSource() bool
	GetSource() []byte
	SetSource([]byte)
}
