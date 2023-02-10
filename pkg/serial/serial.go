package serial

import (
	"encoding/binary"
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
	// The binary format: {data_len(uint32)}{data(bytes)}
	// First, we need to read the dataLen, then read dataLen of bytes as data to return.
	// The dataLen is formatted as LittleEndian, see writeData() method.
	var dataLen uint32
	err := binary.Read(src, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, nil
	}
	size := int(dataLen)

	data := make([]byte, 0, size)
	remain := size
	for remain > 0 {
		buffer := make([]byte, remain)
		readSize, err := src.Read(buffer)
		if err != nil {
			return nil, err
		}
		// According to the doc of io.Reader, the readSize returned by Read method
		// will be less or equal than the buffer size. But won't be bigger than the
		// buffer's scratch space.
		// So here we slice out the valid data from the buffer. If there is sill data
		// has not been read, go to the next loop to continue.
		data = append(data, buffer[:readSize]...)
		remain -= readSize
	}

	return data, nil
}
