package logger

import (
	"log"
	"os"
)

const logFlags = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile | log.Lmsgprefix

var (
	infoLogger    *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", logFlags)
	debugLogger = log.New(os.Stdout, "DEBUG: ", logFlags)
	warningLogger = log.New(os.Stdout, "WARNING: ", logFlags)
	errorLogger = log.New(os.Stderr, "ERROR: ", logFlags)
}

func InfoLogger() *log.Logger {
	return infoLogger
}

func DebugLogger() *log.Logger {
	return debugLogger
}

func WarningLogger() *log.Logger {
	return warningLogger
}

func ErrorLogger() *log.Logger {
	return errorLogger
}
