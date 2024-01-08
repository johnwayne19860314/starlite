package utils

import (
	"io"
	"net/http"
)

func CloseAndDiscardRespBody(resp *http.Response) {
	if resp != nil {
		defer resp.Body.Close()
		_, err := io.Copy(io.Discard, resp.Body)
		if err != nil {
			return
		}
	}
}
