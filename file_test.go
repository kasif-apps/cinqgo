package cinqgo

import (
	"os"
	"path"
	"strings"
	"testing"
)

type Data struct {
	Count int `json:"count"`
}

func TestIOWrite(t *testing.T) {
	slice := NewSlice(Data{Count: 10})

	dir, _ := os.Getwd()
	target := path.Join(dir, "data.json")
	transactor := NewFileTransactor(&slice, target)

	defer transactor.Init()()

	slice.Set(Data{Count: 11})

	raw, err := os.ReadFile(target)

	if err != nil {
		t.Fatalf(err.Error())
	}

	comparison := "{\"count\":11}"
	assertEqual(t, string(raw), comparison)

	os.Remove(target)
}

func TestEncoding(t *testing.T) {
	slice := NewSlice(Data{Count: 10})

	assertEqual(t, slice.Get().Count, 10)

	dir, _ := os.Getwd()
	target := path.Join(dir, "data-corrupt.json")

	transactor := NewFileTransactor(&slice, target)

	transactor.EncodeParadigm = func(value []byte) ([]byte, error) {
		v := string(value)
		v += " hi"

		return []byte(v), nil
	}

	transactor.DecodeParadigm = func(value []byte) ([]byte, error) {
		v := string(value)

		return []byte(strings.TrimSuffix(v, " hi")), nil
	}

	defer transactor.Init()()
	transactor.Load()

	slice.Set(Data{Count: 5})

	assertEqual(t, slice.Get().Count, 5)

	os.Remove(target)
}
