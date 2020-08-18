package zapLogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"shebinbin.com/alertSyslog/config"
	"time"
)

var zlog *zap.SugaredLogger

func LoggerFactory() *zap.SugaredLogger {
	return zlog
}

func logLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.DebugLevel
	}
}

func init() {

	//fileName := "nero.log"

	/*writerSy := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1 << 30, // megabytes 1G
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
		Compress:   true,
	})*/

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // writerSy),
		logLevel(config.LogLevel),
	)
	//logger := zap.New(core, zap.AddCaller())
	//logger.Info("info")
	zlog = zap.New(core, zap.AddCaller()).Sugar()
	zlog.Info("日志组件初始化完成,日志等级为：", config.LogLevel)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
