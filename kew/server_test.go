package kew

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(&Config{DbPath: "./server1"})
	defer s.Store.Drop()

	// Enqueue
	body := strings.NewReader("bar")
	req, _ := http.NewRequest("POST", "/foo", body)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 201)

	// Invalid complete (must dequeue first)
	req, _ = http.NewRequest("DELETE", "/foo/1", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 400)

	// Info
	req, _ = http.NewRequest("GET", "/foo/1", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)

	// Invalid dequeue (still being enqueued)
	req, _ = http.NewRequest("GET", "/foo/head", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)

	// Dequeue
	req, _ = http.NewRequest("GET", "/foo/head?wait=30", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)

	// Complete
	req, _ = http.NewRequest("DELETE", "/foo/1", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)

	// Info not found
	req, _ = http.NewRequest("GET", "/foo/1", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)
}

func TestServerAuth(t *testing.T) {
	s := NewServer(&Config{DbPath: "./server2", Auth: "secret"})
	defer s.Store.Drop()

	body := strings.NewReader("bar")

	req, _ := http.NewRequest("POST", "/foo", body)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 401)

	req, _ = http.NewRequest("POST", "/foo", body)
	req.SetBasicAuth("", "secret")
	w = httptest.NewRecorder()
	s.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 201)
}