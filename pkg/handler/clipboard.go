package handler

import "github.com/fioncat/clipee/pkg/clipboard"

func init() {
	Register("clipboard", &Clipboard{})
}

type Clipboard struct {}

func (*Clipboard) Handle(meta, data []byte, from string) error {
	return clipboard.Handle(meta, data)
}
