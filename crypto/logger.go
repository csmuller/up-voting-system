package crypto

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// logger used to log execution times to a file.
	logger *log.Logger
)

// SetupLogger sets up the logger to log to the given file. The file is created in the current
// user's home directory.
func SetupLogger(logfile string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("couldn't find user's home directory for creating logging file: %v", err)
	}
	f, err := os.OpenFile(filepath.Join(homeDir, logfile),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()
	logger = log.New(io.Writer(f), "", log.LstdFlags)
}

// LogExecutionTime logs the time past since the given start time.
func LogExecutionTime(start time.Time, opDesc string) {
	elapsed := time.Since(start)
	if logger != nil {
		logger.Printf("%s took %s", opDesc, elapsed)
	}
}
