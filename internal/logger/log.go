package logger

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers map[string]*log.Logger
	once    sync.Once
)

func getFileName(e exchange.Exchange, s exchange.Symbol) string {
	return e.Name() + "_" + string(s) + ".log"
}

// InitLogger 初始化日志
func InitLogger(e exchange.Exchange, symbol exchange.Symbol, maxSize, maxBackups, maxAge int, compress bool) {
	once.Do(func() {
		filename := "app.log"
		if e != nil && symbol != "" {
			filename = getFileName(e, symbol)
		}
		path := path.Join("log", filename)

		// 创建日志目录（如果不存在）
		dir := getDir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("无法创建日志目录: %v", err)
			}
		}

		// 配置日志文件滚动
		logWriter := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    maxSize,    // 单个日志文件的最大大小（MB）
			MaxBackups: maxBackups, // 保留的旧日志文件数量
			MaxAge:     maxAge,     // 日志文件保留天数
			Compress:   compress,   // 启用压缩
		}

		// 初始化日志记录器
		loggers[filename] = log.New(logWriter, "", log.Ldate|log.Ltime)
	})
}

// getDir 获取日志文件的目录
func getDir(filePath string) string {
	return filepath.Dir(filePath)
}

// Info 记录信息级别日志
func Info(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := loggers[getFileName(e, s)]
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.Println(v...)
}

// Error 记录错误级别日志
func Error(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := loggers[getFileName(e, s)]
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.SetPrefix("ERROR: ")
	logger.Println(v...)
}

// Debug 记录调试级别日志
func Debug(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := loggers[getFileName(e, s)]
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.SetPrefix("DEBUG: ")
	logger.Println(v...)
}

func Print(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := loggers[getFileName(e, s)]
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.Println(v...)
}
