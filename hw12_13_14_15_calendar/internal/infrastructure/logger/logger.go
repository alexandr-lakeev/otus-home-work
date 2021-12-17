package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logg *zap.Logger
}

func New(level string) (*Logger, error) {
	logg, err := initLogger(level)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logg: logg,
	}, nil
}

func (l Logger) Debug(msg string) {
	l.logg.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.logg.Info(msg)
}

func (l Logger) Warning(msg string) {
	l.logg.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logg.Error(msg)
}

func (l Logger) Panic(msg string) {
	l.logg.Panic(msg)
}

func initLogger(level string) (*zap.Logger, error) {
	levelMap := map[string]zapcore.Level{
		"DEBUG":   zap.DebugLevel,
		"INFO":    zap.InfoLevel,
		"WARNING": zap.WarnLevel,
		"ERROR":   zap.ErrorLevel,
		"PANIC":   zap.PanicLevel,
	}

	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewDevelopmentEncoderConfig()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	zapLevel, ok := levelMap[level]
	if !ok {
		return nil, fmt.Errorf("wrong log level: %s", level)
	}

	atom.SetLevel(zapLevel)

	return logger, nil
}
