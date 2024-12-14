package utils

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() {
	logger = log.New(os.Stdout, "info: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string) {
	logger.Println("info: " + message)
}

func LogError(message string) {
	logger.Println("error: " + message)
}
