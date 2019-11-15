package proton

import (
	"encoding/base64"
	"fmt"
)

const (
	FilterBase64 FilterType = "base64"
)

type FilterType string

func Filter(b []byte, t FilterType) ([]byte, error) {
	switch t {
	case FilterBase64:
		return Base64Filter(b)
	}

	return nil, fmt.Errorf("proton: unknown filter: %s", t)
}

type FilterFunc func(b []byte) ([]byte, error)

func Base64Filter(b []byte) ([]byte, error) {
	out := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(out, b)
	return out, nil
}
