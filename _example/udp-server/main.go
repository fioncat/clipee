package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fioncat/clipee/pkg/serial"
)

func main() {
	listener, err := net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 9981,
		})
	if err != nil {
		errExit(err)
	}

	fmt.Println("server start to listen")
	for {
		packet, err := serial.Read(listener)
		if err != nil {
			fmt.Printf("failed to read udp data: %v\n", err)
			time.Sleep(time.Second * 3)
			continue
		}

		handler := string(packet.Handler)
		meta := string(packet.Meta)
		str := string(packet.Data)

		fmt.Printf("recv packet: {handler: %q, meta: %q, data: %q}\n", handler, meta, str)
	}
}

func errExit(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}
