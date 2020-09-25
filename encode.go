package json

import (
    "bytes"
    "io"

    "github.com/francoispqt/gojay"
)

func (v *valueImpl)Encode(dst []byte) []byte {
    var buf bytes.Buffer
    v.EncodeTo(&buf)
    dst = append(dst, buf.Bytes()...)
    return dst
}

func (v *valueImpl)EncodeTo(w io.Writer) {
    enc := gojay.BorrowEncoder(w)
    defer enc.Release()

    switch v.typ {
    case Object: enc.EncodeObject(v)
    case Array: enc.EncodeArray(v)
    default:
        panic("bug")
    }
}

func encodeObject(key string, val Value, enc *gojay.Encoder) {
    switch val.Type() {
        case Null   : enc.NullKey(key)
        case Bool   : enc.BoolKey(key, val.Bool())
        case Int    : enc.Int64Key(key, val.Int())
        case Float  : enc.Float64Key(key, val.Float())
        case String : enc.StringKey(key, string(val.String(nil)))

    case Object:
        vv := val.(*valueImpl)
        enc.ObjectKey(key, vv)

    case Array:
        vv := val.(*valueImpl)
        enc.ArrayKey(key, vv)
    }
}

func (v *valueImpl)MarshalJSONObject(enc *gojay.Encoder) {
    if v.typ != Object {
        panic("bug")
    }

    if v.HasSource() {
        enc.AppendBytes(v.GetSource())
    }

    for key, val := range v.objCompKV {
        l := uint8(key[0])
        k := key[1:l+1]
        encodeObject(string(k), val, enc)
    }
    for key, val := range v.objKV {
        encodeObject(key, val, enc)
    }
}

func (v *valueImpl)IsNil() bool {
    return v.typ == Null
}

func (v *valueImpl)MarshalJSONArray(enc *gojay.Encoder) {
    if v.typ != Array {
        panic("bug")
    }

    if v.HasSource() {
        enc.AppendBytes(v.GetSource())
    }
    for _, val := range v.arrayV {
        switch val.Type() {
            case Null   : enc.Null()
            case Bool   : enc.Bool(val.Bool())
            case Int    : enc.Int64(val.Int())
            case Float  : enc.Float64(val.Float())
            case String : enc.String(string(val.String(nil)))

        case Object:
            vv := val.(*valueImpl)
            enc.Object(vv)

        case Array:
            vv := val.(*valueImpl)
            enc.Array(vv)
        }
    }
}

