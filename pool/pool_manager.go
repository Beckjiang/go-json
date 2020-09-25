package pool

import (
    "sync"
    "sync/atomic"
)

type Pool struct {
    p sync.Pool
    ctor func() RefCountable

    Name string
    NewCount, GetCount, PutCount uint64
}

func (p *Pool)newFunc() interface{} {
    atomic.AddUint64(&p.NewCount, 1)
    o := p.ctor()
    o.setPool(p)
    o.setInstance(o)
    return o
}

func (p *Pool)Acquire() RefCountable {
    atomic.AddUint64(&p.GetCount, 1)
    o := p.p.Get().(RefCountable)
    o.IncRef()
    return o
}

func (p *Pool)put(o RefCountable) {
    atomic.AddUint64(&p.PutCount, 1)
    p.p.Put(o)
}

type PoolManager struct {
    lock sync.Mutex
    pools []*Pool
}

func (pm *PoolManager)NewPool(name string, ctor func() RefCountable, initNum int) *Pool {
    pm.lock.Lock()
    defer pm.lock.Unlock()

    p := &Pool{
        ctor : ctor,
        Name : name,
    }
    p.p = sync.Pool{
        New: p.newFunc,
    }

    for i:=0; i < initNum; i++ {
        p.p.Put(p.newFunc())
    }

    pm.pools = append(pm.pools, p)
    return p
}

func (pm *PoolManager)Each(fn func(p *Pool)) {
    if fn == nil {
        return
    }

    pm.lock.Lock()
    defer pm.lock.Unlock()

    for _, p := range pm.pools {
        fn(p)
    }
}

var PM PoolManager

