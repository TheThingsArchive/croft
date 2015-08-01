package lora

import (
	"bytes"
	"testing"
)

func TestParseMessage(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x1, 0x10, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	buf.WriteString(`{"stat":{"lati":100}}`)

	c := NewConn(nil)
	msg, err := c.parseMessage(nil, buf.Bytes(), buf.Len())
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

	payload := msg.Payload.(PushMessagePayload)
	t.Logf("%#v", payload)
	if payload.Stat == nil || payload.RXPK != nil {
		t.Error("The payload is not parsed correctly")
	}
	if payload.Stat.Lati != 100 {
		t.Error("The tmst field is not parsed correctly")
	}
}
