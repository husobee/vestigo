//+build go1.7
package vestigo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParamNotExists(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test?group=2", nil)
	// shouldnt exist
	val := Param(r, "location")
	assert.Equal(t, "", val)
}

func TestParamNames(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test?group=2", nil)
	AddParam(r, "user", "test")
	AddParam(r, "location", "San Francisco, CA")
	actual := ParamNames(r)

	var foundLocation bool
	var foundUser bool
	for _, v := range actual {
		if v == "user" {
			foundUser = true
		}
		if v == "location" {
			foundLocation = true
		}
	}

	assert.Equal(t, foundUser, true)
	assert.Equal(t, foundLocation, true)
}
