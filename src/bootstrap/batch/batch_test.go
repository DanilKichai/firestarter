package batch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	var testCases = []struct {
		caseName    string
		batch       *Batch
		expectedErr error
	}{
		{
			caseName: "valid",
			batch: &[]Batch{{
				{
					Path: filepath.Join(t.TempDir(), "w_mode"),
					Mode: &[]int{0644}[0],
				},
				{
					Path: filepath.Join(t.TempDir(), "wo_mode"),
				},
			}}[0],
		},
		{
			caseName: "valid with bad path",
			batch: &[]Batch{{
				{
					Path: "fixtures/not_exists_dir/file",
					Mode: &[]int{0644}[0],
				},
			}}[0],
			expectedErr: os.ErrNotExist,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			err := testCase.batch.Write()

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, testCase.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
