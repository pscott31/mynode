package proto

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/netip"
)

const USER_AGENT = "/pscott31-mynode:0.0.1/"

type Version struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    NetAddress
	AddrFrom    NetAddress // for version > 106
	Nonce       uint64
	UserAgent   VarString
	StartHeight int32
	Relay       bool //  for version > 70001
}

func NewVersion(protocolVersion int32, nodeServices uint64, timestamp int64, remoteAddrPort netip.AddrPort) (Version, error) {
	// Give up on timestamp overflow
	if timestamp < 0 || timestamp > math.MaxUint32 {
		return Version{}, fmt.Errorf("timestamp out of range")
	}

	return Version{
		Version:   protocolVersion,
		Services:  nodeServices,
		Timestamp: timestamp,
		AddrRecv: NetAddress{
			Time:     uint32(timestamp),
			Services: nodeServices,
			IP:       remoteAddrPort,
		},
		AddrFrom:  NetAddress{},
		Nonce:     rand.Uint64(),
		UserAgent: VarString(USER_AGENT),
	}, nil
}

func (vp Version) MarshalToWriter(w io.Writer) error {
	var err error
	if binary.Write(w, binary.LittleEndian, vp.Version) != nil {
		return fmt.Errorf("unable to write version: %w", err)
	}

	if binary.Write(w, binary.LittleEndian, vp.Services) != nil {
		return fmt.Errorf("unable to write services: %w", err)
	}

	if binary.Write(w, binary.LittleEndian, vp.Timestamp) != nil {
		return fmt.Errorf("unable to write timestamp: %w", err)
	}

	if err = vp.AddrRecv.MarshalToWriter(w); err != nil {
		return fmt.Errorf("unable to write recieve address: %w", err)
	}

	if vp.Version < 106 {
		return nil
	}

	if err = vp.AddrFrom.MarshalToWriter(w); err != nil {
		return fmt.Errorf("unable to write from address: %w", err)
	}

	if binary.Write(w, binary.LittleEndian, vp.Nonce) != nil {
		return fmt.Errorf("unable to write nonce: %w", err)
	}

	if err = vp.UserAgent.MarshalToWriter(w); err != nil {
		return fmt.Errorf("unable to write user agent: %w", err)
	}

	if err = binary.Write(w, binary.LittleEndian, vp.StartHeight); err != nil {
		return fmt.Errorf("unable to write start height: %w", err)
	}

	if vp.Version < 70001 {
		return nil
	}

	if err = binary.Write(w, binary.LittleEndian, vp.Relay); err != nil {
		return fmt.Errorf("unable to write relay: %w", err)
	}

	return nil
}

func (vp *Version) UnmarshalFromReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &vp.Version); err != nil {
		return fmt.Errorf("unable to read version: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &vp.Services); err != nil {
		return fmt.Errorf("unable to read services: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &vp.Timestamp); err != nil {
		return fmt.Errorf("unable to read timestamp: %w", err)
	}

	if err := vp.AddrRecv.UnmarshalFromReader(r); err != nil {
		return fmt.Errorf("unable to read receive address: %w", err)
	}

	if vp.Version < 106 {
		return nil
	}

	if err := vp.AddrFrom.UnmarshalFromReader(r); err != nil {
		return fmt.Errorf("unable to read from address: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &vp.Nonce); err != nil {
		return fmt.Errorf("unable to read nonce: %w", err)
	}

	if err := vp.UserAgent.UnmarshalFromReader(r); err != nil {
		return fmt.Errorf("unable to read user agent: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &vp.StartHeight); err != nil {
		return fmt.Errorf("unable to read start height: %w", err)
	}

	if vp.Version < 70001 {
		return nil
	}

	if err := binary.Read(r, binary.LittleEndian, &vp.Relay); err != nil {
		return fmt.Errorf("unable to read start height: %w", err)
	}
	return nil
}
