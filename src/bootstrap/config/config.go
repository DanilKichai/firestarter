package config

import (
	"bootstrap/efi/efidevicepath"
	"bootstrap/efi/efivarfs"
	"fmt"
	"net"
	"slices"
)

type IPv4 struct {
	Static  bool
	Address string
	Gateway string
}

type IPv6 struct {
	Static      bool
	SolicitDHCP bool
	Address     string
	Gateway     string
}

type Config struct {
	MAC               *string
	VLAN              *int
	IPv4              *IPv4
	IPv6              *IPv6
	DNS               []string
	URI               *string
	PartitionSelector *string
	FilePath          *string
}

func Load(efivars string) (*Config, error) {
	current, err := efivarfs.ParseVar[*efivarfs.BootCurrent](efivars, "BootCurrent", efivarfs.GlobalVariable)
	if err != nil {
		return nil, fmt.Errorf("get boot current value: %w", err)
	}

	entry, err := efivarfs.ParseVar[*efivarfs.LoadOption](efivars, fmt.Sprintf("Boot%04X", *current), efivarfs.GlobalVariable)
	if err != nil {
		return nil, fmt.Errorf("get current load option: %w", err)
	}

	cfg := Config{}

	for _, fp := range entry.FilePathList {
		switch fp.Type {
		case efidevicepath.MACAddressType:
			mac, err := efidevicepath.ParsePath[*efidevicepath.MACAddress](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse MAC from current load option: %w", err)
			}

			cfg.MAC = &[]string{mac.MACAddress.String()}[0]

		case efidevicepath.VLANType:
			vlan, err := efidevicepath.ParsePath[*efidevicepath.VLAN](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse VLAN from current load option: %w", err)
			}

			cfg.VLAN = &[]int{int(vlan.Vlanid)}[0]

		case efidevicepath.IPv4Type:
			ipv4, err := efidevicepath.ParsePath[*efidevicepath.IPv4](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse IPv4 from current load option: %w", err)
			}

			addr := ipv4.LocalIPAddress.String()
			prefix, _ := net.IPMask(net.ParseIP(ipv4.SubnetMask.String()).To4()).Size()

			cfg.IPv4 = &IPv4{
				Static:  ipv4.StaticIPAddress,
				Address: fmt.Sprintf("%s/%d", addr, prefix),
				Gateway: ipv4.GatewayIPAddress.String(),
			}

		case efidevicepath.IPv6Type:
			ipv6, err := efidevicepath.ParsePath[*efidevicepath.IPv6](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse IPv6 from current load option: %w", err)
			}

			static := false
			if ipv6.IPAddressOrigin == efidevicepath.IPv6ManualOrigin {
				static = true
			}

			dhcp := false
			if ipv6.IPAddressOrigin == efidevicepath.IPv6StatefullAutoOrigin {
				dhcp = true
			}

			addr := ipv6.LocalIPAddress.String()
			prefix := ipv6.PrefixLength

			cfg.IPv6 = &IPv6{
				Static:      static,
				SolicitDHCP: dhcp,
				Address:     fmt.Sprintf("%s/%d", addr, prefix),
				Gateway:     ipv6.GatewayIPAddress.String(),
			}

		case efidevicepath.DNSType:
			dns, err := efidevicepath.ParsePath[*efidevicepath.DNS](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse DNS from current load option: %w", err)
			}

			for _, addr := range dns.Instances {
				cfg.DNS = append(cfg.DNS, addr.String())
			}

		case efidevicepath.URIType:
			uri, err := efidevicepath.ParsePath[*efidevicepath.URI](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse URI from current load option: %w", err)
			}

			cfg.URI = &[]string{uri.Data}[0]

		case efidevicepath.HardDriveType:
			hd, err := efidevicepath.ParsePath[*efidevicepath.HardDrive](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse HardDrive from current load option: %w", err)
			}

			var selector string

			if hd.SignatureType == efidevicepath.HardDriveMBRSignature {
				selector =
					fmt.Sprintf(
						"PARTUUID=%08x-%02d",
						hd.PartitionSignature[0:4],
						hd.PartitionNumber,
					)
			}

			if hd.SignatureType == efidevicepath.HardDriveGUIDSignature {
				s1 := hd.PartitionSignature[0:4]
				s2 := hd.PartitionSignature[4:6]
				s3 := hd.PartitionSignature[6:8]

				slices.Reverse(s1)
				slices.Reverse(s2)
				slices.Reverse(s3)

				selector = fmt.Sprintf(
					"PARTUUID=%04x-%02x-%02x-%02x-%06x",
					s1,
					s2,
					s3,
					hd.PartitionSignature[8:10],
					hd.PartitionSignature[10:16],
				)
			}

			cfg.PartitionSelector = &selector

		case efidevicepath.FilePathType:
			file, err := efidevicepath.ParsePath[*efidevicepath.FilePath](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse FilePath from current load option: %w", err)
			}

			cfg.FilePath = &[]string{file.PathName}[0]
		}
	}

	return &cfg, nil
}
