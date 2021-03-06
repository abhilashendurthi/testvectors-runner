/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package logger implements utilities to instantiate and manipulate a new logger
*/
package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	//logger is instance of logrus
	logger                      = logrus.New()
	logDir                      = "/tmp/"
	fileName                    = "logfile_" + time.Now().Format(time.RFC3339) + ".log"
	fileOptions                 = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	filePermissions os.FileMode = 0666
	writer          io.Writer
)

//StandardLogger is a representation of logrus logger
type StandardLogger struct {
	*logrus.Logger
}

//NewLogger method to initialize logger
func NewLogger() *StandardLogger {
	var standardLogger = &StandardLogger{logger}
	//TODO - write inputs to config file and parse
	//output folder - default /tmp
	//also log to stdout - default true
	//log level - default warn

	//To display calling file and function
	standardLogger.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			return fmt.Sprintf("%s()", funcname), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	//Timestamp setting
	formatter.TimestampFormat = "2006-01-02T15:04:05.999Z07:00"

	//alternatively,
	//formatter.TimestampFormat = time.RFC3339Nano
	formatter.FullTimestamp = true
	formatter.ForceColors = true

	standardLogger.SetFormatter(formatter)

	//default output to stdout & /tmp
	file, _ := os.OpenFile(logDir+"/"+fileName, fileOptions, filePermissions)
	writer = io.MultiWriter(os.Stdout, file)
	standardLogger.SetOutput(writer)

	//default log level
	standardLogger.SetLevel(logrus.ErrorLevel)
	return standardLogger
}

//SetLogLevel changes the log verbosity
func (sl *StandardLogger) SetLogLevel(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warnf("Invalid log level")
	} else {
		sl.SetLevel(logLevel)
	}
}

//SetLogFolder sets the output folder
func (sl *StandardLogger) SetLogFolder(logDirPath string) {
	logFolder, err := os.Stat(logDirPath)
	if os.IsNotExist(err) || !logFolder.IsDir() {
		logrus.Warnf("Error opening log directory %s", err)
		writer = io.Writer(os.Stdout)
	} else {
		file, _ := os.OpenFile(logDirPath+"/"+fileName, fileOptions, filePermissions)
		writer = io.MultiWriter(os.Stdout, file)
	}
	sl.SetOutput(writer)
}
