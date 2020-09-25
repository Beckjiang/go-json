package json

import (
    "github.com/beckjiang/go-json/pool"
)

const objCompactKeySize = 0
type objCompKey [objCompactKeySize+1]byte
func makeCompKey(k []byte) (key objCompKey) {
    key[0] = byte(int8(len(k)))
    copy(key[1:], k)
    for i:=len(k)+1; i<objCompactKeySize+1; i++ {
        key[i] = 0
    }
    return key
}

type valueImpl struct {
    pool.RefCount

    typ ValueType

    intV int64
    floatV float64
    stringV []byte

    objCompKV map[objCompKey]Value
    objKV map[string]Value

    arrayV []Value
    source []byte
}

func (v *valueImpl)Type() ValueType { return v.typ }

// getter
func (v *valueImpl)Bool() bool { return v.intV == 1 }
func (v *valueImpl)Int() int64 { return v.intV }
func (v *valueImpl)Float() float64 { return v.floatV }
func (v *valueImpl)String(dst []byte) []byte{ return append(dst, v.stringV...) }

func (v *valueImpl)ObjectVal(key []byte) (found bool, val Value) {
    return v.objectVal(key, false)
}
func (v *valueImpl)GrabObjectVal(key []byte) (found bool, val Value) {
    return v.objectVal(key, true)
}
func (v *valueImpl)objectVal(key []byte, grab bool) (found bool, val Value) {
    if len(key) <= objCompactKeySize {
        if v.objCompKV == nil {
            return
        }
        k := makeCompKey(key)
        val, found = v.objCompKV[k]
        if found && grab {
            delete(v.objCompKV, k)
        }
    } else {
        if v.objKV == nil {
            return
        }
        val, found = v.objKV[string(key)]
        if found && grab {
            delete(v.objKV, string(key))
        }
    }
    return
}

func (v *valueImpl)ArrayElem(index int) Value { return v.arrayV[index] }
func (v *valueImpl)SetArrayElem(index int, val Value) {  v.arrayV[index] = val }
func (v *valueImpl)ArraySize() int { return len(v.arrayV) }

// setter

func (v *valueImpl)SetBool(val bool) {
    v.typ = Bool
    if val {
        v.intV = 1
    } else {
        v.intV = 0
    }
}

func (v *valueImpl)SetInt(val int64) {
    v.typ = Int
    v.intV = val
}

func (v *valueImpl)SetFloat(val float64) {
    v.typ = Float
    v.floatV = val
}

func (v *valueImpl)SetString(val []byte) {
    v.typ = String
    v.stringV = append(v.stringV, val...)
}

func (v *valueImpl)SetNull() {
    v.typ = Null
}

func (v *valueImpl)SetObjectKV(key []byte, val Value) {
    v.typ = Object

    if len(key) <= objCompactKeySize {
        k := makeCompKey(key)
        if v.objCompKV == nil {
            v.objCompKV = make(map[objCompKey]Value)
        }
        v.objCompKV[k] = val
    } else {
        if v.objKV == nil {
            v.objKV = make(map[string]Value)
        }
        v.objKV[string(key)] = val
    }

    return
}

func (v *valueImpl)AddArrayElem(val Value) {
    v.typ = Array
    v.arrayV = append(v.arrayV, val.(*valueImpl))
}

func (v *valueImpl)SetType(t ValueType) {
    v.typ = t
}

func (v *valueImpl)AppendSource(b []byte) {
    if len(v.source) > 0 {
        v.source = append(v.source, ',')
    }
    v.source = append(v.source, b ...)
}
func (v *valueImpl)HasSource() bool {
    return len(v.source) > 0
}
func (v *valueImpl)GetSource() []byte {
    return v.source
}
func (v *valueImpl)SetSource(s []byte)  {
    v.source = s
}

func (v *valueImpl)DelArrayElem(index int) Value {
    if index < len(v.arrayV)-1 {
        copy(v.arrayV[index:], v.arrayV[index+1:])
    }
    v.arrayV[len(v.arrayV)-1] = nil
    v.arrayV = v.arrayV[:len(v.arrayV)-1]
    return nil
}

func (v *valueImpl)Reset() {
    v.typ = None
    v.intV = 0
    v.floatV = 0
    if v.stringV != nil {
        v.stringV = v.stringV[:0]
    }

    if v.source != nil {
        v.source = v.source[:0]
    }

    for _, val := range v.objCompKV {
       val.DecRef()
       //delete(v.objCompKV, key)
    }
    v.objCompKV = nil

    for key, val := range v.objKV {
        val.DecRef()
        delete(v.objKV, key)
    }

    if len(v.arrayV) > 0 {
        for i, val := range v.arrayV {
            val.DecRef()
            v.arrayV[i] = nil
        }
        v.arrayV = v.arrayV[:0]
    }
}
