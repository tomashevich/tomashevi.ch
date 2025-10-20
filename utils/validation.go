package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func UnmarshalJSON(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	return decoder.Decode(&v)
}

func UnmarshalJSONString(rawJson string, v any) error {
	decoder := json.NewDecoder(bytes.NewReader([]byte(rawJson)))
	decoder.DisallowUnknownFields()

	return decoder.Decode(&v)
}
