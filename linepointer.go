package bithacking

import "fmt"

// LinePointer is a recreation of ItemIdData in postgres. It's a compact
// representation of three pieces of data: offset, state and length; 15, 2 and
// 15 bits respectively. The sum is 32 bits, which is an important cache
// alignment optimization.
type LinePointer uint32

const (
	maxLinePointerOffset = uint16(1<<15) - 1
	maxState             = uint8(1<<2) - 1
	maxLinePointerLength = uint16(1<<15) - 1
)

// linePointerOffsets is the number of bits from right-most bit of a LinePointer
// to left-most bit.
var linePointerOffsets = [3]uint8{
	15, // length, 15 bits, [0, 32768)
	17, // state,   2 bits, [0,     4)
	32, // offset, 15 bits, [0, 32768)
}

// NewLinePointer constructs a LinePointer value. It returns an error if the
// offset or length inputs cannot fit into a 15-bit unsigned integer. It will
// also return an error if state cannot fit into a 2-bit integer. These
// restrictions are to enforce the same widths of the fields that are in
// postgres' ItemIdData struct.
func NewLinePointer(offset uint16, state uint8, length uint16) (out LinePointer, err error) {
	if offset > maxLinePointerOffset {
		err = fmt.Errorf("offset must be <= %d", maxLinePointerOffset)
		return
	}
	if state > maxState {
		err = fmt.Errorf("state must be <= %d", maxState)
		return
	}
	if length > maxLinePointerLength {
		err = fmt.Errorf("length must be <= %d", maxLinePointerLength)
		return
	}
	var p, q uint32

	q = uint32(offset)
	p |= q << (32 - linePointerOffsets[0])

	q = uint32(state)
	p |= q << (32 - linePointerOffsets[1])

	q = uint32(length)
	p |= q << (32 - linePointerOffsets[2])

	out = LinePointer(p)
	return
}

func (p LinePointer) String() string {
	type linePointer LinePointer
	q := fmt.Sprintf("%032b", linePointer(p))

	// offset + " " + state + " " + length
	return q[:32-linePointerOffsets[1]] + " " + q[linePointerOffsets[0]:linePointerOffsets[1]] + " " + q[32-linePointerOffsets[0]:]
}

// Offset retrieves the leftmost 15 bits.
func (p LinePointer) Offset() uint16 {
	q := uint32(p)
	r := uint16(q >> linePointerOffsets[1])
	return r & maxLinePointerOffset
}

// State retrieves the bits in between the length, offset bits.
func (p LinePointer) State() uint8 {
	q := uint32(p)
	// distance from lowest bit of middle (the middle two bits represent the
	// state we want to extract) to least significant bit of entire LinePointer.
	r := uint8(q >> linePointerOffsets[0])
	// remove higher order bits from middle's left side
	return r & 3
}

// Length retrieves the rightmost 15 bits.
func (p LinePointer) Length() uint16 {
	q := uint16(p) // do not need the most significant 16 bits (leftmost 16)
	return q & maxLinePointerLength
}
