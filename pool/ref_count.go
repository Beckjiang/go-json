package pool

import (
    "reflect"
    "sync/atomic"
)

type RefCountable interface {
    Reset()
    IncRef() RefCountable
    DecRef()
    OwnedMe(RefCountable)

    setInstance(i RefCountable)
    setPool(p *Pool)
}

type RefCount struct {
    pool *Pool
    count int32
    instance RefCountable

    owned []RefCountable
}

func (r *RefCount) setInstance(i RefCountable) {
    r.instance = i
}

func (r *RefCount) setPool(p *Pool) {
    r.pool = p
}

func (r *RefCount) OwnedMe(o RefCountable) {
    r.owned = append(r.owned, o)
}

func (r *RefCount) IncRef() RefCountable {
    atomic.AddInt32(&r.count, 1)
    return r.instance
}

func (r *RefCount) DecRef() {
    if atomic.LoadInt32(&r.count) == 0 {
        panic("bug! The object has been released: " + reflect.TypeOf(r.instance).String())
    }
    if atomic.AddInt32(&r.count, -1) == 0 {
        for i, o := range r.owned {
            o.DecRef()
            r.owned[i] = nil
        }
        if r.owned != nil {
            r.owned = r.owned[:0]
        }

        r.instance.Reset()
        r.pool.put(r.instance)
    }
}

