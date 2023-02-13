package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fioncat/clipee/pkg/serial"
)

func main() {
	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9981,
	}
	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		errExit(err)
	}
	defer conn.Close()

	encoder := serial.NewEncoder("test-handler", conn)

	var idx int
	for {
		meta := fmt.Sprintf("meta-%d", idx)
		msg := fmt.Sprintf("message-%d", idx)

		err = encoder.Encode([]byte(meta), []byte(msg))
		if err != nil {
			fmt.Printf("failed to send udp data: %v\n", err)
		} else {
			fmt.Printf("send message %q to server\n", msg)
		}
		idx++
		time.Sleep(time.Second * 3)
	}
}

func errExit(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}
