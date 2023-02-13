package serial

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var testCases = []struct {
		handler string
		meta    string
		data    string
	}{
		{
			handler: "handler-0",
			meta:    "meta-0",
			data:    "data-0",
		},
		{
			handler: "clipboard",
			meta:    "0",
			data:    "This is a simple message from our clipboard",
		},
		{
			handler: "file",
			meta:    "hello-world.txt,6678",
			data:    "This is a simple file content\nPlease donot modify it.",
		},
		{
			handler: "text",
			meta:    "",
			data:    "Packet without metadata~",
		},
	}

	for _, testCase := range testCases {
		var buff bytes.Buffer
		srcPacket := &Packet{
			Handler: testCase.handler,
			Meta:    []byte(testCase.meta),
			Data:    []byte(testCase.data),
		}
		if len(srcPacket.Meta) == 0 {
			srcPacket.Meta = nil
		}
		err := Write(&buff, srcPacket)
		if err != nil {
			t.Fatal(err)
		}

		encoded := buff.Bytes()
		var src bytes.Buffer
		src.Write(encoded)

		decoded, err := Read(&src)
		if err != nil {
			t.Fatal(err)
		}
		if len(decoded.Meta) == 0 {
			decoded.Meta = nil
		}

		if !reflect.DeepEqual(decoded, srcPacket) {
			t.Fatalf("decoded packet is not equal: %+v vs %+v", decoded, srcPacket)
		}
	}
}
