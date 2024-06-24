package proto

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const MAX_PROTOCOL_MESSAGE_LENGTH = 4 * 1024 * 1024

// Each message sent between nodes is wrapped by this header, with the command-specific data
// stored in the payload.
type Message struct {
	Magic    uint32
	Command  MessageType
	Length   uint32
	Checksum uint32
	Payload  []byte
}

// NewMessage wraps a payload up with a message header & calculates the checksum.
func NewMessage[T Marshallable](magic uint32, messageType MessageType, payload T) Message {
	payloadBytes, err := MarshalToBytes(payload)
	if err != nil {
		log.Fatalln("error marshalling version message payload: ", err.Error())
	}

	hash := sha256.Sum256(payloadBytes)
	hash = sha256.Sum256(hash[:])

	return Message{
		Magic:    magic,
		Command:  messageType,
		Length:   uint32(len(payloadBytes)),
		Checksum: binary.LittleEndian.Uint32(hash[:4]),
		Payload:  payloadBytes,
	}
}

func (m Message) MarshalToWriter(w io.Writer) error {
	// Write magic
	if err := binary.Write(w, binary.LittleEndian, m.Magic); err != nil {
		return fmt.Errorf("unable to write magic: %w", err)
	}

	// Write command variable length string
	if err := m.Command.MarshalToWriter(w); err != nil {
		return err
	}

	// Write payload length
	if err := binary.Write(w, binary.LittleEndian, m.Length); err != nil {
		return fmt.Errorf("unable to write payload length: %w", err)
	}

	// Write payload checksum
	if err := binary.Write(w, binary.LittleEndian, m.Checksum); err != nil {
		return fmt.Errorf("unable to write payload checksum: %w", err)
	}

	// Write payload
	if err := binary.Write(w, binary.LittleEndian, m.Payload); err != nil {
		return fmt.Errorf("unable to write payload: %w", err)
	}

	return nil
}

func (m *Message) UnmarshalFromReader(r io.Reader) error {
	// Read and unmarshal the magic value
	if err := binary.Read(r, binary.LittleEndian, &m.Magic); err != nil {
		return fmt.Errorf("unable to read magic: %w", err)
	}

	// Unmarshal the command
	if err := m.Command.UnmarshalFromReader(r); err != nil {
		return err
	}

	// Read and unmarshal the payload length
	if err := binary.Read(r, binary.LittleEndian, &m.Length); err != nil {
		return fmt.Errorf("unable to read payload length: %w", err)
	}

	// Ensure message is a sane size
	if m.Length > MAX_PROTOCOL_MESSAGE_LENGTH {
		return fmt.Errorf("message length %d exceeds maximum protocol message length %d", m.Length, MAX_PROTOCOL_MESSAGE_LENGTH)
	}

	// Read and unmarshal the checksum
	if err := binary.Read(r, binary.LittleEndian, &m.Checksum); err != nil {
		return fmt.Errorf("unable to read payload checksum: %w", err)
	}

	// Read the payload
	m.Payload = make([]byte, m.Length)
	if err := binary.Read(r, binary.LittleEndian, &m.Payload); err != nil {
		return fmt.Errorf("unable to read payload: %w", err)
	}

	// Check the checksum matches
	hash := sha256.Sum256(m.Payload)
	hash = sha256.Sum256(hash[:])
	calculatedChecksum := binary.LittleEndian.Uint32(hash[:4])

	if m.Checksum != calculatedChecksum {
		return fmt.Errorf("payload checksum mismatch: %x (computed) != %x (in message)", calculatedChecksum, m.Checksum)
	}

	return nil
}
