package efidevicepath

import (
	"archshell/pkg/efi/common"
	"net/url"
)

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#uniform-resource-identifiers-uri-device-path

const URIType = 3 + 24*0x100

type URI struct {
	Data string
}

func (u *URI) UnmarshalBinary(data []byte) error {
	s := string(data)

	if len(s) == 0 {
		return nil
	}

	_, err := url.ParseRequestURI(s)
	if err != nil {
		return common.ErrDataRepresentation
	}

	u.Data = s

	return nil
}
