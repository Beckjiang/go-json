package pool

func init() {
    BS = &ByteSlicePool {
        pool : PM.NewPool("pool.ByteSlicePool", func() RefCountable {
                return &ByteSlice{}
            }, 2048),
    }

    BSZC = &ByteSliceZCPool {
        pool : PM.NewPool("pool.ByteSliceZCPool", func() RefCountable {
            return &ByteSliceZC{}
        }, 2048),
    }

    TwoDimBS = &TwoDimByteSlicePool {
        pool : PM.NewPool("pool.TwoDimByteSlicePool", func() RefCountable {
            return &TwoDimByteSlice{}
        }, 2048),
    }
}
