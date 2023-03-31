package cinqgo

import (
	"os"
)

type IOTransactor[T any] struct {
	Transactor[T]
	Slice  *Slice[T]
	Source string
}

func (t IOTransactor[T]) Init() func() {
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

func (t *IOTransactor[T]) Load() error {
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

func NewIOTransactor[T any](slice *Slice[T], source string) IOTransactor[T] {
	return IOTransactor[T]{
		Transactor: Transactor[T]{},
		Slice:      slice,
		Source:     source,
	}
}
