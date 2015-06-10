package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	port = flag.Int("port", 8090, "The server port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	fmt.Printf("proxy: %v\n", err)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "set,"+r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10)+"\n")

	// body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, "ioutil.ReadAll: %v", err)
		return
	}

	fmt.Fprintf(w, "set body:\n%v\n", string(body))
	fmt.Printf("set body:\n%v\n", string(body))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get,"+r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10)+"\n")

	// body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Fprintf(w, "ioutil.ReadAll: %v", err)
		return
	}

	fmt.Fprintf(w, "get body:\n%v\n", string(body))
	fmt.Printf("get body:\n%v\n", string(body))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	test_3_Handler(w, r)
}
