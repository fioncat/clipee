package serial

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Packet struct {
	Handler string
	Meta    []byte
	Data    []byte
}

func Write(dst io.Writer, packet *Packet) error {
	err := writeData(dst, []byte(packet.Handler))
	if err != nil {
		return err
	}

	err = writeData(dst, packet.Meta)
	if err != nil {
		return err
	}

	return writeData(dst, packet.Data)
}

func writeData(dst io.Writer, data []byte) error {
	dataLen := uint32(len(data))
	err := binary.Write(dst, binary.LittleEndian, dataLen)
	if err != nil {
		return err
	}
	if dataLen == 0 {
		return nil
	}
	_, err = dst.Write(data)
	return err
}

func Read(src io.Reader) (*Packet, error) {
	var packet Packet

	data, err := readData(src)
	if err != nil {
		return nil, err
	}
	packet.Handler = string(data)

	data, err = readData(src)
	if err != nil {
		return nil, err
	}
	packet.Meta = data

	data, err = readData(src)
	if err != nil {
		return nil, err
	}
	packet.Data = data

	return &packet, nil
}

func readData(src io.Reader) ([]byte, error) {
	var dataLen uint32
	err := binary.Read(src, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, nil
	}

	data := make([]byte, int(dataLen))
	readLen, err := src.Read(data)
	if err != nil {
		return nil, err
	}

	if readLen != int(dataLen) {
		return nil, fmt.Errorf("packet bad format: expect to read %d data, but read %d", dataLen, readLen)
	}

	return data, nil
}
