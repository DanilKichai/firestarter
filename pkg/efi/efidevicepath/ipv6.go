package efidevicepath

import (
	"encoding/binary"
	"archshell/pkg/efi/common"
	"net/netip"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#ipv6-device-path

const IPv6Type = 3 + 13*0x100

const (
	IPv6ManualOrigin        = 0x00
	IPv6StatelessAutoOrigin = 0x01
	IPv6StatefulAutoOrigin  = 0x02
)

type IPv6 struct {
	LocalIPAddress   netip.Addr
	RemoteIPAddress  netip.Addr
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint16
	IPAddressOrigin  byte
	PrefixLength     uint8
	GatewayIPAddress netip.Addr
}

func (ip *IPv6) UnmarshalBinary(data []byte) error {
	if len(data) != 56 {
		return common.ErrDataSize
	}

	ip.LocalIPAddress = netip.AddrFrom16([16]byte(data[0:16]))
	ip.RemoteIPAddress = netip.AddrFrom16([16]byte(data[16:32]))
	ip.LocalPort = binary.LittleEndian.Uint16(data[32:34])
	ip.RemotePort = binary.LittleEndian.Uint16(data[34:36])
	ip.Protocol = binary.LittleEndian.Uint16(data[36:38])

	if origin := data[38:39][0]; origin > 0x02 {
		return common.ErrDataRepresentation
	} else {
		ip.IPAddressOrigin = origin
	}

	ip.PrefixLength = uint8(data[39:40][0])
	ip.GatewayIPAddress = netip.AddrFrom16([16]byte(data[40:56]))

	return nil
}
