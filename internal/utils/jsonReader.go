package utils

import (
	"bytes"
	"net/http"
)

func ReadJSONBody(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		Log.Error("read json body err: ", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
