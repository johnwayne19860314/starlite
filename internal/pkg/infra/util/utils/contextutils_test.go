package utils

import (
	"context"
	"net/http"
	"testing"
)

func TestNewServiceContextFromHTTPRequest(t *testing.T) {
	const requestId = "request id"
	req := http.Request{}
	req.Header = http.Header{}
	req.Header.Add("x-txid", requestId)
	ctx := NewServiceContextWithValue(context.Background(), RequestIDCtxID, GetRequestIdByKey(&req, "x-txid"))
	if rId := GetRequestIdFromContext(ctx); rId != requestId {
		t.Fatal(rId, "does not match expected", requestId)
	}
	header := ServiceHeadersFromContext(ctx, http.Header{})
	if val := header.Get("x-txid"); val != requestId {
		t.Fatal(val, "does not match expected", requestId)
	}
}
