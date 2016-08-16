//+build !go1.7

package vestigo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
