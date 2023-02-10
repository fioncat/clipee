package handler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fioncat/clipee/config"
)

func init() {
	Register("file", &File{})
}

type File struct{}

func (*File) Handle(meta, data []byte, from string) error {
	metaBuffer := bytes.NewBuffer(meta)

	var mode uint32
	err := binary.Read(metaBuffer, binary.LittleEndian, &mode)
	if err != nil {
		return fmt.Errorf("failed to read mode from meta: %v", err)
	}

	name := string(metaBuffer.Bytes())
	if name == "" {
		return errors.New("filename is empty")
	}

	dir := config.Get().Share

	stat, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to mkdir for share directory: %v", err)
		}

	case err == nil:
		if !stat.IsDir() {
			return fmt.Errorf("share path %s is not a directory", dir)
		}

	default:
		return fmt.Errorf("failed to read share directory: %v", err)
	}

	path := filepath.Join(dir, name)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	payload := bytes.NewBuffer(data)
	_, err = io.Copy(file, payload)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %v", path, err)
	}

	return nil
}
