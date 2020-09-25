package pool

import (
    "strings"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestByteSlice(t *testing.T) {
    var bs ByteSlice
    {
        bs.Grow(512, 10)
        assert.Equal(t, 0  , bs.Len())
        assert.Equal(t, 0  , len(bs.data))
        assert.Equal(t, 512, cap(bs.data))

        assert.Equal(t, 0 , len(bs.elems))
        assert.Equal(t, 10, cap(bs.elems))

        bs.Append([]byte("a"), []byte("b"), []byte("c"))
        bs.AppendConcat([]byte("d"), []byte("e"))
        n, err := bs.AppendFromReaderN(strings.NewReader("f"), 1)
        assert.Equal(t, 1, n)
        assert.Equal(t, nil, err)

        assert.Equal(t, 5  , bs.Len())
        assert.Equal(t, 11  , len(bs.data))
        assert.Equal(t, 512, cap(bs.data))

        assert.Equal(t, 5 , len(bs.elems))
        assert.Equal(t, 10, cap(bs.elems))

        {
            // growing
            bs.Grow(1024, 20)
            assert.Equal(t, 5  , bs.Len())
            assert.Equal(t, 11  , len(bs.data))
            assert.Equal(t, 1024, cap(bs.data))

            assert.Equal(t, 5 , len(bs.elems))
            assert.Equal(t, 20, cap(bs.elems))
        }

        var bs2 [50][]byte
        s := bs.ToBytes(bs2[:0])
        assert.Equal(t, 5, len(s))
        assert.Equal(t, []byte("a") , s[0])
        assert.Equal(t, []byte("b") , s[1])
        assert.Equal(t, []byte("c") , s[2])
        assert.Equal(t, []byte("de"), s[3])
        assert.Equal(t, []byte("f") , s[4])
    }

    bs.Reset()
    assert.Equal(t, 0   , bs.Len())
    assert.Equal(t, 0   , len(bs.data))
    assert.Equal(t, 1024, cap(bs.data))

    assert.Equal(t, 0  , len(bs.elems))
    assert.Equal(t, 20 , cap(bs.elems))
}

func TestByteSliceEmptyAndNil(t *testing.T) {
    var bs ByteSlice

    bs.Append([]byte("0"))
    bs.Append(nil)

    bs.AppendConcat([]byte("2"), []byte("3"))
    bs.AppendConcat(nil, nil)

    n, err := bs.AppendFromReaderN(strings.NewReader("5"), 1)
    assert.Equal(t, 1, n)
    assert.Equal(t, nil, err)

    assert.Equal(t, 5, bs.Len())

    assert.Equal(t, []byte("0"), bs.Index(0))
    assert.Equal(t, false      , bs.IsNil(0))

    assert.Equal(t, []byte(nil), bs.Index(1))
    assert.Equal(t, true       , bs.IsNil(1))

    assert.Equal(t, []byte("23"), bs.Index(2))
    assert.Equal(t, false       , bs.IsNil(2))

    assert.Equal(t, []byte(nil), bs.Index(3))
    assert.Equal(t, true       , bs.IsNil(3))

    assert.Equal(t, []byte("5"), bs.Index(4))
    assert.Equal(t, false      , bs.IsNil(4))
}
