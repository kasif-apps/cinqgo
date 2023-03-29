package cinqgo

import (
	"os"
)

type IOTransactor[T any] struct {
	Slice  *Slice[T]
	Source string
}

func (t IOTransactor[T]) Init() func() {
	slice := *t.Slice

	return slice.Subscribe(func(e Event) {
		data, err := Encode(e.Detail)

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

	value, err := Decode[T](raw)

	if err != nil {
		return err
	}

	t.Slice.Assign(value)
	return nil
}
