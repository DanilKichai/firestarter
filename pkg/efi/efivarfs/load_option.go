package efivarfs

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
	"fmt"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#efi-device-path-protocol

type FilePath struct {
	Type uint16
	Data []byte
}

type FilePathList []FilePath

func (fpl *FilePathList) AppendFilePathFromBinary(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, common.ErrDataSize
	}

	t := binary.LittleEndian.Uint16(data[0:2])
	l := binary.LittleEndian.Uint16(data[2:4])

	if l < 4 {
		return nil, common.ErrFilePathLength
	}

	if len(data) < int(l) {
		return nil, common.ErrDataSize
	}

	fp := FilePath{
		Type: t,
		Data: data[4:l],
	}

	*fpl = append(*fpl, fp)

	return data[l:], nil
}

func (fpl *FilePathList) UnmarshalBinary(data []byte) (err error) {
	for len(data) > 0 {
		if data, err = fpl.AppendFilePathFromBinary(data); err != nil {
			return fmt.Errorf("the \"FilePath\" data parse error occured: %w", err)
		}
	}

	return nil
}

// https://uefi.org/specs/UEFI/2.10/03_Boot_Manager.html#load-options

type LoadOption struct {
	Attributes         uint32
	FilePathListLength uint16
	Description        string
	FilePathList       FilePathList
	OptionalData       []uint8
}

func (lo *LoadOption) UnmarshalBinary(data []byte) error {
	if len(data) < 10 {
		return common.ErrDataSize
	}

	lo.Attributes = binary.LittleEndian.Uint32(data[4:8])
	lo.FilePathListLength = binary.LittleEndian.Uint16(data[8:10])

	desc, offset, err := common.GetNullTerminatedUnicodeString(data, 10)
	if err != nil {
		return fmt.Errorf("the \"Description\" data parse error occured: %w", err)
	}

	lo.Description = desc

	if len(data) < offset+int(lo.FilePathListLength) {
		return common.ErrDataSize
	}

	err = lo.FilePathList.UnmarshalBinary(data[offset : offset+int(lo.FilePathListLength)])
	if err != nil {
		return fmt.Errorf("the \"FilePathList\" data parse error occured: %w", err)
	}

	lo.OptionalData = data[offset+int(lo.FilePathListLength):]

	return nil
}
