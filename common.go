// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.
package vestigo

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
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

// GetWithCors - Helper method to add HTTP GET Method to router
func (r *Router) GetWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("GET", path, handler, cors)
}

// Get - Helper method to add HTTP GET Method to router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.Add("GET", path, handler)
}

// PostWithCors - Helper method to add HTTP POST Method to router
func (r *Router) PostWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("POST", path, handler, cors)
}

// Post - Helper method to add HTTP POST Method to router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.Add("POST", path, handler)
}

// ConnectWithCors - Helper method to add HTTP CONNECT Method to router
func (r *Router) ConnectWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("CONNECT", path, handler, cors)
}

// Connect - Helper method to add HTTP CONNECT Method to router
func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.Add("CONNECT", path, handler)
}

// DeleteWithCors - Helper method to add HTTP DELETE Method to router
func (r *Router) DeleteWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("DELETE", path, handler, cors)
}

// Delete - Helper method to add HTTP DELETE Method to router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.Add("DELETE", path, handler)
}

// HeadWithCors - Helper method to add HTTP HEAD Method to router
func (r *Router) HeadWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("HEAD", path, handler, cors)
}

// Head - Helper method to add HTTP HEAD Method to router
func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.Add("HEAD", path, handler)
}

// PatchWithCors - Helper method to add HTTP PATCH Method to router
func (r *Router) PatchWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("PATCH", path, handler, cors)
}

// Patch - Helper method to add HTTP PATCH Method to router
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.Add("PATCH", path, handler)
}

// PutWithCors - Helper method to add HTTP PUT Method to router
func (r *Router) PutWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("PUT", path, handler, cors)
}

// Put - Helper method to add HTTP PUT Method to router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.Add("PUT", path, handler)
}

// TraceWithCors - Helper method to add HTTP TRACE Method to router
func (r *Router) TraceWithCors(path string, handler http.HandlerFunc, cors CorsOptionsInterface) {
	r.AddWithCors("TRACE", path, handler, cors)
}

// Trace - Helper method to add HTTP TRACE Method to router
func (r *Router) Trace(path string, handler http.HandlerFunc) {
	r.Add("TRACE", path, handler)
}

// Param - Get a url parameter by name
func Param(r *http.Request, name string) string {
	return r.FormValue(":" + name)
}

// ParamNames - Get a url parameter name list
func ParamNames(r *http.Request) []string {
	r.ParseForm()
	names := []string{}
	for k, _ := range r.Form {
		names = append(names, k)
	}
	return names
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

// CorsOptionsInterface - Interface which defines what a CORS Option
// should be able to perform
type CorsOptionsInterface interface {
	GetAllowOrigin() []string
	GetAllowCredentials() bool
	GetExposeHeaders() []string
	GetMaxAge() time.Duration
	GetAllowMethods() []string
	GetAllowHeaders() []string
	Merge(CorsOptionsInterface) CorsOptionsInterface
}

// CorsAccessControl - Default implementation of CorsOptionsInterface
type CorsAccessControl struct {
	AllowOrigin      []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           time.Duration
	AllowMethods     []string
	AllowHeaders     []string
}

// AllowOrigin - returns the allow-origin string representation
func (c *CorsAccessControl) GetAllowOrigin() []string {
	return c.AllowOrigin
}

// AllowCredentials - returns the allow-credentials string representation
func (c *CorsAccessControl) GetAllowCredentials() bool {
	return c.AllowCredentials
}

// ExposeHeaders - returns the expose-headers string representation
func (c *CorsAccessControl) GetExposeHeaders() []string {
	return c.ExposeHeaders
}

// MaxAge - returns the max-age string representation
func (c *CorsAccessControl) GetMaxAge() time.Duration {
	return c.MaxAge
}

// AllowMethods - returns the allow-methods string representation
func (c *CorsAccessControl) GetAllowMethods() []string {
	return c.AllowMethods
}

// AllowHeaders - returns the allow-headers string representation
func (c *CorsAccessControl) GetAllowHeaders() []string {
	return c.AllowHeaders
}

func (c *CorsAccessControl) Merge(c2 CorsOptionsInterface) CorsOptionsInterface {
	result := new(CorsAccessControl)
	if c2 == nil {
		result.AllowOrigin = c.GetAllowOrigin()
		result.AllowCredentials = c.GetAllowCredentials()
		result.ExposeHeaders = c.GetExposeHeaders()
		result.MaxAge = c.GetMaxAge()
		result.AllowMethods = c.GetAllowMethods()
		result.AllowHeaders = c.GetAllowHeaders()
		return result
	}

	if allowOrigin := c2.GetAllowOrigin(); len(allowOrigin) != 0 {
		result.AllowOrigin = append(c.GetAllowOrigin(), c2.GetAllowOrigin()...)
	} else {
		result.AllowOrigin = c.GetAllowOrigin()
	}
	if allowCredentials := c2.GetAllowCredentials(); allowCredentials == true {
		result.AllowCredentials = c2.GetAllowCredentials()
	} else {
		result.AllowCredentials = c.GetAllowCredentials()
	}
	if exposeHeaders := c2.GetExposeHeaders(); len(exposeHeaders) != 0 {
		result.ExposeHeaders = append(c.GetExposeHeaders(), c2.GetExposeHeaders()...)
	} else {
		result.ExposeHeaders = c.GetExposeHeaders()
	}
	if maxAge := c2.GetMaxAge(); maxAge.Seconds() != 0 {
		result.MaxAge = c2.GetMaxAge()
	} else {
		result.MaxAge = c.GetMaxAge()
	}
	if allowMethods := c2.GetAllowMethods(); len(allowMethods) != 0 {
		result.AllowMethods = append(c.GetAllowMethods(), c2.GetAllowMethods()...)
	} else {
		result.AllowMethods = c.GetAllowMethods()
	}
	if allowHeaders := c2.GetAllowHeaders(); len(allowHeaders) != 0 {
		result.AllowHeaders = append(c.GetAllowHeaders(), c2.GetAllowHeaders()...)
	} else {
		result.AllowHeaders = c.GetAllowHeaders()
	}
	fmt.Println("result: ", result)

	return result
}

func CorsPreflight(gcors CorsOptionsInterface, lcors CorsOptionsInterface, allowedMethods string, w http.ResponseWriter, r *http.Request) error {

	cors := gcors.Merge(lcors)

	if origin := r.Header.Get("Origin"); cors != nil && origin != "" {
		// validate origin is in list of acceptable allow-origins
		allowedOrigin := false
		allowedOriginExact := false
		for _, v := range cors.GetAllowOrigin() {
			if v == origin {
				w.Header().Add("Access-Control-Allow-Origin", origin)
				allowedOriginExact = true
				allowedOrigin = true
				break
			}
		}
		if !allowedOrigin {
			for _, v := range cors.GetAllowOrigin() {
				if v == "*" {
					w.Header().Add("Access-Control-Allow-Origin", v)
					allowedOrigin = true
					break
				}
			}
		}

		if !allowedOrigin {
			// other option headers needed
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
			return errors.New("quick cors end")

		}

		// if the request includes access-control-request-method
		if method := r.Header.Get("Access-Control-Request-Method"); method != "" {
			// if there are no cors settings for this resource, use the allowedMethods,
			// if there are settings for cors, use those
			responseMethods := []string{}
			if methods := cors.GetAllowMethods(); len(methods) != 0 {
				for _, x := range methods {
					if x == method {
						responseMethods = append(responseMethods, x)
					}
				}
			} else {
				for _, x := range strings.Split(allowedMethods, ", ") {
					if x == method {
						responseMethods = append(responseMethods, x)
					}
				}
			}
			if len(responseMethods) > 0 {
				w.Header().Add("Access-Control-Allow-Methods", strings.Join(responseMethods, ", "))
			} else {
				// other option headers needed
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
				return errors.New("quick cors end")
			}
		}

		// if allow credentials is allowed on this resource respond with true
		if allowCredentials := cors.GetAllowCredentials(); allowedOriginExact && allowCredentials {
			w.Header().Add("Access-Control-Allow-Credentials", "true")
		}

		if exposeHeaders := cors.GetExposeHeaders(); len(exposeHeaders) != 0 {
			// if we have expose headers, send them
			w.Header().Add("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ", "))
		}
		if maxAge := cors.GetMaxAge(); maxAge.Seconds() != 0 {
			// optional, if we have a max age, send it
			sec := fmt.Sprint(int64(maxAge.Seconds()))
			w.Header().Add("Access-Control-Max-Age", sec)
		}

		if header := r.Header["Access-Control-Request-Headers"]; len(header) != 0 {
			allowHeaders := cors.GetAllowHeaders()

			goodHeaders := []string{}

			for _, x := range header {
				for _, y := range allowHeaders {
					if strings.ToLower(x) == strings.ToLower(y) {
						goodHeaders = append(goodHeaders, x)
					}
				}
			}

			if len(goodHeaders) == len(header) {
				w.Header().Add("Access-Control-Allow-Headers", strings.Join(goodHeaders, ", "))
			}
		}
	}
	return nil
}

var (
	// OptionsHandler - Generic Options Handler to handle when method isn't allowed for a resource
	OptionsHandler = func(gcors CorsOptionsInterface, lcors CorsOptionsInterface, allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", allowedMethods)

			if err := CorsPreflight(gcors, lcors, allowedMethods, w, r); err != nil {
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		}
	}
	// MethodNotAllowedHandler - Generic Handler to handle when method isn't allowed for a resource
	MethodNotAllowedHandler = func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", allowedMethods)
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
func registerVar(r *http.Request, fmtpname string, pvalue string) {
	if r.URL.RawQuery != "" {
		r.URL.RawQuery += "&" + fmtpname + pvalue
	} else {
		r.URL.RawQuery += fmtpname + pvalue
	}
}
