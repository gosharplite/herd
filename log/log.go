package log

import (
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
	slog.Info(s)
	log.Printf(s)
	return s
}

func Err(f string, v ...interface{}) string {
	s := fmt.Sprintf(f, v...)
	slog.Err(s)
	log.Printf(s)
	return s
}
