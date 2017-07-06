package log

import (
	"errors"
	"os"
	"path"

	logging "github.com/op/go-logging"
)

var logFormatter = logging.MustStringFormatter(
	` %{level:.1s}%{time:0102 15:04:05.999999} %{pid} %{shortfile}] %{message}`,
)

type FileLog struct {
	Name    string
	Logger  *logging.Logger
	LogFile *os.File
}

type LoggerConfig interface {
	GetString(key string) string
}

// NewLogger creates a named log under the files path
func NewLogger(filesPath string, name string) (*FileLog, error) {
	log := logging.MustGetLogger(name)

	logDirPath := path.Join(filesPath, "log")
	if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
		os.Mkdir(logDirPath, 0777)
	}

	logFilePath := path.Join(logDirPath, name+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.New("Unable to create log file:" + err.Error())
	}

	fileLog := logging.NewLogBackend(logFile, "["+name+"]", 0)
	consoleLog := logging.NewLogBackend(os.Stdout, "["+name+"]", 0)

	fileLogLevel := logging.AddModuleLevel(fileLog)
	fileLogLevel.SetLevel(logging.INFO, "")

	consoleLogBackend := logging.NewBackendFormatter(consoleLog, logFormatter)
	fileLogBackend := logging.NewBackendFormatter(fileLog, logFormatter)

	log.SetBackend(logging.SetBackend(fileLogBackend, consoleLogBackend))

	return &FileLog{
		Name:    name,
		Logger:  log,
		LogFile: logFile,
	}, nil
}
