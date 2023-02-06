package daemon

import (
	"fmt"
	"net"
	"time"

	"github.com/fioncat/clipee/pkg/handler"
	"github.com/fioncat/clipee/pkg/serial"
	"github.com/sirupsen/logrus"
)

func Server(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen tcp: %v", err)
	}

	logrus.Infof("[server] begin to listen %s", addr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			logrus.Errorf("[server] failed to establish tcp connection: %v", err)
			time.Sleep(time.Second)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	from := conn.RemoteAddr().String()
	defer conn.Close()

	logrus.Infof("[server] begin to receive data from %s", from)
	for {
		packet, err := serial.Read(conn)
		if err != nil {
			logrus.Errorf("[server] failed to read data from client %s: %v, close connection", from, err)
			return
		}

		logrus.Infof("[server] received %v data from %s", len(packet.Data), from)
		err = handler.Do(packet, from)
		if err != nil {
			logrus.Errorf("[server] failed to handle packet: %v", err)
		}
	}
}
