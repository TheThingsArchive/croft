package lora

import (
	"bytes"
	"testing"
	"time"
)

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

	if data.MHDR != 0x80 {
		t.Fatalf("The MAC header should be 0x80 but is %X", data.MHDR)
	}

	if !bytes.Equal(data.DevAddr, []byte{0x8F, 0x77, 0xBB, 0x07}) {
		t.Fatalf("The node EUI should be 8F77BB07 but is %X", data.DevAddr)
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
