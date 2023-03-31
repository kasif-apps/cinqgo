package cinqgo

type Transactor[T any] struct {
	Slice          *Slice[T]
	Encode         func(value T) ([]byte, error)
	Decode         func(record []byte) (T, error)
	EncodeParadigm func(value []byte) ([]byte, error)
	DecodeParadigm func(value []byte) ([]byte, error)
}
