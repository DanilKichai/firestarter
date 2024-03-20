package main

import (
	efivarfs "bootstrap/efi/efivarfs"
	"fmt"
)

func main() {
	current, err := efivarfs.ParseVar[efivarfs.BootCurrent]("BootCurrent", efivarfs.GlobalVariable)
	if err != nil {
		return
	}

	fmt.Println(current.Value)

	entry, err := efivarfs.ParseVar[efivarfs.BootEntry](fmt.Sprintf("Boot%04X", current.Value), efivarfs.GlobalVariable)
	if err != nil {
		return
	}

	fmt.Println(entry.FilePathListLength)
}

// https://uefi.org/specs/UEFI/2.10/24_Network_Protocols_SNP_PXE_BIS.html#device-path
