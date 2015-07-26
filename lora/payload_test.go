package lora

import (
	"testing"
)

func TestParsePayload(t *testing.T) {
	input := `{"stat":{"some":"val"},"rxpk":{"some":"other"},"rxpk":{"and":"again"},"rxpk":[{"in":"the"},{"result":"here"}]}`
	res, err := parsePayload(input)
	if err != nil {
		t.Errorf("Error Parsing: %s", err.Error())
	}
	sl := *res
	if sl[0].Key != "stat" {
		t.Error("First key not parsed")
	}
	if sl[0].Value != `{"some":"val"}` {
		t.Error("First value not parsed")
	}
	if sl[1].Key != "rxpk" {
		t.Error("Second key not parsed")
	}
	if sl[1].Value != `{"some":"other"}` {
		t.Error("Second value not parsed")
	}
	if sl[2].Key != "rxpk" {
		t.Error("Third key not parsed")
	}
	if sl[2].Value != `{"and":"again"}` {
		t.Error("Third value not parsed")
	}
	if sl[3].Key != "rxpk" {
		t.Error("Fourth key not parsed")
	}
	if sl[3].Value != `[{"in":"the"},{"result":"here"}]` {
		t.Error("Fourth value not parsed")
	}
}
