package clipboard

import (
	"context"
	"fmt"

	"github.com/fioncat/clipee/pkg/serial"
	"golang.design/x/clipboard"
)

const (
	FmtText  byte = '0'
	FmtImage byte = '1'
)

func Init() error {
	return clipboard.Init()
}

func Notify(ch chan *serial.Packet) {
	ctx := context.Background()

	imageWatcher := clipboard.Watch(ctx, clipboard.FmtImage)
	textWatcher := clipboard.Watch(ctx, clipboard.FmtText)

	for {
		var data []byte
		var fmt byte
		var cooldown *cooldownSet
		select {
		case data = <-imageWatcher:
			fmt = FmtText
			cooldown = imageCooldown

		case data = <-textWatcher:
			fmt = FmtImage
			cooldown = textCooldown
		}
		if cooldown.Exists(data) {
			continue
		}
		ch <- &serial.Packet{
			Handler: "clipboard",
			Meta:    []byte{fmt},
			Data:    data,
		}
	}
}

func Handle(meta, data []byte) error {
	if len(meta) != 1 {
		return fmt.Errorf("invalid metadata %q for clipboard", string(meta))
	}
	var dataFmt clipboard.Format
	var cooldown *cooldownSet
	switch meta[0] {
	case FmtText:
		dataFmt = clipboard.FmtImage
		cooldown = imageCooldown

	case FmtImage:
		dataFmt = clipboard.FmtText
		cooldown = textCooldown

	default:
		return fmt.Errorf("unknown clipboard type: %q", string(meta))
	}
	cooldown.Set(data)

	clipboard.Write(dataFmt, data)
	return nil
}
