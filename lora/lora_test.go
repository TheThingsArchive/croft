package lora

import (
	"bytes"
	"testing"
)

func TestParseMessage(t *testing.T) {
	buf := []byte{0x1, 0x10, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x16, 0x31, 0x3f, 0xe3}
	c := NewConn(nil)
	msg, err := c.parseMessage(nil, buf, len(buf))
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", msg.Header)
	if msg.Header.ProtocolVersion != 0x1 {
		t.Error("The protocol version is not parsed correctly")
	}
	if msg.Header.Token != 0x1020 {
		t.Error("The token is not parsed correctly")
	}
	if msg.Header.Identifier != PUSH_DATA {
		t.Error("The identifier is not parsed correctly")
	}
	if !bytes.Equal(msg.Payload, []byte{0x16, 0x31, 0x3f, 0xe3}) {
		t.Error("The payload is not parsed correctly")
	}
}
