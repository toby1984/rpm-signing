package common

import (
	"fmt"
)

// DynamicArray is an array one can write to
// at arbitrary indices, and it will automatically
// grow to accommodate the given number of bytes.
//
// Any write overlapping with a previous write will
// be reported as an error.
type DynamicArray interface {
	// Write data from a []byte to a specific offset.
	// Any write overlapping with a previous write will
	// be reported as an error.
	Write(p []byte, offset uint32) error
	// GetData returns all data written so far.
	// Any gaps between writes will be filled with zeroes.
	GetData() []byte
}

type dynamicArrayImpl struct {
	bitmap                  Bitmap
	data                    []byte
	lastUsedOffsetExclusive uint32
}

var _ DynamicArray = (*dynamicArrayImpl)(nil)

func (da *dynamicArrayImpl) Write(p []byte, offset uint32) error {
	if len(p) == 0 {
		return nil
	}

	if !da.bitmap.Set(offset, uint32(len(p))) {
		return fmt.Errorf("write to overlapping range")
	}

	maxOffsetExclusive := offset + uint32(len(p))
	GrowSliceIfNecessary(maxOffsetExclusive, &da.data)

	copy(da.data[offset:], p)

	if maxOffsetExclusive > da.lastUsedOffsetExclusive {
		da.lastUsedOffsetExclusive = maxOffsetExclusive
	}
	return nil
}

func (da *dynamicArrayImpl) GetData() []byte {
	result := make([]byte, da.lastUsedOffsetExclusive)
	copy(result, da.data[0:da.lastUsedOffsetExclusive])
	return result
}

func NewDynamicArray() DynamicArray {
	return &dynamicArrayImpl{
		bitmap: NewBitmap(),
	}
}
