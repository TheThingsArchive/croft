package lora

import (
	"bytes"
	"testing"
	"time"
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

func TestParseMessageWithGarbage(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x1, 0x10, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0})

	c := NewConn(nil)
	_, err := c.parseMessage(nil, buf.Bytes(), buf.Len())
	if err == nil {
		t.Error("Parse message did not validate input length")
	}
}

func TestParseMessageWithInvalidPayload(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x1, 0x10, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	buf.WriteString("invalid payload")

	c := NewConn(nil)
	_, err := c.parseMessage(nil, buf.Bytes(), buf.Len())
	if err == nil {
		t.Error("Parse message did not validate payload")
	}
}

func TestConvertRXPK(t *testing.T) {
	// Arrange
	key := []byte{0x2B, 0x7E, 0x15, 0x16, 0x28, 0xAE, 0xD2, 0xA6, 0xAB, 0xF7, 0x15, 0x88, 0x09, 0xCF, 0x4F, 0x3C}
	rxpk := &RXPK{
		Time: time.Now(),
		Data: "gI93uwcAAgAGvTNCoZ/MPI1ry1/bBUjbTchQFK7r/gtUscmY3vU+l5twHauwRTAO+GmcOPwaNNU=",
	}

	// Act
	data, err := rxpk.ParseData()
	if err != nil {
		t.Fatalf("Failed to parse data: %s", err.Error())
	}

	// Assert
	if data.MHDR != 0x80 {
		t.Fatalf("The MAC header should be 0x80 but is %X", data.MHDR)
	}

	if data.DevAddr != 0x07BB778F {
		t.Fatalf("The node EUI should be 0x07BB778F but is %X", data.DevAddr)
	}

	if data.FCtrl != 0 {
		t.Fatalf("The control should be 0 but is %d", data.FCtrl)
	}

	if len(data.FOpts) != 0 {
		t.Fatalf("The options should be empty but is %#v", data.FOpts)
	}

	if data.FCnt != 2 {
		t.Fatalf("The counter should be 2 but is %d", data.FCnt)
	}

	if data.FPort != 6 {
		t.Fatalf("The port should be 6 but is %d", data.FPort)
	}

	ok, err := data.TestIntegrity(key)
	if err != nil {
		t.Fatalf("Failed to test integrity: %s", err.Error())
	}
	if !ok {
		t.Fatal("Integrity test failed")
	}

	payload, err := data.DecryptPayload(key)
	if err != nil {
		t.Fatalf("Error decrypting data: %s", err.Error())
	}

	if string(payload) != `{"name":"Turiphro","count":13,"water":true}` {
		t.Fatal("The decrypted data does not match the expected result")
	}
}
