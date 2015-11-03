package main

import (
	"testing"
)

func TestFmtDevAddr(t *testing.T) {
	// Arrange
	devAddr := uint32(0xABC)

	// Act
	fmt := fmtDevAddr(devAddr)

	// Assert
	if fmt != "00000ABC" {
		t.Fatalf("The format %s is not correct", fmt)
	}
}
