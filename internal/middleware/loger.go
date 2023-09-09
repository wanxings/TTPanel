package middleware

import (
	"TTPanel/internal/global"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
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

func RouterLog() gin.HandlerFunc {
	// 创建日志对象
	log := logrus.New()

	// 解析日志级别
	logLevel, err := logrus.ParseLevel(global.Config.Logger.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	// 设置日志级别
	log.SetLevel(logLevel)

	// 创建 Hook，将日志存储到文件中
	hook := NewFileHook(global.Config.Logger.RootPath+"/routers", "2006-01-02-15", 2)

	// 将 Hook 添加到日志对象的 Hooks 列表中
	log.AddHook(hook)

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 返回 gin 中间件函数，用于将日志写入文件
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		var url string
		var bodyParams any
		//是否记录参数
		if global.Config.Logger.RouterLogParams {
			// URL
			url = c.Request.URL.String()
			// 获取请求的Body参数
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			//bodyParams = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			var isJson map[string]any
			if err := json.Unmarshal(bodyBytes, &isJson); err != nil {
				bodyParams = "Body is not json"
			} else {
				bodyParams = isJson
			}
		}
		// 处理请求
		c.Next()

		//是否记录路由日志
		if !global.Config.Logger.RouterLog {
			return
		}

		// 获取请求的Header参数
		//headerParams := c.Request.Header
		//fmt.Println("Header params:", headerParams)

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		//reqMethod := c.Request.Method

		// 请求路由
		//reqUri := c.Request.RequestURI

		// 状态码
		//statusCode := c.Writer.Status()

		// 请求IP
		//clientIP := c.ClientIP()
		headerFields := make(map[string]interface{})
		for k, v := range c.Request.Header {
			headerFields[k] = v
		}
		uid, _ := c.Get("UID")
		uname, _ := c.Get("USERNAME")
		// 日志格式
		log.WithFields(logrus.Fields{
			"uid":        uid,
			"uname":      uname,
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latencyTime,
			"user-agent": c.Request.UserAgent(),
			"header":     headerFields,
			"body":       bodyParams,
			"url":        url,
		}).Info("Request completed")
		return
	}
}
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
		filename := hook.dir + "/" + time.Now().Format(hook.nameFormat) + "-" + strconv.Itoa(hook.fileNum) + ".log"
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
