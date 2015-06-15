package log

import (
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
)

var slog *syslog.Writer

func init() {
	var err error
	slog, err = syslog.New(syslog.LOG_INFO, "[herd]")
	if err != nil {
		log.Printf("syslog.New err: %s", err)
	}
}

func Info(f string, v ...interface{}) string {
	s := fmt.Sprintf(f, v...)
	log.Printf(s)
	return s
}

func Err(f string, v ...interface{}) string {
	s := fmt.Sprintf(f, v...)
	slog.Err(s)
	log.Printf(s)
	return s
}

func Marshal(v interface{}) string {

	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		e := "Err: json.MarshalIndent()"
		slog.Err(e)
		log.Printf(e)
		return ""
	}

	s := string(j)

	log.Printf(s)

	return s
}
