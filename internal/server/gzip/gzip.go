package gzip

import (
	"bytes"
	"compress/gzip"
)

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	gw.Close()
	return buf.Bytes(), nil
}
