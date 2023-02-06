package handler

import (
	"fmt"

	"github.com/fioncat/clipee/pkg/serial"
)

type Handler interface {
	Handle(meta, data []byte, from string) error
}

var handlers = map[string]Handler{}

func Register(name string, h Handler) {
	if _, ok := handlers[name]; ok {
		panic(fmt.Sprintf("handler %q is duplicate", name))
	}
	handlers[name] = h
}

func Do(packet *serial.Packet, from string) error {
	handler := handlers[packet.Handler]
	if handler == nil {
		return fmt.Errorf("packet error: unknown handler %q", packet.Handler)
	}

	err := handler.Handle(packet.Meta, packet.Data, from)
	if err != nil {
		return fmt.Errorf("failed to handle %q: %v", packet.Handler, err)
	}
	return nil
}
