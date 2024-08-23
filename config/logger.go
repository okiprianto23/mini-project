package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
)

// NewLogger creates and returns a new zap.Logger instance with a specific configuration.
func NewLogger() (*LoggerCustom, error) {
	config := zap.Config{
		Encoding:         "json",                                                    // Format output log
		EncoderConfig:    zap.NewProductionEncoderConfig(),                          // Konfigurasi encoder default
		Level:            zap.NewAtomicLevelAt(Int8ToZapLevel(AppConfig.Log.Level)), // Level log
		OutputPaths:      []string{"stdout"},                                        // Output log ke stdout
		ErrorOutputPaths: []string{"stderr"},                                        // Output error ke stderr
	}

	// Kustomisasi format encoder
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.StacktraceKey = ""

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	//set with json
	defaultLogger := DefautLogger(logger)
	defaultLogger.Set("ip", "-")
	defaultLogger.Set("pid", os.Getpid())
	defaultLogger.Set("thread", "-")
	defaultLogger.Set("request_id", "-")
	defaultLogger.Set("source", "-")
	defaultLogger.Set("access_token", "-")
	defaultLogger.Set("application", AppConfig.Server.Application)
	defaultLogger.Set("version", AppConfig.Server.Version)
	defaultLogger.Set("code", "-")
	defaultLogger.Set("user_id", "-")
	defaultLogger.Set("url", "-")
	defaultLogger.Set("method", "-")

	return defaultLogger, nil
}

type LoggerCustom struct {
	defaultLogger *zap.Logger
	Logger        *zap.Logger
	keys          map[string]interface{}
	ModelLogger   *modelLogger
}

type modelLogger struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id" `

	Version     string `json:"version"`
	Application string `json:"application"`
	IP          string `json:"ip"`
	PID         int    `json:"pid"`
	Thread      string `json:"thread"`

	RequestID string `json:"request_id"`
	Source    string `json:"source"`

	ProcessingTime int64  `json:"processing_time"`
	ByteIn         int    `json:"byte_in"`
	ByteOut        int    `json:"byte_out"`
	Status         int    `json:"status" `
	Code           string `json:"code"`
	Url            string `json:"url"`
	Method         string `json:"method"`
}

func DefautLogger(lg *zap.Logger) *LoggerCustom {
	return &LoggerCustom{
		defaultLogger: lg,
		keys:          make(map[string]interface{}),
		ModelLogger:   &modelLogger{},
	}
}

func (log *LoggerCustom) Set(key string, value interface{}) *LoggerCustom {
	//set Logger by default
	log.Logger = log.defaultLogger

	// Perbarui atau tambahkan field ke map keys
	log.keys[key] = value

	// Buat Logger baru dengan field yang diperbarui
	fields := make([]zap.Field, 0, len(log.keys))
	for k, v := range log.keys {
		fields = append(fields, zap.Any(k, v))
	}

	// Buat Logger baru dengan field-field yang ada
	log.Logger = log.Logger.With(fields...)

	// Update modelLogger jika kunci ada dalam tag JSON
	v := reflect.ValueOf(log.ModelLogger).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == key {
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(value))
			}
		}
	}

	return log
}

// Int8ToZapLevel converts int8 to zapcore.Level.
func Int8ToZapLevel(level int8) zapcore.Level {
	switch level {
	case 0:
		return zapcore.DebugLevel
	case 1:
		return zapcore.InfoLevel
	case 2:
		return zapcore.WarnLevel
	case 3:
		return zapcore.ErrorLevel
	case 4:
		return zapcore.DPanicLevel
	case 5:
		return zapcore.PanicLevel
	case 6:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel // Default to InfoLevel if out of range
	}
}
