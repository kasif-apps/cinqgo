package cinqgo

type Transactor[T any] struct {
	Encode         func(value T) ([]byte, error)
	Decode         func(record []byte) (T, error)
	EncodeParadigm func(value []byte) ([]byte, error)
	DecodeParadigm func(value []byte) ([]byte, error)
}
