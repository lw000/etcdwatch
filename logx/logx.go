package logx

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var L *zap.Logger

func Init(logPath string, logLevel string) {
	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    100,     // 每个日志文件保存10M，默认 100M
		MaxAge:     30,      // 保留30个备份，默认不限
		MaxBackups: 7,       // 保留7天，默认不限
		Compress:   true,    // 是否压缩，默认不压缩
	}

	fileWrite := zapcore.AddSync(&hook)

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWrite, level)

	consoleCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level)

	allCore := zapcore.NewTee(fileCore, consoleCore)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	development := zap.Development()
	filed := zap.Fields(zap.String("serviceName", "etcdwatch"))
	L = zap.New(allCore, caller, development, filed)
}
