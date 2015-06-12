package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	port = flag.Int("port", 8090, "The server port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/setscale", setScaleHandler)
	http.HandleFunc("/getscale", getScaleHandler)
	http.HandleFunc("/get", getHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	fmt.Printf("proxy: %v\n", err)
}
