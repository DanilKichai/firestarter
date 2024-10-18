package efidevicepath

// https://uefi.org/specs/UEFI/2.10/10_Protocols_Device_Path_Protocol.html#uniform-resource-identifiers-uri-device-path

const URIType = 3 + 24*0x100

type URI struct {
	Data string
}

func (u *URI) UnmarshalBinary(data []byte) error {
	u.Data = string(data)

	return nil
}
