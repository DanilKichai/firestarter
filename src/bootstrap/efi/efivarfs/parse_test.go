package efivarfs

import (
	"bootstrap/efi/common"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVar(t *testing.T) {
	var testCases_BootCurrent = []struct {
		caseName           string
		efivars            string
		name               string
		guid               string
		expectedResult     interface{}
		expectedResultType string
		expectedErr        error
	}{
		{
			caseName:           "correct BootCurrent",
			efivars:            "fixtures",
			name:               "BootCurrent",
			guid:               GlobalVariable,
			expectedResult:     &[]BootCurrent{0x1002}[0],
			expectedResultType: "*efivarfs.BootCurrent",
		},
		{
			caseName:           "incorrect BootCurrent with too short(BootCurrent,5)",
			efivars:            "fixtures",
			name:               "BootCurrent_short(BootCurrent,5)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.BootCurrent",
			expectedErr:        common.ErrDataIsTooShort,
		},
		{
			caseName:           "not exists BootCurrent",
			efivars:            "fixtures",
			name:               "BootCurrent_not_exists",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.BootCurrent",
			expectedErr:        os.ErrNotExist,
		},
		{
			caseName: "correct LoadOption",
			efivars:  "fixtures",
			name:     "Boot1002",
			guid:     GlobalVariable,
			expectedResult: &[]LoadOption{{
				Attributes:         0x01,
				FilePathListLength: 0x6c,
				Description:        "ArchLinux (LTS kernel)",
				FilePathList: FilePathList{
					FilePath{
						Type: 0x104,
						Data: []uint8{
							0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00,
							0x00, 0x00, 0x00, 0x00, 0x5b, 0x13, 0xa5, 0x65,
							0xd7, 0xec, 0xde, 0x4b, 0xa7, 0xa2, 0x25, 0x1e,
							0x7b, 0x74, 0x65, 0x9d, 0x02, 0x02,
						},
					},
					FilePath{
						Type: 0x404,
						Data: []uint8{
							0x45, 0x00, 0x46, 0x00, 0x49, 0x00, 0x5c, 0x00,
							0x4c, 0x00, 0x69, 0x00, 0x6e, 0x00, 0x75, 0x00,
							0x78, 0x00, 0x5c, 0x00, 0x61, 0x00, 0x72, 0x00,
							0x63, 0x00, 0x68, 0x00, 0x2d, 0x00, 0x6c, 0x00,
							0x69, 0x00, 0x6e, 0x00, 0x75, 0x00, 0x78, 0x00,
							0x2d, 0x00, 0x6c, 0x00, 0x74, 0x00, 0x73, 0x00,
							0x2e, 0x00, 0x65, 0x00, 0x66, 0x00, 0x69, 0x00,
							0x00, 0x00,
						},
					},
					FilePath{
						Type: 0xff7f,
						Data: []uint8{},
					},
				},
				OptionalData: []uint8{},
			}}[0],
			expectedResultType: "*efivarfs.LoadOption",
		},
		{
			caseName:           "empty LoadOption with empty(FilePathList)",
			efivars:            "fixtures",
			name:               "Boot1002_empty(FilePathList)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedResult: &[]LoadOption{{
				Attributes:         0x1,
				FilePathListLength: 0x0,
				Description:        "ArchLinux (LTS kernel)",
				FilePathList:       FilePathList(nil),
				OptionalData:       []uint8{},
			}}[0],
		},
		{
			caseName:           "incorrect LoadOption with unterminated(Description)",
			efivars:            "fixtures",
			name:               "Boot1002_unterminated(Description)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        common.ErrDataIsTooShort,
		},
		{
			caseName:           "incorrect LoadOption with too short(FilePath,3)",
			efivars:            "fixtures",
			name:               "Boot1002_short(FilePath,3)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        ErrIncorrectFilePathLength,
		},
		{
			caseName:           "incorrect LoadOption with too long(FilePath,out_of_range)",
			efivars:            "fixtures",
			name:               "Boot1002_long(FilePath,out_of_range)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        common.ErrDataIsTooShort,
		},
		{
			caseName:           "incorrect LoadOption with too long(FilePathList,out_of_range)",
			efivars:            "fixtures",
			name:               "Boot1002_long(FilePathList,out_of_range)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        common.ErrDataIsTooShort,
		},
		{
			caseName:           "incorrect LoadOption with too short(LoadOption,9)",
			efivars:            "fixtures",
			name:               "Boot1002_short(LoadOption,9)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        common.ErrDataIsTooShort,
		},
		{
			caseName:           "incorrect LoadOption with too short(FilePathList,3)",
			efivars:            "fixtures",
			name:               "Boot1002_short(FilePathList,3)",
			guid:               GlobalVariable,
			expectedResultType: "*efivarfs.LoadOption",
			expectedErr:        common.ErrDataIsTooShort,
		},
	}

	for _, testCase := range testCases_BootCurrent {
		t.Run(testCase.caseName, func(t *testing.T) {
			var (
				err error
				res interface{}
			)

			switch testCase.expectedResultType {
			case "*efivarfs.BootCurrent":
				res, err = ParseVar[*BootCurrent](testCase.efivars, testCase.name, testCase.guid)
			case "*efivarfs.LoadOption":
				res, err = ParseVar[*LoadOption](testCase.efivars, testCase.name, testCase.guid)
			default:
				t.Fatalf("Unexpected type provided: %s", testCase.expectedResultType)
			}

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, testCase.expectedResult, res)
			}
		})
	}

}
