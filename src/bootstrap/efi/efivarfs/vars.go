package efivarfs

import (
	"encoding/binary"
	"errors"
)

type BootCurrent struct {
	Attributes uint32
	Value      uint16
}

// https://uefi.org/specs/UEFI/2.10/03_Boot_Manager.html#load-options
type BootEntry struct {
	Attributes         uint32
	FilePathListLength uint16
}

func (b *BootCurrent) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return errors.New("unmarshal error")
	}

	b.Attributes = binary.LittleEndian.Uint32(data[0:4])
	b.Value = binary.LittleEndian.Uint16(data[4:6])

	return nil
}

func (b *BootEntry) UnmarshalBinary(data []byte) error {
	b.Attributes = binary.LittleEndian.Uint32(data[0:4])
	b.FilePathListLength = binary.LittleEndian.Uint16(data[4:6])
	return nil
}
