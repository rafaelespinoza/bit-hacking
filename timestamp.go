package bithacking

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type Timestamp uint32

// name,	 	length, 	range, 		possible vals
// (empty),     1 bit, 		(0,   1), 	none
// year,        8 bits, 	(1,   8), 	0-255 (num of years after 1900, highest year is 2155)
// month,       4 bits, 	(9,  12), 	1-12
// day,         5 bits, 	(13, 17), 	1-31
// hour of day, 5 bits, 	(18, 22), 	0-23
// min of hour, 6 bits, 	(23, 28), 	0-59
// weekday,     3 bits, 	(29, 31), 	0-6

// distance (in bits) from right-most bit of a time unit to left-most bit of the
// entire 32-bit timestamp
var bitOffsets = [6]uint8{
	8,  // year
	12, // month
	17, // day of month
	22, // hour of day
	28, // min of hour
	31, // weekday
}

// similar to link: https://stackoverflow.com/a/40309527
func decodeTimestamp(s string) []uint32 {
	vals := make([]uint32, 6)
	binStr, _ := binStrToNum(s)

	n := uint32(31) // distance from right of HSB in binStr

	for i, offset := range bitOffsets {
		var timeVal uint32
		offsetEnd := uint32(31 - offset) // 23

		var bit uint32
		for ; n >= offsetEnd; n-- {
			mask := (1 << uint32(n)) - 0
			bit = binStr & uint32(mask)
			timeVal |= bit
			// fmt.Printf("%02d\t%032b\n\t%032b\n\t%032b\n", n, mask, binStr, timeVal)
			if n == 0 {
				break // decrementing a zero-value uint is a bad idea
			}
		}

		timeVal >>= offsetEnd
		vals[i] = timeVal
	}

	return vals
}

// Algorithm for cutting out from start to end is a two-step process. Shift
// original value from start bits to the right. Then perform bit-wise AND with
// mask of (end - start) ones. Input bit start is inclusive, end is exclusive.
// Bits are numbered from 0.  Bit mask with N ones at the end. `1<<end` is
// `2^end`, which has a single `1` at position `end+1`, then has all zeroes to
// the right. Get the needed bitmask by subtracting 1.
// https://stackoverflow.com/a/10090443, https://stackoverflow.com/a/10090450
func extractMid(val, start, end uint32) uint32 {
	bits := val >> start // drop LSB's
	dist := end - start
	mask := (1 << dist) - 1
	return bits & uint32(mask)
}

func crappy_extractTimeVals(s string) []uint32 {
	exp := []uint32{118, 5, 12, 15, 53, 6} // only care about checking May 12 2018
	vals := make([]uint32, 6)
	n, _ := binStrToNum(s)

	fmt.Printf("%2s %8s %4s %16s %6s %33s %33s\n", "i", "offset", "dist", "act", "exp", "act bin", "exp bin")

	for i, offset := range bitOffsets {
		d := 31 - offset
		var a uint32
		a = n >> d
		a = a << d
		c := uint32(extractMid(n, uint32(d), 32))
		fmt.Printf("%2d %8d %4d %16d %6d %033b %033b\n", i, offset, d, c, exp[i], c, exp[i])
		vals[i] |= c
	}

	return vals
}

func newTimestamp(t time.Time) Timestamp {
	r := uint32(0)

	n := t.Year() - 1900
	r |= uint32(n) << (31 - bitOffsets[0])

	n = int(t.Month())
	r |= uint32(n) << (31 - bitOffsets[1])

	n = t.Day()
	r |= uint32(n) << (31 - bitOffsets[2])

	n = t.Hour()
	r |= uint32(n) << (31 - bitOffsets[3])

	n = t.Minute()
	r |= uint32(n) << (31 - bitOffsets[4])

	n = int(t.Weekday())
	r |= uint32(n) << (31 - bitOffsets[5])

	return Timestamp(r)
}

func binStrToNum(b string) (uint32, error) {
	if len(b) != 32 {
		m := fmt.Sprintf("Input must be length of 32. Length is %d", len(b))
		e := errors.New(m)
		return uint32(0), e
	}

	s := 0

	for i, c := range b {
		if c != '0' && c != '1' {
			e := errors.New("each character must represent a bit as 0 or 1")
			return uint32(0), e
		}

		s += valueAt(i, c)
	}

	return uint32(s), nil
}

func valueAt(i int, c rune) int {
	p := exponentiate2(31 - i)
	n := int(c)
	return p * (n - 48) // 48 is ASCII for `0`
}

func exponentiate2(n int) int {
	m := float64(n)
	p := math.Pow(2, m)
	return int(p)
}
