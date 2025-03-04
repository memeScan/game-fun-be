package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	// LevelError 错误
	LevelError = iota
	// LevelWarning 警告
	LevelWarning
	// LevelInformational 提示
	LevelInformational
	// LevelDebug 除错
	LevelDebug
)

var (
	loggerOnce sync.Once
	logger     *Logger
)

// 缓冲区对象池
var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// Logger 日志结构体
type Logger struct {
	level     int
	logFile   *os.File
	writer    *bufio.Writer
	mu        sync.Mutex
	logDir    string
	logChan   chan string
	done      chan struct{}
	closeOnce sync.Once
}

// 初始化日志目录
func initLogDir() string {
	// 优先使用环境变量
	if dir := os.Getenv("LOG_DIR"); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("无法创建日志目录: %v", err))
		}
		return dir
	}

	// 使用当前项目下的 logs 文件夹
	logDir := filepath.Join(".", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("无法创建日志目录: %v", err))
	}
	return logDir
}

// 获取当前日志文件路径
func getLogFilePath(logDir string) string {
	return filepath.Join(logDir, fmt.Sprintf("app.%s.log", time.Now().Format("20060102")))
}

// 清理旧日志
func cleanOldLogs(logDir string) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			threshold := now.AddDate(0, 0, -30)

			files, err := os.ReadDir(logDir)
			if err != nil {
				fmt.Printf("Error reading log directory: %v\n", err)
				continue
			}

			for _, file := range files {
				if info, err := file.Info(); err == nil {
					if !info.IsDir() && info.ModTime().Before(threshold) {
						filePath := filepath.Join(logDir, info.Name())
						if err := os.Remove(filePath); err != nil {
							fmt.Printf("Error removing old log file %s: %v\n", filePath, err)
						}
					}
				}
			}
		}
	}()
}

func (ll *Logger) openLogFile() error {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if ll.writer != nil {
		ll.writer.Flush()
	}
	if ll.logFile != nil {
		ll.logFile.Close()
	}

	f, err := os.OpenFile(getLogFilePath(ll.logDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	ll.logFile = f
	ll.writer = bufio.NewWriterSize(f, 256*1024)
	return nil
}

func (ll *Logger) formatLog(level, format string, v ...interface{}) string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// 时间戳
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05.000 +0800"))
	buf.WriteByte(' ')
	buf.WriteString(level)
	buf.WriteByte(' ')

	if len(v) > 0 {
		fmt.Fprintf(buf, format, v...)
	} else {
		buf.WriteString(format)
	}
	buf.WriteByte('\n')

	return buf.String()
}

func (ll *Logger) writeLog(msg string) {
	select {
	case ll.logChan <- msg:
		// 成功发送到通道
	default:
		// 通道已满，直接写入文件
		ll.mu.Lock()
		ll.writer.WriteString(msg)
		ll.writer.Flush()
		ll.mu.Unlock()
	}
}

func (ll *Logger) Panic(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := ll.formatLog("[Panic]", format, v...)
	ll.writeLog(msg)
	panic(msg)
}

func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	ll.writeLog(ll.formatLog("[E]", format, v...))
}

func (ll *Logger) Warning(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	ll.writeLog(ll.formatLog("[W]", format, v...))
}

func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	ll.writeLog(ll.formatLog("[I]", format, v...))
}

func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	msg := ll.formatLog("[D]", format, v...)
	if ll.level == LevelDebug {
		fmt.Print(msg)
	}
	ll.writeLog(msg)
}

// BuildLogger 构建logger
func BuildLogger(level string) {
	intLevel := LevelError
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "info":
		intLevel = LevelInformational
	case "debug":
		intLevel = LevelDebug
	}

	logDir := initLogDir()
	l := &Logger{
		level:   intLevel,
		logDir:  logDir,
		logChan: make(chan string, 50000),
		done:    make(chan struct{}),
	}

	if err := l.openLogFile(); err != nil {
		panic(fmt.Sprintf("无法打开日志文件: %v", err))
	}

	go l.startAsyncWriter()
	cleanOldLogs(logDir)

	// 启动日志文件轮转
	go func() {
		for {
			select {
			case <-l.done:
				return
			default:
				now := time.Now()
				next := now.Add(24 * time.Hour)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				time.Sleep(next.Sub(now))
				l.openLogFile()
			}
		}
	}()

	logger = l
}

func (ll *Logger) startAsyncWriter() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var msgCount int

	for {
		select {
		case msg := <-ll.logChan:
			ll.mu.Lock()
			ll.writer.WriteString(msg)
			msgCount++
			if msgCount >= 1000 {
				ll.writer.Flush()
				msgCount = 0
			}
			ll.mu.Unlock()

		case <-ticker.C:
			if msgCount > 0 {
				ll.mu.Lock()
				ll.writer.Flush()
				ll.mu.Unlock()
				msgCount = 0
			}

		case <-ll.done:
			ll.mu.Lock()
			ll.writer.Flush()
			ll.mu.Unlock()
			return
		}
	}
}

// Log 返回日志对象
func Log() *Logger {
	loggerOnce.Do(func() {
		// 从环境变量获取日志级别，默认为 "debug"
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "debug"
		}
		BuildLogger(logLevel)
	})
	return logger
}

// GetGoroutineID returns the goroutine ID of the calling goroutine
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func (ll *Logger) Close() error {
	var err error
	ll.closeOnce.Do(func() {
		close(ll.done)
		// 等待所有消息写入
		for len(ll.logChan) > 0 {
			time.Sleep(time.Millisecond * 100)
		}

		ll.mu.Lock()
		defer ll.mu.Unlock()

		if ll.writer != nil {
			err = ll.writer.Flush()
		}
		if ll.logFile != nil {
			if closeErr := ll.logFile.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}
	})
	return err
}
