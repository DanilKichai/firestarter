package main

import (
	"bootstrap/efi/efidevicepath"
	"bootstrap/efi/efivarfs"
	"fmt"
)

func main() {
	current, err := efivarfs.ParseVar[*efivarfs.BootCurrent]("/sys/firmware/efi/efivars", "BootCurrent", efivarfs.GlobalVariable)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	entry, err := efivarfs.ParseVar[*efivarfs.LoadOption]("/sys/firmware/efi/efivars", fmt.Sprintf("Boot%04X", *current), efivarfs.GlobalVariable)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	mac := &efidevicepath.MACAddress{}
	vlan := &efidevicepath.VLAN{}
	ipv4 := &efidevicepath.IPv4{}
	ipv6 := &efidevicepath.IPv6{}
	dns := &efidevicepath.DNS{}
	uri := &efidevicepath.URI{}

	for _, fp := range entry.FilePathList {
		switch fp.Type {
		case efidevicepath.MACAddressType:
			mac, err = efidevicepath.ParsePath[*efidevicepath.MACAddress](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}

		case efidevicepath.VLANType:
			vlan, err = efidevicepath.ParsePath[*efidevicepath.VLAN](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}

		case efidevicepath.IPv4Type:
			ipv4, err = efidevicepath.ParsePath[*efidevicepath.IPv4](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}

		case efidevicepath.IPv6Type:
			ipv6, err = efidevicepath.ParsePath[*efidevicepath.IPv6](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}

		case efidevicepath.DNSType:
			dns, err = efidevicepath.ParsePath[*efidevicepath.DNS](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}

		case efidevicepath.URIType:
			uri, err = efidevicepath.ParsePath[*efidevicepath.URI](fp.Data)
			if err != nil {
				fmt.Println(err.Error())

				return
			}
		}

	}

	if mac != nil && vlan != nil && ipv4 != nil && ipv6 != nil && dns != nil && uri != nil {

	}

	fmt.Println(entry.Description)
}
