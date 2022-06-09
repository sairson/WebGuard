package main

import (
	"fmt"
	"github.com/sairson/WebGuard"
	"net/http"
)

type baseHandle struct {
}

func main() {
	guard := WebGuard.New("cfg.yml", true, handler)
	http.HandleFunc("/", guard.RunGuard())
	http.ListenAndServe(":8555", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ok")
}
