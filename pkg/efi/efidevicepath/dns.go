package efidevicepath

import (
	"archshell/pkg/efi/common"
	"net/netip"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#dns-device-path

const DNSType = 3 + 31*0x100

type DNS struct {
	IsIPv6    bool
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
		return common.ErrDataSize
	}

	switch data[0:1][0] {
	case 0x00:
		d.IsIPv6 = false
	case 0x01:
		d.IsIPv6 = true
	default:
		return common.ErrDataRepresentation
	}

	resolvers := data[1:]
	if len(resolvers)%16 != 0 {
		return common.ErrDataSize
	}

	for i := 0; i < len(resolvers); i += 16 {
		a := EFIIPAddress(resolvers[i : i+16]).Addr(d.IsIPv6)
		d.Instances = append(d.Instances, a)
	}

	return nil
}
