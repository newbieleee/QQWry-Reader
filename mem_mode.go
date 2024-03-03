package ip

import (
    `errors`
    `io`
    `os`
)

type memoryMode struct {
    filename string
    data     []byte
    pos      int64
}

func (m *memoryMode) Open() error {
    if len(m.data) == 0 {
        data, err := os.ReadFile(m.filename)
        if err != nil {
            return err
        }
        m.data = data
    } else {
        m.pos = 0
    }
    return nil
}

func (m *memoryMode) Read(p []byte) (n int, err error) {
    n = copy(p, m.data[m.pos:])
    return
}

func (m *memoryMode) ReadAt(p []byte, off int64) (n int, err error) {
    n = copy(p, m.data[off:])
    return
}

var ErrOutOfRange = errors.New("偏移量超出范围")

func (m *memoryMode) Seek(offset int64, whence int) (n int64, err error) {
    var r int64
    switch whence {
    case io.SeekEnd:
        r = int64(len(m.data)-1) + offset

    case io.SeekCurrent:
        r = m.pos + offset
    case io.SeekStart:
        r = offset
    }
    if r < 0 || r > int64(len(m.data)-1) {
        return 0, ErrOutOfRange
    }
    m.pos = r
    return r, nil
}

func (m *memoryMode) Close() error {
    m.data = m.data[:0]
    return nil
}
