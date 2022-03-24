package main

import (
	"log"
	"net/http"
)

func HelloHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello!"))
}

func main() {
	http.HandleFunc("/", HelloHandler)
	err := http.ListenAndServe(":80", http.DefaultServeMux)
	if err != nil {
		log.Fatal(err)
	}
}
