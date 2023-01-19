package hlogger

import (
	"log"
	"os"
	"sync"
)

// hydraLogger struct for the hydra logger
type hydraLogger struct {
	*log.Logger
	filename string
}

var (
	hlogger *hydraLogger
	once    sync.Once
)

// GetInstance creates a singleton instance of the hydra logger
func GetInstance() *hydraLogger {
	once.Do(func() {
		hlogger = createLogger("hydralogger.log")
	})
	return hlogger
}

// Create a logger instance
func createLogger(fname string) *hydraLogger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	// return a initialized hydraLogger pointer
	return &hydraLogger{
		filename: fname,
		Logger:   log.New(file, "Hydra ", log.Lshortfile),
	}

}
