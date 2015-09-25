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
	root       *node
	globalCors *CorsAccessControl
}

// NewRouter - Create a new vestigo router
func NewRouter() *Router {
	return &Router{
		root: &node{
			resource: NewResource(),
		},
	}
}

// GlobalCors - Settings for Global Cors Options
func (r *Router) SetGlobalCors(c *CorsAccessControl) {
	r.globalCors = c
}

// SetCors - Set per resource Cors Policy
func (r *Router) SetCors(path string, c *CorsAccessControl) {
	r.AddWithCors("CORS", path, nil, c)
}

// ServeHTTP - implementation of a http.Handler, making Router a http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := r.Find(req)
	h(w, req)
}

// Add - Add a method/handler combination to the router
func (r *Router) AddWithCors(method, path string, h http.HandlerFunc, cors *CorsAccessControl) {
	r.add(method, path, h, cors)
}

// Add - Add a method/handler combination to the router
func (r *Router) Add(method, path string, h http.HandlerFunc) {
	r.add(method, path, h, nil)
}

// Add - Add a method/handler combination to the router
func (r *Router) add(method, path string, h http.HandlerFunc, cors *CorsAccessControl) {
	pnames := []string{} // Param names

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, stype, nil, cors)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, ptype, pnames, cors)
				return
			}
			r.insert(method, path[:i], nil, ptype, pnames, cors)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, stype, nil, cors)
			pnames = append(pnames, "_name")
			r.insert(method, path[:i+1], h, mtype, pnames, cors)
			return
		}
	}

	r.insert(method, path, h, stype, pnames, cors)
}

// Find - Find A route within the router tree
func (r *Router) Find(req *http.Request) (h http.HandlerFunc) {
	// get tree base node from the router
	cn := r.root

	h = NotFoundHandler

	if !validMethod(req.Method) {
		// if the method is completely invalid
		h = MethodNotAllowedHandler(cn.resource.allowedMethods)
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
			if cn.resource != nil {
				// Found route, check if method is applicable
				theHandler, allowedMethods := cn.resource.GetMethodHandler(req.Method)
				if theHandler == nil {
					if uint16(req.Method[0])<<8|uint16(req.Method[1]) == 0x4f50 {
						h = OptionsHandler(r.globalCors, cn.resource.Cors, allowedMethods)
						return
					}
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

			registerVar(req, cn.fmtpnames[n], search[:i])
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
			registerVar(req, cn.fmtpnames[len(cn.pnames)-1], search)
			search = "" // End search
			continue
		}

		// Not found
		return
	}
}

// insert - insert a route into the router tree
func (r *Router) insert(method, path string, h http.HandlerFunc, t ntype, pnames []string, cors *CorsAccessControl) {
	// Adjust max param

	cn := r.root

	if !validMethod(method) && method != "CORS" {
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
				cn.resource = NewResource()
				cn.resource.Cors = cn.resource.Cors.Merge(cors)
				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				cn.pnames = pnames
			}
		} else if l < pl {
			// Split node
			newResource := NewResource()
			cn.resource.CopyTo(newResource)
			n := newNode(cn.typ, cn.prefix[l:], cn, cn.children, newResource, cn.pnames)

			// Reset parent node
			cn.typ = stype
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.resource = NewResource()
			cn.pnames = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.typ = t
				cn.resource.Cors = cn.resource.Cors.Merge(cors)

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				cn.pnames = pnames
			} else {
				// Create child node
				newResource := NewResource()
				newResource.Cors = newResource.Cors.Merge(cors)
				if method != "CORS" {
					newResource.AddMethodHandler(method, h)
				}
				n = newNode(t, search[l:], cn, nil, newResource, pnames)
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
			newResource := NewResource()
			if method != "CORS" {
				newResource.AddMethodHandler(method, h)
			}
			newResource.Cors = newResource.Cors.Merge(cors)
			n := newNode(t, search, cn, nil, newResource, pnames)
			cn.addChild(n)
		} else {
			if cors != nil {
				cn.resource.Cors = cn.resource.Cors.Merge(cors)
			}
			// Node already exists
			if h != nil {
				// add the handler to the node's map of methods to handlers

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				cn.pnames = pnames
			}
		}
		return
	}
}
