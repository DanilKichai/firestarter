package unyaml

import (
	"fmt"
	"path/filepath"
)

type UnYAML []struct {
	Path  string `yaml:"path"`
	Type  string `yaml:"type"`
	User  string `yaml:"user"`
	Group string `yaml:"group"`
	Mode  string `yaml:"mode"`
	Data  string `yaml:"data"`
}

func (u *UnYAML) Write(destination string, allowAbsolutePath bool, allowFileOverride bool, allowDirectoryExists bool) error {
	for i, e := range *u {
		path := filepath.Clean(e.Path)

		if path == "" {
			return fmt.Errorf("extract %d entry of YAML: %w", i, ErrEmptyPath)
		}

		if !allowAbsolutePath || !filepath.IsAbs(path) {
			path = filepath.Join(destination, path)
		}

		if len(e.Data) == 0 && e.Type == "symlink" {
			return fmt.Errorf("extract %d entry of YAML: %w", i, ErrEmptyData)
		}

		/*
			user :=
			group :=
			mode :=
		*/

		switch e.Type {
		case "file":
			//

		case "catalog":
			//

		case "symlink":
			//

		default:
			return fmt.Errorf("extract %d entry of YAML: %w", i, ErrTypeRepresentation)
		}
	}

	return nil
}
