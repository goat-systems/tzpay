package logging

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

func GetLogging(file string) (*logging.Logger, *os.File) {
	log := logging.MustGetLogger("example")
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	var logger *logging.LogBackend
	var f *os.File
	if file != "" {
		logger = logging.NewLogBackend(os.Stdout, "", 0)
	} else {
		f, err := openFile(file)
		if err != nil {
			fmt.Println("[logging] could not open file " + file + " for logging")
			logger = logging.NewLogBackend(os.Stdout, "", 0)
		} else {
			logger = logging.NewLogBackend(f, "", 0)
		}
	}
	backendFormatter := logging.NewBackendFormatter(logger, format)
	logging.SetBackend(backendFormatter)

	return log, f
}

func openFile(file string) (*os.File, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}
