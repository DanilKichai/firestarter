package common

import (
	"encoding/binary"
	"unicode/utf16"
)

func GetNullTerminatedUnicodeString(data []byte, offset int) (string, int, error) {
	if offset+2 > len(data) {
		return "", offset, ErrDataSize
	}

	term := false
	var ustr []uint16

	for ; offset < len(data)-1; offset += 2 {
		char := binary.LittleEndian.Uint16(data[offset : offset+2])

		if char == 0 {
			term = true
			offset += 2

			break
		}

		ustr = append(ustr, char)
	}

	if !term {
		return "", offset, ErrDataSize
	}

	return string(utf16.Decode(ustr)), offset, nil
}
