package transactor

import (
	"encoding/json"
)

func Encode[T any](value T) ([]byte, error) {
	encoded, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func EncodeToString[T any](value T) (string, error) {
	encoded, err := json.Marshal(value)

	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func Decode[T any](record []byte) (T, error) {
	var value T

	err := json.Unmarshal(record, &value)

	if err != nil {
		var temp T
		return temp, err
	}

	return value, nil
}

func DecodeFromString[T any](record string) (T, error) {
	return Decode[T]([]byte(record))
}
