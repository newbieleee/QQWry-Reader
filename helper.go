package ip

import (
    `io`
)

type concurrencyReader struct {
    pos    int64
    reader io.ReaderAt
}

func (r *concurrencyReader) Read(b []byte) (n int, err error) {
    n, err = r.reader.ReadAt(b, r.pos)
    r.pos += int64(n)
    return
}
