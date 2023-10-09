package main

import "log"

type Logger interface {
	Error(v ...interface{})
	Printf(format string, args ...interface{})
	Print(v ...interface{})
}

type DefaultLogger struct{}

func (l DefaultLogger) Error(args ...interface{}) {
	log.Println(append(args, "ERROR: "))
}

func (l DefaultLogger) Printf(format string, args ...interface{}) {
	log.Println(append(args, "INFO: "))
}

func (l DefaultLogger) Print(args ...interface{}) {
	log.Println(append(args, "DEBUG: "))
}
