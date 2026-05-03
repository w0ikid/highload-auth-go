package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создает новый экземпляр логгера.
// Поддерживает уровни: debug, info, warn, error, fatal.
// Если передать "prod", будет использован стандартный Production конфиг (Info).
func New(level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		// Если это не стандартный уровень, проверяем на "prod"
		if level == "prod" {
			zapLevel = zap.InfoLevel
		} else {
			zapLevel = zap.InfoLevel
		}
	}

	var cfg zap.Config
	if level == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z07:00")

		cfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(zapLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoderCfg,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	return cfg.Build(
		zap.WithCaller(false),
		zap.AddStacktrace(zap.ErrorLevel + 1),
	)
}