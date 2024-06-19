package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/pscott31/mynode/config"
	"github.com/pscott31/mynode/proto"
)

func main() {
	config := config.Default()

	addrPort, err := netip.ParseAddrPort(config.RemoteAddr)
	if err != nil {
		log.Fatalln("error parsing remote address: ", err.Error())
	}

	conn, err := net.Dial("tcp", addrPort.String())
	if err != nil {
		fmt.Println("error dialing {}", err.Error())
		return
	}
	defer conn.Close()

	payload, err := proto.NewVersion(config.Version, config.Services, time.Now().Unix(), addrPort)
	if err != nil {
		log.Fatalf("error creating version message payload: %v ", err.Error())
	}

	payloadBytes, err := proto.MarshalToBytes(payload)
	if err != nil {
		log.Fatalln("error marshalling version message payload: ", err.Error())
	}

	message := proto.NewMessage(config.Magic, proto.MSG_VERSION, payloadBytes)
	messageBytes, err := proto.MarshalToBytes(message)
	if err != nil {
		fmt.Println("Error marshalling message: ", err.Error())
	}

	fmt.Println(hex.EncodeToString(messageBytes))

	if n, err := conn.Write(messageBytes); err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	} else {
		fmt.Println("Wrote", n, "bytes")
	}

	err = os.WriteFile("out.bin", messageBytes, 0o644)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	// TODO!
	// time.Sleep(1 * time.Second)

	var reply proto.Message
	if err = reply.UnmarshalFromReader(conn); err != nil {
		log.Fatalf("error unmarshalling version response message: %v", err)
	}

	fmt.Printf("reply: %+v\n", reply)

	if reply.Magic != config.Magic {
		log.Fatalf("mismatched node networks - sent %x, got %x", message.Magic, reply.Magic)
	}

	if reply.Command != proto.MSG_VERSION {
		log.Fatalf("expected 'version' message in respose, got %s", reply.Command)
	}

	var replyVersion proto.Version
	if err = replyVersion.UnmarshalFromReader(bytes.NewBuffer(reply.Payload)); err != nil {
		log.Fatalf("error unmarshalling version response payload: %v", err)
	}
	fmt.Printf("reply version: %+v\n", replyVersion)
	// dudger := make([]byte, 24)
	// n, err := conn.Read(dudger)
	// if err != nil {
	// 	fmt.Println("Error reading from conn: ", err.Error())
	// }

	// fmt.Println("Read", n, "bytes")
}
