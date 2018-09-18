package log

import (
	"errors"
	"fmt"
	_log "log"
	"os"
)

// Constants for logger.
const (
	LevelWrong = 0
	LevelDebug = 1
	LevelInfo  = 2
	LevelWarn  = 3
	LevelError = 4
	LevelCrit  = 5

	OptDebug = "DEBUG"
	OptInfo  = "INFO"
	OptWarn  = "WARN"
	OptError = "ERROR"
	OptCrit  = "CRIT"

	PrefixDebug = "[Debg] : "
	PrefixInfo  = "[Info] : "
	PrefixWarn  = "[Warn] : "
	PrefixError = "[Erro] : "
	PrefixCrit  = "[Crit] : "
)

// log is grobal Logger instance
var log *Logger

// Logger contains logger's info and basic logger instances.
type Logger struct {
	path *string
	fp   *os.File

	level int

	logDebug *_log.Logger
	logInfo  *_log.Logger
	logWarn  *_log.Logger
	logError *_log.Logger
	logCrit  *_log.Logger
}

// Init sets log path and logger's level.
func Init(path *string, level *string) error {
	logLevel := MapLevel(level)
	if logLevel == LevelWrong {
		return errors.New("Wrong log level")
	}

	if path != nil {
		logFp, err := os.OpenFile(*path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		logDebug := _log.New(logFp, PrefixDebug, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logInfo := _log.New(logFp, PrefixInfo, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logWarn := _log.New(logFp, PrefixWarn, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logError := _log.New(logFp, PrefixError, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logCrit := _log.New(logFp, PrefixCrit, _log.Ldate|_log.Ltime|_log.Lshortfile)

		log = &Logger{path: path, fp: logFp, level: logLevel, logDebug: logDebug,
			logInfo: logInfo, logWarn: logWarn, logError: logError, logCrit: logCrit}
	} else {
		logDebug := _log.New(os.Stdout, PrefixDebug, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logInfo := _log.New(os.Stdout, PrefixInfo, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logWarn := _log.New(os.Stdout, PrefixWarn, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logError := _log.New(os.Stdout, PrefixError, _log.Ldate|_log.Ltime|_log.Lshortfile)
		logCrit := _log.New(os.Stdout, PrefixCrit, _log.Ldate|_log.Ltime|_log.Lshortfile)

		log = &Logger{path: nil, fp: nil, level: logLevel, logDebug: logDebug,
			logInfo: logInfo, logWarn: logWarn, logError: logError, logCrit: logCrit}
	}

	return nil
}

// Clean clears the logger
func Clean() {
	if log.fp != nil {
		log.fp.Close()
	}
}

// MapLevel maps option to log level.
func MapLevel(opt *string) int {
	switch *opt {
	case OptDebug:
		return LevelDebug
	case OptInfo:
		return LevelInfo
	case OptWarn:
		return LevelWarn
	case OptError:
		return LevelError
	case OptCrit:
		return LevelCrit
	default:
		return LevelWrong
	}
}

// Debugf works the same as printf with debug prefix
func Debugf(format string, v ...interface{}) {
	if log.level > LevelDebug {
		return
	}

	log.logDebug.Output(2, fmt.Sprintf(format, v...))
}

// Debug works the same as print with debug prefix
func Debug(v ...interface{}) {
	if log.level > LevelDebug {
		return
	}

	log.logDebug.Output(2, fmt.Sprint(v...))
}

// Debugln works the same as println with debug prefix
func Debugln(v ...interface{}) {
	if log.level > LevelDebug {
		return
	}

	log.logDebug.Output(2, fmt.Sprint(v...))
}

// Infof works the same as printf with information prefix
func Infof(format string, v ...interface{}) {
	if log.level > LevelInfo {
		return
	}

	log.logInfo.Output(2, fmt.Sprintf(format, v...))
}

// Info works the same as print with information prefix
func Info(v ...interface{}) {
	if log.level > LevelInfo {
		return
	}

	log.logInfo.Output(2, fmt.Sprint(v...))
}

// Infoln works the same as println with information prefix
func Infoln(v ...interface{}) {
	if log.level > LevelInfo {
		return
	}

	log.logInfo.Output(2, fmt.Sprint(v...))
}

// Warnf works the same as printf with warning prefix
func Warnf(format string, v ...interface{}) {
	if log.level > LevelWarn {
		return
	}

	log.logWarn.Output(2, fmt.Sprintf(format, v...))
}

// Warn works the same as print with warning prefix
func Warn(v ...interface{}) {
	if log.level > LevelWarn {
		return
	}

	log.logWarn.Output(2, fmt.Sprint(v...))
}

// Warnln works the same as println with warning prefix
func Warnln(v ...interface{}) {
	if log.level > LevelWarn {
		return
	}

	log.logWarn.Output(2, fmt.Sprint(v...))
}

// Errorf works the same as printf with error prefix
func Errorf(format string, v ...interface{}) {
	if log.level > LevelError {
		return
	}

	log.logError.Output(2, fmt.Sprintf(format, v...))
}

// Error works the same as print with error prefix
func Error(v ...interface{}) {
	if log.level > LevelError {
		return
	}

	log.logError.Output(2, fmt.Sprint(v...))
}

// Errorln works the same as println with error prefix
func Errorln(v ...interface{}) {
	if log.level > LevelError {
		return
	}

	log.logError.Output(2, fmt.Sprint(v...))
}

// Critf works the same as printf with critical prefix
func Critf(format string, v ...interface{}) {
	if log.level > LevelCrit {
		return
	}

	log.logCrit.Output(2, fmt.Sprintf(format, v...))
}

// Crit works the same as print with critical prefix
func Crit(v ...interface{}) {
	if log.level > LevelCrit {
		return
	}

	log.logCrit.Output(2, fmt.Sprint(v...))
}

// Critln works the same as println with critical prefix
func Critln(v ...interface{}) {
	if log.level > LevelCrit {
		return
	}

	log.logCrit.Output(2, fmt.Sprint(v...))
}
