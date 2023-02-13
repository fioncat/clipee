package daemon

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/fioncat/clipee/config"
	"github.com/fioncat/clipee/pkg/clipboard"
	"github.com/fioncat/clipee/pkg/daemon"
	"github.com/fioncat/clipee/pkg/handler"
	"github.com/fioncat/clipee/pkg/serial"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagNoDaemon bool

func init() {
	Start.PersistentFlags().BoolVarP(&flagNoDaemon, "no-daemon", "n", false, "no daemon")
}

var Start = &cobra.Command{
	Use:   "start [-n]",
	Short: "Start clipee",

	RunE: func(_ *cobra.Command, _ []string) error {
		if !flagNoDaemon {
			d, err := daemon.New()
			if err != nil {
				return fmt.Errorf("failed to init daemon: %v", err)
			}
			return d.Start(func() error {
				return StartAll()
			})
		}

		return StartAll()
	},

	Args: cobra.ExactArgs(0),
}

var Stop = &cobra.Command{
	Use:   "stop",
	Short: "Stop daemon",

	RunE: func(_ *cobra.Command, _ []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		return d.Stop()
	},
}

var Status = &cobra.Command{
	Use:   "status",
	Short: "Show daemon status",

	RunE: func(_ *cobra.Command, _ []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		return d.ShowStatus()
	},
}

var Logs = &cobra.Command{
	Use:   "logs",
	Short: "Show daemon logs",

	DisableFlagParsing: true,

	RunE: func(_ *cobra.Command, args []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		path := d.LogPath()
		args = append(args, path)
		cmd := exec.Command("tail", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	},
}

func StartClient() {
	remotes := config.Get().Remotes
	conns := make(map[string]net.Conn, len(remotes))
	packets := make(chan *serial.Packet, 1000)

	go clipboard.Notify(packets)

	var err error
	for packet := range packets {
		for _, remote := range remotes {
			conn := conns[remote]
			if conn == nil {
				conn, err = net.Dial("tcp", remote)
				if err != nil {
					logrus.Errorf("[client] failed to dial %s: %v", remote, err)
					continue
				}
				conns[remote] = conn
				logrus.Infof("[client] dial to %s success", remote)
			}
			err = serial.Write(conn, packet)
			if err != nil {
				logrus.Errorf("[client] failed to send data to %s: %v", remote, err)
				delete(conns, remote)
				continue
			}
			logrus.Infof("[client] handler %s send %v data", packet.Handler, len(packet.Data))
		}
	}
}

func StartServer() error {
	listen, err := net.Listen("tcp", config.Get().Listen)
	if err != nil {
		return fmt.Errorf("failed to listen tcp: %v", err)
	}

	logrus.Infof("[server] begin to listen %s", config.Get().Listen)
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

func StartAll() error {
	if len(config.Get().Remotes) > 0 {
		go StartClient()
	}
	return StartServer()
}
