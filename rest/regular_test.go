package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewRequest(t *testing.T, path, method, body string) *http.Request {
	var rd io.Reader

	if body != "" {
		rd = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, path, rd)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func DoRequest(t *testing.T, srv http.Handler, path, method, body string) *httptest.ResponseRecorder {
	req := NewRequest(t, path, method, body)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func CheckStatus(t *testing.T, res *httptest.ResponseRecorder, code int) {
	if res.Code != code {
		t.Errorf("wrong response code, want %v, got %v", code, res.Code)
	}
}

func CheckBody(t *testing.T, res *httptest.ResponseRecorder, body string) {
	if res.Body.String() != body {
		t.Errorf("wrong body returned for request: want %q, got %q", body, res.Body.String())
	}
}

func TestCreateDeleteFile(t *testing.T) {
	srv := NewServer(false)
	filepath := fmt.Sprintf("%s/%s", "foo", "bar")

	res := DoRequest(t, srv, filepath, "POST", "test content")
	CheckStatus(t, res, http.StatusCreated)
	CheckBody(t, res, "")

	res = DoRequest(t, srv, filepath, "GET", "")
	CheckStatus(t, res, http.StatusOK)
	CheckBody(t, res, "test content")

	res = DoRequest(t, srv, filepath, "DELETE", "")
	CheckStatus(t, res, http.StatusOK)
	CheckBody(t, res, "")

	res = DoRequest(t, srv, filepath, "GET", "")
	CheckStatus(t, res, http.StatusNotFound)
	CheckBody(t, res, "")
}
