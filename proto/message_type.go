package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MessageType string

const (
	MAX_COMMAND_LENGTH int         = 12
	MSG_VERSION        MessageType = "version"
	MSG_VERACK         MessageType = "verack"
)

func (mt MessageType) MarshalToWriter(w io.Writer) error {
	command_length := len(mt)
	if command_length > MAX_COMMAND_LENGTH {
		return fmt.Errorf("command name too long (%v when protocol max is %v)", command_length, MAX_COMMAND_LENGTH)
	}

	if err := binary.Write(w, binary.LittleEndian, []byte(mt)); err != nil {
		return fmt.Errorf("unable to write command: %w", err)
	}

	for i := 0; i < MAX_COMMAND_LENGTH-command_length; i++ {
		if err := binary.Write(w, binary.LittleEndian, byte(0)); err != nil {
			return fmt.Errorf("unable to write address padding: %w", err)
		}
	}

	return nil
}

func (mt *MessageType) UnmarshalFromReader(r io.Reader) error {
	// Create a buffer to hold the command, which has a fixed size of MAX_COMMAND_LENGTH
	buf := make([]byte, MAX_COMMAND_LENGTH)

	if _, err := io.ReadFull(r, buf); err != nil {
		return fmt.Errorf("unable to read command: %w", err)
	}

	trimmedBuf := bytes.Trim(buf, "\x00")

	*mt = MessageType(string(trimmedBuf))

	return nil
}
