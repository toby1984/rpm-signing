package common

import (
	"fmt"
	"strings"
)

type Bitmap interface {
	// Set populates the given number of bits starting
	// at bit offset with a '1'
	//
	// it returns true if none of those bits had
	// already been set before, otherwise false
	Set(bitOffset uint32, bitCount uint32) bool
	String() string
}

const elementSizeInBytes = 8
const bitsPerElement = elementSizeInBytes * 8

type BitmapImpl struct {
	data []uint64
}

var _ Bitmap = (*BitmapImpl)(nil)

func (b *BitmapImpl) String() string {

	if len(b.data) == 0 {
		return ""
	}

	lastIdx := 0
	for i := len(b.data) - 1; i >= 0; i-- {
		if b.data[i] != 0 {
			lastIdx = i
			break
		}
	}
	var result strings.Builder
	for offset := 0; offset <= lastIdx; offset++ {
		value := b.data[offset]
		var mask = uint64(1) << 63
		for ; mask != 0; mask = mask >> 1 {
			if value&mask != 0 {
				result.WriteRune('1')
			} else {
				result.WriteRune('0')
			}
		}
		if offset+1 < lastIdx {
			result.WriteString("\n")
		}
	}
	return strings.TrimRight(result.String(), "0")
}

func (b *BitmapImpl) simpleWrite(bitOffset uint32, bitCount uint32) bool {
	startIndexInclusive := bitOffset / bitsPerElement
	endIndexInclusive := (bitOffset + bitCount - 1) / bitsPerElement

	if startIndexInclusive != endIndexInclusive {
		panic(fmt.Sprintf("simpleWrite() called with wrong arguments, start index %d is not equal to end index %d", startIndexInclusive, endIndexInclusive))
	}

	// bitCount
	var mask uint64
	maskLen := bitCount
	if maskLen > 1 {
		mask = (1 << maskLen) - 1
	} else {
		mask = 1
	}
	bitoffsetInLongWord := bitOffset % bitsPerElement
	leftShift := 64 - maskLen - bitoffsetInLongWord
	mask <<= leftShift
	GrowSliceIfNecessary(startIndexInclusive, &b.data)
	if b.data[startIndexInclusive]&mask != 0 {
		return false
	}
	b.data[startIndexInclusive] |= mask
	return true
}

func (b *BitmapImpl) Set(bitOffset uint32, bitCount uint32) bool {

	if bitCount == 0 {
		return true
	}

	startIndexInclusive := bitOffset / bitsPerElement
	endIndexInclusive := (bitOffset + bitCount - 1) / bitsPerElement

	if startIndexInclusive == endIndexInclusive {
		// trivial case -> write within one element only
		return b.simpleWrite(bitOffset, bitCount)
	}

	// write spreads across at least two elements

	bitsRemaining := bitCount

	bitOffsetAtStartSegment := startIndexInclusive * bitsPerElement
	cnt := 64 - (bitOffset - bitOffsetAtStartSegment)
	if !b.simpleWrite(bitOffset, cnt) {
		return false
	}

	bitsRemaining -= cnt
	if bitsRemaining == 0 {
		return true
	}

	// write intermediate segments (if any)
	currentIndex := startIndexInclusive + 1
	var mask uint64 = 0xffffffffffffffff
	for bitsRemaining > bitsPerElement {
		GrowSliceIfNecessary(currentIndex, &b.data)
		if b.data[currentIndex]&mask != 0 {
			return false
		}
		b.data[currentIndex] |= mask
		currentIndex++
		bitsRemaining -= bitsPerElement
	}

	// write last segment (if any)
	if bitsRemaining > 0 {
		return b.simpleWrite(endIndexInclusive*bitsPerElement, bitsRemaining)
	}
	return true
}

func NewBitmap() Bitmap {
	return &BitmapImpl{}
}
