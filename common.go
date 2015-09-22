// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.
package vestigo

import (
	"net/http"
	"strings"
)

// methods - a list of methods that are allowed
var methods = []string{
	"CONNECT",
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"PATCH",
	"POST",
	"PUT",
	"TRACE",
}

// Get - Helper method to add HTTP GET Method to router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.Add("GET", path, handler)
}

// Post - Helper method to add HTTP POST Method to router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.Add("POST", path, handler)
}

// Connect - Helper method to add HTTP CONNECT Method to router
func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.Add("CONNECT", path, handler)
}

// Delete - Helper method to add HTTP DELETE Method to router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.Add("DELETE", path, handler)
}

// Head - Helper method to add HTTP HEAD Method to router
func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.Add("HEAD", path, handler)
}

// Options - Helper method to add HTTP OPTIONS Method to router
func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.Add("OPTIONS", path, handler)
}

// Patch - Helper method to add HTTP PATCH Method to router
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.Add("PATCH", path, handler)
}

// Put - Helper method to add HTTP PUT Method to router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.Add("PUT", path, handler)
}

// Trace - Helper method to add HTTP TRACE Method to router
func (r *Router) Trace(path string, handler http.HandlerFunc) {
	r.Add("TRACE", path, handler)
}

// Param - Get a url parameter by name
func Param(r *http.Request, name string) string {
	return r.FormValue(":" + name)
}

//validMethod - validate that the http method is valid.
func validMethod(method string) bool {
	var ok = false
	for _, v := range methods {
		if v == method {
			ok = true
			break
		}
	}
	return ok
}

var (
	// MethodNotAllowedHandler - Generic Handler to handle when method isn't allowed for a resource
	MethodNotAllowedHandler = func(allowedMethods ...string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", strings.Join(allowedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
		}
	}
	// NotFoundHandler - Generic Handler to handle when resource isn't found
	NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}
)

// registerVar - Put the URL Parameter into the request's RawQuery, very PAT like.
func registerVar(r *http.Request, pname string, pvalue string) {
	if r.URL.RawQuery != "" {
		r.URL.RawQuery += "&" + ":" + pname + "=" + pvalue
	} else {
		r.URL.RawQuery += ":" + pname + "=" + pvalue
	}
}
