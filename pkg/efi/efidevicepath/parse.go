package efidevicepath

import (
	"encoding"
	"archshell/pkg/efi/common"
	"fmt"
)

func ParsePath[T encoding.BinaryUnmarshaler](data []byte) (T, error) {
	t := common.New[T]()
	if err := t.UnmarshalBinary(data); err != nil {
		return common.Nil[T](), fmt.Errorf("unmarshal data: %w", err)
	}

	return t, nil
}
