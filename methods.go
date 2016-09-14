// +build !go1.5,!go1.4,!go1.3,!go1.2,!go1.1,!go1.0

// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

import "net/http"

// methods - a list of methods that are allowed
var methods = map[string]bool{
	http.MethodConnect: true,
	http.MethodDelete:  true,
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodPatch:   true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodTrace:   true,
}

var (
	httpConnect = http.MethodConnect
	httpDelete  = http.MethodDelete
	httpGet     = http.MethodGet
	httpHead    = http.MethodHead
	httpOptions = http.MethodOptions
	httpPatch   = http.MethodPatch
	httpPost    = http.MethodPost
	httpPut     = http.MethodPut
	httpTrace   = http.MethodTrace
)
