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
		t.Helper() // OMIT
		if res.Code != code {
			t.Errorf("wrong response code, want %v, got %v",
				http.StatusText(code), http.StatusText(res.Code))
		}
	}
}

func WantBody(body string) CheckFunc {
	return func(t *testing.T, res *httptest.ResponseRecorder) {
		t.Helper() // OMIT
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
			t.Helper() // OMIT
			t.Errorf("HTTP response header %v has wrong value: want %q, got %q",
				name,
				value, res.Header().Get(name))
		}
	}
}

// START INTRO OMIT
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
			Requests: Requests{
				{
					NewRequest(t, "/text/file.txt", "POST", "test content"),
					[]CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantBody("test content"),
					},
				},
			},
		},
		// END INTRO OMIT
		// START REGULAR OMIT
		{
			Requests: Requests{
				{
					NewRequest(t, "/text/file.txt", "POST", "test content"),
					[]CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantBody("test content"),
					},
				}, {
					NewRequest(t, "/text/file.txt", "DELETE", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusNotFound),
					},
				},
			},
		},
		// END OMIT
		// START APPEND1 OMIT
		{
			AppendOnly: true, // HL
			Requests: Requests{
				{
					NewRequest(t, "/text/file.txt", "POST", "test content"),
					[]CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantBody("test content"),
					},
				}, {
					NewRequest(t, "/text/file.txt", "DELETE", ""),
					[]CheckFunc{
						WantStatus(http.StatusMethodNotAllowed), // HL
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantBody("test content"),
					},
				},
			},
		},
		// END OMIT
		// START APPEND2 OMIT
		{
			AppendOnly: true,
			Requests: Requests{
				{
					NewRequest(t, "/lock/bar", "POST", "test content"),
					[]CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					NewRequest(t, "/lock/bar", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantBody("test content"),
					},
				}, {
					NewRequest(t, "/lock/bar", "DELETE", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK), // HL
					},
				},
			},
		},
		// END OMIT
		// START HEADER OMIT
		{
			Requests: Requests{
				{
					NewRequest(t, "/text/file.txt", "POST", "test content"),
					[]CheckFunc{
						WantStatus(http.StatusCreated),
					},
				}, {
					NewRequest(t, "/text/file.txt", "GET", ""),
					[]CheckFunc{
						WantStatus(http.StatusOK),
						WantHeader("Content-Type", "application/octet-stream"), // HL
						WantBody("test content"),
					},
				},
			},
		},
		// END OMIT
	}

	// START FUNC OMIT
	for _, test := range tests {
		t.Run("", func(st *testing.T) {
			srv := NewServer(test.AppendOnly)

			for _, r := range test.Requests {
				// execute a single request
				res := httptest.NewRecorder()
				srv.ServeHTTP(res, r.Request)

				// run all checks on the response
				for _, fn := range r.Checks {
					st.Run("", func(sst *testing.T) {
						fn(sst, res)
					})
				}
			}
		})
	}
	// END FUNC OMIT
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
		CreateRandomFile(t, "text"),
		CreateRandomFile(t, "lock"),
	}

	for _, seq := range tests {
		t.Run("", func(st *testing.T) {
			srv := NewServer(false)

			for _, request := range seq {
				// execute a single request
				res := httptest.NewRecorder()
				srv.ServeHTTP(res, request.Request)

				// run all checks on the response
				for _, fn := range request.Checks {
					st.Run("", func(sst *testing.T) {
						fn(st, res)
					})
				}
			}
		})
	}
}
