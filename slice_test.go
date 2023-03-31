package cinqgo

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	count := NewSlice(10)
	assertEqual(t, count.Get(), 10)

	err := count.Set(15)
	assertEqual(t, count.Get(), 15)
	assertEqual(t, err, nil)

	err = count.SetFunc(func(old_value int) (int, error) {
		return old_value + 1, nil
	})
	assertEqual(t, count.Get(), 16)
	assertEqual(t, err, nil)

	channel := make(chan int)

	go count.SetChan(channel)

	channel <- 20

	assertEqual(t, count.Get(), 20)

	channelSetter := make(chan SetterFunc[int])

	go count.SetChanFunc(channelSetter)

	channelSetter <- func(old_value int) (int, error) {
		return old_value + 1, nil
	}

	assertEqual(t, count.Get(), 21)
}

func TestReadonlySlice(t *testing.T) {
	immutable := NewReadonlySlice("you cannot change me")

	err := immutable.Set("hi")
	assertEqual(t, immutable.Get(), "you cannot change me")
	assertEqual(t, err.Error(), "cannot set a readonly slice")

	err = immutable.SetFunc(func(old_value string) (string, error) {
		return "hello", nil
	})
	assertEqual(t, immutable.Get(), "you cannot change me")
	assertEqual(t, err.Error(), "cannot set a readonly slice")
}

func TestDerivedSlice(t *testing.T) {
	count := NewSlice(10)
	label, kill := Derive(count, func(value int) string {
		return fmt.Sprintf("%d apples", value)
	})

	defer kill()
	assertEqual(t, label.Get(), "10 apples")

	err := count.Set(15)
	assertEqual(t, label.Get(), "15 apples")
	assertEqual(t, err, nil)

	err = count.SetFunc(func(old_value int) (int, error) {
		return old_value + 1, nil
	})
	assertEqual(t, label.Get(), "16 apples")
	assertEqual(t, err, nil)

	channel := make(chan int)

	go count.SetChan(channel)

	channel <- 20

	assertEqual(t, label.Get(), "20 apples")
}
