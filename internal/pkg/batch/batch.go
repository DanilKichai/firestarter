package batch

import (
	"fmt"
	"os"
)

type Batch []struct {
	Path string `yaml:"path"`
	Mode *int   `yaml:"mode"`
	Type string `yaml:"type"`
	Data string `yaml:"data"`
}

func (b *Batch) Write() error {
	for i, e := range *b {
		switch e.Type {
		case "file":
			m := os.FileMode(0o644)
			if e.Mode != nil {
				m = os.FileMode(*e.Mode)
			}

			if err := os.WriteFile(e.Path, []byte(e.Data), m); err != nil {
				return fmt.Errorf("write %d file of batch: %w", i, err)
			}

		case "directory":
			m := os.FileMode(0o755)
			if e.Mode != nil {
				m = os.FileMode(*e.Mode)
			}

			if err := os.Mkdir(e.Path, m); err != nil {
				return fmt.Errorf("write %d directory of batch: %w", i, err)
			}

		case "rdirectory":
			m := os.FileMode(0o755)
			if e.Mode != nil {
				m = os.FileMode(*e.Mode)
			}

			if err := os.MkdirAll(e.Path, m); err != nil {
				return fmt.Errorf("recursive write %d directory of batch: %w", i, err)
			}

		case "symlink":
			if err := os.Symlink(e.Data, e.Path); err != nil {
				return fmt.Errorf("write %d symlink of batch: %w", i, err)
			}

		default:
			return fmt.Errorf("write %d entry of batch: %w", i, ErrUnsupportedType)
		}

	}

	return nil
}
