package vestigo

import "net/http"

type Interceptor interface {
	// returns true if the interceptor should run before handler
	Before() bool
	// returns true if the interceptor should run after handler
	After() bool
	// the actual intercept function, returns true if the request should continue to handler and/or
	// chained interceptors, false if the execution should terminate
	Intercept(w http.ResponseWriter, r *http.Request) bool
}
