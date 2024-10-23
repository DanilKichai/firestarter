package efidevicepath

import (
	"bootstrap/efi/common"
	"encoding/binary"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#vlan-device-path-node

const VLANType = 3 + 20*0x100

type VLAN struct {
	Vlanid uint16
}

func (v *VLAN) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return common.ErrDataSize
	}

	id := binary.LittleEndian.Uint16(data[0:2])
	if id > 4095 || id == 0 {
		return common.ErrDataRepresentation
	}

	v.Vlanid = id

	return nil
}
