package config

import (
	"archshell/pkg/efi/efidevicepath"
	"archshell/pkg/efi/efivarfs"
	"fmt"
	"net"
	"slices"
	"strings"
)

type IPv4 struct {
	Static  bool
	Address string
	Gateway string
	DNS     []string
}

type IPv6 struct {
	Static   bool
	Stateful bool
	Address  string
	Gateway  string
	DNS      []string
}

type Config struct {
	MAC           *string
	VLAN          *int
	IPv4          *IPv4
	IPv6          *IPv6
	URI           *string
	PartitionUUID *string
	FilePath      *string
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

			stateful := false
			if ipv6.IPAddressOrigin == efidevicepath.IPv6StatefulAutoOrigin {
				stateful = true
			}

			addr := ipv6.LocalIPAddress.String()
			prefix := ipv6.PrefixLength

			cfg.IPv6 = &IPv6{
				Static:   static,
				Stateful: stateful,
				Address:  fmt.Sprintf("%s/%d", addr, prefix),
				Gateway:  ipv6.GatewayIPAddress.String(),
			}

		case efidevicepath.DNSType:
			dns, err := efidevicepath.ParsePath[*efidevicepath.DNS](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse DNS from current load option: %w", err)
			}

			if dns.IsIPv6 {
				for _, addr := range dns.Instances {
					cfg.IPv6.DNS = append(cfg.IPv6.DNS, addr.String())
				}

			} else {
				for _, addr := range dns.Instances {
					cfg.IPv4.DNS = append(cfg.IPv4.DNS, addr.String())
				}

			}

		case efidevicepath.URIType:
			uri, err := efidevicepath.ParsePath[*efidevicepath.URI](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse URI from current load option: %w", err)
			}

			s := string(uri.Data)

			if len(s) == 0 {
				break
			}

			cfg.URI = &s

		case efidevicepath.HardDriveType:
			hd, err := efidevicepath.ParsePath[*efidevicepath.HardDrive](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse HardDrive from current load option: %w", err)
			}

			var partuuid string

			if hd.SignatureType == efidevicepath.HardDriveMBRSignature {
				s := hd.PartitionSignature[0:4]
				slices.Reverse(s)

				partuuid =
					fmt.Sprintf(
						"%08x-%02d",
						s,
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

				partuuid = fmt.Sprintf(
					"%04x-%02x-%02x-%02x-%06x",
					s1,
					s2,
					s3,
					hd.PartitionSignature[8:10],
					hd.PartitionSignature[10:16],
				)
			}

			if len(partuuid) == 0 {
				break
			}

			cfg.PartitionUUID = &partuuid

		case efidevicepath.FilePathType:
			file, err := efidevicepath.ParsePath[*efidevicepath.FilePath](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse FilePath from current load option: %w", err)
			}

			s := file.PathName
			s = strings.Replace(s, `\`, `/`, -1)

			cfg.FilePath = &s
		}
	}

	return &cfg, nil
}
