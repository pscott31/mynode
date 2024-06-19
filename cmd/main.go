package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/pscott31/mynode/config"
	"github.com/pscott31/mynode/proto"
)

func main() {
	// TODO: Seed the random number generator
	// rand.Seed(time.Now().UnixNano())

	config := config.Default()

	addrPort, err := netip.ParseAddrPort(config.RemoteAddr)
	if err != nil {
		log.Fatalln("error parsing remote address: ", err.Error())
	}

	// Connect to remote node
	conn, err := net.Dial("tcp", addrPort.String())
	if err != nil {
		log.Fatalln("error dialing: ", err.Error())
	}
	defer conn.Close()

	log.Printf("Connected to %s", addrPort.String())

	// Make our version message
	ourVersion, err := proto.NewVersion(config.Version, config.Services, time.Now().Unix(), addrPort)
	if err != nil {
		log.Fatalf("error creating version message payload: %v ", err.Error())
	}

	// Wrap it up with a ourVersionMsg header
	ourVersionMsg := proto.NewMessage(config.Magic, proto.MSG_VERSION, ourVersion)
	ourVersionBytes, err := proto.MarshalToBytes(ourVersionMsg)
	if err != nil {
		fmt.Println("error marshalling message: ", err.Error())
	}

	// Send it down the pipe
	log.Printf("sending our version %+v", ourVersion)
	if _, err := conn.Write(ourVersionBytes); err != nil {
		log.Fatalln("Error writing to connection: ", err.Error())
	}

	// Handy for debugging
	// err = os.WriteFile("out.bin", messageBytes, 0o644)
	// if err != nil {
	// 	log.Fatalf("failed writing to file: %s", err)
	// }

	var theirVersionMsg proto.Message
	if err = theirVersionMsg.UnmarshalFromReader(conn); err != nil {
		log.Fatalf("error unmarshalling version response message: %v", err)
	}

	// First check stuff in the header
	// Magic should match
	if theirVersionMsg.Magic != config.Magic {
		log.Fatalf("mismatched node networks - sent %x, got %x", ourVersionMsg.Magic, theirVersionMsg.Magic)
	}

	// And they should be sending us a 'version' message in response
	if theirVersionMsg.Command != proto.MSG_VERSION {
		log.Fatalf("expected 'version' message in respose, got %s", theirVersionMsg.Command)
	}

	// Unmarshal the message payload into a version message
	var theirVersion proto.Version
	if err = theirVersion.UnmarshalFromReader(bytes.NewBuffer(theirVersionMsg.Payload)); err != nil {
		log.Fatalf("error unmarshalling version response payload: %v", err)
	}
	fmt.Printf("received their version: %+v\n", theirVersion)

	// Check we're not connected to ourselves
	if theirVersion.Nonce == ourVersion.Nonce {
		log.Fatalf("nonce in response matches nonce in request (connected to self?)")
	}

	// Check that the receiving address in the response matches our connection's sending address
	if theirVersion.AddrRecv.IP.String() != conn.LocalAddr().String() {
		log.Fatalf("address in response (%s) does not match address of connected peer (%s)", theirVersion.AddrRecv.IP, conn.LocalAddr())
	}

	// TODO: Send VERACK
}
