package pool

import (
    "io"
)

var BS *ByteSlicePool

type sliceHeader struct {
    offset, len int
}

// ByteSlice 用于存储一维数组，可以区分 空数组和nil
type ByteSlice struct {
    RefCount

    data []byte
    elems []sliceHeader
}

func (b *ByteSlice)Grow(dataCap, elemCap int) {
    dataLen := len(b.data)
    dataExtend := dataCap - cap(b.data)
    if dataExtend > 0 {
        b.data = append(b.data, make([]byte, dataExtend)...)
        b.data = b.data[:dataLen]
    }

    elemLen := len(b.elems)
    elemExtend := elemCap - cap(b.elems)
    if elemExtend > 0 {
        b.elems = append(b.elems, make([]sliceHeader, elemExtend)...)
        b.elems = b.elems[:elemLen]
    }
}

func (b *ByteSlice)Len() int {
    return len(b.elems)
}

func (b *ByteSlice)Index(index int) []byte {
    start := b.elems[index].offset
    if b.data[start] == '-' {
        return nil
    }

    end := start + b.elems[index].len
    return b.data[start+1:end]
}

func (b *ByteSlice)IsNil(index int) bool {
    return b.Index(index) == nil
}

func (b *ByteSlice)CopyTo(index int, buf []byte) []byte {
    if b.IsNil(index) {
        return nil
    }
    start := b.elems[index].offset
    end := start + b.elems[index].len
    buf = append(buf, b.data[start+1:end]...)
    return buf
}

func (b *ByteSlice)Reset() {
    if len(b.data) > 0 {
        b.data  = b.data[:0]
    }

    if len(b.elems) > 0 {
        b.elems = b.elems[:0]
    }
}


// AppendConcat 将所有参数拼接成1个数据块
func (b *ByteSlice)AppendConcat(bs ...[]byte) *ByteSlice{
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
    b.elems = append(b.elems, sliceHeader{ used, length+1 })
    return b
}

func (b *ByteSlice)AppendFromReaderN(rd io.Reader, expect int) (n int, err error) {
    used := len(b.data)
    b.data = append(b.data, make([]byte, expect+1)...)

    n, err = io.ReadFull(rd, b.data[used+1:])
    if err != nil {
        b.data = b.data[:used]
        return n, err
    }

    b.data[used] = '+'
    b.elems = append(b.elems, sliceHeader{ used, expect+1 })
    return n, err
}

func (b *ByteSlice)Append(bs ...[]byte) *ByteSlice {
    for _, bts := range bs {
        used := len(b.data)
        if bts == nil {
            b.data = append(b.data, '-')
        } else {
            b.data = append(b.data, '+')
        }
        b.data  = append(b.data, bts...)
        b.elems = append(b.elems, sliceHeader{ used, len(bts)+1 })
    }
    return b
}

func (b *ByteSlice)ToBytes(bs [][]byte) [][]byte {
    for i, h := range b.elems {
        if b.IsNil(i) {
            bs = append(bs, nil)
        } else {
            bs = append(bs, b.data[h.offset+1:h.offset+h.len])
        }
    }
    return bs
}

type ByteSlicePool struct {
    pool *Pool
}

func (p *ByteSlicePool)Get() *ByteSlice {
    bs := p.pool.Acquire().(*ByteSlice)
    return bs
}
