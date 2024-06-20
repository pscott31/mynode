package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/netip"
)

var (
	IPV4_IPV6_PREFIX = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}
	INVALID_IP_ADDR  = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type NetAddress struct {
	Time     uint32
	Services uint64
	IP       netip.AddrPort
}

func (na NetAddress) MarshalToWriter(w io.Writer) error {
	var err error

	// Timestamp is not written in the version message, which is all we're using NetAddress for at the moment
	// if err = binary.Write(w, binary.LittleEndian, na.Time); err != nil {
	// 	return fmt.Errorf("unable to write time: %w", err)
	// }

	// Write the services bitmask
	if err = binary.Write(w, binary.LittleEndian, na.Services); err != nil {
		return fmt.Errorf("unable to write services: %w", err)
	}

	// Marshall the IP address
	// NetAttr.MarshalBinary marshalls to a []byte differing lenths:
	//   0  (for Addr{}),
	//   4  (IPv4)
	//   16 (IPv6) bytes.
	// The bitcoin protocol always wants 16 bytes, with a special prefix for IPv4 addresses
	ipBytes, err := na.IP.Addr().MarshalBinary()
	if err != nil {
		return err
	}

	switch n := len(ipBytes); n {
	case 0:
		err = binary.Write(w, binary.BigEndian, make([]byte, 16))
	case 4:
		err = binary.Write(w, binary.BigEndian, IPV4_IPV6_PREFIX)
		if err != nil {
			break
		}
		err = binary.Write(w, binary.BigEndian, ipBytes)
	case 16:
		err = binary.Write(w, binary.BigEndian, ipBytes)
	}
	if err != nil {
		return fmt.Errorf("unable to write ip address: %w", err)
	}

	if err := binary.Write(w, binary.BigEndian, na.IP.Port()); err != nil {
		return fmt.Errorf("unable to write port: %w", err)
	}

	return nil
}

func (na *NetAddress) UnmarshalFromReader(r io.Reader) error {
	// // Timestamp is not written in the version message, which is all we're using NetAddress for at the moment
	// if err := binary.Read(r, binary.LittleEndian, &na.Time); err != nil {
	// 	return fmt.Errorf("unable to read time: %w", err)
	// }

	if err := binary.Read(r, binary.LittleEndian, &na.Services); err != nil {
		return fmt.Errorf("unable to read services: %w", err)
	}

	// Prepare a buffer to read the IP address (16 bytes)
	ipBytes := make([]byte, 16)
	if _, err := io.ReadFull(r, ipBytes); err != nil {
		return fmt.Errorf("unable to read IP address: %w", err)
	}

	// If the 16 byte buffer is all zeros, marshal to zero value Addr{}
	addr := netip.Addr{}
	if !bytes.Equal(ipBytes, INVALID_IP_ADDR) {
		// If we have the magic 12 byte prefix, marshal to IPv4 Addr{}
		if bytes.Equal(ipBytes[0:12], IPV4_IPV6_PREFIX) {
			ipBytes = ipBytes[12:]
		}

		// Else marshal to IPv6 Addr{}
		var ok bool
		addr, ok = netip.AddrFromSlice(ipBytes)
		if !ok {
			return fmt.Errorf("unable to parse IP address: %s", ipBytes)
		}
	}

	// Read the port (2 bytes)
	var port uint16
	if err := binary.Read(r, binary.BigEndian, &port); err != nil {
		return fmt.Errorf("unable to read port: %w", err)
	}

	na.IP = netip.AddrPortFrom(addr, port)
	return nil
}
