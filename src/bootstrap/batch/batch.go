package batch

import (
	"fmt"
	"os"
)

type Batch []struct {
	Path string `yaml:"path"`
	Mode *int   `yaml:"mode"`
	Data string `yaml:"data"`
}

func (u *Batch) Write() error {
	for i, e := range *u {
		m := os.FileMode(0o644)
		if e.Mode != nil {
			m = os.FileMode(*e.Mode)
		}

		if err := os.WriteFile(e.Path, []byte(e.Data), m); err != nil {
			return fmt.Errorf("write %d file of batch: %w", i, err)
		}

	}

	return nil
}
