package ip

import (
	`os`
)

type Builder interface {
	build(*Data) error
}

type builder func(*Data) error

func (f builder) build(data *Data) error {
	return f(data)
}

func WithFileMode(filename string) Builder {
	return builder(func(data *Data) error {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		data.source = f
		return nil
	})
}

func WithMemoryMode(filename string) Builder {
	return builder(func(data *Data) error {
		f, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		data.source = &memoryMode{data: f}
		return nil
	})
}
