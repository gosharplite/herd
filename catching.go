package main

import (
	"errors"
	"time"
)

var catching = make(map[string]kvt)

type kvt struct {
	prefix string
	key    string
	value  string
	t      time.Time
}

func addCatch(prefix, key, value string) {
	catching[prefix+","+key] = kvt{
		prefix: prefix,
		key:    key,
		value:  value,
		t:      time.Now(),
	}

	go checkCatch()
}

func getCatch(prefix, key string) (string, error) {
	v, ok := catching[prefix+","+key]
	if !ok {
		return "", errors.New("no catche")
	}

	return v.value, nil
}

func checkCatch() {
	if len(catching) > 1000 {
		for _, c := range catching {
			if time.Now().Sub(c.t) > 10*time.Minute {
				delete(catching, c.key)
			}
		}
	}
}
