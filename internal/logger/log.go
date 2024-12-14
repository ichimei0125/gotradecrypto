package logger

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers map[string]*log.Logger
	onces   map[string]*sync.Once
	mutex   sync.Mutex
)

func getFileName(e exchange.Exchange, s exchange.Symbol) string {
	filename := "app.log"
	if e != nil && s != "" {
		filename = common.GetUniqueName(e.Name(), string(s)) + ".log"
	}
	return filename
}

// InitLogger
func InitLogger(e exchange.Exchange, symbol exchange.Symbol, maxSize, maxBackups, maxAge int, compress bool) {
	filename := getFileName(e, symbol)
	log_folder := config.GetEnvVar(common.ENV_LOG_PATH[0], common.ENV_LOG_PATH[1])
	path := path.Join(log_folder, filename)
	mutex.Lock()
	if onces == nil {
		onces = make(map[string]*sync.Once)
	}
	if loggers == nil {
		loggers = make(map[string]*log.Logger)
	}
	if _, exists := onces[filename]; !exists {
		onces[filename] = &sync.Once{}
	}
	mutex.Unlock()

	once := onces[filename]
	once.Do(func() {

		dir := getDir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Error on create dir: %v", err)
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

		loggers[filename] = log.New(logWriter, "", log.Ldate|log.Ltime)
	})
}

// getDir 获取日志文件的目录
func getDir(filePath string) string {
	return filepath.Dir(filePath)
}

func getLogger(e exchange.Exchange, s exchange.Symbol) *log.Logger {
	filename := getFileName(e, s)
	mutex.Lock()
	defer mutex.Unlock()
	logger, exists := loggers[filename]
	if !exists {
		log.Printf("Logger not initialized for Exchange: %s, Symbol: %s. Please call InitLogger first.", e.Name(), string(s))
	}
	return logger
}

// Info 记录信息级别日志
func Info(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := getLogger(e, s)
	if logger == nil {
		return
	}
	logger.SetPrefix("INFO: ")
	logger.Println(v...)
}

// Error 记录错误级别日志
func Error(v ...interface{}) {
	var e exchange.Exchange = nil
	var s exchange.Symbol = ""

	if len(v) >= 3 {
		e, s = v[0].(exchange.Exchange), v[1].(exchange.Symbol)
		v = v[2:]
	}

	logger := getLogger(e, s)
	if logger == nil {
		return
	}
	logger.SetPrefix("ERROR: ")
	logger.Println(v...)
}

// Debug 记录调试级别日志
func Debug(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := getLogger(e, s)
	if logger == nil {
		return
	}
	logger.SetPrefix("DEBUG: ")
	logger.Println(v...)
}

func Print(e exchange.Exchange, s exchange.Symbol, v ...interface{}) {
	logger := getLogger(e, s)
	if logger == nil {
		return
	}
	logger.SetPrefix("")
	logger.Println(v...)
}
