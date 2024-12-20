package config

import (
	"archshell/pkg/efi/common"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	var testCases = []struct {
		caseName       string
		efivars        string
		expectedResult *Config
		expectedErr    error
	}{
		{
			caseName:    "not exists efivars",
			efivars:     "fixtures/not_exists",
			expectedErr: os.ErrNotExist,
		},
		{
			caseName: "valid efivars with dynamic DHCP for IPv4 and IPv6",
			efivars:  "fixtures/valid_dhcp",
			expectedResult: &[]Config{{
				MAC:  &[]string{"3c:ec:ef:c4:45:80"}[0],
				VLAN: &[]int{1}[0],
				IPv4: &IPv4{
					Static:  false,
					Address: "0.0.0.0/0",
					Gateway: "0.0.0.0",
					DNS: []string{
						"192.168.0.1",
						"192.168.0.2",
					},
				},
				IPv6: &IPv6{
					Static:   false,
					Stateful: false,
					Address:  "::/0",
					Gateway:  "::",
				},
				URI: &[]string{"http://www.google.com/"}[0],
			}}[0],
		},
		{
			caseName: "valid efivars with dynamic DHCP for IPv4 and IPv6 (stateful IPv6)",
			efivars:  "fixtures/valid_dhcp_stateful6",
			expectedResult: &[]Config{{
				MAC:  &[]string{"3c:ec:ef:c4:45:80"}[0],
				VLAN: &[]int{1}[0],
				IPv4: &IPv4{
					Static:  false,
					Address: "0.0.0.0/0",
					Gateway: "0.0.0.0",
				},
				IPv6: &IPv6{
					Static:   false,
					Stateful: true,
					Address:  "::/0",
					Gateway:  "::",
					DNS: []string{
						"2001:4860:4860::8888",
						"2001:4860:4860::8844",
					},
				},
				URI: &[]string{"http://www.google.com/"}[0],
			}}[0],
		},
		{
			caseName: "valid efivars with static addresses for IPv4 and IPv6",
			efivars:  "fixtures/valid_static",
			expectedResult: &[]Config{{
				MAC:  &[]string{"3c:ec:ef:c4:45:80"}[0],
				VLAN: &[]int{1}[0],
				IPv4: &IPv4{
					Static:  true,
					Address: "192.168.0.11/24",
					Gateway: "192.168.0.254",
					DNS: []string{
						"192.168.0.1",
						"192.168.0.2",
					},
				},
				IPv6: &IPv6{
					Static:   true,
					Stateful: false,
					Address:  "bdf9:564f:9126:fcfa:a1c8:818:d700:44f2/64",
					Gateway:  "::",
				},
				URI: &[]string{"http://www.google.com/"}[0],
			}}[0],
		},
		{
			caseName:    "invalid efivars with invalid load option",
			efivars:     "fixtures/invalid_lo",
			expectedErr: common.ErrDataSize,
		},
		{
			caseName:    "invalid efivars with invalid DNS",
			efivars:     "fixtures/invalid_dns",
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName:    "invalid efivars with invalid IPv4",
			efivars:     "fixtures/invalid_ipv4",
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName:    "invalid efivars with invalid IPv6",
			efivars:     "fixtures/invalid_ipv6",
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName:    "invalid efivars with invalid MAC",
			efivars:     "fixtures/invalid_mac",
			expectedErr: common.ErrDataSize,
		},
		{
			caseName:    "invalid efivars with invalid VLAN",
			efivars:     "fixtures/invalid_vlan",
			expectedErr: common.ErrDataSize,
		},
		{
			caseName: "valid efivars with empty URI",
			efivars:  "fixtures/empty_uri",
			expectedResult: &[]Config{{
				MAC:  &[]string{"3c:ec:ef:c4:45:80"}[0],
				VLAN: &[]int{1}[0],
				IPv4: &IPv4{
					Static:  false,
					Address: "0.0.0.0/0",
					Gateway: "0.0.0.0",
					DNS: []string{
						"192.168.0.1",
						"192.168.0.2",
					},
				},
				IPv6: &IPv6{
					Static:   false,
					Stateful: false,
					Address:  "::/0",
					Gateway:  "::",
				},
			}}[0],
		},
		{
			caseName:    "invalid efivars with invalid URI",
			efivars:     "fixtures/invalid_uri",
			expectedErr: common.ErrDataRepresentation,
		},
		{
			caseName: "valid efivars with GPT partitioned HardDrive + FilePath",
			efivars:  "fixtures/valid_harddrive(gpt)+filepath",
			expectedResult: &[]Config{{
				PartitionUUID: &[]string{"65a5135b-ecd7-4bde-a7a2-251e7b74659d"}[0],
				FilePath:      &[]string{`EFI\Linux\arch-linux-lts.efi`}[0],
			}}[0],
		},
		{
			caseName: "valid efivars with MBR partitioned HardDrive + FilePath",
			efivars:  "fixtures/valid_harddrive(mbr)+filepath",
			expectedResult: &[]Config{{
				PartitionUUID: &[]string{"65a5135b-01"}[0],
				FilePath:      &[]string{`EFI\Linux\arch-linux-lts.efi`}[0],
			}}[0],
		},
		{
			caseName: "valid efivars with unsupported HardDrive Signature + FilePath",
			efivars:  "fixtures/unsupported_harddrive(signature)+filepath",
			expectedResult: &[]Config{{
				FilePath: &[]string{`EFI\Linux\arch-linux-lts.efi`}[0],
			}}[0],
		},
		{
			caseName:    "invalid efivars with invalid HardDrive",
			efivars:     "fixtures/invalid_harddrive",
			expectedErr: common.ErrDataSize,
		},
		{
			caseName:    "invalid efivars with invalid FilePath",
			efivars:     "fixtures/invalid_filepath",
			expectedErr: common.ErrDataSize,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			cfg, err := Load(testCase.efivars)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, cfg)
				assert.Equal(t, testCase.expectedResult, cfg)
			}
		})
	}
}
