package proto

import (
	"encoding/binary"
	"fmt"
	"io"
)

// A wrapper around a string for marshalling/unmarshalling in the BTC protocol
type VarString string

// Marshalled as var_int for length, followed by the string itself
func (vs VarString) MarshalToWriter(w io.Writer) error {
	var err error

	if err = VarInt(len(vs)).MarshalToWriter(w); err != nil {
		return fmt.Errorf("unable to write var string length: %w", err)
	}

	if err = binary.Write(w, binary.LittleEndian, []uint8(vs)); err != nil {
		return fmt.Errorf("unable to write var string: %w", err)
	}

	return nil
}

func (vs *VarString) UnmarshalFromReader(r io.Reader) error {
	var length VarInt
	if err := length.UnmarshalFromReader(r); err != nil {
		return fmt.Errorf("unable to read var string length: %w", err)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("unable to read var string: %w", err)
	}

	*vs = VarString(buf)
	return nil
}
