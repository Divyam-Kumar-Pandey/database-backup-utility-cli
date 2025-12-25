package utils

import (
	"fmt"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	logFile     *os.File
)

func InitLogger() {
	var err error
	logFile, err = os.OpenFile(
		"db-backup-cli.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	InfoLogger = log.New(
		logFile,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	ErrorLogger = log.New(
		logFile,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

// func CloseLogger() {
// 	if logFile != nil {
// 		_ = logFile.Close()
// 	}
// }

func LogInfo(msg string) {
	fmt.Println("INFO: " + msg)
	if InfoLogger != nil {
		InfoLogger.Println(msg)
	}
}

func LogError(msg string) {
	fmt.Println("ERROR: " + msg)
	if ErrorLogger != nil {
		ErrorLogger.Println(msg)
	}
}
