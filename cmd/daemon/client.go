package daemon

import (
	"net"

	"github.com/fioncat/clipee/pkg/clipboard"
	"github.com/fioncat/clipee/pkg/serial"
	"github.com/sirupsen/logrus"
)

func Client(dsts []string) {
	conns := make(map[string]net.Conn, len(dsts))
	packets := make(chan *serial.Packet, 1000)

	go clipboard.Notify(packets)

	var err error
	for packet := range packets {
		for _, dst := range dsts {
			conn := conns[dst]
			if conn == nil {
				conn, err = net.Dial("tcp", dst)
				if err != nil {
					logrus.Errorf("[client] failed to dial %s: %v", dst, err)
					continue
				}
				conns[dst] = conn
				logrus.Infof("[client] dial to %s success", dst)
			}
			err = serial.Write(conn, packet)
			if err != nil {
				logrus.Errorf("[client] failed to send data to %s: %v", dst, err)
				delete(conns, dst)
				continue
			}
			logrus.Infof("[client] handler %s send %v data", packet.Handler, len(packet.Data))
		}
	}
}
