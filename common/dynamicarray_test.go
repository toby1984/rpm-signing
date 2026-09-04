package common

import (
	"encoding/hex"
	"slices"
	"testing"
)

func TestDynamicArray(t *testing.T) {
	array := NewDynamicArray()

	expected1 := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8}
	if array.Write(expected1, 0) != nil {
		t.Fatal("Write() failed")
	}
	if !slices.Equal(array.GetData(), expected1) {
		t.Fatalf("GetData() failed. Expected %s but got %s", hex.Dump(expected1), hex.Dump(array.GetData()))
	}
	if array.Write(expected1, 0) == nil {

	}
}
