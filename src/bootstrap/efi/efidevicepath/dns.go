package efidevicepath

import (
	"bootstrap/efi/common"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#dns-device-path

const DNSType = 3 + 31*0x100

type DNS struct {
	IsIPv6    bool
	Instances []EFIIPAddress
}

type EFIIPAddress []byte

func (d *DNS) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return common.ErrDataIsTooShort
	}

	switch data[0:1][0] {
	case 0x00:
		d.IsIPv6 = false
	case 0x01:
		d.IsIPv6 = true
	default:
		return ErrInvalidBooleanRepresentation
	}

	for i := 0; i < (len(data)-1)/128; i++ {
		eia := EFIIPAddress(data[128*i+1 : 128*(i+1)+1])
		d.Instances = append(d.Instances, eia)
	}

	return nil
}
