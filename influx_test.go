package main

import (
	"fmt"
	"testing"
)

func TestFormatJsonBytesForInflux(t *testing.T) {
	input := []byte(`{"stat":{"time":"2015-07-23 10:22:03 GMT","lati":52.37376,"long":4.88663,"alti":-11,"rxnb":0,"rxok":0,"rxfw":0,"ackr":0.0,"dwnb":0,"txnb":0}}`)
	output, err := FormatJsonBytesForInflux(input)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Print(string(output))
}
