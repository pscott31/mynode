package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/netip"
)

var IPV4_IPV6_PREFIX = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}

type NetAddress struct {
	Time     uint32
	Services uint64
	IP       netip.AddrPort
}

func (na NetAddress) MarshalToWriter(w io.Writer) error {
	var err error

	if err = binary.Write(w, binary.LittleEndian, na.Time); err != nil {
		return fmt.Errorf("unable to write time: %w", err)
	}

	if err = binary.Write(w, binary.LittleEndian, na.Services); err != nil {
		return fmt.Errorf("unable to write services: %w", err)
	}

	ipBytes, err := na.IP.MarshalBinary()
	if err != nil {
		return err
	}

	for i := 0; i < 16-len(ipBytes); i++ {
		if err = binary.Write(w, binary.LittleEndian, byte(0)); err != nil {
			return fmt.Errorf("unable to write address padding: %w", err)
		}
	}

	if err = binary.Write(w, binary.BigEndian, ipBytes); err != nil {
		return fmt.Errorf("unable to write address: %w", err)
	}

	if err := binary.Write(w, binary.BigEndian, na.IP.Port()); err != nil {
		return fmt.Errorf("unable to write port: %w", err)
	}

	return nil
}

func (na *NetAddress) UnmarshalFromReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &na.Time); err != nil {
		return fmt.Errorf("unable to read time: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &na.Services); err != nil {
		return fmt.Errorf("unable to read services: %w", err)
	}

	// Prepare a buffer to read the IP address (16 bytes for IPv6)
	ipBytes := make([]byte, 16)
	if _, err := io.ReadFull(r, ipBytes); err != nil {
		return fmt.Errorf("unable to read IP address: %w", err)
	}

	// Read the port (2 bytes)
	var port uint16
	if err := binary.Read(r, binary.BigEndian, &port); err != nil {
		return fmt.Errorf("unable to read port: %w", err)
	}

	// Convert IP bytes and port into netip.AddrPort

	if bytes.Equal(ipBytes[0:12], IPV4_IPV6_PREFIX) {
		ipBytes = ipBytes[12:]
	}

	addr, ok := netip.AddrFromSlice(ipBytes)
	if !ok {
		return fmt.Errorf("unable to parse IP address: %s", ipBytes)
	}
	na.IP = netip.AddrPortFrom(addr, port)

	return nil
}
