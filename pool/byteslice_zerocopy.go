package pool

import (
)

var BSZC *ByteSliceZCPool

// ByteSliceZC 用于存储一维数组，可以区分 空数组和nil
type ByteSliceZC struct {
    RefCount

    data [][]byte
}

func (b *ByteSliceZC)Grow(c int) {
    dataLen := len(b.data)
    dataExtend := c - cap(b.data)
    if dataExtend > 0 {
        b.data = append(b.data, make([][]byte, dataExtend)...)
        b.data = b.data[:dataLen]
    }
}

func (b *ByteSliceZC)Len() int {
    return len(b.data)
}

func (b *ByteSliceZC)Index(index int) []byte {
    return b.data[index]
}

func (b *ByteSliceZC)IsNil(index int) bool {
    return b.Index(index) == nil
}

func (b *ByteSliceZC)CopyTo(index int, buf []byte) []byte {
    if b.IsNil(index) {
        return nil
    }
    buf = append(buf, b.Index(index)...)
    return buf
}

func (b *ByteSliceZC)Reset() {
    for i := range b.data {
        b.data[i] = nil
    }

    if len(b.data) > 0 {
        b.data  = b.data[:0]
    }
}

func (b *ByteSliceZC)Append(bs ...[]byte) *ByteSliceZC {
    b.data = append(b.data, bs...)
    return b
}

func (b *ByteSliceZC)ToBytes(bs [][]byte) [][]byte {
    return b.data
}

type ByteSliceZCPool struct {
    pool *Pool
}

func (p *ByteSliceZCPool)Get() *ByteSliceZC {
    bs := p.pool.Acquire().(*ByteSliceZC)
    return bs
}
