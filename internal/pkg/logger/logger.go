package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(conf Config, appName, version, buildDate, commit string) (*zap.Logger, error) {
	zapConf := zap.Config{
		Level:       zapLevel(conf.Level),
		Development: conf.Debug,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig(),
		OutputPaths:      conf.GetOutput(),
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConf.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.Any("application", struct {
			Name      string `json:"name"`
			Version   string `json:"version"`
			BuildDate string `json:"buildDate"`
			Commit    string `json:"commit"`
		}{
			Name:      appName,
			Version:   version,
			BuildDate: buildDate,
			Commit:    commit,
		}),
	)
	return logger, nil
}

func zapLevel(level string) (l zap.AtomicLevel) {
	switch level {
	case DebugLevel:
		l = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		panic("unknown log level")
	}
	return l
}

func zapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "@log_name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stackTrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
