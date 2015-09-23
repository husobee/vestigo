# Vestigo - A Vestige Of Echo's URL Router

## Abstract

[Echo][echo-main] has a very fast URL router.  This repo is a vestige of just the URL Router,
broken out into a stand alone module.  There is such an abundance of parts and pieces that can be fit
together for go web services, it seems like a shame to have a very fast URL router require the use
of one framework, and one context model.  This library aims to give the world a fast, and featureful
URL router that can stand on it's own, without being forced into a particular web framework.

## Design

1. Radix Tree Based
2. Attach URL Parameters into Request (PAT style) instead of context

### TODOs for V1

- [ ] Fix bug in router where handler.allowedMethods is getting populated where it shouldn't be
- [ ] Valiators for URL params
- [ ] Validate with Tests RFC 2616 Compliance (OPTIONS, etc)

### Long Term TODOs
- [ ] Implement RFC 6570 URI Parameters

## Performance

Initial implementation on a fork of [standard http performance testing libary][http-perf-test] shows the following:

```
BenchmarkVestigo_GithubAll         20000             75763 ns/op            9280 B/op        339 allocs/op
```

I should mention that the above performance is about 2x slower then the fastest URL router I have tested (Echo/Gin), and
is slightly worse than HTTPRouter, but I am happy with this performance considering this implementation is the fastest 
implementation that can handle standard http.HandlerFunc handlers, without forcing end users to use a particular context, 
or use a non-standard handler function, locking them into an implementation.

## Examples

```go

package main

import (
	"log"
	"net/http"

	"github.com/husobee/vestigo"
)

func main () {
    router := vestigo.NewRouter()

    router.Get("/welcome", GetWelcomeHandler)
    router.Post("/welcome/:name", PostWelcomeHandler)

	log.Fatal(http.ListenAndServe(":1234", router))

}

func PostWelcomeHandler(w http.ResponseWriter, r *http.Request) {
    name := vestigo.Param(r, "name") // url params live in the request
    w.WriteHeader(200)
    w.Write([]byte("wecome " + name +"!"))
}
func GetWelcomeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("wecome!"))
}

```

## Licensing

Portions of the URL Router were taken from [Echo][echo-main] and are covered under their [License][echo-main-license].

The rest of the implementation is covered under The MIT License covered under this [License][vestigo-main-license].

# Contributing

If you wish to contribute, please fork this repository, submit an issue, or pull request with your suggestions.  
Please use gofmt and golint before trying to contribute.


[echo-main]: https://github.com/labstack/echo
[echo-main-license]: https://github.com/labstack/echo/blob/master/LICENSE
[vestigo-main-license]: https://github.com/husobee/vestigo/blob/master/LICENSE
[http-perf-test]: https://github.com/julienschmidt/go-http-routing-benchmark
