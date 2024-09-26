package common

import (
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

// 定义配置类信息，注意里面变量名称首字母需大写，才可获取配置信息
type Config struct {

	// Config 配置文件结构体
	LogFile struct {
		LogDir     string `default:"logs"`
		FileName   string `default:"server.log"`
		LogLevel   string `default:"info"`
		MaxSize    int    `default:128`
		MaxBackups int    `default:180`
		MaxAge     int    `default:1`
		Compress   bool   `default:false`
		PrintTag   bool   `default:false`
	}

	//请求信息
	ReqInfo struct {
		Url string
	}
}

var Conf = Config{}

// 读取配置参数
func InitConfig() error {
	//打包使用配置文件路径
	//filePath := "./config.yml"
	//本地测试配置文件路径
	filePath := "src/config.yml"
	err := configor.Load(&Conf, filePath)
	if err != nil {
		log.Fatal("加载配置文件失败.......，错误信息：", err)
		return err
	}
	return nil
}

// String returns a lower-case ASCII representation of the log level.
func Level2String(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
