package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func DoRequest(t *testing.T, srv *Server, path, method, body string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	return res
}

func CheckStatus(t *testing.T, res *httptest.ResponseRecorder, code int) {
	if res.Code != code {
		t.Errorf("wrong response code, want %v, got %v",
			http.StatusText(code), http.StatusText(res.Code))
	}
}

func CheckBody(t *testing.T, res *httptest.ResponseRecorder, body string) {
	if res.Body.String() != body {
		t.Errorf("wrong body returned for request: want %q, got %q", body, res.Body.String())
	}
}

func TestCreateDeleteFile(t *testing.T) {
	srv := NewServer(false)
	filepath := "/text/file.txt"

	res := DoRequest(t, srv, filepath, "POST", "test content")
	CheckStatus(t, res, http.StatusCreated)

	res = DoRequest(t, srv, filepath, "GET", "")
	CheckStatus(t, res, http.StatusOK)
	CheckBody(t, res, "test content")

	res = DoRequest(t, srv, filepath, "DELETE", "")
	CheckStatus(t, res, http.StatusOK)

	res = DoRequest(t, srv, filepath, "GET", "")
	CheckStatus(t, res, http.StatusNotFound)
}

func TestAppendOnlyCreateDeleteFile(t *testing.T) {
	srv := NewServer(true) // HL
	filepath := "/text/file.txt"

	res := DoRequest(t, srv, filepath, "POST", "test content")
	CheckStatus(t, res, http.StatusCreated)

	res = DoRequest(t, srv, filepath, "GET", "")
	CheckStatus(t, res, http.StatusOK)
	CheckBody(t, res, "test content")

	res = DoRequest(t, srv, filepath, "DELETE", "")
	CheckStatus(t, res, http.StatusMethodNotAllowed) // HL
}
