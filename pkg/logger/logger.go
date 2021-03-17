package logger

import (
	"IntelligentTransfer/config"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var ZapLogger *zap.Logger

// 初始化zap logger
func init() {
	//从日志中读取路径
	now := time.Now()
	logPath := config.GetConfig().GetString("log.log_path")
	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf(logPath+"%04d-"+"%02d-"+"%02d.log", now.Year(), now.Month(), now.Day()),
		MaxSize:    128,
		MaxAge:     7,
		MaxBackups: 30,
		Compress:   false,
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	// 如果是开发环境，同时在控制台上也输出
	writes = append(writes, zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	ZapLogger = zap.New(core, caller, development)
	ZapLogger.Info("log 初始化成功")
}

//封装Sugar print方法
func Debugf(format string, v ...interface{}) {
	ZapLogger.Sugar().Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	ZapLogger.Sugar().Infof(format, v...)
}

func Errorf(format string, v ...interface{}) {
	ZapLogger.Sugar().Errorf(format, v...)
}

func Debug(args ...interface{}) {
	ZapLogger.Sugar().Debug(args...)
}

func Info(args ...interface{}) {
	ZapLogger.Sugar().Info(args...)
}

func Error(args ...interface{}) {
	ZapLogger.Sugar().Error(args...)
}
