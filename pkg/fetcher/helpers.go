package fetcher

import (
	"encoding/json"
	"io"
)

func DecodeFromJson[T any](input io.Reader, dest *T) error {
	decoder := json.NewDecoder(input)
	err := decoder.Decode(dest)

	if err != nil {
		return err
	}

	return nil
}

func ToPointer[T any](n T) *T {
	return &n
}
