package bithacking

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimestamp(t *testing.T) {
	tables := []struct {
		input    time.Time
		expected uint32
	}{
		{time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), 540673},
		{
			time.Date(2018, 5, 12, 15, 53, 0, 0, time.UTC),
			992681902,
		},
	}

	for _, test := range tables {
		actual := newTimestamp(test.input)
		cmp := uint32(actual)
		// fmt.Printf("%032b\n", cmp)

		if cmp != test.expected {
			t.Errorf("%d != %d\n", cmp, test.expected)
		}
	}
}

func TestDecodeTimestamp(t *testing.T) {
	tables := []struct {
		input    string
		expected []uint32
	}{
		{"00000000000010000100000000000001", []uint32{0, 1, 1, 0, 0, 1}},
		{"00111011001010110001111110101110", []uint32{118, 5, 12, 15, 53, 6}},
	}

	for _, test := range tables {
		actual := decodeTimestamp(test.input)

		for i, v := range actual {
			if v != test.expected[i] {
				t.Errorf("%d != %d\n", v, test.expected[i])
			}
		}
	}
}

func TestExtractMid(t *testing.T) {
	var actual uint32
	var cmp string

	actual = extractMid(0xdeadbeef, 4, 16)
	cmp = fmt.Sprintf("%x", actual)

	if cmp != "bee" {
		t.Errorf("%v != %v\n", cmp, "bee")
	}
}

func TestBinStrToNum(t *testing.T) {
	tables := []struct {
		input    string
		expected uint32
	}{
		{"00000000000010000100000000000001", uint32(540673)},
		{"00011111001000110001111011110110", uint32(522395382)},
	}

	for _, test := range tables {
		actual, err := binStrToNum(test.input)

		if err != nil {
			t.Errorf("%v\n", err)
		} else if actual != test.expected {
			t.Errorf("\nactual:  %32d\nexpected %32d\n", actual, test.expected)
		}
	}
}
