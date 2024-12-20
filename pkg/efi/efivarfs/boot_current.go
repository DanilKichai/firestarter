package efivarfs

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
)

// https://uefi.org/specs/UEFI/2.10/03_Boot_Manager.html#globally-defined-variables

type BootCurrent uint16

func (bc *BootCurrent) UnmarshalBinary(data []byte) error {
	if len(data) != 6 {
		return common.ErrDataSize
	}

	*bc = BootCurrent(binary.LittleEndian.Uint16(data[4:6]))

	return nil
}
