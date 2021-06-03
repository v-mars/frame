package logger

import (
	"fmt"
	"frame/pkg/utils"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Logger     = logrus.New()
)
var logLevels = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
}

func Initial(level string) {
	formatter := &Formatter{
		LogFormat:       "",
		//LogFormat:       "%time% [%lvl%] %msg%",
		TimestampFormat: "2006-01-02 15:04:05",
	}
	//ginFormatter := &Formatter{
	//	LogFormat:       "%msg%",
	//	TimestampFormat: "2006-01-02 15:04:05",
	//}
	InitLog("DEBUG","std.log","./logs",formatter, Logger)

}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func InitOutToFile(logPath, logFileName string) error {
	logFilePath := path.Join(logPath, logFileName)
	logDirPath := path.Dir(logFilePath)
	if !utils.FileExists(logDirPath) {
		err := os.MkdirAll(logDirPath, os.ModePerm)
		if err != nil {
			msg := fmt.Sprintf("K8sDeployment log dir %s error: %s\n", logDirPath, err)
			fmt.Println(msg)
			return err
		}
	}
	return nil
}

func InitLog(LogLevel, logFileName,LogPath string, formatter *Formatter, loggerObj *logrus.Logger){

	level, ok := logLevels[strings.ToUpper(LogLevel)]
	if !ok {
		level = logrus.InfoLevel
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	loggerObj.SetFormatter(formatter)
	loggerObj.SetOutput(os.Stdout)
	loggerObj.SetLevel(level)

	// Output to file
	logFilePath := path.Join(LogPath, logFileName)
	if err := InitOutToFile(LogPath, logFileName); err!=nil{
		log.Fatal(err)
	}
	rotateFileHook, err := NewRotateFileHook(RotateFileConfig{
		Filename:   logFilePath,
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     7,
		LocalTime:  true,
		Level:      level,
		Formatter:  formatter,
	})
	if err != nil {
		fmt.Printf("log rotate hook error: %s\n", err)
		return
	}
	loggerObj.AddHook(rotateFileHook)
}

