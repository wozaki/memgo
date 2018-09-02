package memgo

import (
	"bytes"
	"compress/zlib"
)

func compress(value []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(value)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompress(value []byte) (string, error) {
	b := bytes.NewReader(value)
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	err = r.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
