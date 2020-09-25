package pool

import (
    "bytes"
    "github.com/tinylib/msgp/msgp"
)

type Buffer struct {
    bytes.Buffer
    RefCount
}

var bufferPool *Pool
func init() {
    bufferPool =  PM.NewPool("pool.Buffer{}", func() RefCountable {
        b := &Buffer{}
        b.Grow(512)
        return b
    }, 2048)
}

func NewBuffer() *Buffer {
    return bufferPool.Acquire().(*Buffer)
}

func NewBufferBytes(s []byte) *Buffer {
    buf := bufferPool.Acquire().(*Buffer)
    buf.Write(s)
    return buf
}

func NewBufferString(s string) *Buffer {
    buf := bufferPool.Acquire().(*Buffer)
    buf.WriteString(s)
    return buf
}

func (buf *Buffer)MarshalMsg(b []byte) (o []byte, err error) {
    return msgp.AppendBytes(b, buf.Bytes()), nil
}

func (buf *Buffer)Msgsize() int {
    return buf.Len()
}

func (buf *Buffer)UnmarshalMsg(b []byte) (o []byte, err error) {
    var v []byte
    v, o, err = msgp.ReadBytesZC(b)
    if err != nil {
        return
    }
    buf.Write(v)
    return
}
