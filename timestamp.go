package bithacking

import (
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

// timestampOffsets is the distance (in bits) from right-most bit of a time unit
// to left-most bit of the entire 32-bit timestamp.
var timestampOffsets = [6]uint8{
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
	r |= uint32(n) << (31 - timestampOffsets[0])

	n = int(t.Month())
	r |= uint32(n) << (31 - timestampOffsets[1])

	n = t.Day()
	r |= uint32(n) << (31 - timestampOffsets[2])

	n = t.Hour()
	r |= uint32(n) << (31 - timestampOffsets[3])

	n = t.Minute()
	r |= uint32(n) << (31 - timestampOffsets[4])

	n = int(t.Weekday())
	r |= uint32(n) << (31 - timestampOffsets[5])

	return Timestamp(r)
}

func BinStrToNum(b string) (num uint32, err error) {
	if len(b) != 32 {
		err = fmt.Errorf("input must have length 32, length is %d", len(b))
		return
	}

	i, err := strconv.ParseUint(b, 2, 32)
	num = uint32(i)
	return
}

func DecodeTimestamp(s string) []uint8 {
	vals := make([]uint8, 6)
	binStr, _ := BinStrToNum(s)

	n := uint32(31)

	for i, offset := range timestampOffsets {
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
