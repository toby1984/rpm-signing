package common

import (
	"fmt"
	"strings"
	"testing"
)

type bitrange struct {
	start uint32
	len   uint32
}

func (bitrange bitrange) String() string {
	return fmt.Sprintf("bitrange{%d, %d}", bitrange.start, bitrange.len)
}

func assertSetSucceeds(ranges []bitrange, t *testing.T) bool {

	bitmap := NewBitmap()
	for _, r := range ranges {
		fmt.Printf("Testing bitmap %s\n", r)
		if !bitmap.Set(r.start, r.len) {
			t.Fatalf("BitmapImpl Set() should have succeeded on %s", r)
		}
	}
	return true
}

func assertSetFails(ranges []bitrange, t *testing.T) bool {

	bitmap := NewBitmap()
	for idx, r := range ranges {
		fmt.Printf("Testing bitmap %s\n", r)
		var failed bool
		if idx == 0 {
			failed = !bitmap.Set(r.start, r.len)
		} else {
			failed = bitmap.Set(r.start, r.len)
		}
		if failed {
			t.Fatalf("BitmapImpl Set() should have failed on %s", r)
			return false
		}
	}
	return true
}

func TestBitmapImpl_ToString(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(0, 1) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := "1"
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

func TestBitmapImpl_ToString2(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(1, 1) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := "01"
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

func TestBitmapImpl_ToString3(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(1, 2) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := "011"
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

func TestBitmapImpl_ToString4(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(63, 1) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := strings.Repeat("0", 63) + "1"
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

func TestBitmapImpl_ToString5(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(63, 2) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := strings.Repeat("0", 63) + "11"
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

func TestBitmapImpl_ToString6(t *testing.T) {
	bitmap := NewBitmap()
	if !bitmap.Set(32, 64) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	expected := strings.Repeat("0", 32) + strings.Repeat("1", 64)
	actual := bitmap.String()
	if actual != expected {
		t.Fatalf("BitmapImpl String() should have returned \"%s\", got \"%s\"", expected, actual)
	}
}

/*
Write() at offset 0x10a4 with length 16
Write() at offset 0x0 with length 41
Write() at offset 0x29 with length 65
*/
func TestBitmapImpl_SpecialCase2(t *testing.T) {

	bitmap := NewBitmap()
	if !bitmap.Set(0x10a4, 16) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	if !bitmap.Set(0x00, 41) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	if !bitmap.Set(0x29, 65) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
}

func TestBitmapImpl_SpecialCase(t *testing.T) {

	bitmap := NewBitmap()
	if !bitmap.Set(0x2c5, 16) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	if !bitmap.Set(0x00, 2) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
	if !bitmap.Set(0x2, 10) {
		t.Fatal("BitmapImpl Set() should have succeeded")
	}
}

func TestBitmapImpl_Set_Success(t *testing.T) {

	// test non-overlapping ranges
	ranges := [][]bitrange{
		{bitrange{start: 0, len: 65}, bitrange{start: 65, len: 65}},
		{bitrange{start: 0, len: 1}, bitrange{start: 1, len: 1}},
		{bitrange{start: 0, len: 0}},
		{bitrange{start: 0, len: 63}, bitrange{start: 64, len: 1}},
		{bitrange{start: 0, len: 64}, bitrange{start: 64, len: 64}},
		{bitrange{start: 30, len: 64}, bitrange{start: 30 + 64, len: 64}},
	}

	for _, rng := range ranges {
		if !assertSetSucceeds(rng, t) {
			return
		}
	}
}

func TestBitmapImpl_Set_Fail(t *testing.T) {

	// test overlapping ranges
	ranges := [][]bitrange{
		{bitrange{start: 0, len: 3}, bitrange{start: 1, len: 1}},
		{bitrange{start: 0, len: 1}, bitrange{start: 0, len: 1}},
		{bitrange{start: 0, len: 3}, bitrange{start: 0, len: 1}},
		{bitrange{start: 0, len: 3}, bitrange{start: 2, len: 1}},
	}

	for _, rng := range ranges {
		if !assertSetFails(rng, t) {
			return
		}
	}
}
