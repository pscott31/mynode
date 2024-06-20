package proto_test

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/pscott31/mynode/proto"
	"github.com/stretchr/testify/assert"
)

func TestNetAddress_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		netAddr proto.NetAddress
	}{
		{
			name: "IPv4 address",
			netAddr: proto.NetAddress{
				Services: 1,
				IP:       netip.MustParseAddrPort("192.0.2.1:80"),
			},
		},
		{
			name: "IPv6 address",
			netAddr: proto.NetAddress{
				Services: 1,
				IP:       netip.MustParseAddrPort("[2001:db8::1]:80"),
			},
		},
		{
			name: "Invalid IP address",
			netAddr: proto.NetAddress{
				Services: 1,
				IP:       netip.AddrPort{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt.netAddr.MarshalToWriter(buf)
			assert.NoError(t, err)

			var gotAddr proto.NetAddress
			err = gotAddr.UnmarshalFromReader(buf)
			if err != nil {
				t.Errorf("UnmarshalFromReader() error = %v", err)
				return
			}

			assert.Equal(t, gotAddr.Services, tt.netAddr.Services)
			assert.Zero(t, gotAddr.IP.Compare(tt.netAddr.IP))
		})
	}
}
