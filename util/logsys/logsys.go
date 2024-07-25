package logsys

import (
	"errors"
	"fmt"
	"log"
)

var IsDebug = false

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
var Info = func(format string, cate string, a ...interface{}) {
	log.Printf("[info]["+cate+"]"+format+"\r\n", a...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
var Error = func(format string, cate string, a ...interface{}) error {
	log.Printf("[Error]["+cate+"]"+format+"\r\n", a...)
	return errors.New(fmt.Sprintf(format, a...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
var Debug = func(format string, cate string, a ...interface{}) {
	if !IsDebug {
		return
	}
	log.Printf("[Debug]["+cate+"]"+format+"\r\n", a...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
var Panicln = func(format string, cate string, a ...interface{}) {
	log.Panicln(fmt.Sprintf("[Err]["+cate+"]"+format+"\r\n", a...))
}
