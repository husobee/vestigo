// Portions Copyright 2015 Labstack.  All rights reserved.
// Portions Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.
package vestigo

import "net/http"

const (
	stype ntype = iota
	ptype
	mtype
)

type (
	ntype    uint8
	children []*node
)

// Router - The main vestigo router data structure
type Router struct {
	root *node
}

// NewRouter - Create a new vestigo router
func NewRouter() *Router {
	return &Router{
		root: &node{
			handler: new(handler),
		},
	}
}

// ServeHTTP - implementation of a http.Handler, making Router a http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := r.Find(req)
	h(w, req)
}

// Add - Add a method/handler combination to the router
func (r *Router) Add(method, path string, h http.HandlerFunc) {
	pnames := []string{} // Param names

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, stype, nil)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, ptype, pnames)
				return
			}
			r.insert(method, path[:i], nil, ptype, pnames)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, stype, nil)
			pnames = append(pnames, "_name")
			r.insert(method, path[:i+1], h, mtype, pnames)
			return
		}
	}

	r.insert(method, path, h, stype, pnames)
}

// Find - Find A route within the router tree
func (r *Router) Find(req *http.Request) (h http.HandlerFunc) {
	// get tree base node from the router
	cn := r.root

	h = NotFoundHandler

	if !validMethod(req.Method) {
		// if the method is completely invalid
		h = MethodNotAllowedHandler(cn.handler.allowedMethods)
		return
	}

	var (
		search = req.URL.Path
		c      *node  // Child node
		n      int    // Param counter
		nt     ntype  // Next type
		nn     *node  // Next node
		ns     string // Next search
	)

	// TODO: Check empty path???

	// Search order static > param > match-any
	for {
		if search == "" {
			if cn.handler != nil {
				// Found route, check if method is applicable
				theHandler, allowedMethods := cn.handler.GetMethodHandler(req.Method)
				if theHandler == nil {
					// route is valid, but method is not allowed, 405
					h = MethodNotAllowedHandler(allowedMethods)
					return
				}
				h = theHandler
			}
			return
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// Continue search
			search = search[l:]
		} else {
			cn = nn
			search = ns
			if nt == ptype {
				goto Param
			} else if nt == mtype {
				goto MatchAny
			} else {
				// Not found
				return
			}
		}

		if search == "" {
			// TODO: Needs improvement
			if cn.findChildWithType(mtype) == nil {
				continue
			}
			// Empty value
			goto MatchAny
		}

		// Static node
		c = cn.findChild(search[0], stype)
		if c != nil {
			// Save next
			if cn.label == '/' {
				nt = ptype
				nn = cn
				ns = search
			}
			cn = c
			continue
		}

		// Param node
	Param:
		c = cn.findChildWithType(ptype)
		if c != nil {
			// Save next
			if cn.label == '/' {
				nt = mtype
				nn = cn
				ns = search
			}
			cn = c
			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}

			registerVar(req, cn.pnames[n], search[:i])
			n++
			search = search[i:]
			continue
		}

		// Match-any node
	MatchAny:
		//		c = cn.getChild()
		c = cn.findChildWithType(mtype)
		if c != nil {
			cn = c
			registerVar(req, cn.pnames[len(cn.pnames)-1], search)
			search = "" // End search
			continue
		}

		// Not found
		return
	}
}

// insert - insert a route into the router tree
func (r *Router) insert(method, path string, h http.HandlerFunc, t ntype, pnames []string) {
	// Adjust max param

	cn := r.root

	if !validMethod(method) {
		panic("invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.typ = t
				cn.handler = new(handler)
				cn.handler.AddMethodHandler(method, h)
				cn.pnames = pnames
			}
		} else if l < pl {
			// Split node
			newHandler := new(handler)
			cn.handler.CopyTo(newHandler)
			n := newNode(cn.typ, cn.prefix[l:], cn, cn.children, newHandler, cn.pnames)

			// Reset parent node
			cn.typ = stype
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.handler = new(handler)
			cn.pnames = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.typ = t
				cn.handler.AddMethodHandler(method, h)
				cn.pnames = pnames
			} else {
				// Create child node
				newHandler := new(handler)
				newHandler.AddMethodHandler(method, h)
				n = newNode(t, search[l:], cn, nil, newHandler, pnames)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			newHandler := new(handler)
			newHandler.AddMethodHandler(method, h)
			n := newNode(t, search, cn, nil, newHandler, pnames)
			cn.addChild(n)
		} else {
			// Node already exists
			if h != nil {
				// add the handler to the node's map of methods to handlers
				cn.handler.AddMethodHandler(method, h)
				cn.pnames = pnames
			}
		}
		return
	}
}
