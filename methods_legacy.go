// +build go1.5 go1.4 go1.3 go1.2 go1.1 go1.0

// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.

package vestigo

// methods - a list of methods that are allowed
var methods = map[string]bool{
	"CONNECT": true,
	"DELETE":  true,
	"GET":     true,
	"HEAD":    true,
	"OPTIONS": true,
	"PATCH":   true,
	"POST":    true,
	"PUT":     true,
	"TRACE":   true,
}

var (
	httpConnect = "CONNECT"
	httpDelete  = "DELETE"
	httpGet     = "GET"
	httpHead    = "HEAD"
	httpOptions = "OPTIONS"
	httpPatch   = "PATCH"
	httpPost    = "POST"
	httpPut     = "PUT"
	httpTrace   = "TRACE"
)
