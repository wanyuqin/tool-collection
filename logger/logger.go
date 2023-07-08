package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/wanyuqin/tool-collection/backend/x/xfile"
	"os"
	"path/filepath"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).DisableColors = true
	log.Out = os.Stdout
	log.AddHook(&FileHook{})

}

func Debug(info string) {
	log.Debug(info)
}

func Error(info string) {
	log.Error(info)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

type FileHook struct {
}

func (f *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FileHook) Fire(e *logrus.Entry) error {
	logPath := ""
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	switch e.Level {
	case logrus.InfoLevel:
		logPath = filepath.Join(homeDir, ".tools_collection", "info.log")
	case logrus.DebugLevel:
		logPath = filepath.Join(homeDir, ".tools_collection", "debug.log")
	case logrus.ErrorLevel:
		logPath = filepath.Join(homeDir, ".tools_collection", "err.log")
	}

	file, err := os.OpenFile(logPath, os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(e.Message))
	return err
}

func InitLogFile(rootPath string) error {
	errLogPath := filepath.Join(rootPath, "err.log")
	infoLogPath := filepath.Join(rootPath, "info.log")
	debugLogPath := filepath.Join(rootPath, "debug.log")

	if err := xfile.CreateFileIfNotExist(errLogPath); err != nil {
		return err
	}

	if err := xfile.CreateFileIfNotExist(infoLogPath); err != nil {
		return err
	}

	if err := xfile.CreateFileIfNotExist(debugLogPath); err != nil {
		return err
	}

	return nil
}
