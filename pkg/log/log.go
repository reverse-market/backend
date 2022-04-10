package log

import (
	"log"
	"os"
)

type Loggers struct {
	info  *log.Logger
	error *log.Logger
}

func (l *Loggers) Info() *log.Logger {
	return l.info
}

func (l *Loggers) Error() *log.Logger {
	return l.error
}

func NewLoggers() *Loggers {
	return &Loggers{
		info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
