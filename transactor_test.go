package cinqgo

import (
	"strconv"
	"strings"
	"testing"
)

func TestTransactor(t *testing.T) {
	transactor := Transactor[int]{
		Encode: func(value int) ([]byte, error) {
			return []byte(strconv.Itoa(value)), nil
		},
		Decode: func(record []byte) (int, error) {
			return strconv.Atoi(string(record))
		},
		EncodeParadigm: func(value []byte) ([]byte, error) {
			return []byte(string(value) + " suffix"), nil
		},
		DecodeParadigm: func(value []byte) ([]byte, error) {
			str := string(value)
			return []byte(strings.TrimSuffix(str, " suffix")), nil
		},
	}

	encoded, err := transactor.Encode(10)
	assertEqual(t, string(encoded), "10")
	assertEqual(t, err, nil)

	decoded, err := transactor.Decode([]byte("15"))
	assertEqual(t, decoded, 15)
	assertEqual(t, err, nil)

	decoded, err = transactor.Decode([]byte("hello"))
	assertEqual(t, decoded, 0)
	assertEqual(t, err.Error(), "strconv.Atoi: parsing \"hello\": invalid syntax")

	paradigmEncoded, err := transactor.EncodeParadigm(encoded)
	assertEqual(t, string(paradigmEncoded), "10 suffix")
	assertEqual(t, err, nil)

	paradigmDecoded, err := transactor.DecodeParadigm([]byte(paradigmEncoded))
	assertEqual(t, string(paradigmDecoded), "10")
	assertEqual(t, err, nil)
}
