// Copyright 2015 Husobee Associates, LLC.  All rights reserved.
// Use of this source code is governed by The MIT License, which
// can be found in the LICENSE file included.
// Portions Copyright 2013 Julien Schmidt.  All rights reserved.
// Use of portions of this source code is governed by a BSD style License
// which can be found in the LICENSE-go-http-routing-benchmark file included.
// Portions Copyright 2015 Labstack.  All rights reserved.

package vestigo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var (
	api = []route{
		// OAuth Authorizations
		{"GET", "/authorizations", nil},
		{"GET", "/authorizations/:id", nil},
		{"POST", "/authorizations", nil},
		{"PUT", "/authorizations/clients/:client_id", nil},
		{"PATCH", "/authorizations/:id", nil},
		{"DELETE", "/authorizations/:id", nil},
		{"GET", "/applications/:client_id/tokens/:access_token", nil},
		{"DELETE", "/applications/:client_id/tokens", nil},
		{"DELETE", "/applications/:client_id/tokens/:access_token", nil},

		// Activity
		{"GET", "/events", nil},
		{"GET", "/repos/:owner/:repo/events", nil},
		{"GET", "/networks/:owner/:repo/events", nil},
		{"GET", "/orgs/:org/events", nil},
		{"GET", "/users/:user/received_events", nil},
		{"GET", "/users/:user/received_events/public", nil},
		{"GET", "/users/:user/events", nil},
		{"GET", "/users/:user/events/public", nil},
		{"GET", "/users/:user/events/orgs/:org", nil},
		{"GET", "/feeds", nil},
		{"GET", "/notifications", nil},
		{"GET", "/repos/:owner/:repo/notifications", nil},
		{"PUT", "/notifications", nil},
		{"PUT", "/repos/:owner/:repo/notifications", nil},
		{"GET", "/notifications/threads/:id", nil},
		{"PATCH", "/notifications/threads/:id", nil},
		{"GET", "/notifications/threads/:id/subscription", nil},
		{"PUT", "/notifications/threads/:id/subscription", nil},
		{"DELETE", "/notifications/threads/:id/subscription", nil},
		{"GET", "/repos/:owner/:repo/stargazers", nil},
		{"GET", "/users/:user/starred", nil},
		{"GET", "/user/starred", nil},
		{"GET", "/user/starred/:owner/:repo", nil},
		{"PUT", "/user/starred/:owner/:repo", nil},
		{"DELETE", "/user/starred/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/subscribers", nil},
		{"GET", "/users/:user/subscriptions", nil},
		{"GET", "/user/subscriptions", nil},
		{"GET", "/repos/:owner/:repo/subscription", nil},
		{"PUT", "/repos/:owner/:repo/subscription", nil},
		{"DELETE", "/repos/:owner/:repo/subscription", nil},
		{"GET", "/user/subscriptions/:owner/:repo", nil},
		{"PUT", "/user/subscriptions/:owner/:repo", nil},
		{"DELETE", "/user/subscriptions/:owner/:repo", nil},

		// Gists
		{"GET", "/users/:user/gists", nil},
		{"GET", "/gists", nil},
		{"GET", "/gists/public", nil},
		{"GET", "/gists/starred", nil},
		{"GET", "/gists/:id", nil},
		{"POST", "/gists", nil},
		{"PATCH", "/gists/:id", nil},
		{"PUT", "/gists/:id/star", nil},
		{"DELETE", "/gists/:id/star", nil},
		{"GET", "/gists/:id/star", nil},
		{"POST", "/gists/:id/forks", nil},
		{"DELETE", "/gists/:id", nil},

		// Git Data
		{"GET", "/repos/:owner/:repo/git/blobs/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/blobs", nil},
		{"GET", "/repos/:owner/:repo/git/commits/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/commits", nil},
		{"GET", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"GET", "/repos/:owner/:repo/git/refs", nil},
		{"POST", "/repos/:owner/:repo/git/refs", nil},
		{"PATCH", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"DELETE", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"GET", "/repos/:owner/:repo/git/tags/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/tags", nil},
		{"GET", "/repos/:owner/:repo/git/trees/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/trees", nil},

		// Issues
		{"GET", "/issues", nil},
		{"GET", "/user/issues", nil},
		{"GET", "/orgs/:org/issues", nil},
		{"GET", "/repos/:owner/:repo/issues", nil},
		{"GET", "/repos/:owner/:repo/issues/:number", nil},
		{"POST", "/repos/:owner/:repo/issues", nil},
		{"PATCH", "/repos/:owner/:repo/issues/:number", nil},
		{"GET", "/repos/:owner/:repo/assignees", nil},
		{"GET", "/repos/:owner/:repo/assignees/:assignee", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/comments", nil},
		{"GET", "/repos/:owner/:repo/issues/comments", nil},
		{"GET", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"POST", "/repos/:owner/:repo/issues/:number/comments", nil},
		{"PATCH", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"DELETE", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/events", nil},
		{"GET", "/repos/:owner/:repo/issues/events", nil},
		{"GET", "/repos/:owner/:repo/issues/events/:id", nil},
		{"GET", "/repos/:owner/:repo/labels", nil},
		{"GET", "/repos/:owner/:repo/labels/:name", nil},
		{"POST", "/repos/:owner/:repo/labels", nil},
		{"PATCH", "/repos/:owner/:repo/labels/:name", nil},
		{"DELETE", "/repos/:owner/:repo/labels/:name", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"POST", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name", nil},
		{"PUT", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"GET", "/repos/:owner/:repo/milestones/:number/labels", nil},
		{"GET", "/repos/:owner/:repo/milestones", nil},
		{"GET", "/repos/:owner/:repo/milestones/:number", nil},
		{"POST", "/repos/:owner/:repo/milestones", nil},
		{"PATCH", "/repos/:owner/:repo/milestones/:number", nil},
		{"DELETE", "/repos/:owner/:repo/milestones/:number", nil},

		// Miscellaneous
		{"GET", "/emojis", nil},
		{"GET", "/gitignore/templates", nil},
		{"GET", "/gitignore/templates/:name", nil},
		{"POST", "/markdown", nil},
		{"POST", "/markdown/raw", nil},
		{"GET", "/meta", nil},
		{"GET", "/rate_limit", nil},

		// Organizations
		{"GET", "/users/:user/orgs", nil},
		{"GET", "/user/orgs", nil},
		{"GET", "/orgs/:org", nil},
		{"PATCH", "/orgs/:org", nil},
		{"GET", "/orgs/:org/members", nil},
		{"GET", "/orgs/:org/members/:user", nil},
		{"DELETE", "/orgs/:org/members/:user", nil},
		{"GET", "/orgs/:org/public_members", nil},
		{"GET", "/orgs/:org/public_members/:user", nil},
		{"PUT", "/orgs/:org/public_members/:user", nil},
		{"DELETE", "/orgs/:org/public_members/:user", nil},
		{"GET", "/orgs/:org/teams", nil},
		{"GET", "/teams/:id", nil},
		{"POST", "/orgs/:org/teams", nil},
		{"PATCH", "/teams/:id", nil},
		{"DELETE", "/teams/:id", nil},
		{"GET", "/teams/:id/members", nil},
		{"GET", "/teams/:id/members/:user", nil},
		{"PUT", "/teams/:id/members/:user", nil},
		{"DELETE", "/teams/:id/members/:user", nil},
		{"GET", "/teams/:id/repos", nil},
		{"GET", "/teams/:id/repos/:owner/:repo", nil},
		{"PUT", "/teams/:id/repos/:owner/:repo", nil},
		{"DELETE", "/teams/:id/repos/:owner/:repo", nil},
		{"GET", "/user/teams", nil},

		// Pull Requests
		{"GET", "/repos/:owner/:repo/pulls", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number", nil},
		{"POST", "/repos/:owner/:repo/pulls", nil},
		{"PATCH", "/repos/:owner/:repo/pulls/:number", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/commits", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/files", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/merge", nil},
		{"PUT", "/repos/:owner/:repo/pulls/:number/merge", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/comments", nil},
		{"GET", "/repos/:owner/:repo/pulls/comments", nil},
		{"GET", "/repos/:owner/:repo/pulls/comments/:number", nil},
		{"PUT", "/repos/:owner/:repo/pulls/:number/comments", nil},
		{"PATCH", "/repos/:owner/:repo/pulls/comments/:number", nil},
		{"DELETE", "/repos/:owner/:repo/pulls/comments/:number", nil},

		// Repositories
		{"GET", "/user/repos", nil},
		{"GET", "/users/:user/repos", nil},
		{"GET", "/orgs/:org/repos", nil},
		{"GET", "/repositories", nil},
		{"POST", "/user/repos", nil},
		{"POST", "/orgs/:org/repos", nil},
		{"GET", "/repos/:owner/:repo", nil},
		{"PATCH", "/repos/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/contributors", nil},
		{"GET", "/repos/:owner/:repo/languages", nil},
		{"GET", "/repos/:owner/:repo/teams", nil},
		{"GET", "/repos/:owner/:repo/tags", nil},
		{"GET", "/repos/:owner/:repo/branches", nil},
		{"GET", "/repos/:owner/:repo/branches/:branch", nil},
		{"DELETE", "/repos/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/collaborators", nil},
		{"GET", "/repos/:owner/:repo/collaborators/:user", nil},
		{"PUT", "/repos/:owner/:repo/collaborators/:user", nil},
		{"DELETE", "/repos/:owner/:repo/collaborators/:user", nil},
		{"GET", "/repos/:owner/:repo/comments", nil},
		{"GET", "/repos/:owner/:repo/commits/:sha/comments", nil},
		{"POST", "/repos/:owner/:repo/commits/:sha/comments", nil},
		{"GET", "/repos/:owner/:repo/comments/:id", nil},
		{"PATCH", "/repos/:owner/:repo/comments/:id", nil},
		{"DELETE", "/repos/:owner/:repo/comments/:id", nil},
		{"GET", "/repos/:owner/:repo/commits", nil},
		{"GET", "/repos/:owner/:repo/commits/:sha", nil},
		{"GET", "/repos/:owner/:repo/readme", nil},
		{"GET", "/repos/:owner/:repo/contents/*path", nil},
		{"PUT", "/repos/:owner/:repo/contents/*path", nil},
		{"DELETE", "/repos/:owner/:repo/contents/*path", nil},
		{"GET", "/repos/:owner/:repo/:archive_format/:ref", nil},
		{"GET", "/repos/:owner/:repo/keys", nil},
		{"GET", "/repos/:owner/:repo/keys/:id", nil},
		{"POST", "/repos/:owner/:repo/keys", nil},
		{"PATCH", "/repos/:owner/:repo/keys/:id", nil},
		{"DELETE", "/repos/:owner/:repo/keys/:id", nil},
		{"GET", "/repos/:owner/:repo/downloads", nil},
		{"GET", "/repos/:owner/:repo/downloads/:id", nil},
		{"DELETE", "/repos/:owner/:repo/downloads/:id", nil},
		{"GET", "/repos/:owner/:repo/forks", nil},
		{"POST", "/repos/:owner/:repo/forks", nil},
		{"GET", "/repos/:owner/:repo/hooks", nil},
		{"GET", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/hooks", nil},
		{"PATCH", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/hooks/:id/tests", nil},
		{"DELETE", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/merges", nil},
		{"GET", "/repos/:owner/:repo/releases", nil},
		{"GET", "/repos/:owner/:repo/releases/:id", nil},
		{"POST", "/repos/:owner/:repo/releases", nil},
		{"PATCH", "/repos/:owner/:repo/releases/:id", nil},
		{"DELETE", "/repos/:owner/:repo/releases/:id", nil},
		{"GET", "/repos/:owner/:repo/releases/:id/assets", nil},
		{"GET", "/repos/:owner/:repo/stats/contributors", nil},
		{"GET", "/repos/:owner/:repo/stats/commit_activity", nil},
		{"GET", "/repos/:owner/:repo/stats/code_frequency", nil},
		{"GET", "/repos/:owner/:repo/stats/participation", nil},
		{"GET", "/repos/:owner/:repo/stats/punch_card", nil},
		{"GET", "/repos/:owner/:repo/statuses/:ref", nil},
		{"POST", "/repos/:owner/:repo/statuses/:ref", nil},

		// Search
		{"GET", "/search/repositories", nil},
		{"GET", "/search/code", nil},
		{"GET", "/search/issues", nil},
		{"GET", "/search/users", nil},
		{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword", nil},
		{"GET", "/legacy/repos/search/:keyword", nil},
		{"GET", "/legacy/user/search/:keyword", nil},
		{"GET", "/legacy/user/email/:email", nil},

		// Users
		{"GET", "/users/:user", nil},
		{"GET", "/user", nil},
		{"PATCH", "/user", nil},
		{"GET", "/users", nil},
		{"GET", "/user/emails", nil},
		{"POST", "/user/emails", nil},
		{"DELETE", "/user/emails", nil},
		{"GET", "/users/:user/followers", nil},
		{"GET", "/user/followers", nil},
		{"GET", "/users/:user/following", nil},
		{"GET", "/user/following", nil},
		{"GET", "/user/following/:user", nil},
		{"GET", "/users/:user/following/:target_user", nil},
		{"PUT", "/user/following/:user", nil},
		{"DELETE", "/user/following/:user", nil},
		{"GET", "/users/:user/keys", nil},
		{"GET", "/user/keys", nil},
		{"GET", "/user/keys/:id", nil},
		{"POST", "/user/keys", nil},
		{"PATCH", "/user/keys/:id", nil},
		{"DELETE", "/user/keys/:id", nil},
	}
)

func TestRouterParam(t *testing.T) {
	r := NewRouter()
	r.Add("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {})
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	h := r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "1", Param(req, "id"))
	}
}

func TestRouterTwoParam(t *testing.T) {
	r := NewRouter()
	r.Add("GET", "/users/:uid/files/:fid", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/users/1/files/1", nil)
	w := httptest.NewRecorder()
	h := r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "1", Param(req, "uid"))
		assert.Equal(t, "1", Param(req, "fid"))
	}
}

func TestAddParamEncode(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test?:user=1", nil)
	AddParam(r, "id", "2 2")
	assert.Equal(t, r.URL.RawQuery, ":user=1&%3Aid=2+2")
}

func TestParamNames(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test?:user=1&group=2", nil)
	AddParam(r, "location", "San Francisco, CA")
	actual := ParamNames(r)

	var foundLocation bool
	var foundUser bool
	for _, v := range actual {
		if v == ":user" {
			foundUser = true
		}
		if v == ":location" {
			foundLocation = true
		}
	}

	assert.Equal(t, foundUser, true)
	assert.Equal(t, foundLocation, true)
}

func TestRouterMicroParam(t *testing.T) {
	r := NewRouter()
	r.Add("GET", "/:a/:b/:c", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/1/2/3", nil)
	h := r.Find(req)
	if assert.NotNil(t, h) {
		assert.Equal(t, "1", Param(req, "a"))
		assert.Equal(t, "2", Param(req, "b"))
		assert.Equal(t, "3", Param(req, "c"))
	}
}

func TestRouterMixParamMatchAny(t *testing.T) {
	r := NewRouter()

	// Route
	r.Add("GET", "/users/:id/*", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/users/joe/comments", nil)
	w := httptest.NewRecorder()
	h := r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "joe", Param(req, "id"))
	}
}

func TestRouterMultiRoute(t *testing.T) {
	r := NewRouter()

	// Routes
	r.Add("GET", "/users", func(w http.ResponseWriter, r *http.Request) {})
	r.Add("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {})
	r.Add("GET", "/users/static", func(w http.ResponseWriter, r *http.Request) {})
	r.Add("GET", "/:id", func(w http.ResponseWriter, r *http.Request) {})

	// Route > /users
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	h := r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
	}

	// Route > /users/:id
	req, _ = http.NewRequest("GET", "/users/1", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "1", Param(req, "id"))
	}

	// Route > /users/static
	req, _ = http.NewRequest("GET", "/users/static", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()
	if assert.NotNil(t, h) {
		h(w, req)
	}

	// Route > /users/static
	req, _ = http.NewRequest("GET", "/users/something", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()
	if assert.NotNil(t, h) {
		h(w, req)
	}

	// Route > /user
	req, _ = http.NewRequest("GET", "/user/1", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()

	h(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)

	// Route > /user
	req, _ = http.NewRequest("GET", "/user", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()

	h(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, "user", Param(req, "id"))

	// Route > /test
	req, _ = http.NewRequest("GET", "/users123", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()

	h(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "users123", Param(req, "id"))

	// Invalid Method for Resource
	// Route > /user
	req, _ = http.NewRequest("INVALID", "/users", nil)
	h = r.Find(req)
	w = httptest.NewRecorder()
	h(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

}

/*
func TestRouterPriority(t *testing.T) {
	r := NewRouter()

	// Routes
	r.Add(GET, "/users", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/new", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/:id", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/dew", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/:id/files", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/newsee", func(w http.ResponseWriter, r *http.Request) {})
	r.Add(GET, "/users/*", func(w http.ResponseWriter, r *http.Request) {})

	// Route > /users
	h, _ := r.Find(GET, "/users", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 1, c.Get("a"))
	}

	// Route > /users/new
	h, _ = r.Find(GET, "/users/new", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 2, c.Get("b"))
	}

	// Route > /users/:id
	h, _ = r.Find(GET, "/users/1", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 3, c.Get("c"))
	}

	// Route > /users/dew
	h, _ = r.Find(GET, "/users/dew", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 4, c.Get("d"))
	}

	// Route > /users/:id/files
	h, _ = r.Find(GET, "/users/1/files", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 5, c.Get("e"))
	}

	// Route > /users/:id
	h, _ = r.Find(GET, "/users/news", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 3, c.Get("c"))
	}

	// Route > /users/*
	h, _ = r.Find(GET, "/users/joe/books", c)
	if assert.NotNil(t, h) {
		h(c)
		assert.Equal(t, 7, c.Get("g"))
		assert.Equal(t, "joe/books", c.Param("_name"))
	}
}
*/

func TestRouterParamNames(t *testing.T) {
	r := NewRouter()

	// Routes
	r.Add("GET", "/users", func(w http.ResponseWriter, r *http.Request) {})
	r.Add("GET", "/users/:uid", func(w http.ResponseWriter, r *http.Request) {})
	r.Add("GET", "/users/:uid/files/:fid", func(w http.ResponseWriter, r *http.Request) {})

	// Route > /users/:id
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	h := r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "1", Param(req, "uid"))
	}

	// Route > /users/:uid/files/:fid
	req, _ = http.NewRequest("GET", "/users/1/files/1", nil)
	w = httptest.NewRecorder()
	h = r.Find(req)
	if assert.NotNil(t, h) {
		h(w, req)
		assert.Equal(t, "1", Param(req, "uid"))
		assert.Equal(t, "1", Param(req, "fid"))
	}
}

func TestRouterParamGet(t *testing.T) {
	r := NewRouter()
	r.Add("GET", "/users/:uid", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "222", r.URL.Query().Get(":uid"))
		assert.Equal(t, "222", Param(r, "uid"))
		assert.Equal(t, "red", r.URL.Query().Get("color"))
		assert.Equal(t, "burger", r.URL.Query().Get("food"))
	})

	req, _ := http.NewRequest("GET", "/users/222?color=red&food=burger", nil)
	h := httptest.NewRecorder()
	r.ServeHTTP(h, req)
}

func TestRouterParamPost(t *testing.T) {
	r := NewRouter()
	r.Add("POST", "/users/:uid", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "123", r.FormValue("id"))
		assert.Equal(t, "123", r.Form.Get("id"))
		assert.Equal(t, "222", r.URL.Query().Get(":uid"))
		assert.Equal(t, "222", Param(r, "uid"))
		assert.Equal(t, "red", r.URL.Query().Get("color"))
		assert.Equal(t, "burger", r.URL.Query().Get("food"))
	})

	form := url.Values{}
	form.Add("id", "123")
	req, _ := http.NewRequest("POST", "/users/222?color=red&food=burger", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h := httptest.NewRecorder()
	r.ServeHTTP(h, req)
}

// TestUnderscoreFirstCall references issues #29
func TestUnderscoreFirstCall(t *testing.T) {
	r := NewRouter()
	h := httptest.NewRecorder()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "health") })
	r.Get("/_/accounts/foo", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "not param") })
	r.Get("/_/:project/bar", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "param") })

	req, _ := http.NewRequest("GET", "/_/a/bar", nil)
	r.ServeHTTP(h, req)
	assert.Equal(t, 200, h.Code)
}

// TestUnderscoreSecondCall references issues #29
func TestUnderscoreSecondCall(t *testing.T) {
	r := NewRouter()
	h := httptest.NewRecorder()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "health") })
	r.Get("/_/accounts/foo", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "not param") })
	r.Get("/_/:project/bar", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "param") })

	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(h, req)
	assert.Equal(t, 200, h.Code)

	req, _ = http.NewRequest("GET", "/_/a/bar", nil)
	r.ServeHTTP(h, req)
	assert.Equal(t, 200, h.Code)
}

func TestRouterAPI(t *testing.T) {
	r := NewRouter()

	for _, route := range api {
		r.Add(route.method, route.path, func(w http.ResponseWriter, req *http.Request) {})
	}

	w := httptest.NewRecorder()

	for _, route := range api {

		req, _ := http.NewRequest(route.method, route.path, nil)
		h := r.Find(req)
		if assert.NotNil(t, h) {
			for _, n := range ParamNames(req) {
				if assert.NotEmpty(t, n) {
					assert.NotNil(t, Param(req, n))
				}
			}
			h(w, req)
		}
	}
}

func TestRouterAddInvalidMethod(t *testing.T) {
	r := NewRouter()
	assert.Panics(t, func() {
		r.Add("INVALID", "/", func(w http.ResponseWriter, req *http.Request) {})
	})
}

func TestMethodSpecificAddRoute(t *testing.T) {
	router := NewRouter()
	m := map[string]func(path string, handler http.HandlerFunc){
		"GET":     router.Get,
		"POST":    router.Post,
		"CONNECT": router.Connect,
		"DELETE":  router.Delete,
		"PATCH":   router.Patch,
		"PUT":     router.Put,
		"TRACE":   router.Trace,
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("success-" + r.Method))
	}
	path := "/test"
	for _, v := range m {
		v(path, f)
	}
	for k := range m {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(k, path, nil)
		if err != nil {
			t.Errorf("Failed to create a new request, method: %s, path: %s", k, path)
		}
		router.ServeHTTP(w, r)
		if w.Code != 200 || w.Body.String() != "success-"+k {
			t.Errorf("Invalid response, method: %s, path: %s, code: %d, body: %s", k, path, w.Code, w.Body.String())
		}
	}

}

func TestHandleAddRoute(t *testing.T) {
	router := NewRouter()
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("success-" + r.Method))
	}
	path := "/test"
	router.Handle(path, http.HandlerFunc(f))
	for k := range methods {
		if k == http.MethodHead || k == http.MethodOptions || k == http.MethodTrace {
			continue
		}
		w := httptest.NewRecorder()
		r, err := http.NewRequest(k, path, nil)
		if err != nil {
			t.Errorf("Failed to create a new request, method: %s, path: %s", k, path)
		}
		router.ServeHTTP(w, r)
		if w.Code != 200 || w.Body.String() != "success-"+k {
			t.Errorf("Invalid response, method: %s, path: %s, code: %d, body: %s", k, path, w.Code, w.Body.String())
		}
	}
}

func TestHandleFuncAddRoute(t *testing.T) {
	router := NewRouter()
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("success-" + r.Method))
	}
	path := "/test"
	router.HandleFunc(path, f)
	for k := range methods {
		if k == http.MethodHead || k == http.MethodOptions || k == http.MethodTrace {
			continue
		}
		w := httptest.NewRecorder()
		r, err := http.NewRequest(k, path, nil)
		if err != nil {
			t.Errorf("Failed to create a new request, method: %s, path: %s", k, path)
		}
		router.ServeHTTP(w, r)
		if w.Code != 200 || w.Body.String() != "success-"+k {
			t.Errorf("Invalid response, method: %s, path: %s, code: %d, body: %s", k, path, w.Code, w.Body.String())
		}
	}
}

func TestRouterServeHTTP(t *testing.T) {
	r := NewRouter()

	r.Add("GET", "/users", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})

	// OK
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Not found
	req, _ = http.NewRequest("GET", "/files", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSameRouteDifferentMethodDifferentPnamesServeHTTP(t *testing.T) {
	r := NewRouter()

	r.Add("GET", "/:var1/test/:var2", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("var1: " + Param(req, "var1")))
	})
	r.Add("POST", "/:var2/test/:var1", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("var2: " + Param(req, "var2")))
	})

	// OK
	req, _ := http.NewRequest("GET", "/one/test/two", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "var1: one")

	req, _ = http.NewRequest("POST", "/two/test/one", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "var2: two")
}

func TestMatchedURLTemplate(t *testing.T) {
	r := NewRouter()

	for _, v := range []string{
		"/users/:test_param",
		"/users",
		"/users/:test_param/:param_two",
		"/:test_params",
	} {

		r.Add("GET", v, func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(""))
		})
	}

	for _, v := range []string{
		"/users/:test_param",
		"/users",
		"/users/:test_param/:param_two",
		"/:test_params",
	} {
		req, _ := http.NewRequest("GET", v, nil)
		templ := r.GetMatchedPathTemplate(req)
		assert.Equal(t, templ, v)
	}
}

func TestIssue49a(t *testing.T) {

	r := NewRouter()

	r.Get("/greet/:name", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome %s\n", Param(r, "name"))
	})
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You opened %s\n", Param(r, "_name"))
	})

	req, _ := http.NewRequest("GET", "/g", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, strings.Contains(w.Body.String(), "You opened"))

}

func TestIssue49b(t *testing.T) {

	r := NewRouter()

	r.Get("/books/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/books/\n")
	})

	r.Get("/books/:book", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/books/%s\n", Param(r, "book"))
	})

	r.Get("/bookcase/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/bookcase/\n")
	})

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "wildcard!! %s\n", Param(r, "_name"))
	})

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, strings.Contains(w.Body.String(), "wildcard!!"))

	req, _ = http.NewRequest("GET", "/bookk", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, strings.Contains(w.Body.String(), "wildcard!!"))

	req, _ = http.NewRequest("GET", "/book", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, strings.Contains(w.Body.String(), "wildcard!!"))

}

func TestIssue51(t *testing.T) {

	fmt.Println("here??")
	r := NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/books/%s\n", Param(r, "book"))
	})

	r.Get("/users/:name", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/users/%s\n", Param(r, "name"))
	})

	r.Get("/admin/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/admin/\n")
	})

	r.Get("/books/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/books/\n")
	})

	r.Get("/books/:book", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/books/%s\n", Param(r, "book"))
	})

	req, _ := http.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	//assert.True(t, strings.Contains(w.Body.String(), "not found"))

}
