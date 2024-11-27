package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *log.Logger
	once   sync.Once
)

// InitLogger 初始化日志
func InitLogger(logFilePath string, maxSize, maxBackups, maxAge int, compress bool) {
	once.Do(func() {
		// 创建日志目录（如果不存在）
		dir := getDir(logFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("无法创建日志目录: %v", err)
			}
		}

		// 配置日志文件滚动
		logWriter := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    maxSize,    // 单个日志文件的最大大小（MB）
			MaxBackups: maxBackups, // 保留的旧日志文件数量
			MaxAge:     maxAge,     // 日志文件保留天数
			Compress:   compress,   // 启用压缩
		}

		// 初始化日志记录器
		logger = log.New(logWriter, "", log.Ldate|log.Ltime)
	})
}

// getDir 获取日志文件的目录
func getDir(filePath string) string {
	return filepath.Dir(filePath)
}

// Info 记录信息级别日志
func Info(v ...interface{}) {
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.Println(v...)
}

// Error 记录错误级别日志
func Error(v ...interface{}) {
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.SetPrefix("ERROR: ")
	logger.Println(v...)
}

// Debug 记录调试级别日志
func Debug(v ...interface{}) {
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.SetPrefix("DEBUG: ")
	logger.Println(v...)
}

func Print(v ...interface{}) {
	if logger == nil {
		log.Println("Use InitLogger")
		return
	}
	logger.Println(v...)
}
