package logging

import (
  "log"
  "os"
)

type Logger struct {
  *log.Logger
}

func NewLogger() *Logger {
  return &Logger{
    Logger: log.New(os.Stdout, "gRPC Server: ", log.Ldate|log.Ltime|log.Lshortfile),
  }
}

func (l *Logger) LogRPC(method string, duration string) {
  l.Printf("Method: %s, Duration: %s", method, duration)
}
