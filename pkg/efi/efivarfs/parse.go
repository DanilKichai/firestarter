package efivarfs

import (
	"encoding"
	"archshell/pkg/efi/common"
	"fmt"
	"os"
	"path/filepath"
)

func ParseVar[T encoding.BinaryUnmarshaler](efivars string, name string, guid string) (T, error) {
	data, err := os.ReadFile(filepath.Join(efivars, fmt.Sprintf("%s-%s", name, guid)))
	if err != nil {
		return common.Nil[T](), fmt.Errorf("read file: %w", err)
	}

	t := common.New[T]()
	if err = t.UnmarshalBinary(data); err != nil {
		return common.Nil[T](), fmt.Errorf("unmarshal data: %w", err)
	}

	return t, nil
}
