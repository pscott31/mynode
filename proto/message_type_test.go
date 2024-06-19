package proto_test

import (
	"testing"

	"github.com/pscott31/mynode/proto"
	"github.com/stretchr/testify/assert"
)

func TestMessageType_MarshalToWriter(t *testing.T) {
	expected := []byte("test\x00\x00\x00\x00\x00\x00\x00\x00")
	actual, err := proto.MarshalToBytes(proto.MessageType("test"))

	assert.Equal(t, actual, expected)
	assert.NoError(t, err)
}

func TestMessageType_MarshalToWriterFails(t *testing.T) {
	expected := []byte(nil)
	actual, err := proto.MarshalToBytes(proto.MessageType("a_message_type_that_is_too_long"))

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
}
