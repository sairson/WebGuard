package main

import (
	"fmt"
	"github.com/sairson/WebGuard"
	"net/http"
)

func main() {
	guard := WebGuard.New("cfg.yml", true, nil, handler)
	http.HandleFunc("/", guard.RunGuard())
	http.ListenAndServe(":8555", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ok")
}

func debug(in interface{}) {
	fmt.Println(in)
}
