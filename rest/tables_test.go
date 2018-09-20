package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewRequest(t *testing.T, path, method, body string) *http.Request {
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

// START INTRO OMIT
func TestTables(t *testing.T) {
	type Requests []struct {
		Request *http.Request
		Code    int
		Body    string
	}

	var tests = []struct {
		AppendOnly bool
		Requests   Requests
	}{
		{
			Requests: Requests{
				{
					Request: NewRequest(t, "/text/test.txt", "POST", "test content"),
					Code:    http.StatusCreated,
				}, {
					Request: NewRequest(t, "/text/test.txt", "GET", ""),
					Code:    http.StatusOK,
					Body:    "test content",
				},
			},
		},
		// END INTRO OMIT
		// START REGULAR OMIT
		{
			Requests: Requests{
				{
					Request: NewRequest(t, "/text/test.txt", "POST", "test content"),
					Code:    http.StatusCreated,
				}, {
					Request: NewRequest(t, "/text/test.txt", "GET", ""),
					Code:    http.StatusOK,
					Body:    "test content",
				}, {
					Request: NewRequest(t, "/text/test.txt", "DELETE", ""),
					Code:    http.StatusOK,
				}, {
					Request: NewRequest(t, "/text/test.txt", "GET", ""),
					Code:    http.StatusNotFound,
				},
			},
		},
		// END OMIT
		// START APPEND1 OMIT
		{
			AppendOnly: true, // HL
			Requests: Requests{
				{
					Request: NewRequest(t, "/text/test.txt", "POST", "test content"),
					Code:    http.StatusCreated,
				}, {
					Request: NewRequest(t, "/text/test.txt", "GET", ""),
					Code:    http.StatusOK,
					Body:    "test content",
				}, {
					Request: NewRequest(t, "/text/test.txt", "DELETE", ""),
					Code:    http.StatusMethodNotAllowed, // HL
				}, {
					Request: NewRequest(t, "/text/test.txt", "GET", ""),
					Code:    http.StatusOK,
					Body:    "test content",
				},
			},
		},
		// END OMIT
		// START APPEND2 OMIT
		{
			AppendOnly: true,
			Requests: Requests{
				{
					Request: NewRequest(t, "/lock/bar", "POST", "test content"),
					Code:    http.StatusCreated,
				}, {
					Request: NewRequest(t, "/lock/bar", "GET", ""),
					Code:    http.StatusOK,
					Body:    "test content",
				}, {
					Request: NewRequest(t, "/lock/bar", "DELETE", ""),
					Code:    http.StatusOK, // HL
				},
			},
		},
		// END OMIT
	}

	// START FUNC OMIT
	for _, test := range tests {
		t.Run("", func(st *testing.T) {
			// use a new server for each sub-test
			srv := NewServer(test.AppendOnly)

			for _, r := range test.Requests {
				res := httptest.NewRecorder()
				srv.ServeHTTP(res, r.Request)

				if r.Code != res.Code {
					st.Errorf("%v %v wrong response code, want %v, got %v",
						r.Request.Method, r.Request.URL,
						http.StatusText(r.Code), http.StatusText(res.Code))
				}

				if r.Body != "" && res.Body.String() != r.Body {
					st.Errorf("%v %v wrong body returned for request: want %q, got %q",
						r.Request.Method, r.Request.URL, r.Body, res.Body.String())
				}
			}
		})
	}
	// END FUNC OMIT
}
