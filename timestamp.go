package bithacking

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Timestamp uint32

// (empty),	none,	1 bit, 	(0,   1)
// year,	0-255,	8 bits, (1,   8), num of years after 1900, highest year is 2155
// month, 	1-12,	4 bits, (9,  12)
// day, 	1-31,	5 bits, (13, 17)
// hour,	0-23, 	5 bits, (18, 22)
// min, 	0-59,	6 bits, (23, 28)
// weekday,	0-6, 	3 bits, (29, 31)

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

func BinStrToNum(b string) (uint32, error) {
	if len(b) != 32 {
		m := fmt.Sprintf("Input must be length of 32. Length is %d", len(b))
		e := errors.New(m)
		return uint32(0), e
	}

	i, e := strconv.ParseUint(b, 2, 32)
	return uint32(i), e
}

func DecodeTimestamp(s string) []uint8 {
	vals := make([]uint8, 6)
	binStr, _ := BinStrToNum(s)

	n := uint32(31)

	for i, offset := range bitOffsets {
		var offsetEnd, mask, bit, timeVal uint32
		offsetEnd = 31 - uint32(offset)

		for ; n >= offsetEnd; n-- {
			mask = 1 << n
			bit = binStr & mask
			timeVal |= bit
			if n < 1 {
				break // decrementing a zero-value uint is a bad idea
			}
		}

		timeVal >>= offsetEnd
		vals[i] = uint8(timeVal)
	}

	return vals
}
