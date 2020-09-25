package json

import (
	"errors"
	"github.com/francoispqt/gojay"
	"github.com/minio/simdjson-go"
	"sync"
	"unsafe"
)

var ErrUnsupportedCPU = errors.New("Unsupported CPU")
var ErrUnknownType = errors.New("unknown type")

func bytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func getValueWithSchema(v Value, iter simdjson.Iter, schema *DecodeSchema) error {
	switch iter.Type() {
	case simdjson.TypeInt:
		val, err := iter.Int()
		if err != nil {
			return err
		}
		v.SetInt(val)

	case simdjson.TypeUint:
		val, err := iter.Uint()
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case simdjson.TypeFloat:
		val, err := iter.Float()
		if err != nil {
			return err
		}
		v.SetFloat(val)

	case simdjson.TypeBool:
		val, err := iter.Bool()
		if err != nil {
			return err
		}
		v.SetBool(val)

	case simdjson.TypeNull:
		v.SetNull()

	case simdjson.TypeString:
		val, err := iter.StringBytes()
		if err != nil {
			return err
		}
		v.SetString(val)

	case simdjson.TypeObject:
		var objBuf simdjson.Object
		obj, err := iter.Object(&objBuf)
		if err != nil {
			return err
		}

		if schema == nil {
			return parseObject(v, obj)
		} else {
			return parseObjectWithSchema(v, obj, schema)
		}

	case simdjson.TypeArray:
		arr, err := iter.Array(nil)
		if err != nil {
			return err
		}
		if schema == nil {
			return parseArray(v, arr)
		} else {
			return parseArrayWithSchema(v, arr, schema)
		}

	default:
		return ErrUnknownType
	}
	return nil
}

//var sourcePool *pool.Pool
//func init() {
//	sourcePool = pool.PM.NewPool("json.DecodeSourceWrapper", func() pool.RefCountable {
//		return new(DecodeSourceWrapper)
//	}, 128)
//}

func parseToSource(enc *gojay.Encoder, name string, iter *simdjson.Iter) error {

	switch iter.Type() {
	case simdjson.TypeInt:
		val, err := iter.Int()
		if err != nil {
			return err
		}
		enc.Int64Key(name, val)

	case simdjson.TypeUint:
		val, err := iter.Int()
		if err != nil {
			return err
		}
		enc.Int64Key(name, val)

	case simdjson.TypeFloat:
		val, err := iter.Float()
		if err != nil {
			return err
		}
		enc.Float64Key(name, val)

	case simdjson.TypeBool:
		val, err := iter.Bool()
		if err != nil {
			return err
		}
		enc.BoolKey(name, val)

	case simdjson.TypeNull:
		enc.NullKey(name)

	case simdjson.TypeString:
		val, err := iter.StringBytes()
		if err != nil {
			return err
		}
		enc.StringKey(name, string(val))

	case simdjson.TypeObject:
		ds := BorrowDecodeSourceWrapper()
		ds.I = iter

		enc.ObjectKey(name, ds)
		ds.DecRef()
		//parseObjectToBytes(enc, obj, iter)

	case simdjson.TypeArray:
		ds := BorrowDecodeSourceWrapper()
		ds.I = iter

		enc.ArrayKey(name, ds)
		ds.DecRef()

	default:
		// 不需要 encode
		return ErrUnknownType
	}

	return nil
}

func parseObject(container Value, obj *simdjson.Object) error {
	container.SetType(Object)

	var it simdjson.Iter
	for {
		name, t, err := obj.NextElementBytes(&it)
		if err != nil {
			return err
		}
		if t == simdjson.TypeNone {
			break
		}
		v := NewNull()
		if err := getValueWithSchema(v, it, nil); err != nil {
			v.DecRef()
			return err
		}
		container.SetObjectKV(name, v)
	}
	return nil
}

func parseObjectWithSchema(container Value, obj *simdjson.Object, schema *DecodeSchema) error {
	container.SetType(Object)

	var it simdjson.Iter
	var currentSchema *DecodeSchema

	enc := gojay.BorrowEncoder(nil)
	defer enc.Release()
	enc.AppendByte('{')

	for {
		name, t, err := obj.NextElementBytes(&it)
		if err != nil {
			return err
		}
		if t == simdjson.TypeNone {
			break
		}

		if schema.NeedsSaveToSource(bytesToString(name)) {
			err := parseToSource(enc, bytesToString(name), &it)
			if err == nil {
				continue
			}
		} else if schema.NeedsSkip(bytesToString(name)) {
			continue
		}

		currentSchema = schema.Property[bytesToString(name)]
		v := NewNull()
		if err := getValueWithSchema(v, it, currentSchema); err != nil {
			v.DecRef()
			return err
		}
		container.SetObjectKV(name, v)
	}

	var b []byte
	b = enc.Buf()
	container.AppendSource(b[1:])

	return nil
}

func parseArray(container Value, arr *simdjson.Array) error {
	container.SetType(Array)

	it := arr.Iter()
	var elem simdjson.Iter
	for i := 0; ; i++ {
		t, err := it.AdvanceIter(&elem)
		if err != nil {
			return err
		}
		if t == simdjson.TypeNone {
			break
		}

		v := NewNull()
		if err := getValueWithSchema(v, elem, nil); err != nil {
			v.DecRef()
			return err
		}
		container.AddArrayElem(v)
	}

	return nil
}

func parseArrayWithSchema(container Value, arr *simdjson.Array, schema *DecodeSchema) error {
	container.SetType(Array)

	enc := gojay.BorrowEncoder(nil)
	defer enc.Release()
	enc.AppendByte('[')

	if schema == nil || schema.ElementSchema == nil {
		// 写入 source
		ds := BorrowDecodeSourceWrapper()
		ds.A = arr

		enc.Array(ds)
		ds.DecRef()
	} else {
		it := arr.Iter()
		var elem simdjson.Iter

		for i := 0; ; i++ {
			t, err := it.AdvanceIter(&elem)
			if err != nil {
				return err
			}
			if t == simdjson.TypeNone {
				break
			}

			v := NewNull()
			if err := getValueWithSchema(v, elem, schema.ElementSchema); err != nil {
				v.DecRef()
				return err
			}
			container.AddArrayElem(v)
		}
	}

	var b []byte
	b = enc.Buf()
	container.AppendSource(b[1:])

	return nil
}

var parsedJsonPool = sync.Pool{
	New: func() interface{} {
		return &simdjson.ParsedJson{}
	},
}

func SimdJsonDecoder(v Value, str []byte) error {
	return SchemaJsonDecoder(v, str, nil)
}

func SchemaJsonDecoder(v Value, str []byte, schema *DecodeSchema) error {
	if !simdjson.SupportedCPU() {
		return ErrUnsupportedCPU
	}

	pj := parsedJsonPool.Get().(*simdjson.ParsedJson)
	parsed, err := simdjson.Parse(str, pj)
	if err != nil {
		return err
	}
	defer func() {
		// 不是 Put(pj), 总之 simdjson.Parse() 就是这么奇葩
		parsed.Reset()
		parsedJsonPool.Put(parsed)
	}()

	it := parsed.Iter()
	it.Advance()

	var iter simdjson.Iter
	if _, _, err := it.Root(&iter); err != nil {
		return err
	}
	if err := getValueWithSchema(v, iter, schema); err != nil {
		return err
	}
	return nil
}
