package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

// CustomLogger untuk menyimpan log ke file dengan rotasi menggunakan zap dan lumberjack
type CustomLogger struct {
	logger.Config
	ZapLogger *zap.Logger
}

// NewCustomLogger membuat instance CustomLogger dengan zap dan rotasi file
func NewCustomLogger(filePath string) (*CustomLogger, error) {
	// Pastikan direktori relatif ada
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Konfigurasi lumberjack untuk rotasi file
	lumberLogger := &lumberjack.Logger{
		Filename:   filePath, // Path relatif, misalnya "logs/gorm.log"
		MaxSize:    100,      // Ukuran maksimum file sebelum rotasi (dalam MB)
		MaxBackups: 3,        // Jumlah maksimum file cadangan
		MaxAge:     28,       // Usia maksimum file cadangan (dalam hari)
		Compress:   true,     // Kompres file cadangan
	}

	// Konfigurasi zap dengan lumberjack sebagai WriteSyncer
	writeSyncer := zapcore.AddSync(lumberLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // Encoder JSON
		writeSyncer,
		zapcore.InfoLevel, // Level minimum untuk zap
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	absPath, _ := filepath.Abs(filePath)
	fmt.Printf("Log file created at: %s\n", absPath)

	return &CustomLogger{
		Config: logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
		ZapLogger: zapLogger,
	}, nil
}

// LogMode mengatur level logging
func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info untuk log level INFO
func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.saveLog(ctx, "INFO", fmt.Sprintf(msg, data...), 0, 0, "")
	}
}

// Warn untuk log level WARN
func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.saveLog(ctx, "WARN", fmt.Sprintf(msg, data...), 0, 0, "")
	}
}

// Error untuk log level ERROR
func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.saveLog(ctx, "ERROR", fmt.Sprintf(msg, data...), 0, 0, fmt.Sprintf("%v", data))
	}
}

// Trace untuk log query SQL
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	// Abaikan query SELECT
	if os.Getenv("ENV") != "dev" {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT") {
			return
		}

		l.saveLog(ctx, "TRACE", sql, float64(elapsed.Nanoseconds())/1e6, rows, errorMsg)
	} else {
		if errorMsg != "" {
			// log.Printf("SQL : %s Error: %s", sql, errorMsg)
		}
		log.Printf("SQL : %s Error: %s", sql, errorMsg)
	}

}

// saveLog menyimpan log menggunakan zap
func (l *CustomLogger) saveLog(ctx context.Context, level, query string, duration float64, rows int64, err string) {

	fields := []zap.Field{
		zap.String("created_at", time.Now().Format(time.RFC3339)),
		zap.String("level", level),
		zap.String("query", query),
		zap.Float64("duration", duration),
		zap.Int64("rows", rows),
		zap.String("error", err),
	}
	switch level {
	case "INFO":
		l.ZapLogger.Info("gorm_log", fields...)
	case "WARN":
		l.ZapLogger.Warn("gorm_log", fields...)
	case "ERROR":
		l.ZapLogger.Error("gorm_log", fields...)
	case "TRACE":
		l.ZapLogger.Info("gorm_log", fields...)
	}
}
