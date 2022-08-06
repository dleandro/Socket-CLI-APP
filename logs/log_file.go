package logs

import (
	"os"
)

const LOG_FILE_NAME = "numbers.log"

type Logger struct {
	LogFileIsSetupChan chan bool
}

func InitLogger() *Logger {
	logger := &Logger{
		LogFileIsSetupChan: make(chan bool),
	}

	go logger.createLogFile()

	return logger
}

func (logger *Logger) createLogFile() {
	file, err := os.Create(LOG_FILE_NAME) // TODO: could've put this in a go routine also not well handled

	if (err != nil) { // send signal to a new channel saying that app and server need to quit because the log file failed to be created

	}

	close(logger.LogFileIsSetupChan)

	// send signal once it finishes creating log file so that we know we can proceed with the logs
}