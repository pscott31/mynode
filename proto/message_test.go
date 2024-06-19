package proto_test

import (
	"bytes"
	"testing"

	"github.com/pscott31/mynode/proto"

	"github.com/stretchr/testify/assert"
)

var (
	EXAMPLE_MESSAGE_COMMAND      = "example"
	EXAMPLE_MESSAGE_PAYLOAD      = "payload"
	EXAMPLE_MESSAGE_PAYLOAD_SIZE = len(EXAMPLE_MESSAGE_PAYLOAD)
	EXAMPLE_MESSAGE              = proto.Message{
		Magic:    42,
		Command:  "example",
		Length:   8,          // Assuming the length of the payload
		Checksum: 0x805a173f, // Example checksum
		Payload:  []byte{0x07, 'p', 'a', 'y', 'l', 'o', 'a', 'd'},
	}
)

var EXAMPLE_MESSAGE_BYTES = []byte{
	0x2a, 0x00, 0x00, 0x00, // Magic (42)
	'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x00, 0x00, 0x00, 0x00, 0x00, // Command ("example" + padding to 12 bytes)
	0x08, 0x00, 0x00, 0x00, // Length of payload, little endian
	0x3f, 0x17, 0x5a, 0x80, // Checksum, little endian
	0x07,                              // Payload Size
	'p', 'a', 'y', 'l', 'o', 'a', 'd', // Payload
}

// This message has a checksum that does not match the payload
var EXAMPLE_INVALID_MESSAGE_BYTES = []byte{
	0x2a, 0x00, 0x00, 0x00, // Magic (42)
	'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x00, 0x00, 0x00, 0x00, 0x00, // Command ("example" + padding to 12 bytes)
	0x08, 0x00, 0x00, 0x00, // Length of payload, little endian
	0xff, 0xff, 0xff, 0xff, // INCORRECT checksum, little endian
	0x07,                              // String Length
	'p', 'a', 'y', 'l', 'o', 'a', 'd', // Payload
}

func TestMessage_New(t *testing.T) {
	newMessage := proto.NewMessage(42, proto.MessageType("example"), proto.VarString("payload"))
	assert.Equal(t, EXAMPLE_MESSAGE, newMessage, "The new message should match the expected message")
}

func TestMessage_Marshal(t *testing.T) {
	marshalledMessage, err := proto.MarshalToBytes(EXAMPLE_MESSAGE)
	assert.NoError(t, err, "Marshaling should not produce an error")
	assert.Equal(t, EXAMPLE_MESSAGE_BYTES, marshalledMessage, "The marshaled bytes should match the expected bytes")
}

func TestMessage_UnMarshal(t *testing.T) {
	var unmarshaledMessage proto.Message

	err := unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(EXAMPLE_MESSAGE_BYTES))
	assert.NoError(t, err)

	assert.Equal(t, EXAMPLE_MESSAGE, unmarshaledMessage, "The unmarshalled message should match the expected message")
}

func TestMessage_UnMarshal_Fail(t *testing.T) {
	var unmarshaledMessage proto.Message

	// Not enough bytes to read magic
	err := unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(make([]byte, 3)))
	assert.ErrorContains(t, err, "magic")

	// Not enough bytes to read command
	err = unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(make([]byte, 5)))
	assert.ErrorContains(t, err, "command")

	// Not enough bytes to read length
	err = unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(make([]byte, 17)))
	assert.ErrorContains(t, err, "payload length")

	// Not enough bytes to read checksum
	err = unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(make([]byte, 21)))
	assert.ErrorContains(t, err, "checksum")

	// Payload doesn't match
	err = unmarshaledMessage.UnmarshalFromReader(bytes.NewBuffer(EXAMPLE_INVALID_MESSAGE_BYTES))
	assert.ErrorContains(t, err, "checksum")
}

func TestMessageMarshalUnmarshal(t *testing.T) {
	originalMessage := proto.NewMessage(42, proto.MessageType("command"), proto.VarString("payload"))

	// Marshal the message to a buffer
	var buf bytes.Buffer
	err := originalMessage.MarshalToWriter(&buf)
	assert.NoError(t, err, "Marshaling should not produce an error")

	// Unmarshal the message from the buffer
	var unmarshaledMessage proto.Message
	err = unmarshaledMessage.UnmarshalFromReader(&buf)
	assert.NoError(t, err, "Unmarshaling should not produce an error")

	// Compare the original and unmarshaled messages
	assert.Equal(t, originalMessage, unmarshaledMessage, "The unmarshaled message should be equal to the original")
}
