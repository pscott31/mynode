package proto

import (
	"encoding/binary"
	"fmt"
	"io"
)

type VarInt uint64

// Unsigned integers are encoded in the BTC protocol depending on the size to save space.
// < 0xFD         1 byte   uint8
// <= 0xFFFF      3 bytes  0xFD followed by the length as uint16
// <= 0xFFFFFFFF  5	bytes  0xFE followed by the length as uint32
// > 0xFFFFFFFF   9 bytes9 0xFF followed by the length as uint64
func (vi VarInt) MarshalToWriter(w io.Writer) error {
	var err error
	switch {
	case vi < 0xFD:
		return binary.Write(w, binary.LittleEndian, uint8(vi))
	case vi <= 0xFFFF:
		if err := binary.Write(w, binary.LittleEndian, uint8(0xFD)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(vi))
	case vi <= 0xFFFFFFFF:
		if err := binary.Write(w, binary.LittleEndian, uint8(0xFE)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(vi))
	default:
		if err := binary.Write(w, binary.LittleEndian, uint8(0xFF)); err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, uint64(vi))
	}
	return err
}

func (vi *VarInt) UnmarshalFromReader(r io.Reader) error {
	var firstByte uint8
	if err := binary.Read(r, binary.LittleEndian, &firstByte); err != nil {
		return fmt.Errorf("unable to read the first byte: %w", err)
	}

	switch firstByte {
	case 0xFD:
		var value uint16
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return fmt.Errorf("unable to read uint16 value: %w", err)
		}
		*vi = VarInt(value)
	case 0xFE:
		var value uint32
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return fmt.Errorf("unable to read uint32 value: %w", err)
		}
		*vi = VarInt(value)
	case 0xFF:
		var value uint64
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return fmt.Errorf("unable to read uint64 value: %w", err)
		}
		*vi = VarInt(value)
	default:
		*vi = VarInt(firstByte)
	}

	return nil
}
