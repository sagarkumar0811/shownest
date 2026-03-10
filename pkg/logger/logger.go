package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shownest/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type Config struct {
	Environment string
	ServiceName string
	LogDir      string
	LogLevel    string
}

func Init(ctx context.Context, provider config.ConfigProvider) error {
	appName, err := provider.Get(ctx, "app")
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	envData, err := provider.Get(ctx, "env")
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}

	logDirData, err := provider.Get(ctx, "logdir")
	if err != nil {
		return fmt.Errorf("failed to get logdir: %w", err)
	}

	logLevelData, err := provider.Get(ctx, "loglevel")
	if err != nil {
		return fmt.Errorf("failed to get loglevel: %w", err)
	}

	var app, env, logDir, logLevel string
	json.Unmarshal(appName, &app)
	json.Unmarshal(envData, &env)
	json.Unmarshal(logDirData, &logDir)
	json.Unmarshal(logLevelData, &logLevel)

	cfg := Config{
		Environment: env,
		ServiceName: app,
		LogDir:      logDir,
		LogLevel:    logLevel,
	}

	return initLogger(cfg)
}

func initLogger(cfg Config) error {
	var cores []zapcore.Core

	level := zapcore.InfoLevel
	if cfg.LogLevel != "" {
		if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
			return fmt.Errorf("invalid log level %q: %w", cfg.LogLevel, err)
		}
	} else if cfg.Environment == config.EnvironmentLocal {
		level = zapcore.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	if cfg.Environment == config.EnvironmentLocal {
		logsDir := cfg.LogDir
		if logsDir == "" {
			return fmt.Errorf("logdir is required for local environment")
		}

		if _, err := os.Stat(logsDir); os.IsNotExist(err) {
			return fmt.Errorf("logs directory does not exist: %s", logsDir)
		}

		dateStr := time.Now().Format(config.DateFormat)
		logFile := filepath.Join(logsDir, fmt.Sprintf("%s_%s.log", cfg.ServiceName, dateStr))
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			level,
		)
		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("app", cfg.ServiceName)),
	)

	globalLogger = logger
	zap.ReplaceGlobals(logger)

	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

func WithContext(ctx context.Context) *zap.Logger {
	logger := Get()
	if requestId := getRequestId(ctx); requestId != "" {
		logger = logger.With(zap.String("requestId", requestId))
	}
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func getRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestId, ok := ctx.Value("requestId").(string); ok {
		return requestId
	}
	return ""
}
