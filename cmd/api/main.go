package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request received")
		_, _ = fmt.Fprint(w, "pong")
	})
	_ = http.ListenAndServe(":8080", nil)
}
