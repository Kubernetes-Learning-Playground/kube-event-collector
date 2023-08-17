package logCollector

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/practice/kube-event/pkg/model"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	Logger *StructLogger
)

type StructLogger struct {
	*logrus.Logger
}

func NewStructLogger(logPath string) *StructLogger {
	return InitLog(logPath)
}

func InitLog(path string) *StructLogger {

	if err := os.MkdirAll(path, 0777); err != nil {
		panic(fmt.Errorf("create log dir: %s error", path))
	}
	// formatter
	formatter := &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return strconv.Itoa(frame.Line), frame.Function
		},
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyFile:  "filename",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "lineno",
			logrus.FieldKeyLevel: "level",
		},
	}
	// 默认每天分割，保存 7 天
	logf, _ := rotatelogs.New(
		filepath.Join(path, "event.%Y%m%d.log"),
		rotatelogs.WithLinkName(filepath.Join(path, "event.log")),
	)
	fiHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logf,
		logrus.InfoLevel:  logf,
		logrus.WarnLevel:  logf,
		logrus.ErrorLevel: logf,
		logrus.FatalLevel: logf,
		logrus.PanicLevel: logf,
	}, formatter)


	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        false,
		TimestampFormat: "2006-01-02 15:04:05",
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return fmt.Sprintf(" (%s:%d)", frame.Function, frame.Line)
		},
	})
	log.AddHook(fiHook)
	Logger = &StructLogger{
		Logger: log,
	}
	return Logger
}

func (sl *StructLogger) EventLog(event *model.Event) {
	sl.Logger.WithFields(logrus.Fields{
		"name":        event.Name,
		"reason":      event.Reason,
		"message":     event.Message,
		"host":        event.Host,
		"kind":        event.Kind,
		"event_level": event.Type,
		"count":       event.Count,
	}).Info(event.Message)
}