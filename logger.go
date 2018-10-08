package main

import (
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (cfg *LoggerConfig) Build() (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = cfg.Level
	zapCfg.OutputPaths = cfg.OutputPaths
	zapCfg.ErrorOutputPaths = cfg.ErrorOutputPaths
	return zapCfg.Build()
}

// SetLogger uses the provided logger to replace zap's global loggers and
// hijack output from the standard library's "log" package. It returns a
// function to undo these changes.
func SetLogger(log *zap.Logger) func() {
	undoGlobals := zap.ReplaceGlobals(log)
	undoHijack := zap.RedirectStdLog(log)
	return func() {
		undoGlobals()
		undoHijack()
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}
	opts = append(opts, zap.Development())
	opts = append(opts, zap.AddCaller())

	stackLevel := zap.ErrorLevel
	//stackLevel = zap.WarnLevel
	opts = append(opts, zap.AddStacktrace(stackLevel))

	opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSampler(core, time.Second, 100, 100)
	}))

	return opts
}

func BuildLogger(cfg *LoggerConfig) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  true,
	})

	encoder_cfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caller",
		MessageKey:     "Msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewJSONEncoder(encoder_cfg),
		w,
		cfg.Level,
	)

	logger := zap.New(core, buildOptions(w)...)

	SetLogger(logger)

	return
}
