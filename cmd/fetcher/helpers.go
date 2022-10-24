package main

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

func Includes[T comparable](slice []T, f func(el T) bool) bool {
	for _, v := range slice {
		if f(v) {
			return true
		}
	}

	return false
}
