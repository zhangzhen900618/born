package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
)

var logger *Logger

func init() {
	New()
}

func New(opts ...Option) {
	o := getDefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	logger = &Logger{
		options: o,
	}
	logger.init()
}

type Logger struct {
	zerolog.Logger
	options
}

func (l *Logger) init() {
	lv, err := zerolog.ParseLevel(l.options.level)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(lv)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	multi := zerolog.MultiLevelWriter(l.getWriters()...)
	l.Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()
}

func (l *Logger) getWriters() []io.Writer {
	var writers []io.Writer
	if l.enableWrite {
		writer := &lumberjack.Logger{
			Filename:   l.fileName,
			MaxSize:    l.maxSize,
			MaxAge:     l.maxAge,
			MaxBackups: l.maxBackups,
		}
		writers = append(writers, writer)
	}
	if l.enableConsole {
		writers = append(writers, os.Stdout)
	}
	return writers
}

func Trace() *zerolog.Event {
	return logger.Trace()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Err(err error) *zerolog.Event {
	return logger.Err(err)
}
func Fatal() *zerolog.Event {
	return logger.Fatal()
}

func Panic() *zerolog.Event {
	return logger.Panic()
}
