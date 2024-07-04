package loggerhlp

import (
	"errors"
	"fmt"
	"log"

	"github.com/anoideaopen/glog"
	"github.com/anoideaopen/glog/logr"
	"github.com/anoideaopen/glog/std"
	"github.com/sirupsen/logrus"
)

const (
	logTypeStd       = "std"
	logTypeLRTxt     = "lr-txt"
	logTypeLRTxtDev  = "lr-txt-dev"
	logTypeLRJson    = "lr-json"
	logTypeLRJsonDev = "lr-json-dev"
	// logTypeGCP is build especially for Google Cloud Engine Logs
	// log in json, 'level' field named 'severity'.
	logTypeGCP = "gcp"
)

// CreateLogger creates logger
func CreateLogger(loggerType, logLvl string) (glog.Logger, error) {
	switch loggerType {
	case logTypeStd:
		return createStdLogger(logLvl)

	case logTypeLRTxt, logTypeLRTxtDev,
		logTypeLRJson, logTypeLRJsonDev, logTypeGCP:
		return createLrLogger(loggerType, logLvl)

	default:
		return nil,
			errors.New("failed to create logger: unknown type " + loggerType)
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
		return nil, errors.New("failed to create logger: unknown log level " + logLvl)
	}

	return std.New(log.Default(), ll), nil
}

func createLrLogger(loggerType, logLvl string) (glog.Logger, error) {
	ll, err := logrus.ParseLevel(logLvl)
	if err != nil {
		return nil,
			fmt.Errorf("failed to create logger: %w", err)
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
	case logTypeGCP:
		lrFormatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "severity",
			},
		}
	}

	if lrFormatter != nil {
		logrus.SetFormatter(lrFormatter)
	}

	return logr.New(lrLogger), nil
}
