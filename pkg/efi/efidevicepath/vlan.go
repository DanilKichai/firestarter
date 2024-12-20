package efidevicepath

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
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

	v.Vlanid = binary.LittleEndian.Uint16(data[0:2])

	return nil
}
