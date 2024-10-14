package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNullTerminatedUnicodeString(t *testing.T) {
	var testCases = []struct {
		caseName       string
		data           []byte
		offset         int
		expectedString string
		expectedOffset int
		expectedErr    error
	}{
		{
			caseName:       "correct (without offset)",
			data:           []byte{116, 0, 101, 0, 115, 0, 116, 0, 0, 0},
			offset:         0,
			expectedString: "test",
			expectedOffset: 10,
		},
		{
			caseName:       "correct (with offset)",
			data:           []byte{116, 0, 101, 0, 115, 0, 116, 0, 0, 0},
			offset:         4,
			expectedString: "st",
			expectedOffset: 10,
		},
		{
			caseName:       "incorrect offset",
			data:           []byte{116, 0, 101, 0, 115, 0, 116, 0, 0, 0},
			offset:         9,
			expectedString: "",
			expectedErr:    ErrDataIsTooShort,
		},
		{
			caseName:       "unterminated",
			data:           []byte{116, 0, 101, 0, 115, 0, 116, 0},
			offset:         0,
			expectedString: "",
			expectedErr:    ErrDataIsTooShort,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			s, offset, err := GetNullTerminatedUnicodeString(testCase.data, testCase.offset)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
				assert.Equal(t, "", s)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedString, s)
				assert.Equal(t, testCase.expectedOffset, offset)
			}
		})
	}
}
