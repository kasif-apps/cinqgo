package cinqgo

import "testing"

func assertEqual[T comparable](t *testing.T, a T, b T) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}
