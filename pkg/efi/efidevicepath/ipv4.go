package efidevicepath

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
	"net/netip"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#ipv4-device-path

const IPv4Type = 3 + 12*0x100

type IPv4 struct {
	LocalIPAddress   netip.Addr
	RemoteIPAddress  netip.Addr
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint16
	StaticIPAddress  bool
	GatewayIPAddress netip.Addr
	SubnetMask       netip.Addr
}

func (ip *IPv4) UnmarshalBinary(data []byte) error {
	if len(data) != 23 {
		return common.ErrDataSize
	}

	ip.LocalIPAddress = netip.AddrFrom4([4]byte(data[0:4]))
	ip.RemoteIPAddress = netip.AddrFrom4([4]byte(data[4:8]))
	ip.LocalPort = binary.LittleEndian.Uint16(data[8:10])
	ip.RemotePort = binary.LittleEndian.Uint16(data[10:12])
	ip.Protocol = binary.LittleEndian.Uint16(data[12:14])

	switch data[14:15][0] {
	case 0x00:
		ip.StaticIPAddress = false
	case 0x01:
		ip.StaticIPAddress = true
	default:
		return common.ErrDataRepresentation
	}

	ip.GatewayIPAddress = netip.AddrFrom4([4]byte(data[15:19]))
	ip.SubnetMask = netip.AddrFrom4([4]byte(data[19:23]))

	return nil
}
