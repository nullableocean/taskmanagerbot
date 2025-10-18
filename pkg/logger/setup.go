package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// SetupFileForLogs добавляет вывод в файл для логгера из пакета "log"
func SetupFileForLogs(logDir, logFile string) (*os.File, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("setup log file: cannot create log directory: %v\n", err)
	}

	logPath := filepath.Join(logDir, logFile)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return nil, fmt.Errorf("setup log file: using stdout only, cannot open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return file, nil
}
