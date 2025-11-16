package utils

import (
	"bytes"
	"encoding/base64"
	"io"
)

func DecodeBase64(src []byte) ([]byte, error) {
	decoder := base64.NewDecoder(base64.RawURLEncoding, bytes.NewReader(src))
	var dst bytes.Buffer

	_, err := io.Copy(&dst, decoder)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return dst.Bytes(), nil
}

func EncodeBase64(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.RawURLEncoding, &buf)

	_, err := encoder.Write(src)
	if err != nil {
		return nil, err
	}
	encoder.Close()

	return buf.Bytes(), nil
}
