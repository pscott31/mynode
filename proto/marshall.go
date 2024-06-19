package proto

import (
	"bytes"
	"io"
)

type Marshallable interface {
	MarshalToWriter(w io.Writer) error
}

func MarshalToBytes[T Marshallable](t T) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := t.MarshalToWriter(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
