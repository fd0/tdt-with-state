package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTables(t *testing.T) {
	type Requests []struct {
		Request    *http.Request
		StatusCode int
		Body       string
	}

	var tests = []struct {
		AppendOnly bool
		Requests   Requests
	}{
		{
			AppendOnly: false,
			Requests: Requests{
				{
					Request:    NewRequest(t, "/foo/bar", "POST", "test content"),
					StatusCode: http.StatusCreated,
				}, {
					Request:    NewRequest(t, "/foo/bar", "GET", ""),
					StatusCode: http.StatusOK,
					Body:       "test content",
				}, {
					Request:    NewRequest(t, "/foo/bar", "DELETE", ""),
					StatusCode: http.StatusOK,
				}, {
					Request:    NewRequest(t, "/foo/bar", "GET", ""),
					StatusCode: http.StatusNotFound,
				},
			},
		},
		{
			AppendOnly: true,
			Requests: Requests{
				{
					Request:    NewRequest(t, "/foo/bar", "POST", "test content"),
					StatusCode: http.StatusCreated,
				}, {
					Request:    NewRequest(t, "/foo/bar", "GET", ""),
					StatusCode: http.StatusOK,
					Body:       "test content",
				}, {
					Request:    NewRequest(t, "/foo/bar", "DELETE", ""),
					StatusCode: http.StatusMethodNotAllowed,
				}, {
					Request:    NewRequest(t, "/foo/bar", "GET", ""),
					StatusCode: http.StatusOK,
					Body:       "test content",
				},
			},
		},
		{
			AppendOnly: true,
			Requests: Requests{
				{
					Request:    NewRequest(t, "/lock/bar", "POST", "test content"),
					StatusCode: http.StatusCreated,
				}, {
					Request:    NewRequest(t, "/lock/bar", "GET", ""),
					StatusCode: http.StatusOK,
					Body:       "test content",
				}, {
					Request:    NewRequest(t, "/lock/bar", "DELETE", ""),
					StatusCode: http.StatusOK,
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
				if request.StatusCode != res.Code {
					st.Errorf("%v %v wrong response code, want %v, got %v",
						request.Request.Method, request.Request.URL, request.StatusCode, res.Code)
				}

				if request.Body != "" && res.Body.String() != request.Body {
					st.Errorf("%v %v wrong body returned for request: want %q, got %q",
						request.Request.Method, request.Request.URL, request.Body, res.Body.String())
				}
			}
		})
	}
}
