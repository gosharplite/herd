package log

import (
	"encoding/json"
	"fmt"
	"log"
)

func Info(f string, v ...interface{}) string {
	s := fmt.Sprintf(f, v...)
	log.Printf(s)
	return s
}

func Err(f string, v ...interface{}) string {
	s := fmt.Sprintf(f, v...)
	log.Printf(s)
	return s
}

func Marshal(v interface{}) string {

	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		e := "Err: json.MarshalIndent()"
		log.Printf(e)
		return ""
	}

	s := string(j)

	log.Printf(s)

	return s
}
