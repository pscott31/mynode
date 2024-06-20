package proto_test

import (
	"bytes"
	"testing"

	"github.com/pscott31/mynode/proto"
	"github.com/stretchr/testify/assert"
)

func TestVarInt_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		varInt       proto.VarInt
		expectedSize int
	}{
		{
			name:         "Small number",
			varInt:       proto.VarInt(0xFC),
			expectedSize: 1,
		},
		{
			name:         "Boundary of 0xFD",
			varInt:       proto.VarInt(0xFD),
			expectedSize: 3,
		},
		{
			name:         "At 0xFFFF",
			varInt:       proto.VarInt(0xFFFF),
			expectedSize: 3,
		},
		{
			name:         "After 0xFFFF",
			varInt:       proto.VarInt(0x10000),
			expectedSize: 5,
		},
		{
			name:         "At 0xFFFFFFFF",
			varInt:       proto.VarInt(0xFFFFFFFF),
			expectedSize: 5,
		},
		{
			name:         "Large number",
			varInt:       proto.VarInt(0xFFFFFFFFFF),
			expectedSize: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt.varInt.MarshalToWriter(buf)
			assert.NoError(t, err)
			assert.Equal(t, buf.Len(), tt.expectedSize)

			var gotVarInt proto.VarInt
			err = gotVarInt.UnmarshalFromReader(buf)
			assert.NoError(t, err)
			assert.Equal(t, gotVarInt, tt.varInt)
		})
	}
}
