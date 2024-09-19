package efidevicepath

import (
	"encoding/binary"
	"fmt"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#vlan-device-path-node

const VLANType = 3 + 20*0x100

type VLAN uint16

func (v *VLAN) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("unmarshal data is too short")
	}

	*v = VLAN(binary.LittleEndian.Uint16(data[0:2]))

	return nil
}
