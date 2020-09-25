package json

import (
	"github.com/beckjiang/go-json/pool"
	"github.com/francoispqt/gojay"
	"github.com/minio/simdjson-go"
)

type DecodeDefaultAction uint8

const (
	Skip DecodeDefaultAction = iota
	SaveToSource
)

const ActionSaveToSource = "save_to_source"
const ActionSkip = "skip"

type DecodeSchema struct {
	Name          []byte
	Type          ValueType
	Property      map[string]*DecodeSchema
	ElementSchema *DecodeSchema
	DefaultAction DecodeDefaultAction
}

func (ds *DecodeSchema) NeedsSaveToSource(name string) bool {
	return (ds.Property[name] == nil && ds.DefaultAction == SaveToSource) ||
		(ds.Property[name] != nil && ds.Property[name].Type == Source)
}

func (ds *DecodeSchema) NeedsSkip(name string) bool {
	return ds.Property[name] == nil && ds.DefaultAction == Skip
}

var decodeSourceWrapperPool *pool.Pool

func init() {
	decodeSourceWrapperPool = pool.PM.NewPool("source.wrapper", func() pool.RefCountable {
		return new(DecodeSourceWrapper)
	}, 128)
}

func BorrowDecodeSourceWrapper() *DecodeSourceWrapper {
	return decodeSourceWrapperPool.Acquire().(*DecodeSourceWrapper)
}

type DecodeSourceWrapper struct {
	pool.RefCount

	I *simdjson.Iter
	A *simdjson.Array
}

func (ds *DecodeSourceWrapper) IsNil() bool {
	return false
}

func (ds *DecodeSourceWrapper) Reset() {
	ds.I = nil
	ds.A = nil
}

func (ds *DecodeSourceWrapper) MarshalJSONObject(enc *gojay.Encoder) {
	var objBuf simdjson.Object
	obj, err := ds.I.Object(&objBuf)
	if err != nil {
		return
	}
	for {
		name, t, _ := obj.NextElementBytes(ds.I)
		if t == simdjson.TypeNone {
			break
		}
		_ = parseToSource(enc, bytesToString(name), ds.I)
	}

}

func (ds *DecodeSourceWrapper) MarshalJSONArray(enc *gojay.Encoder) {
	if ds.I == nil && ds.A == nil {
		panic("bug")
	}
	var arr *simdjson.Array
	if ds.A != nil {
		arr = ds.A
	} else {
		arr, _ = ds.I.Array(nil)
	}

	it := arr.Iter()
	var elem simdjson.Iter

	for i := 0; ; i++ {
		t, err := it.AdvanceIter(&elem)
		if err != nil {
			break
		}
		if t == simdjson.TypeNone {
			break
		}

		switch t {
		case simdjson.TypeInt:
			val, err := elem.Int()
			if err != nil {
				break
			}
			enc.Int64(val)

		case simdjson.TypeUint:
			val, err := elem.Int()
			if err != nil {
				break
			}
			enc.Int64(val)

		case simdjson.TypeFloat:
			val, err := elem.Float()
			if err != nil {
				break
			}
			enc.Float64(val)

		case simdjson.TypeBool:
			val, err := elem.Bool()
			if err != nil {
				break
			}
			enc.Bool(val)

		case simdjson.TypeNull:
			enc.Null()

		case simdjson.TypeString:
			val, err := elem.StringBytes()
			if err != nil {
				break
			}
			enc.String(string(val))

		case simdjson.TypeObject:
			ds := BorrowDecodeSourceWrapper()
			ds.I = &elem

			enc.Object(ds)
			ds.DecRef()

		case simdjson.TypeArray:
			ds := BorrowDecodeSourceWrapper()
			ds.I = &elem

			enc.Array(ds)
			ds.DecRef()
		}
	}
}

func NewDecodeSchemaByMap(m map[string]interface{}) *DecodeSchema {
	if m == nil {
		return nil
	}
	base := &DecodeSchema{DefaultAction: Skip}

	switch m["type"].(string) {
	case "int":
		base.Type = Int
	case "string":
		base.Type = String
	case "float":
		base.Type = Float
	case "object":
		base.Type = Object
		property := m["property"].(map[string]interface{})
		if base.Property == nil {
			base.Property = make(map[string]*DecodeSchema)
		}
		for k, s := range property {
			base.Property[k] = NewDecodeSchemaByMap(s.(map[string]interface{}))
		}

	case "array":
		base.Type = Array
		base.ElementSchema = NewDecodeSchemaByMap(m["element_schema"].(map[string]interface{}))
	}

	if m["default_action"] != nil {
		switch m["default_action"].(string) {
		case ActionSaveToSource:
			base.DefaultAction = SaveToSource
		case ActionSkip:
			base.DefaultAction = Skip
		}
	}

	return base
}