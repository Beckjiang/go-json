package json

import (
    "github.com/beckjiang/go-json/pool"
)

var valuePool *pool.Pool
func init() {
    valuePool = pool.PM.NewPool("json.value", func() pool.RefCountable {
        return new(valueImpl)
    }, 128)
}

func NewBool(val bool) Value {
    v := valuePool.Acquire().(Value)
    v.SetBool(val)
    return v
}

func NewInt(val int64) Value {
    v := valuePool.Acquire().(Value)
    v.SetInt(val)
    return v
}

func NewFloat(val float64) Value {
    v := valuePool.Acquire().(Value)
    v.SetFloat(val)
    return v
}

func NewString(val []byte) Value {
    v := valuePool.Acquire().(Value)
    v.SetString(val)
    return v
}

func NewObject() Value {
    v := valuePool.Acquire().(Value)
    v.SetType(Object)
    return v
}

func NewArray() Value {
    v := valuePool.Acquire().(Value)
    v.SetType(Array)
    return v
}

func NewNull() Value {
    v := valuePool.Acquire().(Value)
    v.SetType(Null)
    return v
}
