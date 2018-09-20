package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type File struct {
	Type string
	Name string
}

type Server struct {
	Data map[File][]byte
	sync.Mutex

	AppendOnly bool
}

func NewServer(appendOnly bool) *Server {
	return &Server{
		Data:       make(map[File][]byte), // OMIT
		AppendOnly: appendOnly,
	}
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		status int
		file   File
	)

	path := strings.Split(req.URL.Path[1:], "/")
	if len(path) != 2 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, "invalid path, must be /{type}/{name}")
		return
	}

	file = File{Type: path[0], Name: path[1]}

	switch req.Method {
	case "GET":
		status, err = s.Get(file, res)
	case "POST":
		status, err = s.Post(file, res, req.Body)
	case "DELETE":
		status, err = s.Delete(file, res)
	default:
		err = errors.New("method not allowed")
	}

	if req.Body != nil {
		req.Body.Close()
	}

	if err != nil {
		if status == 0 {
			status = http.StatusInternalServerError
		}
		res.WriteHeader(status)
		fmt.Fprintln(res, err)
		return
	}

	if status != 0 {
		res.WriteHeader(status)
	}
}

func (s *Server) Get(file File, res http.ResponseWriter) (status int, err error) {
	s.Lock()
	content, ok := s.Data[file]
	s.Unlock()

	if !ok {
		return http.StatusNotFound, nil
	}

	res.Header().Add("content-type", "application/octet-stream")

	res.WriteHeader(http.StatusOK)
	_, err = res.Write(content)
	if err != nil {
		return http.StatusBadGateway, err
	}

	return 0, nil
}

func (s *Server) Post(file File, res http.ResponseWriter, body io.Reader) (status int, err error) {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	s.Lock()
	_, ok := s.Data[file]
	if ok {
		s.Unlock()
		return http.StatusBadRequest, errors.New("file already exists")
	}
	s.Data[file] = buf

	s.Unlock()
	return http.StatusCreated, nil
}

func (s *Server) Delete(file File, res http.ResponseWriter) (status int, err error) {
	if file.Type != "lock" && s.AppendOnly {
		return http.StatusMethodNotAllowed, errors.New("server is in append-only mode")
	}

	s.Lock()
	_, ok := s.Data[file]
	if !ok {
		s.Unlock()
		return http.StatusNotFound, nil
	}
	delete(s.Data, file)
	s.Unlock()

	return http.StatusOK, nil
}

func main() {
	var appendOnly bool
	if len(os.Args) > 1 && os.Args[1] == "--append-only" {
		appendOnly = true
	}

	addr := "localhost:1234"
	fmt.Printf("listen on %v (append-only mode: %v)\n", addr, appendOnly)

	http.ListenAndServe(addr, NewServer(appendOnly))
}
