package WebGuard

import (
	"net/http"
	"testing"
)

type baseHandle struct {
}

func _Test(t *testing.T) {
	handle := &baseHandle{}
	http.Handle("/", handle)
	server := &http.Server{
		Addr:    ":8555", // proxy port
		Handler: handle,  // Cache structure
	}

	server.ListenAndServe()
}

func (b *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
