package g

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

// InitLog 初始化log
func InitLog(level string) (err error) {
	switch level {
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	default:
		Log.Fatal("[logger.go] log conf only allow [info, debug, warn], please check your confguire")
	}
	Log.SetFormatter(&nested.Formatter{TimestampFormat: "2006-01-02 15:04:05"})

	return
}
