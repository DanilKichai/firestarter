package config

import (
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
			caseName: "valid efivars with static addresses for IPv4 and IPv6",
			efivars:  "valid_static",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "valid efivars with dynamic DHCP for IPv4 and IPv6",
			efivars:  "valid_dhcp",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "valid efivars with dynamic DHCP for IPv4 and IPv6 (force IPv6 DHCP)",
			efivars:  "valid_dhcp_force6",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid DNS",
			efivars:  "invalid_dns",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid IPv4",
			efivars:  "invalid_ipv4",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid IPv6",
			efivars:  "invalid_ipv6",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid MAC",
			efivars:  "invalid_mac",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid URI",
			efivars:  "invalid_uri",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
		},
		{
			caseName: "invalid efivars with invalid VLAN",
			efivars:  "invalid_vlan",
			expectedResult: &[]Config{{
				MAC:  "",
				VLAN: 0,
				IPv4: IPv4{},
				IPv6: IPv6{},
				DNS:  []string{},
				URI:  "",
			}}[0],
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
