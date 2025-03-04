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
	LogMaxSize      = 1024 // Max log file size in MB before rotation
	LogMaxBackups   = 30   // Max number of old log files to keep
	LogMaxAge       = 90   // Max days to retain old log files
	LogCompress     = true // Compress old log files
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

// InitLogger initializes a single logger for the entire application
func InitLogger(logLevel, logFile string, dev bool) error {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return err
	}

	// Log rotation configuration
	logWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    LogMaxSize,
		MaxBackups: LogMaxBackups,
		MaxAge:     LogMaxAge,
		Compress:   LogCompress,
	}

	consoleEncoder := buildConsoleEncoder(dev)
	fileEncoder := buildFileEncoder()

	// Define output cores
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(logWriter), level)

	// Merge cores so logs go to both file and console
	core := zapcore.NewTee(consoleCore, fileCore)

	// Create the final logger
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)

	// Create a sugared logger for structured logging
	sugaredLogger = logger.Sugar()

	zap.L().Info("Logger initialized", zap.String("level", logLevel))
	return nil
}

// parseLogLevel converts string log level to zapcore.Level
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

// buildConsoleEncoder sets up console output formatting
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

// buildFileEncoder sets up file output formatting
func buildFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

// ZapLoggerMiddleware logs all HTTP requests in a single log file
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

		if !logger.Core().Enabled(lvl) {
			return
		}

		// Log fields
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

		// Log request content length if available
		if c.Request.ContentLength > 0 {
			fields = append(fields, zap.Int64("req_size", c.Request.ContentLength))
		}

		msg := fmt.Sprintf("%s %s", method, path)
		logMessage(msg, lvl, fields)
	}
}

// logMessage logs HTTP requests at the appropriate level
func logMessage(msg string, lvl zapcore.Level, fields []zap.Field) {
	switch lvl {
	case zap.InfoLevel:
		logger.Info(msg, fields...)
	case zap.WarnLevel:
		logger.Warn(msg, fields...)
	case zap.ErrorLevel:
		logger.Error(msg, fields...)
	}
}

// GetZapSugaredLogger returns a global sugared logger instance
func GetZapSugaredLogger() *zap.SugaredLogger {
	return sugaredLogger
}
