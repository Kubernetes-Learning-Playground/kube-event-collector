package logCollector

import (
	"github.com/practice/kube-event/pkg/model"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Logger *StructLogger
)

type StructLogger struct {
	*logrus.Logger
}

func init() {
	lg := logrus.New()
	Logger = &StructLogger{
		Logger: lg,
	}

	Logger.SetFormatter(&logrus.JSONFormatter{}) // 设置 format json
	// Output to stdout instead of the default stderr
	logfile, _ := os.OpenFile("./logrus.mylog", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	Logger.SetOutput(logfile)
	//Logger.SetOutput(os.Stdout)
}

func (sl *StructLogger) EventLog(event *model.Event) {
	sl.Logger.WithFields(logrus.Fields{
		//"namespace":      namespace,
		"pod":         event.Name,
		"reason":      event.Reason,
		"message":     event.Message,
		"host":        event.Host,
		"kind":        event.Kind,
		"event_level": event.Type,
		"count":       event.Count,
	}).Info(event.Message)
}
