package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/helloworld", helloHandler)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe("193.233.202.119:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/helloworld" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }


    fmt.Fprintf(w, "Hello World!\n")
}
