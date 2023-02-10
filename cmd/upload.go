package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/fioncat/clipee/config"
	"github.com/fioncat/clipee/pkg/serial"
	"github.com/spf13/cobra"
)

var Upload = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to other machines",

	Args: cobra.ExactArgs(1),

	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]

		file, err := os.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		base := filepath.Base(name)
		mode := stat.Mode()

		var metaBuffer bytes.Buffer
		err = binary.Write(&metaBuffer, binary.LittleEndian, mode)
		if err != nil {
			return err
		}
		metaBuffer.WriteString(base)

		for _, remote := range config.Get().Remotes {
			conn, err := net.Dial("tcp", remote)
			if err != nil {
				return fmt.Errorf("failed to dial %q: %v", remote, err)
			}
			err = serial.Write(conn, &serial.Packet{
				Handler: "file",
				Meta:    metaBuffer.Bytes(),
				Data:    data,
			})
			if err != nil {
				return fmt.Errorf("failed to send data to %q: %v", remote, err)
			}
			fmt.Printf("Send file %s to %s successed\n", base, remote)
		}

		return nil
	},
}
