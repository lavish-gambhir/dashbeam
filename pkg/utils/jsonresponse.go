package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func PrettyPrint(v any) (string, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ToJson(w io.Writer, sl any) error {
	return json.NewEncoder(w).Encode(sl)
}

func FromJson(r io.Reader, sl any) error {
	return json.NewDecoder(r).Decode(sl)
}
