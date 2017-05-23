package vestigo

import "net/http"

type Interceptor interface {
	// returns if the interceptor should run before handler function
	Before() bool
	// returns if the interceptor should run after handler function
	After() bool
	// the actual intercept function, returns bool indicating if the request should continue
	Intercept(w http.ResponseWriter, r *http.Request) bool
}
