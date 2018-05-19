package svrkit

import (
	"io"
	"log"
	"os"
)

//Log2File log 同时写到 stdout 和文件
func Log2File(file string) {
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("open log file error!", err.Error())
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile)) //同时写到stdout和日志文件
}
