package efidevicepath

import (
	"bootstrap/efi/common"
	"net"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#mac-address-device-path

const MACAddressType = 3 + 11*0x100

type MACAddress struct {
	MACAddress net.HardwareAddr
	IfType     byte
}

func (ma *MACAddress) UnmarshalBinary(data []byte) error {
	if len(data) != 33 {
		return common.ErrDataSize
	}

	addr := data[0:6]
	if [6]byte(addr) == [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00} || [6]byte(addr) == [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff} {
		return common.ErrDataRepresentation
	}

	ma.MACAddress = addr
	ma.IfType = data[32:33][0]

	return nil
}
