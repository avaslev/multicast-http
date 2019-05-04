package main

import (
	"log"
    "net/http"
    "fmt"
    "github.com/avaslev/multicast-http/src/multicast"
)

func handleRequest(res http.ResponseWriter, req *http.Request) {
    multicast.HandleRequest(res, req)
    fmt.Fprintf(res, "Ok")
}

func main() {
    log.Printf("Server will run \n")

	// start server on port 80
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
    }
}