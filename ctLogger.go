package ctLogger

import (
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kdar/logrus-cloudwatchlogs"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

var log *logrus.Logger

func InitLogging(level string, sndToCloudWatch bool, sess *session.Session, group string) {

	log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	switch strings.ToLower(level) {
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warn", "warning":
		log.SetLevel(logrus.WarnLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	default:
		log.SetLevel(logrus.TraceLevel)
	}

	if sndToCloudWatch == true {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		pc, _, _, _ := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		module := details.Name()[0:strings.Index(details.Name(), "/")]
		stream := fmt.Sprintf(time.Now().Format("2006/01/02/15/04/05/001/") + module + "/%x", b)
		hook, err := logrus_cloudwatchlogs.NewHook(group, stream, sess)
		if err != nil {
			log.Fatal(err)
		}
		log.Hooks.Add(hook)
		log.Out = ioutil.Discard
		log.Formatter = logrus_cloudwatchlogs.NewProdFormatter()
	}
}

func Error(err error){
	pc, fn, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	log.WithFields(logrus.Fields{
		"line":   line,
		"function": details.Name(),
		"filename": fn,
	}).Error(err.Error())
}

func Info(msg string){
	pc, fn, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	log.WithFields(logrus.Fields{
		"line":   line,
		"function": details.Name(),
		"filename": fn,
	}).Info(msg)
}

func Debug(msg string){
	pc, fn, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	log.WithFields(logrus.Fields{
		"line":   line,
		"function": details.Name(),
		"filename": fn,
	}).Debug(msg)
}

func Trace(msg string){
	pc, fn, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	log.WithFields(logrus.Fields{
		"line":   line,
		"function": details.Name(),
		"filename": fn,
	}).Trace(msg)
}