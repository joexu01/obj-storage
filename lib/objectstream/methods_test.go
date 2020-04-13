package objectstream

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test method GET

func getHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		log.Println("error writing byte stream: ", err)
	}
}

func TestGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getHandler))
	defer s.Close()

	gs, _ := NewGetStream(s.URL[7:], "any")
	b, _ := ioutil.ReadAll(gs)
	if string(b) != "hello world" {
		t.Error(b)
	}
}

// Test method PUT

func putHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	if string(b) == "test" {
		return
	}
	w.WriteHeader(http.StatusForbidden)
}

func TestPut(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(putHandler))
	defer s.Close()

	ps := NewPutStream(s.URL[7:], "any")
	_, err := io.WriteString(ps, "test")
	if err != nil {
		t.Error(err)
		return
	}
	err = ps.Close()
	if err != nil {
		t.Error(err)
	}
}
