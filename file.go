package cinqgo

import (
	"os"
)

type FileTransactor[T any] struct {
	Transactor[T]
	Source string
}

func (t FileTransactor[T]) Init() func() {
	slice := *t.Slice

	return slice.Subscribe(func(e Event) {
		data := []byte{}
		var err error

		if t.Encode != nil {
			data, err = t.Encode(e.Detail.(T))
		} else {
			data, err = Encode(e.Detail)
		}

		if t.EncodeParadigm != nil {
			data, err = t.EncodeParadigm(data)
		}

		if err != nil {
			panic(err)
		}

		os.WriteFile(t.Source, data, 0644)
	})
}

func (t *FileTransactor[T]) Load() error {
	raw, err := os.ReadFile(t.Source)

	if err != nil {
		return err
	}

	if t.DecodeParadigm != nil {
		raw, err = t.DecodeParadigm(raw)
	}

	var value T

	if t.Decode != nil {
		value, err = t.Decode(raw)
	} else {
		value, err = Decode[T](raw)
	}

	if err != nil {
		return err
	}

	t.Slice.Set(value)
	return nil
}

func NewFileTransactor[T any](slice *Slice[T], source string) FileTransactor[T] {
	return FileTransactor[T]{
		Transactor: Transactor[T]{
			Slice:  slice,
			Encode: Encode[T],
			Decode: Decode[T],
		},
		Source: source,
	}
}
