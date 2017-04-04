package vestigo

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimmedParamNames(t *testing.T) {
	router := NewRouter()

	// Wildcard should be _name
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/*")
		assert.Equal(t, []string([]string{"_name"}), TrimmedParamNames(r), "Should be _name")
	})

	// Some random name should be b
	router.Get("/a/:b", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/a/:b")
		assert.Equal(t, []string([]string{"b"}), TrimmedParamNames(r), "Should be b")
	})

	// Some random name
	router.Get("/:a", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/:a")
		assert.Equal(t, []string([]string{"a"}), TrimmedParamNames(r), "Should be a")
	})

	// Multiple parameters random name
	router.Get("/a/:b/c/:d", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/a/:b/c/:d")
		var foundB bool
		var foundD bool
		for _, v := range TrimmedParamNames(r) {
			if v == "b" {
				foundB = true
			}
			if v == "d" {
				foundD = true
			}
		}
		assert.True(t, foundB, "Result should contain b")
		assert.True(t, foundD, "Result should contain d")
	})

	rec := httptest.NewRecorder()

	req1, _ := http.NewRequest("GET", "/a/", nil)
	req4, _ := http.NewRequest("GET", "/a", nil)
	req2, _ := http.NewRequest("GET", "/a/b", nil)
	req3, _ := http.NewRequest("GET", "/a/b/c/d", nil)

	router.ServeHTTP(rec, req1)
	router.ServeHTTP(rec, req2)
	router.ServeHTTP(rec, req3)
	router.ServeHTTP(rec, req4)
}
