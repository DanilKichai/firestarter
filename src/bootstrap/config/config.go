package config

import (
	"bootstrap/efi/efidevicepath"
	"bootstrap/efi/efivarfs"
	"fmt"
	"net"
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
	MAC  *string
	VLAN *int
	IPv4 *IPv4
	IPv6 *IPv6
	DNS  []string
	URI  *string
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

			m := mac.MACAddress.String()
			cfg.MAC = &m

		case efidevicepath.VLANType:
			vlan, err := efidevicepath.ParsePath[*efidevicepath.VLAN](fp.Data)
			if err != nil {
				return nil, fmt.Errorf("parse VLAN from current load option: %w", err)
			}

			v := int(vlan.Vlanid)
			cfg.VLAN = &v

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

			u := uri.Data
			cfg.URI = &u
		}

	}

	return &cfg, nil
}
