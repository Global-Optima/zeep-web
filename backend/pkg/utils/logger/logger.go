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
)

var (
	globalLogger *zap.Logger
)

func Init(logLevel, logFile string, dev bool) error {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return err
	}

	rotatingLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     90,
		Compress:   true,
	}

	encoder := buildEncoder(dev)

	// Console output core
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout), // console output
		level,
	)

	// File output core
	fileCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(rotatingLogger), // file output
		level,
	)

	// Tee both cores so logs go to console and file
	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
	globalLogger = logger

	// Test log entry
	zap.L().Info("Logger initialized", zap.String("file", logFile), zap.String("level", logLevel))

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

func buildEncoder(dev bool) zapcore.Encoder {
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

func NewSugar(name string) *zap.SugaredLogger {
	if globalLogger == nil {
		logger := zap.NewNop()
		return logger.Sugar()
	}
	return globalLogger.Named(name).Sugar()
}

func ZapRequestLogger() gin.HandlerFunc {
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

		logger := globalLogger // or zap.L() if you prefer the global
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

		// Log the request line with the appropriate level
		msg := fmt.Sprintf("%s %s", method, path)
		switch lvl {
		case zap.InfoLevel:
			logger.Info(msg, fields...)
		case zap.WarnLevel:
			logger.Warn(msg, fields...)
		case zap.ErrorLevel:
			logger.Error(msg, fields...)
		}
	}
}
