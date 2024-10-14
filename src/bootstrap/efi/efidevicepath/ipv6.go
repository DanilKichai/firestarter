package efidevicepath

import (
	"bootstrap/efi/common"
	"encoding/binary"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#ipv6-device-path

const IPv6Type = 3 + 13*0x100

const (
	IPv6ManualOrigin        = 0x00
	IPv6StatelessAutoOrigin = 0x01
	IPv6StatefullAutoOrigin = 0x02
)

type IPv6 struct {
	LocalIPAddress   []byte
	RemoteIPAddress  []byte
	LocalPort        uint16
	RemotePort       uint16
	Protocol         []byte
	IPAddressOrigin  byte
	PrefixLength     uint8
	GatewayIPAddress []byte
}

func (ip *IPv6) UnmarshalBinary(data []byte) error {
	if len(data) < 56 {
		return common.ErrDataIsTooShort
	}

	ip.LocalIPAddress = data[0:16]
	ip.RemoteIPAddress = data[16:32]
	ip.LocalPort = binary.LittleEndian.Uint16(data[32:34])
	ip.RemotePort = binary.LittleEndian.Uint16(data[34:36])
	ip.Protocol = data[36:38]
	ip.PrefixLength = uint8(data[39:40][0])
	ip.GatewayIPAddress = data[40:56]

	return nil
}
