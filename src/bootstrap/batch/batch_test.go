package batch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var testCases = []struct {
		caseName          string
		batch             *Batch
		expectedErr       error
		expecterErrSubstr *string
	}{
		{
			caseName: "valid",
			batch: &[]Batch{{
				{
					Path: filepath.Join(t.TempDir(), "file_w_mode"),
					Mode: &[]int{0644}[0],
					Type: "file",
				},
				{
					Path: filepath.Join(t.TempDir(), "file_wo_mode"),
					Type: "file",
				},
				{
					Path: filepath.Join(t.TempDir(), "directory_w_mode"),
					Mode: &[]int{0755}[0],
					Type: "directory",
				},
				{
					Path: filepath.Join(t.TempDir(), "directory_wo_mode"),
					Type: "directory",
				},
				{
					Path: filepath.Join(t.TempDir(), "r/directory_w_mode"),
					Mode: &[]int{0755}[0],
					Type: "rdirectory",
				},
				{
					Path: filepath.Join(t.TempDir(), "r/directory_wo_mode"),
					Type: "rdirectory",
				},
				{
					Path: filepath.Join(t.TempDir(), "symlink"),
					Type: "symlink",
					Data: "/dev/null",
				},
			}}[0],
		},
		{
			caseName: "valid file with not exists path",
			batch: &[]Batch{{
				{
					Path: "fixtures/not_exists_dir/file",
					Mode: &[]int{0644}[0],
					Type: "file",
				},
			}}[0],
			expectedErr: os.ErrNotExist,
		},
		{
			caseName: "valid directory with not exists path",
			batch: &[]Batch{{
				{
					Path: "fixtures/not_exists_dir/directory",
					Mode: &[]int{0644}[0],
					Type: "directory",
				},
			}}[0],
			expectedErr: os.ErrNotExist,
		},
		{
			caseName: "valid symlink with not exists path",
			batch: &[]Batch{{
				{
					Path: "fixtures/not_exists_dir/symlink",
					Type: "symlink",
					Data: "/dev/null",
				},
			}}[0],
			expectedErr: os.ErrNotExist,
		},
		{
			caseName: "valid rdirectory with invalid path",
			batch: &[]Batch{{
				{
					Path: "fixtures/file",
					Type: "rdirectory",
				},
			}}[0],
			expecterErrSubstr: &[]string{"not a directory"}[0],
		},
		{
			caseName: "invalid with unsupported type",
			batch: &[]Batch{{
				{
					Path: filepath.Join(t.TempDir(), "invalid"),
					Type: "bad",
				},
			}}[0],
			expectedErr: ErrUnsupportedType,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			err := testCase.batch.Write()

			if testCase.expectedErr != nil || testCase.expecterErrSubstr != nil {
				if testCase.expectedErr != nil {
					assert.ErrorIs(t, err, testCase.expectedErr)
				}

				if testCase.expecterErrSubstr != nil {
					assert.ErrorContains(t, err, *testCase.expecterErrSubstr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
