package main

import (
	"io/ioutil"
	"net/http"
)

func getApi() (string, error) {

	resp, err := http.Get("http://192.168.4.54:8080/api")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getVersion() (string, error) {

	resp, err := http.Get("http://192.168.4.54:8080/version")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
