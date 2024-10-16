package efidevicepath

import (
	"bootstrap/efi/common"
	"net/netip"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#dns-device-path

const DNSType = 3 + 31*0x100

type DNS struct {
	Instances []netip.Addr
}

type EFIIPAddress [16]byte

func (eip EFIIPAddress) Addr(isIPv6 bool) netip.Addr {
	if isIPv6 {
		return netip.AddrFrom16(eip)
	} else {
		return netip.AddrFrom4([4]byte(eip[0:4]))
	}
}

func (d *DNS) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return common.ErrDataIsTooShort
	}

	var isIPv6 bool
	switch data[0:1][0] {
	case 0x00:
		isIPv6 = false
	case 0x01:
		isIPv6 = true
	default:
		return ErrInvalidBooleanRepresentation
	}

	resolvers := data[1:]
	for i := 0; i < len(resolvers); i += 128 {
		a := EFIIPAddress(resolvers[i : i+128]).Addr(isIPv6)
		d.Instances = append(d.Instances, a)
	}

	return nil
}
