package extractor

import (
	"encoding/binary"
	"errors"
	"io"
)

type Endian int

const (
	LittleEndian Endian = iota
	BigEndian
)

var ErrUnsupportedEndian = errors.New("unsupported endian")

func readU16(reader io.Reader, endianness Endian) (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}

	var result uint16
	switch endianness {
	case LittleEndian:
		result = binary.LittleEndian.Uint16(buf)
	case BigEndian:
		result = binary.BigEndian.Uint16(buf)
	default:
		return 0, ErrUnsupportedEndian
	}

	return result, nil
}

func readU24(reader io.Reader, endianness Endian) (uint32, error) {
	buf := make([]byte, 3)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}

	var result uint32
	switch endianness {
	case LittleEndian:
		result = uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16
	case BigEndian:
		result = uint32(buf[2]) | uint32(buf[1])<<8 | uint32(buf[0])<<16
	default:
		return 0, ErrUnsupportedEndian
	}

	return result, nil
}

func readU32(reader io.Reader, endianness Endian) (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}

	var result uint32
	switch endianness {
	case LittleEndian:
		result = binary.LittleEndian.Uint32(buf)
	case BigEndian:
		result = binary.BigEndian.Uint32(buf)
	default:
		return 0, ErrUnsupportedEndian
	}

	return result, nil
}
