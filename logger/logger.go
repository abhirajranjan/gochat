package logger

import (
	"gochat/config"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	GetLoggerLevel() zapcore.Level
	AddWriter(newWriter io.Writer) *logger
	WithName(name string)
	InitLogger() *logger

	// logger methods
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Printf(template string, args ...interface{})
	Warn(args ...interface{})
	WarnMsg(msg string, err error)
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Sync() error
}

type logger struct {
	level       string
	dev         bool
	encoder     string
	writer      []zapcore.WriteSyncer
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func NewLogger(conf config.LoggerConfig) *logger {
	l := logger{level: conf.Level, dev: conf.Isdev, encoder: conf.Encoder}
	return &l
}

func (l *logger) GetLoggerLevel() zapcore.Level {
	return loggerLevelMap[l.level]
}

func (l *logger) AddWriter(newWriter io.Writer) *logger {
	l.writer = append(l.writer, zapcore.AddSync(newWriter))
	return l
}

func (l *logger) getWriteSyncer() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(l.writer...)
}

func (l *logger) getEncoderConfig() (conf zapcore.EncoderConfig) {
	if l.dev {
		conf = zap.NewDevelopmentEncoderConfig()
	} else {
		conf = zap.NewProductionEncoderConfig()
	}
	conf.NameKey = "[SERVICE]"
	conf.TimeKey = "[TIME]"
	conf.LevelKey = "[LEVEL]"
	conf.CallerKey = "[LINE]"
	conf.MessageKey = "[MESSAGE]"
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	conf.EncodeCaller = zapcore.ShortCallerEncoder
	conf.EncodeDuration = zapcore.StringDurationEncoder

Select:
	switch l.encoder {
	case "console":
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		conf.EncodeCaller = zapcore.FullCallerEncoder
		conf.ConsoleSeparator = " | "
	case "json":
		conf.FunctionKey = "[CALLER]"
		conf.EncodeName = zapcore.FullNameEncoder
	default:
		l.encoder = "console"
		goto Select
	}
	return
}

func (l *logger) getEncoder() (encoder zapcore.Encoder) {
Select:
	switch l.encoder {
	case "console":
		encoder = zapcore.NewConsoleEncoder(l.getEncoderConfig())
	case "json":
		encoder = zapcore.NewJSONEncoder(l.getEncoderConfig())
	default:
		l.encoder = "console"
		goto Select
	}
	return
}
func (l *logger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

func (l *logger) InitLogger() *logger {
	level := zap.NewAtomicLevelAt(l.GetLoggerLevel())
	core := zapcore.NewCore(l.getEncoder(), l.getWriteSyncer(), level)
	l.logger = zap.New(core)
	l.sugarLogger = l.logger.Sugar()
	return l
}
