package efidevicepath

import (
	"archshell/pkg/efi/common"
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePath(t *testing.T) {
	var testCases_DNS = []struct {
		caseName       string
		data           []byte
		expectedResult *DNS
		expectedErr    error
	}{
		{
			caseName: "valid IPv4 DNS",
			data: []byte{
				0,
				192, 168, 0, 1,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				192, 168, 0, 2,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expectedResult: &[]DNS{{
				IsIPv6: false,
				Instances: []netip.Addr{
					netip.AddrFrom4([4]byte{192, 168, 0, 1}),
					netip.AddrFrom4([4]byte{192, 168, 0, 2}),
				},
			}}[0],
		},
		{
			caseName: "valid IPv6 DNS",
			data: []byte{
				1,

				0x66, 0x0a, 0xde, 0x97,
				0xc2, 0x43, 0xf7, 0xf7,
				0x93, 0x5a, 0xc1, 0x8a,
				0xee, 0xfb, 0xa5, 0x01,

				0x66, 0x0a, 0xde, 0x97,
				0xc2, 0x43, 0xf7, 0xf7,
				0x93, 0x5a, 0xc1, 0x8a,
				0xee, 0xfb, 0xa5, 0x02,
			},
			expectedResult: &[]DNS{{
				IsIPv6: true,
				Instances: []netip.Addr{
					netip.AddrFrom16([16]byte{
						0x66, 0x0a, 0xde, 0x97,
						0xc2, 0x43, 0xf7, 0xf7,
						0x93, 0x5a, 0xc1, 0x8a,
						0xee, 0xfb, 0xa5, 0x01,
					}),
					netip.AddrFrom16([16]byte{
						0x66, 0x0a, 0xde, 0x97,
						0xc2, 0x43, 0xf7, 0xf7,
						0x93, 0x5a, 0xc1, 0x8a,
						0xee, 0xfb, 0xa5, 0x02,
					}),
				},
			}}[0],
		},
		{
			caseName: "invalid IPv4 DNS with short data",
			data: []byte{
				0,

				192, 168, 0, 1,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				192, 168, 0, 2,
			},
			expectedErr: common.ErrDataSize,
		},
		{
			caseName: "invalid IPv4 DNS with incorrect boolean representation",
			data: []byte{
				2,

				192, 168, 0, 1,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,

				192, 168, 0, 2,
			},
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName:    "empty IPv4 DNS",
			data:        []byte{},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_DNS {
		t.Run(testCase.caseName, func(t *testing.T) {
			dns, err := ParsePath[*DNS](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, dns)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, dns)
				assert.Equal(t, testCase.expectedResult, dns)
			}
		})
	}

	var testCases_FilePath = []struct {
		caseName       string
		data           []byte
		expectedResult *FilePath
		expectedErr    error
	}{
		{
			caseName: "valid FilePath",
			data: []byte{
				0x45, 0x00, 0x46, 0x00, 0x49, 0x00, 0x5c, 0x00,
				0x4c, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x75, 0x00,
				0x78, 0x00, 0x5c, 0x00, 0x61, 0x00, 0x72, 0x00,
				0x63, 0x00, 0x68, 0x00, 0x2d, 0x00, 0x6c, 0x00,
				0x69, 0x00, 0x6e, 0x00, 0x75, 0x00, 0x78, 0x00,
				0x2d, 0x00, 0x6c, 0x00, 0x74, 0x00, 0x73, 0x00,
				0x2e, 0x00, 0x65, 0x00, 0x66, 0x00, 0x69, 0x00,
				0x00, 0x00,
			},
			expectedResult: &[]FilePath{{
				PathName: `EFI\Linux\arch-linux-lts.efi`,
			}}[0],
		},
		{
			caseName: "invalid FilePath with too short data",
			data: []byte{
				0x45, 0x00, 0x46, 0x00, 0x49, 0x00, 0x5c, 0x00,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_FilePath {
		t.Run(testCase.caseName, func(t *testing.T) {
			file, err := ParsePath[*FilePath](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, file)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, file)
				assert.Equal(t, testCase.expectedResult, file)
			}
		})
	}

	var testCases_HardDrive = []struct {
		caseName       string
		data           []byte
		expectedResult *HardDrive
		expectedErr    error
	}{
		{
			caseName: "valid HardDrive",
			data: []byte{
				0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x5b, 0x13, 0xa5, 0x65,
				0xd7, 0xec, 0xde, 0x4b, 0xa7, 0xa2, 0x25, 0x1e,
				0x7b, 0x74, 0x65, 0x9d, 0x02, 0x02,
			},
			expectedResult: &[]HardDrive{{
				PartitionNumber: 0x1,
				PartitionStart:  0x800,
				PartitionSize:   0x200000,
				PartitionSignature: [16]uint8{
					0x5b, 0x13, 0xa5, 0x65, 0xd7, 0xec, 0xde, 0x4b,
					0xa7, 0xa2, 0x25, 0x1e, 0x7b, 0x74, 0x65, 0x9d,
				},
				PartitionFormat: 0x02,
				SignatureType:   0x02,
			}}[0],
		},
		{
			caseName: "invalid HardDrive with too short data",
			data: []byte{
				0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_HardDrive {
		t.Run(testCase.caseName, func(t *testing.T) {
			hd, err := ParsePath[*HardDrive](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, hd)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, hd)
				assert.Equal(t, testCase.expectedResult, hd)
			}
		})
	}

	var testCases_IPv4 = []struct {
		caseName       string
		data           []byte
		expectedResult *IPv4
		expectedErr    error
	}{
		{
			caseName: "valid static IPv4",
			data: []byte{
				192, 168, 0, 101,
				192, 168, 0, 1,
				0x00, 0x00,
				0x00, 0x00,
				0x06, 0x00,
				0x01,
				192, 168, 0, 1,
				255, 255, 255, 0,
			},
			expectedResult: &[]IPv4{{
				LocalIPAddress:   netip.AddrFrom4([4]byte{192, 168, 0, 101}),
				RemoteIPAddress:  netip.AddrFrom4([4]byte{192, 168, 0, 1}),
				LocalPort:        0,
				RemotePort:       0,
				Protocol:         0x06,
				StaticIPAddress:  true,
				GatewayIPAddress: netip.AddrFrom4([4]byte{192, 168, 0, 1}),
				SubnetMask:       netip.AddrFrom4([4]byte{255, 255, 255, 0}),
			}}[0],
		},
		{
			caseName: "valid dynamic IPv4",
			data: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x06, 0x00,
				0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expectedResult: &[]IPv4{{
				LocalIPAddress:   netip.AddrFrom4([4]byte{0, 0, 0, 0}),
				RemoteIPAddress:  netip.AddrFrom4([4]byte{0, 0, 0, 0}),
				LocalPort:        0,
				RemotePort:       0,
				Protocol:         0x06,
				StaticIPAddress:  false,
				GatewayIPAddress: netip.AddrFrom4([4]byte{0, 0, 0, 0}),
				SubnetMask:       netip.AddrFrom4([4]byte{0, 0, 0, 0}),
			}}[0],
		},
		{
			caseName: "invalid IPv4 with invalid boolean representation",
			data: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00,
				0x00, 0x00,
				0x06, 0x00,
				0x02,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName: "invalid IPv4 with too short data",
			data: []byte{
				0x00, 0x00, 0x00, 0x00,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_IPv4 {
		t.Run(testCase.caseName, func(t *testing.T) {
			ip4, err := ParsePath[*IPv4](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, ip4)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, ip4)
				assert.Equal(t, testCase.expectedResult, ip4)
			}
		})
	}

	var testCases_IPv6 = []struct {
		caseName       string
		data           []byte
		expectedResult *IPv6
		expectedErr    error
	}{
		{
			caseName: "valid static IPv6",
			data: []byte{
				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x02,

				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x01,

				0x00, 0x00,

				0x00, 0x00,

				0x06, 0x00,

				0x00,

				0x40,

				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x01,
			},
			expectedResult: &[]IPv6{{
				LocalIPAddress: netip.AddrFrom16([16]byte{
					0xd2, 0xba, 0x55, 0x04,
					0xeb, 0x1c, 0xcf, 0x7d,
					0x45, 0x6d, 0x6d, 0x68,
					0x5e, 0x70, 0xd2, 0x02,
				}),
				RemoteIPAddress: netip.AddrFrom16([16]byte{
					0xd2, 0xba, 0x55, 0x04,
					0xeb, 0x1c, 0xcf, 0x7d,
					0x45, 0x6d, 0x6d, 0x68,
					0x5e, 0x70, 0xd2, 0x01,
				}),
				LocalPort:       0,
				RemotePort:      0,
				Protocol:        0x06,
				IPAddressOrigin: 0x00,
				PrefixLength:    64,
				GatewayIPAddress: netip.AddrFrom16([16]byte{
					0xd2, 0xba, 0x55, 0x04,
					0xeb, 0x1c, 0xcf, 0x7d,
					0x45, 0x6d, 0x6d, 0x68,
					0x5e, 0x70, 0xd2, 0x01,
				}),
			}}[0],
		},
		{
			caseName: "invalid static IPv6 with invalid IPAddressOrigin representation",
			data: []byte{
				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x02,

				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x01,

				0x00, 0x00,

				0x00, 0x00,

				0x06, 0x00,

				0x04,

				0x40,

				0xd2, 0xba, 0x55, 0x04,
				0xeb, 0x1c, 0xcf, 0x7d,
				0x45, 0x6d, 0x6d, 0x68,
				0x5e, 0x70, 0xd2, 0x01,
			},
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName: "invalid static IPv6 with short data",
			data: []byte{
				0xd2, 0xba, 0x55, 0x04,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_IPv6 {
		t.Run(testCase.caseName, func(t *testing.T) {
			ip6, err := ParsePath[*IPv6](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, ip6)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, ip6)
				assert.Equal(t, testCase.expectedResult, ip6)
			}
		})
	}

	var testCases_MACAddress = []struct {
		caseName       string
		data           []byte
		expectedResult *MACAddress
		expectedErr    error
	}{
		{
			caseName: "valid MAC Address",
			data: []byte{
				0x3c, 0x55, 0x76, 0xbd, 0xfb, 0x97, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

				0x01,
			},
			expectedResult: &[]MACAddress{{
				MACAddress: net.HardwareAddr{0x3c, 0x55, 0x76, 0xbd, 0xfb, 0x97},
				IfType:     0x01,
			}}[0],
		},
		{
			caseName: "invalid MAC Address with too short data",
			data: []byte{
				0x3c, 0x55, 0x76, 0xbd, 0xfb, 0x97, 0x00, 0x00,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_MACAddress {
		t.Run(testCase.caseName, func(t *testing.T) {
			mac, err := ParsePath[*MACAddress](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, mac)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, mac)
				assert.Equal(t, testCase.expectedResult, mac)
			}
		})
	}

	var testCases_URI = []struct {
		caseName       string
		data           []byte
		expectedResult *URI
		expectedErr    error
	}{
		{
			caseName: "valid URI",
			data: []byte{
				0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x77,
				0x77, 0x77, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
				0x2e, 0x63, 0x6f, 0x6d, 0x2f,
			},
			expectedResult: &[]URI{{
				Data: "https://www.google.com/",
			}}[0],
		},
		{
			caseName: "empty URI",
			data:     []byte{},
			expectedResult: &[]URI{{
				Data: "",
			}}[0],
		},
		{
			caseName: "invalid URI",
			data: []byte{
				0x68, 0x74, 0x74, 0x70, 0x73, 0x2f, 0x2f, 0x2f, 0x77,
				0x77, 0x77, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
				0x2e, 0x63, 0x6f, 0x6d, 0x2f,
			},
			expectedErr: common.ErrDataRepresentation,
		},
	}

	for _, testCase := range testCases_URI {
		t.Run(testCase.caseName, func(t *testing.T) {
			uri, err := ParsePath[*URI](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, uri)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, uri)
				assert.Equal(t, testCase.expectedResult, uri)
			}
		})
	}

	var testCases_VLAN = []struct {
		caseName       string
		data           []byte
		expectedResult *VLAN
		expectedErr    error
	}{
		{
			caseName: "valid VLAN",
			data: []byte{
				0x00, 0x04,
			},
			expectedResult: &[]VLAN{{
				Vlanid: 1024,
			}}[0],
		},
		{
			caseName: "short VLAN",
			data: []byte{
				0x00,
			},
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases_VLAN {
		t.Run(testCase.caseName, func(t *testing.T) {
			vlan, err := ParsePath[*VLAN](testCase.data)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, vlan)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, vlan)
				assert.Equal(t, testCase.expectedResult, vlan)
			}
		})
	}
}
