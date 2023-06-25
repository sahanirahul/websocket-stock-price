package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

// this is a copy of zap level, can use zap.Level instead
const (
	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	PanicLevel
)

type Fields map[string]interface{}

type ILogger interface {
	WriteLogs(ctx context.Context, msg string, level Level, fields Fields)
}

type zlogger struct {
	logger *zap.Logger
}

var once sync.Once

// making logger exportable so that it can be used from anywhere
var Logger *zlogger

func NewLogger() ILogger {
	once.Do(func() {
		initializeLogger()
	})
	return Logger
}

func (l *zlogger) normalizeFields(fields Fields) {
	for key := range fields {
		if fields[key] == nil {
			delete(fields, key)
			continue
		}
		switch val := fields[key].(type) {
		case fmt.Stringer:
			fields[key] = val.String()
		case error:
			fields[key] = val.Error()
		default:
			b, _ := json.Marshal(val)
			fields[key] = string(b)
		}
	}
}

func (l *zlogger) zapFields(fields Fields) []zapcore.Field {
	if len(fields) == 0 {
		return nil
	}
	var zapFields []zapcore.Field
	for key, val := range fields {
		zapFields = append(zapFields, zap.Any(key, val))
	}
	return zapFields
}

func (l *zlogger) WriteLogs(ctx context.Context, msg string, level Level, fields Fields) {
	l.normalizeFields(fields)
	zapFields := l.zapFields(fields)
	switch level {
	case InfoLevel:
		//do info logging
		l.logger.Info(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case PanicLevel:
		l.logger.Panic(msg, zapFields...)
	}
}

func createFile(filepath string) {
	// Create a file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func initializeLogger() error {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	today := time.Now().Format("2006-01-02")
	logdir := os.Getenv("LOGDIR")

	if len(logdir) == 0 {
		pwd, _ := os.Getwd()
		logdir := path.Join(pwd, "logs")
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			log.Fatal(err)
		}

	}
	logfilePath := path.Join(path.Join(logdir, today))
	createFile(logfilePath)
	logFile, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		// will log debuglevel in file
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		// will log infolevel on console
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)
	Logger = &zlogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel)),
	}
	return nil
}
