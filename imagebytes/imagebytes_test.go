package imagebytes_test

import (
	"bytes"
	"testing"

	"github.com/pillowskiy/imagesize/imagebytes"
)

func TestReadU16(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		buf        []byte
		endianness imagebytes.Endian
		expected   uint16
		expectErr  bool
	}{
		{
			name:       "LittleEndian U16",
			buf:        []byte{0x01, 0x02},
			endianness: imagebytes.LittleEndian,
			expected:   0x0201,
			expectErr:  false,
		},
		{
			name:       "BigEndian U16",
			buf:        []byte{0x01, 0x02},
			endianness: imagebytes.BigEndian,
			expected:   0x0102,
			expectErr:  false,
		},
		{
			name:       "Invalid Endian",
			buf:        []byte{0x01, 0x02},
			endianness: 99,
			expected:   0,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.buf)
			result, err := imagebytes.ReadU16(reader, tt.endianness)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadU24(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		buf        []byte
		endianness imagebytes.Endian
		expected   uint32
		expectErr  bool
	}{
		{
			name:       "LittleEndian U24",
			buf:        []byte{0x01, 0x02, 0x03},
			endianness: imagebytes.LittleEndian,
			expected:   0x030201,
			expectErr:  false,
		},
		{
			name:       "BigEndian U24",
			buf:        []byte{0x01, 0x02, 0x03},
			endianness: imagebytes.BigEndian,
			expected:   0x010203,
			expectErr:  false,
		},
		{
			name:       "Invalid Endian",
			buf:        []byte{0x01, 0x02, 0x03},
			endianness: 99,
			expected:   0,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.buf)
			result, err := imagebytes.ReadU24(reader, tt.endianness)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadU32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		buf        []byte
		endianness imagebytes.Endian
		expected   uint32
		expectErr  bool
	}{
		{
			name:       "LittleEndian U32",
			buf:        []byte{0x01, 0x02, 0x03, 0x04},
			endianness: imagebytes.LittleEndian,
			expected:   0x04030201,
			expectErr:  false,
		},
		{
			name:       "BigEndian U32",
			buf:        []byte{0x01, 0x02, 0x03, 0x04},
			endianness: imagebytes.BigEndian,
			expected:   0x01020304,
			expectErr:  false,
		},
		{
			name:       "Invalid Endian",
			buf:        []byte{0x01, 0x02, 0x03, 0x04},
			endianness: 99,
			expected:   0,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.buf)
			result, err := imagebytes.ReadU32(reader, tt.endianness)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
