package utils

import (
	"bytes"
	"io"
)

func ReadJSONBody(body io.ReadCloser) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)
	defer body.Close()
	if err != nil {
		Log.Error("read json body err: ", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadPlainTextBody(body io.ReadCloser) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)
	defer body.Close()
	if err != nil {
		Log.Error("read text body err: ", err)
		return "", err
	}
	return string(buf.Bytes()), nil
}
