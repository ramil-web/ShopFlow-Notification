package logger

import "log"

type Logger struct{}

func (l *Logger) Infof(format string, v ...interface{})  { log.Printf("[INFO] "+format, v...) }
func (l *Logger) Debugf(format string, v ...interface{}) { log.Printf("[DEBUG] "+format, v...) }
func (l *Logger) Errorf(format string, v ...interface{}) { log.Printf("[ERROR] "+format, v...) }
