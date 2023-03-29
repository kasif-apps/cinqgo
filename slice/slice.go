package slice

import (
	"encoding/json"
	"errors"

	et "github.com/kasif-apps/cinqgo/event_target"
)

type Slice[T any] struct {
	et.EventTarget
	value T
}

type SetterFunc[T any] func(old_value T) (T, error)

type GetterFunc[T any, K any] func(value T) K

func (slice *Slice[T]) Assign(value T) error {
	slice.commit(value)
	return nil
}

func (slice *Slice[T]) Set(setter SetterFunc[T]) error {
	newValue, err := setter(slice.value)

	if err != nil {
		return err
	}

	slice.commit(newValue)
	return nil
}

func (slice *Slice[T]) commit(value T) {
	slice.value = value
	slice.DispatchEvent(et.NewEvent("update", value))
}

func (slice Slice[T]) Get() T {
	return slice.value
}

func (slice Slice[T]) Subscribe(callback_fn func(e et.Event)) func() {
	slice.AddEventListener("update", &callback_fn)

	return func() {
		slice.RemoveEventListener("update", &callback_fn)
	}
}

func Derive[T any, K any](parent Slice[T], getter GetterFunc[T, K]) ReadonlySlice[K] {
	intrinsic_value := getter(parent.value)
	sub_slice := NewReadonlySlice(intrinsic_value)

	parent.Subscribe(func(e et.Event) {
		new_value := getter(e.Detail.(T))
		sub_slice.assign(new_value)
	})

	return sub_slice
}

func NewSlice[T any](value T) Slice[T] {
	return Slice[T]{
		EventTarget: et.NewEventTarget(),
		value:       value,
	}
}

type ReadonlySlice[T any] struct {
	*Slice[T]
}

func (slice *ReadonlySlice[T]) assign(value T) {
	slice.commit(value)
}

func (slice ReadonlySlice[T]) Get() T {
	return slice.Slice.Get()
}

func (slice *ReadonlySlice[T]) Assign(value T) error {
	return errors.New("cannot assign to a readonly slice")
}

func (slice *ReadonlySlice[T]) Set(setter SetterFunc[T]) error {
	return errors.New("cannot set a readonly slice")
}

func NewReadonlySlice[T any](value T) ReadonlySlice[T] {
	return ReadonlySlice[T]{
		Slice: &Slice[T]{
			EventTarget: et.NewEventTarget(),
			value:       value,
		},
	}
}

func Upsert[T any](slice Slice[T], upsert T) T {
	root_raw, _ := json.Marshal(slice.Get())
	root := map[string]interface{}{}
	json.Unmarshal(root_raw, &root)

	update_raw, _ := json.Marshal(upsert)
	update := map[string]interface{}{}
	json.Unmarshal(update_raw, &update)

	for key, value := range update {
		if value != nil && value != "" {
			root[key] = value
		}
	}

	result_raw, _ := json.Marshal(root)
	var result T
	json.Unmarshal(result_raw, &result)

	return result
}