package proto_test

import (
	"bytes"
	"testing"

	"github.com/pscott31/mynode/proto"
	"github.com/stretchr/testify/assert"
)

func TestVarString_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		varStr proto.VarString
	}{
		{
			name:   "Empty string",
			varStr: proto.VarString(""),
		},
		{
			name:   "Short string",
			varStr: proto.VarString("hello"),
		},
		{
			name:   "Long string",
			varStr: proto.VarString("The quick brown fox jumps over the lazy dog"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt.varStr.MarshalToWriter(buf)
			assert.NoError(t, err)

			var gotVarStr proto.VarString
			err = gotVarStr.UnmarshalFromReader(buf)
			assert.NoError(t, err)
			assert.Equal(t, gotVarStr, tt.varStr)
		})
	}
}
