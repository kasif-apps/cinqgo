package cinqgo

import (
	"os"
)

type IOTransactor[T any] struct {
	Slice          *Slice[T]
	Source         string
	Encode         func(value T) ([]byte, error)
	Decode         func(record []byte) (T, error)
	EncodeParadigm func(value []byte) ([]byte, error)
	DecodeParadigm func(value []byte) ([]byte, error)
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

	t.Slice.Assign(value)
	return nil
}

func NewIOTransactor[T any](slice *Slice[T], source string) IOTransactor[T] {
	return IOTransactor[T]{
		Slice:  slice,
		Source: source,
	}
}
