package cinqgo

type NetworkTransactor[T any] struct {
	Transactor[T]
	Query  func() T
	Mutate func(value T)
}

func (t *NetworkTransactor[T]) Init() {
	t.Slice.Subscribe(func(e Event) {
		t.Mutate(e.Detail.(T))
	})
}

func (t *NetworkTransactor[T]) Load() {
	result := t.Query()

	t.Slice.Set(result)
}

func NewNetworkTransactor[T any](slice *Slice[T], query func() T, mutate func(value T)) NetworkTransactor[T] {
	return NetworkTransactor[T]{
		Transactor: Transactor[T]{
			Slice:  slice,
			Encode: Encode[T],
			Decode: Decode[T],
		},
		Query:  query,
		Mutate: mutate,
	}
}
