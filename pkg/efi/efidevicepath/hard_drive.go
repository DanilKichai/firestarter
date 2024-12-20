package efidevicepath

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#hard-drive

const HardDriveType = 4 + 1*0x100

const (
	HardDriveNoDiskSignatureType = 0x00
	HardDriveMBRSignatureType    = 0x01
	HardDriveGUIDSignatureType   = 0x02
)

const (
	HardDriveMBRPartitionFormat  = 0x01
	HardDriveGUIDPartitionFormat = 0x01
)

const (
	HardDriveNoDiskSignature = 0x00
	HardDriveMBRSignature    = 0x01
	HardDriveGUIDSignature   = 0x02
)

type HardDrive struct {
	PartitionNumber    uint32
	PartitionStart     uint64
	PartitionSize      uint64
	PartitionSignature [16]byte
	PartitionFormat    uint8
	SignatureType      uint8
}

func (hd *HardDrive) UnmarshalBinary(data []byte) error {
	if len(data) != 38 {
		return common.ErrDataSize
	}

	hd.PartitionNumber = binary.LittleEndian.Uint32(data[0:4])
	hd.PartitionStart = binary.LittleEndian.Uint64(data[4:12])
	hd.PartitionSize = binary.LittleEndian.Uint64(data[12:20])
	hd.PartitionSignature = [16]byte(data[20:36])
	hd.PartitionFormat = data[36:37][0]
	hd.SignatureType = data[37:38][0]
	return nil
}
