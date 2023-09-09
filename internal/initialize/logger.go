package initialize

import (
	"TTPanel/internal/global"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strconv"
	"time"
)

// FileHook 用于将日志存储到文件中
type FileHook struct {
	// 存储日志的目录
	dir string
	// 文件名格式
	nameFormat string
	// 每个日志文件最大大小，单位为 MB
	maxSize int
	// 当前日志文件的序号
	fileNum int
	// 当前日志文件的大小
	fileSize int
	// 当前日志文件的写入器
	writer *lumberjack.Logger
}

type CustomizeFormatter struct{}

func InitLogger() *logrus.Logger {
	// 创建日志对象
	log := logrus.New()

	// 解析日志级别
	logLevel, err := logrus.ParseLevel(global.Config.Logger.LogLevel)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	// 设置日志级别
	log.SetLevel(logLevel)

	// 创建 Hook，将日志存储到文件中
	hook := NewFileHook(global.Config.Logger.RootPath+"/panel", "2006-01-02-15", 2)

	// 将 Hook 添加到日志对象的 Hooks 列表中
	log.AddHook(hook)

	//log.SetFormatter(&logrus.TextFormatter{
	//	DisableColors:   true,
	//	ForceQuote:      true,
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})
	log.SetReportCaller(true)
	log.SetFormatter(&CustomizeFormatter{})
	return log
}

//func newFileLogger() io.Writer {
//	return &lumberjack.Logger{
//		Filename:  global.Config.Logger.SavePath + "/" + global.Config.Logger.FileName + global.Config.Logger.FileExt,
//		MaxSize:   2,
//		MaxAge:    10,
//		LocalTime: true,
//	}
//}

func NewFileHook(dir, nameFormat string, maxSize int) *FileHook {
	return &FileHook{
		dir:        dir,
		nameFormat: nameFormat,
		maxSize:    maxSize,
		fileNum:    0,
		fileSize:   0,
		writer:     nil,
	}
}

// Fire 实现 logrus.Hook 接口的方法，将日志写入文件
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	// 获取日志内容
	msg, err := entry.String()
	if err != nil {
		return err
	}

	// 检查是否需要切换到下一个日志文件
	if hook.writer == nil || hook.fileSize+len(msg) > hook.maxSize*1024*1024 {
		// 关闭当前日志文件的写入器
		if hook.writer != nil {
			_ = hook.writer.Close()
		}

		// 计算下一个日志文件的序号
		hook.fileNum++

		// 创建新的日志文件的写入器
		filename := hook.dir + "/" + time.Now().Format(hook.nameFormat) + "-" + strconv.Itoa(hook.fileNum) + ".json"
		hook.writer = &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    hook.maxSize,
			MaxBackups: 0,
			MaxAge:     0,
			LocalTime:  true,
			Compress:   true,
		}

		// 重置当前日志文件的大小
		hook.fileSize = 0
	}

	// 写入日志内容
	n, err := hook.writer.Write([]byte(msg))
	if err != nil {
		return err
	}

	// 更新当前日志文件的大小
	hook.fileSize += n

	return nil
}

// Levels 实现 logrus.Hook 接口的方法，指定日志级别
func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (s *CustomizeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
