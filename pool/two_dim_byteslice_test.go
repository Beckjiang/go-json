package pool

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTwoDimByteSlice(t *testing.T) {
    var bs TwoDimByteSlice

    {
        bs.Grow(512, 10, 10)
        assert.Equal(t, 0  , bs.Dim())
        assert.Equal(t, 0  , len(bs.data))
        assert.Equal(t, 512, cap(bs.data))

        assert.Equal(t, 0 , len(bs.flat))
        assert.LessOrEqual(t, 100, cap(bs.flat))

        assert.Equal(t, 0 , len(bs.dim))
        assert.Equal(t, 10, cap(bs.dim))

        {
            // dim 1
            bs.NewDim()

            bs.Append([]byte("a"))
            bs.Append([]byte("b"))
            bs.Append([]byte("c"))
            bs.AppendConcat([]byte("d"), []byte("e"))

            assert.Equal(t, 1  , bs.Dim())
            assert.Equal(t, 4  , bs.Len(0))
            assert.Equal(t, 9  , len(bs.data))
            assert.Equal(t, 512, cap(bs.data))

            assert.Equal(t, 4 , len(bs.flat))
            assert.LessOrEqual(t, 100, cap(bs.flat))
        }

        {
            // dim 2
            bs.NewDim()

            bs.Append([]byte("j"))
            bs.Append([]byte("k"))

            assert.Equal(t, 2  , bs.Dim())
            assert.Equal(t, 2  , bs.Len(1))
            assert.Equal(t, 13 , len(bs.data))
            assert.Equal(t, 512, cap(bs.data))

            assert.Equal(t, 6 , len(bs.flat))
            assert.LessOrEqual(t, 100, cap(bs.flat))
        }

        {
            // growing
            bs.Grow(1024, 20, 20)
            assert.Equal(t, 2   , bs.Dim())
            assert.Equal(t, 13  , len(bs.data))
            assert.Equal(t, 1024, cap(bs.data))

            assert.Equal(t, 6 , len(bs.flat))
            //assert.LessOrEqual(t, 400, cap(bs.flat)) // go机制限制不预留太多空间

            assert.Equal(t, 2 , len(bs.dim))
            assert.LessOrEqual(t, 20, cap(bs.dim))
        }

        var bs2 [50][]byte
        {
            // dim 1
            s := bs.ToBytes(0, bs2[:0])
            assert.Equal(t, 4, len(s))
            assert.Equal(t, []byte("a") , s[0])
            assert.Equal(t, []byte("b") , s[1])
            assert.Equal(t, []byte("c") , s[2])
            assert.Equal(t, []byte("de"), s[3])
        }
        {
            // dim 2
            s := bs.ToBytes(1, bs2[:0])
            assert.Equal(t, 2, len(s))
            assert.Equal(t, []byte("j") , s[0])
            assert.Equal(t, []byte("k") , s[1])
        }

        assert.Equal(t, []byte("a") , bs.Index(0, 0))
        assert.Equal(t, []byte("b") , bs.Index(0, 1))
        assert.Equal(t, []byte("c") , bs.Index(0, 2))
        assert.Equal(t, []byte("de"), bs.Index(0, 3))

        assert.Equal(t, []byte("j") , bs.Index(1, 0))
        assert.Equal(t, []byte("k") , bs.Index(1, 1))
    }

    bs.Reset()
    assert.Equal(t, 0   , bs.Dim())
    assert.Equal(t, 0   , len(bs.data))
    assert.Equal(t, 1024, cap(bs.data))

    assert.Equal(t, 0  , len(bs.flat))
    //assert.LessOrEqual(t, 400, cap(bs.flat)) // go机制限制不预留太多空间

    assert.Equal(t, 0  , len(bs.dim))
    assert.Equal(t, 20 , cap(bs.dim))
}

func TestTwoDimByteSliceEmptyAndNil(t *testing.T) {
    var bs TwoDimByteSlice

    bs.NewDim()

    bs.Append([]byte("0"))
    bs.Append(nil)

    bs.AppendConcat([]byte("2"), []byte("3"))
    bs.AppendConcat(nil, nil)

    assert.Equal(t, 4, bs.Len(0))

    assert.Equal(t, []byte("0"), bs.Index(0, 0))
    assert.Equal(t, false      , bs.IsNil(0, 0))

    assert.Equal(t, []byte(nil), bs.Index(0, 1))
    assert.Equal(t, true       , bs.IsNil(0, 1))

    assert.Equal(t, []byte("23"), bs.Index(0, 2))
    assert.Equal(t, false       , bs.IsNil(0, 2))

    assert.Equal(t, []byte(nil), bs.Index(0, 3))
    assert.Equal(t, true       , bs.IsNil(0, 3))
}
