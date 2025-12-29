package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewProduction creates a production logger (JSON format) with specified level
func NewProduction(level string) (Logger, error) {
	zapLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build production logger: %w", err)
	}

	return &zapLogger{
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}, nil
}

// NewDevelopment creates a development logger (console format) with specified level
func NewDevelopment(level string) (Logger, error) {
	zapLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build development logger: %w", err)
	}

	return &zapLogger{
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}, nil
}

// NewLogger creates a logger based on environment and level
// env: "development" or "production"
// level: "debug", "info", "warn", "error"
func NewLogger(env, level string) (Logger, error) {
	if env == "development" {
		return NewDevelopment(level)
	}

	return NewProduction(level)
}

// NewNop returns a no-op logger for testing
func NewNop() Logger {
	nopLogger := zap.NewNop()
	return &zapLogger{
		logger:        nopLogger,
		sugaredLogger: nopLogger.Sugar(),
	}
}

// Must panics if logger creation fails
func Must(logger Logger, err error) Logger {
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}

	return logger
}

// parseLevel converts string level to zapcore.Level
func parseLevel(level string) (zapcore.Level, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return zapLevel, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	return zapLevel, nil
}
