package logs

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	//"panic"	"fatal"	"warn", "warning"	"info"	"debug"	"trace"
	Level        string
	ReportCaller bool
	Format       string
}

func Setup(c Config) {
	logrus.SetReportCaller(c.ReportCaller)
	if c.Level != "" {
		if level, err := logrus.ParseLevel(c.Level); err != nil {
			panic(err)
		} else {
			logrus.SetLevel(level)
		}
	}
	switch c.Format {
	case "json":
		logrus.SetFormatter(JSONFormatter)
	//case "text":
	default:
		logrus.SetFormatter(TextFormatter)
	}
}
