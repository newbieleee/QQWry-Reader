package ip

type Builder interface {
    build(*Data) error
}

type builder func(*Data) error

func (f builder) build(data *Data) error {
    return f(data)
}

func WithFileMode(filename string) Builder {
    return builder(func(data *Data) error {
        data.source = &fileMode{filename: filename}
        return nil
    })
}

func WithMemoryMode(filename string) Builder {
    return builder(func(data *Data) error {
        data.source = &memoryMode{filename: filename}
        return nil
    })
}
