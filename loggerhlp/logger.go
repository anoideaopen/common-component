package loggerhlp

import (
	"log"

	"github.com/newity/glog"
	"github.com/newity/glog/logr"
	"github.com/newity/glog/std"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	logTypeStd       = "std"
	logTypeLRTxt     = "lr-txt"
	logTypeLRTxtDev  = "lr-txt-dev"
	logTypeLRJson    = "lr-json"
	logTypeLRJsonDev = "lr-json-dev"
)

// CreateLogger creates logger
func CreateLogger(loggerType, logLvl string) (glog.Logger, error) {
	switch loggerType {
	case logTypeStd:
		return createStdLogger(logLvl)

	case logTypeLRTxt, logTypeLRTxtDev,
		logTypeLRJson, logTypeLRJsonDev:
		return createLrLogger(loggerType, logLvl)

	default:
		return nil,
			errors.WithStack(
				errors.Errorf("failed to create logger: unknown type %s", loggerType))
	}
}

func createStdLogger(logLvl string) (glog.Logger, error) {
	var ll std.Level
	switch logLvl {
	case "trace":
		ll = std.LevelTrace
	case "debug":
		ll = std.LevelDebug
	case "info":
		ll = std.LevelInfo
	case "warning":
		ll = std.LevelWarning
	case "error":
		ll = std.LevelError
	default:
		return nil, errors.WithStack(errors.Errorf(
			"failed to create logger: unknown log level %s",
			logLvl,
		))
	}

	return std.New(log.Default(), ll), nil
}

func createLrLogger(loggerType, logLvl string) (glog.Logger, error) {
	ll, err := logrus.ParseLevel(logLvl)
	if err != nil {
		return nil,
			errors.WithStack(
				errors.Wrap(err, "failed to create logger"))
	}

	lrLogger := logrus.StandardLogger()
	lrLogger.SetLevel(ll)

	var lrFormatter logrus.Formatter
	switch loggerType {
	case logTypeLRTxt, logTypeLRTxtDev:
		lrFormatter = &logrus.TextFormatter{
			ForceColors:               loggerType == logTypeLRTxtDev,
			DisableColors:             loggerType != logTypeLRTxtDev,
			EnvironmentOverrideColors: false,
			FullTimestamp:             true,
			PadLevelText:              loggerType == logTypeLRTxtDev,
			QuoteEmptyFields:          false,
		}
	case logTypeLRJson, logTypeLRJsonDev:
		lrFormatter = &logrus.JSONFormatter{
			PrettyPrint: loggerType == logTypeLRJsonDev,
		}
	}

	if lrFormatter != nil {
		logrus.SetFormatter(lrFormatter)
	}

	return logr.New(lrLogger), nil
}
