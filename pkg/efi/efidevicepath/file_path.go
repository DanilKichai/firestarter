package efidevicepath

import (
	"archshell/pkg/efi/common"
	"fmt"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#file-path-media-device-path

const FilePathType = 4 + 4*0x100

type FilePath struct {
	PathName string
}

func (fp *FilePath) UnmarshalBinary(data []byte) error {
	file, _, err := common.GetNullTerminatedUnicodeString(data, 0)
	if err != nil {
		return fmt.Errorf("the \"FilePath\" data parse error occured: %w", err)
	}

	fp.PathName = file

	return nil
}
