package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/netip"
	"time"
)

const (
	REMOTE_ADDR         = "127.0.0.1:8333"
	MAGIC_MAIN   uint32 = 0xD9B4BEF9
	OUR_VERSION  int32  = 31900 // TODO
	OUR_SERVICES uint64 = 1     // NODE_NETWORK
	START_HEIGHT int32  = 0
)

// TODO: VarString/VarInt types?

type NetAddress struct {
	Time     uint32
	Services uint64
	IP       netip.AddrPort
}

func (na NetAddress) MarshalBTC(w io.Writer) error {
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

	for i := 0; i < 10-len(ipBytes); i++ {
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

type VersionPayload struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    NetAddress
	AddrFrom    NetAddress // for version > 106
	Nonce       uint64
	UserAgent   string // Must go to size, char[] on the wire
	StartHeight int32
	Relay       bool //  for version > 70001
}

func (vp VersionPayload) MarshalBTC(w io.Writer) error {
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

	if err = vp.AddrRecv.MarshalBTC(w); err != nil {
		return fmt.Errorf("unable to write recieve address: %w", err)
	}

	if err = vp.AddrFrom.MarshalBTC(w); err != nil {
		return fmt.Errorf("unable to write from address: %w", err)
	}

	// TODO: Assumes user agent < 0xFD bytes
	if binary.Write(w, binary.LittleEndian, uint8(len(vp.UserAgent))) != nil {
		return fmt.Errorf("unable to write user agent length: %w", err)
	}

	if binary.Write(w, binary.LittleEndian, vp.UserAgent); err != nil {
		return fmt.Errorf("unable to write user agent: %w", err)
	}

	if binary.Write(w, binary.LittleEndian, START_HEIGHT) != nil {
		return fmt.Errorf("unable to write start height: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("Hello, world!")
	conn, err := net.Dial("tcp", REMOTE_ADDR)
	if err != nil {
		fmt.Println("Error dialing {}", err.Error())
		return
	}
	defer conn.Close()

	remoteAddr, err := netip.ParseAddrPort(REMOTE_ADDR)
	if err != nil {
		fmt.Println("Error parsing remote address: ", err.Error())
		return
	}

	payload := VersionPayload{
		Version:   OUR_VERSION,
		Services:  OUR_SERVICES,
		Timestamp: time.Now().Unix(),
		AddrRecv: NetAddress{
			Time:     uint32(time.Now().Unix()),
			Services: OUR_SERVICES,
			IP:       remoteAddr,
		},
		AddrFrom: NetAddress{},
	}
	payloadBuf := new(bytes.Buffer)
	if err = payload.MarshalBTC(payloadBuf); err != nil {
		fmt.Println("Error marshalling payload: ", err.Error())
		return
	}
	payloadBytes := payloadBuf.Bytes()

	buf := new(bytes.Buffer)
	if binary.Write(buf, binary.LittleEndian, MAGIC_MAIN) != nil {
		fmt.Println("Error writing magic")
		return
	}

	if binary.Write(buf, binary.LittleEndian, []byte("version\000\000\000\000\000")) != nil {
		fmt.Println("Error writing command")
		return
	}

	if binary.Write(buf, binary.LittleEndian, uint32(len(payloadBytes))) != nil {
		fmt.Println("Error writing payload length")
		return
	}

	// payloadHash := sha256.New()
	// payloadHash.Write(payloadBytes)
	payloadHash := sha256.Sum256(payloadBytes)
	if binary.Write(buf, binary.LittleEndian, payloadHash[:3]) != nil {
		fmt.Println("Error writing payload checksum")
		return
	}

	if binary.Write(buf, binary.LittleEndian, payloadBytes) != nil {
		fmt.Println("Error writing payload")
		return
	}

	fmt.Println(hex.EncodeToString(buf.Bytes()))

	if n, err := conn.Write(buf.Bytes()); err != nil {
		fmt.Println("Error writing to conn: ", err.Error())
	} else {
		fmt.Println("Wrote", n, "bytes")
	}

	// TODO!
	time.Sleep(1 * time.Second)

	dudger := make([]byte, 24)
	n, err := conn.Read(dudger)
	if err != nil {
		fmt.Println("Error reading from conn: ", err.Error())
	}

	fmt.Println("Read", n, "bytes")
}
