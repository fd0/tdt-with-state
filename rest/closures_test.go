package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// CheckFunc checks an HTTP response and marks the test as errored if something
// is incorrect.
type CheckFunc func(*testing.T, *httptest.ResponseRecorder)

func WantStatus(code int) CheckFunc {
	return func(t *testing.T, res *httptest.ResponseRecorder) {
		if res.Code != code {
			t.Errorf("wrong response code, want %v, got %v", code, res.Code)
		}
	}
}

func WantBody(body string) CheckFunc {
	return func(t *testing.T, res *httptest.ResponseRecorder) {
		if res.Body.String() != body {
			t.Errorf("wrong body returned for request: want %q, got %q", body, res.Body.String())
		}
	}
}

func WantHeader(name, value string) CheckFunc {
	return func(t *testing.T, res *httptest.ResponseRecorder) {
		if _, ok := res.Header()[name]; !ok {
			t.Errorf("expected HTTP response header %q not found", name)
			return
		}

		if res.Header().Get(name) != value {
			t.Errorf("HTTP response header %v has wrong value: want %q, got %q",
				name, value, res.Header().Get(name))
		}
	}
}

func TestClosures(t *testing.T) {
	type Requests []struct {
		Request *http.Request
		Checks  []CheckFunc
	}

	var tests = []struct {
		AppendOnly bool
		Requests   Requests
	}{
		{
			AppendOnly: false,
			Requests: Requests{
				{
					Request: NewRequest(t, "/foo/bar", "POST", "test content"),
					Checks: []CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					Request: NewRequest(t, "/foo/bar", "GET", ""),
					Checks: []CheckFunc{
						WantStatus(http.StatusOK),
						WantHeader("Content-Type", "application/octet-stream"),
						WantBody("test content"),
					},
				}, {
					Request: NewRequest(t, "/foo/bar", "DELETE", ""),
					Checks: []CheckFunc{
						WantStatus(http.StatusOK),
					},
				}, {
					Request: NewRequest(t, "/foo/bar", "GET", ""),
					Checks: []CheckFunc{
						WantStatus(http.StatusNotFound),
					},
				},
			},
		},
		{
			AppendOnly: true,
			Requests: Requests{
				{
					Request: NewRequest(t, "/locks/bar", "POST", "test content"),
					Checks: []CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					Request: NewRequest(t, "/locks/bar", "GET", ""),
					Checks: []CheckFunc{
						WantStatus(http.StatusOK),
						WantHeader("Content-Type", "application/octet-stream"),
						WantBody("test content"),
					},
				}, {
					Request: NewRequest(t, "/locks/bar", "DELETE", ""),
					Checks: []CheckFunc{
						WantStatus(http.StatusMethodNotAllowed),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(st *testing.T) {
			// use a new server for each sub-test
			srv := NewServer(test.AppendOnly)

			for _, request := range test.Requests {
				// execute a single request
				res := httptest.NewRecorder()
				srv.ServeHTTP(res, request.Request)

				// run all checks on the response
				for _, fn := range request.Checks {
					fn(st, res)
				}
			}
		})
	}
}

// RandomFileName returns a random string with 10 to 20 characters.
func RandomFileName(t *testing.T) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, rand.Intn(20)+10)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(buf)
}

type TestSequence []struct {
	Request *http.Request
	Checks  []CheckFunc
}

func CreateRandomFile(t *testing.T, tpe string) TestSequence {
	path := fmt.Sprintf("/%v/%v", tpe, RandomFileName(t))
	return TestSequence{
		{
			Request: NewRequest(t, path, "POST", "test content"),
			Checks: []CheckFunc{
				WantStatus(http.StatusCreated),
			},
		}, {
			Request: NewRequest(t, path, "GET", ""),
			Checks: []CheckFunc{
				WantStatus(http.StatusOK),
				WantHeader("Content-Type", "application/octet-stream"),
				WantBody("test content"),
			},
		}, {
			Request: NewRequest(t, path, "DELETE", ""),
			Checks: []CheckFunc{
				WantStatus(http.StatusOK),
			},
		}, {
			Request: NewRequest(t, path, "GET", ""),
			Checks: []CheckFunc{
				WantStatus(http.StatusNotFound),
			},
		},
	}
}

func TestMoreClosures(t *testing.T) {
	var tests = []TestSequence{
		CreateRandomFile(t, "foo"),
		CreateRandomFile(t, "data"),
		CreateRandomFile(t, "locks"),
	}

	for _, seq := range tests {
		t.Run("", func(st *testing.T) {
			// use a new server for each sub-test
			srv := NewServer(false)

			for _, request := range seq {
				// execute a single request
				res := httptest.NewRecorder()
				srv.ServeHTTP(res, request.Request)

				// run all checks on the response
				for _, fn := range request.Checks {
					fn(st, res)
				}
			}
		})
	}
}
