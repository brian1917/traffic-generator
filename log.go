package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger is the global logger for Workloader
var Logger log.Logger
var logFile string

func SetUpLogging() {

	// First check env variable, then config file, then use default
	logFile = os.Getenv("TRAFFIC_GENERATOR_LOG_FILE")
	if logFile == "" {
		logFile = "traffic-generator.log"
	}
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	Logger.SetOutput(f)

}

// LogError writes the error the workloader.log and always prints an error to stdout.
func LogError(msg string) {
	Logger.SetPrefix(time.Now().Format("2006-01-02 15:04:05 "))
	fmt.Printf("%s [ERROR] - %s see log file for potentially more information.\r\n", time.Now().Format("2006-01-02 15:04:05 "), msg)
	Logger.Printf("[ERROR] - %s\r\n", msg)
}

// LogErrorf uses string formatting to write to log to workloader.log and always prints msg to stdout.
func LogErrorf(format string, a ...any) {
	LogError(fmt.Sprintf(format, a...))
}

// LogWarning writes the log to workloader.log and optionally prints msg to stdout.
func LogWarning(msg string, stdout bool) {
	Logger.SetPrefix(time.Now().Format("2006-01-02 15:04:05 "))
	if stdout {
		fmt.Printf("%s [WARNING] - %s\r\n", time.Now().Format("2006-01-02 15:04:05 "), msg)
	}
	Logger.Printf("[WARNING] - %s\r\n", msg)
}

// LogWarningf uses string formatting to write to log to workloader.log and optionally prints msg to stdout.
func LogWarningf(stdout bool, format string, a ...any) {
	LogWarning(fmt.Sprintf(format, a...), stdout)
}

// LogInfo writes the log to workloader.log and optionally prints msg to stdout.
func LogInfo(msg string, stdout bool) {
	Logger.SetPrefix(time.Now().Format("2006-01-02 15:04:05 "))
	if stdout {
		fmt.Printf("%s [INFO] - %s\r\n", time.Now().Format("2006-01-02 15:04:05 "), msg)
	}
	Logger.Printf("[INFO] - %s\r\n", msg)
}

// LogInfof uses string formatting to write to log to workloader.log and optionally prints msg to stdout.
func LogInfof(stdout bool, format string, a ...any) {
	LogInfo(fmt.Sprintf(format, a...), stdout)
}
