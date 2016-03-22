// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodNotAllowedDifferentMethodAllowed(t *testing.T) {
	router := NewRouter()
	path := "/test"
	router.Add("GET", path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("some return body"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "POST", path)
	}
	router.ServeHTTP(w, r)

	if w.Code != 405 || w.Body.String() != "Method Not Allowed" {
		t.Errorf("Invalid POST response, method: %s, path: %s, code: %d, body: %s", "POST", path, w.Code, w.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	router := NewRouter()
	path := "/test"
	router.Add("GET", path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("some return body"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("FAKEMETHOD", path, nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "GETFAKEMETHOD", path)
	}
	router.ServeHTTP(w, r)

	if w.Code != 405 || w.Body.String() != "Method Not Allowed" {
		t.Errorf("Invalid GETFAKEMETHOD response, method: %s, path: %s, code: %d, body: %s", "GETFAKEMETHOD", path, w.Code, w.Body.String())
	}
}

func TestEmptyBodyTrace(t *testing.T) {
	router := NewRouter()
	AllowTrace = true
	defer func() {
		AllowTrace = false
	}()
	path := "/test"
	router.Patch(path+"/split", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("TRACE", path+"/split", nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "TRACE", path)
	}
	router.ServeHTTP(w, r)
	if w.Code != 200 || w.Body.String() != "" || w.Header().Get("Content-Type") != "message/http" {
		t.Errorf("Invalid TRACE response, method: %s, path: %s, code: %d, body: %s", "TRACE", path, w.Code, w.Body.String())
	}

}

func TestTrace(t *testing.T) {
	router := NewRouter()
	AllowTrace = true
	defer func() {
		AllowTrace = false
	}()
	path := "/test"
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})
	router.Patch(path+"/split", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})
	router.Connect(path+"/split/again", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("TRACE", path+"/split", bytes.NewBufferString("awesome trace"))
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "TRACE", path)
	}
	router.ServeHTTP(w, r)
	if w.Code != 200 || w.Body.String() != "awesome trace" || w.Header().Get("Content-Type") != "message/http" {
		t.Errorf("Invalid TRACE response, method: %s, path: %s, code: %d, body: %s", "TRACE", path, w.Code, w.Body.String())
	}

}

func TestHead(t *testing.T) {
	router := NewRouter()
	path := "/test"
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-TestHeader", "true")
		w.WriteHeader(200)
		w.Write([]byte("some return body"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "HEAD", path)
	}
	router.ServeHTTP(w, r)
	if w.Code != 200 || w.Body.String() != "" || w.Header().Get("X-TestHeader") != "true" {
		t.Errorf("Invalid HEAD response, method: %s, path: %s, code: %d, body: %s", "HEAD", path, w.Code, w.Body.String())
	}
}

func TestNotFound(t *testing.T) {
	router := NewRouter()
	path := "/test"
	router.Add("GET", path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("some return body"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", path+"broken", nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "GET", path)
	}
	router.ServeHTTP(w, r)

	if w.Code != 404 || w.Body.String() != "Not Found" {
		t.Errorf("Invalid response, method: %s, path: %s, code: %d, body: %s", "GET", path, w.Code, w.Body.String())
	}
}

func TestCustomNotFound(t *testing.T) {
	router := NewRouter()
	path := "/test"
	router.Add("GET", path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("some return body"))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", path+"broken", nil)
	if err != nil {
		t.Errorf("Failed to create a new request, method: %s, path: %s", "GET", path)
	}
	CustomNotFoundHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("custom not found"))

	})
	router.ServeHTTP(w, r)

	if w.Code != 404 || w.Body.String() != "custom not found" {
		t.Errorf("Invalid response, method: %s, path: %s, code: %d, body: %s", "GET", path, w.Code, w.Body.String())
	}
}
