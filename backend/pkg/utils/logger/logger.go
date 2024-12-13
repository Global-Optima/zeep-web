package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevelStr   = "debug"
	InfoLevelStr    = "info"
	WarningLevelStr = "warning"
	ErrorLevelStr   = "error"
	LogMaxSize      = 1024
	LogMaxBackups   = 30
	LogMaxAge       = 90
	LogCompress     = true
)

var (
	loggerGin            *zap.Logger
	sugaredLoggerService *zap.SugaredLogger
)

func InitLoggers(logLevel, ginLogFile, serviceLogFile string, dev bool) error {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return err
	}

	ginLog := &lumberjack.Logger{
		Filename:   ginLogFile,
		MaxSize:    LogMaxSize,
		MaxBackups: LogMaxBackups,
		MaxAge:     LogMaxAge,
		Compress:   LogCompress,
	}

	serviceLog := &lumberjack.Logger{
		Filename:   serviceLogFile,
		MaxSize:    LogMaxSize,
		MaxBackups: LogMaxBackups,
		MaxAge:     LogMaxAge,
		Compress:   LogCompress,
	}

	consoleEncoder := buildConsoleEncoder(dev)
	fileEncoder := buildFileEncoder()

	// Console output core
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout), // console output
		level,
	)

	// File output core
	ginFileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(ginLog), // file output
		level,
	)

	serviceFileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(serviceLog), // file output
		level,
	)

	// Tee both cores so logs go to console and file
	serviceCore := zapcore.NewTee(consoleCore, serviceFileCore)
	ginCore := zapcore.NewTee(consoleCore, ginFileCore)

	loggerGin = zap.New(ginCore, zap.AddCaller(), zap.AddCallerSkip(1))

	loggerService := zap.New(serviceCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(loggerService)
	sugaredLoggerService = NewSugar("service", loggerService)

	// Test log entry
	zap.L().Info("Logger initialized", zap.String("level", logLevel))

	return nil
}

func parseLogLevel(logLevel string) (zapcore.Level, error) {
	switch logLevel {
	case DebugLevelStr:
		return zapcore.DebugLevel, nil
	case InfoLevelStr:
		return zapcore.InfoLevel, nil
	case WarningLevelStr:
		return zapcore.WarnLevel, nil
	case ErrorLevelStr:
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level %s", logLevel)
	}
}

func buildConsoleEncoder(dev bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if dev {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func buildFileEncoder() zapcore.Encoder {
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	fileEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(fileEncoderConfig)
}

func NewSugar(name string, logger *zap.Logger) *zap.SugaredLogger {
	if logger == nil {
		logger = zap.NewNop()
		return logger.Sugar()
	}

	sugaredLoggerService = logger.Named(name).Sugar()

	return sugaredLoggerService
}

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path // fallback if no named route
		}

		c.Next() // process the request

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		requestID := c.GetString("request_id") // if you set request_id in another middleware

		// Determine log level based on status code
		var lvl zapcore.Level
		switch {
		case status >= 500:
			lvl = zap.ErrorLevel
		case status >= 400:
			lvl = zap.WarnLevel
		default:
			lvl = zap.InfoLevel
		}

		logger := loggerGin // or zap.L() if you prefer the global
		if !logger.Core().Enabled(lvl) {
			return
		}

		// Add fields
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Request content length (if known)
		if c.Request.ContentLength > 0 {
			fields = append(fields, zap.Int64("req_size", c.Request.ContentLength))
		}
		msg := fmt.Sprintf("%s %s", method, path)
		logMessage(msg, lvl, fields)
	}
}

func logMessage(msg string, lvl zapcore.Level, fields []zap.Field) {
	// Log the request line with the appropriate level

	switch lvl {
	case zap.InfoLevel:
		loggerGin.Info(msg, fields...)
	case zap.WarnLevel:
		loggerGin.Warn(msg, fields...)
	case zap.ErrorLevel:
		loggerGin.Error(msg, fields...)
	}
}

func GetZapSugaredLogger() *zap.SugaredLogger {
	return sugaredLoggerService
}
