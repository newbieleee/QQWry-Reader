package ip

import (
    `os`
)

type fileMode struct {
    filename string
    data     *os.File
}

func (f *fileMode) Open() error {
    if f.data == nil {
        file, err := os.OpenFile(f.filename, os.O_RDONLY, 0644)
        if err != nil {
            return err
        }
        f.data = file
    }
    _, err := f.data.Seek(0, 0)
    return err
}

func (f *fileMode) Read(p []byte) (n int, err error) {
    return f.data.Read(p)
}

func (f *fileMode) ReadAt(p []byte, off int64) (n int, err error) {
    return f.data.ReadAt(p, off)
}

func (f *fileMode) Seek(offset int64, whence int) (int64, error) {
    return f.data.Seek(offset, whence)
}

func (f *fileMode) Close() error {
    err := f.data.Close()
    f.data = nil
    return err
}
