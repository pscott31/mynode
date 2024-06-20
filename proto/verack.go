package proto

import "io"

type VerAck struct{}

func (va VerAck) MarshalToWriter(w io.Writer) error {
	return nil
}

func (va *VerAck) UnmarshalFromReader(r io.Reader) error {
	return nil
}
