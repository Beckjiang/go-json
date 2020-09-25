package pool

import (
)

var TwoDimBS *TwoDimByteSlicePool

type TwoDimByteSlicePool struct {
    pool *Pool
}

// TwoDimByteSlice 二维递增数组
type TwoDimByteSlice struct {
    RefCount

    data []byte
    flat []sliceHeader
    dim  []sliceHeader
}

func (p *TwoDimByteSlicePool)Get() *TwoDimByteSlice {
    bs := p.pool.Acquire().(*TwoDimByteSlice)
    return bs
}

func (b *TwoDimByteSlice)Reset() {
    if len(b.data) > 0 {
        b.data = b.data[:0]
    }

    if len(b.flat) > 0 {
        b.flat = b.flat[:0]
    }

    if len(b.dim) > 0 {
        b.dim  = b.dim[:0]
    }
}

func (b *TwoDimByteSlice)Grow(dataCap, flatCap, dimCap int) {
    dataLen := len(b.data)
    dataExtend := dataCap - cap(b.data)
    if dataExtend > 0 {
        b.data = append(b.data, make([]byte, dataExtend)...)
        b.data = b.data[:dataLen]
    }

    flatLen := len(b.flat)
    flatExtend := flatCap * dimCap - cap(b.flat)
    if flatExtend > 0 {
        b.flat = append(b.flat, make([]sliceHeader, flatExtend)...)
        b.flat = b.flat[:flatLen]
    }

    dimLen := len(b.dim)
    dimExtend := dimCap - cap(b.dim)
    if dimExtend > 0 {
        b.dim = append(b.dim, make([]sliceHeader, dimExtend)...)
        b.dim = b.dim[:dimLen]
    }
}

func (b *TwoDimByteSlice)NewDim() *TwoDimByteSlice {
    b.dim = append(b.dim, sliceHeader{ len(b.flat), 0 })
    return b
}

func (b *TwoDimByteSlice)Dim() int {
    return len(b.dim)
}

func (b *TwoDimByteSlice)Len(dim int) int {
    return b.dim[dim].len
}

// Index 不要修改返回的数组内容
func (b *TwoDimByteSlice)Index(dim, index int) []byte {
    flat := b.dim[dim].offset
    start := b.flat[flat+index].offset
    if b.data[start] == '-' {
        return nil
    }

    end   := start + b.flat[flat+index].len
    return b.data[start+1:end]
}

func (b *TwoDimByteSlice)IsNil(dim, index int) bool {
    return b.Index(dim, index) == nil
}

func (b *TwoDimByteSlice)CopyTo(dim, index int, buf []byte) []byte {
    bts := b.Index(dim, index)
    if bts == nil {
        return buf
    }
    buf = append(buf, bts...)
    return buf
}

func (b *TwoDimByteSlice)AppendConcat(bs ...[]byte) {
    used := len(b.data)
    b.data = append(b.data, '+')

    allNil := true
    length := 0
    for _, bts := range bs {
        if bts == nil {
            continue
        }
        allNil = false
        length = length + len(bts)
        b.data = append(b.data, bts...)
    }
    if allNil {
        b.data[used] = '-'
    }
    b.flat = append(b.flat, sliceHeader{ used, length+1 })
    b.dim[len(b.dim)-1].len++
}

func (b *TwoDimByteSlice)Append(bs ...[]byte) {
    for _, bts := range bs {
        used := len(b.data)
        if bts == nil {
            b.data = append(b.data, '-')
        } else {
            b.data = append(b.data, '+')
        }
        b.data  = append(b.data, bts...)
        b.flat = append(b.flat, sliceHeader{ used, len(bts)+1 })
        b.dim[len(b.dim)-1].len++
    }
}

func (b *TwoDimByteSlice)ToBytes(dim int, bs [][]byte) [][]byte {
    start := b.dim[dim].offset
    end   := start + b.dim[dim].len
    for i, h := range b.flat[start:end] {
        if b.IsNil(dim, i) {
            bs = append(bs, nil)
        } else {
            bs = append(bs, b.data[h.offset+1:h.offset+h.len])
        }
    }
    return bs
}

func (b *TwoDimByteSlice)ToByteSlice(dim int) (bs *ByteSlice) {
    bs = BS.Get()
    for i := 0; i < b.Len(dim); i++ {
        bs.Append(b.Index(dim, i))
    }
    return bs
}


