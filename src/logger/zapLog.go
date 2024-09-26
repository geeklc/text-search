package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"text-search/src/common"
	"time"
)

// 创建logger
func CreateLogger() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

// Logger 全局变量，导出给调用者使用
var Logger *zap.Logger

// Init 初始化日志相关目录
func Init() error {
	conf := common.Conf
	config := conf.LogFile
	// 创建日志根目录
	if _, err := os.Stat(config.LogDir); os.IsNotExist(err) {
		err := os.MkdirAll(config.LogDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory, err: %v", err)
		}
	}

	// 自定义初始化zap库
	initLogger()

	return nil
}

// getWriter 获取wirter文件写入
func getWriter(logBasePath, LogFileName string, config common.Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logBasePath, LogFileName),
		MaxBackups: config.LogFile.MaxBackups,
		MaxSize:    config.LogFile.MaxSize,
		MaxAge:     config.LogFile.MaxAge,
		Compress:   config.LogFile.Compress,
	}
}

// initLog 初始化日志
func initLogger() {
	conf := common.Conf
	c := conf.LogFile

	// 获取io.Writer实现
	logWriter := getWriter(c.LogDir, c.FileName, conf)

	// 获取日志默认配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// 自定义日志输出格式
	// 修改TimeKey
	encoderConfig.TimeKey = "time"
	// 修改MessageKey
	encoderConfig.MessageKey = "message"
	// 时间格式符合人类观看
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 日志等级大写INFO
	encoderConfig.EncodeLevel = func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(lvl.CapitalString())
	}
	// 日志打印时所处代码位置
	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}
	// 加载自定义配置为json格式
	//encoder := zapcore.NewJSONEncoder(encoderConfig)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 自定义日志级别 info
	logLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zap.FatalLevel && lvl >= common.Level2String(c.LogLevel)
	})

	// 日志文件输出位置
	var core zapcore.Core
	if c.PrintTag {
		//同时在文件和终端输出日志
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel), // info级别日志
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
		)
	} else {
		//只在文件输出日志
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel), // info级别日志
		)
	}

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
