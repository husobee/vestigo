// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import (
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
		t.Errorf("Invalid POST response, method: %s, path: %s, code: %s, body: %s", "POST", path, w.Code, w.Body.String())
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
		t.Errorf("Invalid GETFAKEMETHOD response, method: %s, path: %s, code: %s, body: %s", "GETFAKEMETHOD", path, w.Code, w.Body.String())
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
		t.Errorf("Invalid response, method: %s, path: %s, code: %s, body: %s", "GET", path, w.Code, w.Body.String())
	}
}
