package efidevicepath

import (
	"bootstrap/efi/common"
	"encoding/binary"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#ipv4-device-path

const IPv4Type = 3 + 12*0x100

type IPv4 struct {
	LocalIPAddress   []byte
	RemoteIPAddress  []byte
	LocalPort        uint16
	RemotePort       uint16
	Protocol         []byte
	StaticIPAddress  bool
	GatewayIPAddress []byte
	SubnetMask       []byte
}

func (ip *IPv4) UnmarshalBinary(data []byte) error {
	if len(data) < 23 {
		return common.ErrDataIsTooShort
	}

	ip.LocalIPAddress = data[0:4]
	ip.RemoteIPAddress = data[4:8]
	ip.LocalPort = binary.LittleEndian.Uint16(data[8:10])
	ip.RemotePort = binary.LittleEndian.Uint16(data[10:12])
	ip.Protocol = data[12:14]

	switch data[14:15][0] {
	case 0x00:
		ip.StaticIPAddress = false
	case 0x01:
		ip.StaticIPAddress = true
	default:
		return ErrInvalidBooleanRepresentation
	}

	ip.GatewayIPAddress = data[15:19]
	ip.SubnetMask = data[19:23]

	return nil
}
