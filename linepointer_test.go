package bithacking

import (
	"testing"
)

func TestNewLinePointer(t *testing.T) {
	out, err := NewLinePointer(
		maxLinePointerOffset,
		maxState,
		maxLinePointerLength,
	)
	const max = uint32(1<<32 - 1)
	if err != nil {
		t.Error(err)
	} else if uint32(out) != max {
		t.Errorf("wrong value, got %d, expected %d", uint32(out), max)
	}

	t.Run("errors", func(t *testing.T) {
		t.Run("offset", func(t *testing.T) {
			var input uint16
			var err error

			input = uint16(32767)
			_, err = NewLinePointer(input, uint8(0), uint16(0))
			if err != nil {
				t.Errorf("got unexpected error %v for input %d", err, input)
			}

			input = uint16(32768)
			_, err = NewLinePointer(input, uint8(0), uint16(0))
			if err == nil {
				t.Errorf("expected error but got none for input %d", input)
			}
		})

		t.Run("state", func(t *testing.T) {
			var input uint8
			var err error

			input = uint8(3)
			_, err = NewLinePointer(uint16(0), input, uint16(0))
			if err != nil {
				t.Errorf("got unexpected error %v for input %d", err, input)
			}

			input = uint8(4)
			_, err = NewLinePointer(uint16(0), input, uint16(0))
			if err == nil {
				t.Errorf("expected error but got none for input %d", input)
			}
		})

		t.Run("length", func(t *testing.T) {
			var input uint16
			var err error

			input = uint16(32767)
			_, err = NewLinePointer(uint16(0), uint8(0), input)
			if err != nil {
				t.Errorf("got unexpected error %v for input %d", err, input)
			}

			input = uint16(32768)
			_, err = NewLinePointer(uint16(0), uint8(0), input)
			if err == nil {
				t.Errorf("expected error but got none for input %d", input)
			}
		})
	})
}

func TestLinePointerString(t *testing.T) {
	tests := []struct {
		offset   uint16
		state    uint8
		length   uint16
		expected string
	}{
		{
			offset:   0,
			state:    0,
			length:   0,
			expected: "000000000000000 00 000000000000000",
		},
		{
			offset:   32767,
			state:    3,
			length:   32767,
			expected: "111111111111111 11 111111111111111",
		},
		{
			offset:   1024,
			state:    1,
			length:   2048,
			expected: "000010000000000 01 000100000000000",
		},
		{
			offset:   2048,
			state:    2,
			length:   1024,
			expected: "000100000000000 10 000010000000000",
		},
	}

	for i, test := range tests {
		out, err := NewLinePointer(test.offset, test.state, test.length)
		if err != nil {
			t.Error(i, err)
			return
		}
		actual := out.String()
		if actual != test.expected {
			t.Errorf("test %d\ngot      %q\nexpected %q", i, actual, test.expected)
		}
	}
}

func TestLinePointerOffset(t *testing.T) {
	inputs := []uint16{0, 1024, 32767}

	for i, input := range inputs {
		line, err := NewLinePointer(
			input,
			maxState-1,
			maxLinePointerLength-1,
		)
		if err != nil {
			t.Error(i, err)
			return
		}
		actual := line.Offset()
		if actual != input {
			t.Errorf("test %d, got %d, expected %d", i, actual, input)
		}
	}
}

func TestLinePointerState(t *testing.T) {
	inputs := []uint8{0, 1, 2, 3}

	for i, input := range inputs {
		line, err := NewLinePointer(
			maxLinePointerOffset-1,
			input,
			maxLinePointerLength-1,
		)
		if err != nil {
			t.Error(i, err)
			return
		}
		actual := line.State()
		if actual != input {
			t.Errorf("test %d, got %d, expected %d", i, actual, input)
		}
	}
}

func TestLinePointerLength(t *testing.T) {
	inputs := []uint16{0, 1024, 32767}

	for i, input := range inputs {
		line, err := NewLinePointer(
			maxLinePointerOffset-1,
			maxState-1,
			input,
		)
		if err != nil {
			t.Error(i, err)
			return
		}
		actual := line.Length()
		if actual != input {
			t.Errorf("test %d, got %d, expected %d", i, actual, input)
		}
	}
}
