package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/maxminddb-golang"
	"io"
	"math"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type LogAuditService struct {
}

var successPattern = regexp.MustCompile(`sshd.*Accepted`)
var failurePattern = regexp.MustCompile(`sshd.*Failed`)
var ipPattern = regexp.MustCompile(`from ([^ ]+) port (\d+)`)
var userPattern = regexp.MustCompile(`password for ([^ ]+)`)
var invalidUserPattern = regexp.MustCompile(`Failed (password|none) for (invalid user )?([^ ]+)`)

// PanelOperationLogList 面板操作日志列表
func (s *LogAuditService) PanelOperationLogList(req *request.OperationLogListR, offset, limit int) ([]*model.OperationLog, int, error) {
	conD := model.ConditionsT{"ORDER": "create_time DESC"}
	if !util.StrIsEmpty(req.Query) && !util.StrIsEmpty(req.QueryField) {
		if util.IsGeneral(req.QueryField) {
			conD[req.QueryField+" LIKE ?"] = "%" + req.Query + "%"
		}
	}
	if req.Type != 0 {
		conD["type = ?"] = req.Type
	}
	logList, total, err := (&model.OperationLog{}).List(global.PanelDB, &conD, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return logList, int(total), nil
}

// ClearPanelOperationLog 清空面板操作日志
func (s *LogAuditService) ClearPanelOperationLog() error {
	return (&model.OperationLog{}).DeleteAll(global.PanelDB)
}

// LogFileOccupancy 日志占用
func (s *LogAuditService) LogFileOccupancy() []response.LogFileOccupancy {
	logPathList := map[string]string{
		fmt.Sprintf("%s/data/logs/compress", global.Config.System.PanelPath):   helper.Message("log_audit.dockerComposeLog"),
		fmt.Sprintf("%s/data/logs/panel", global.Config.System.PanelPath):      helper.Message("log_audit.PanelUpdateDebugLog"),
		fmt.Sprintf("%s/data/logs/project", global.Config.System.PanelPath):    helper.Message("log_audit.ProjectRunLog"),
		fmt.Sprintf("%s/data/logs/queue_task", global.Config.System.PanelPath): helper.Message("log_audit.QueueTaskLog"),
		fmt.Sprintf("%s/data/logs/routers", global.Config.System.PanelPath):    helper.Message("log_audit.PanelRouterLog"),
		fmt.Sprintf("%s/data/logs/ssl", global.Config.System.PanelPath):        helper.Message("log_audit.SSLCertificateApplicationLog"),
		fmt.Sprintf("%s/data/logs/mysql", global.Config.System.PanelPath):      helper.Message("log_audit.MysqlOperationLog"),
	}
	var logFileOccupancyList []response.LogFileOccupancy
	for dirPath, description := range logPathList {
		count, size := util.GetFolderStats(dirPath)
		logFileOccupancyList = append(logFileOccupancyList, response.LogFileOccupancy{
			LogoPath:    dirPath,
			Description: description,
			Count:       count,
			Size:        size,
		})
	}
	return logFileOccupancyList
}

func (s *LogAuditService) SSHLoginLogList(query string, status, offset, limit int) ([]*response.SSHLoginLog, int, error) {
	sshLoginLogList := make([]*response.SSHLoginLog, 0)
	//判断日志文件位置
	logPathList := []string{"/var/log/auth.log", "/var/log/secure", "/var/log/messages"}
	for _, logPath := range logPathList {
		if util.IsFile(logPath) {
			//读取日志文件
			processedLogs, err := ParseSSHLogs(logPath, query, status)
			if err != nil {
				return nil, 0, err
			}
			// 计算起始索引和结束索引
			total := len(processedLogs)
			start := offset
			end := offset + limit
			if start > total {
				start = total
			}
			if end > total {
				end = total
			}
			for _, processedLog := range processedLogs[start:end] {
				continent, country := GetIPAttributionByGeoCountry(processedLog.IP)
				processedLog.IPAttribution = fmt.Sprintf("%s-%s", continent, country)
				sshLoginLogList = append(sshLoginLogList, processedLog)
			}
			return sshLoginLogList, total, nil
		}
	}
	return nil, 0, errors.New("not found log")
}
func ParseSSHLogs(logPath, query string, status int) ([]*response.SSHLoginLog, error) {
	processedLogs := make([]*response.SSHLoginLog, 0) // 用于存储处理后的SSH登录日志的切片
	file, err := os.Open(logPath)                     // 打开日志文件
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStat, err := file.Stat() // 获取文件信息
	if err != nil {
		return nil, err
	}
	fileSize := fileStat.Size() // 获取文件大小
	offset := fileSize - 1      // 设置偏移量为文件末尾
	lastLineSize := 0           // 用于记录最后一行的大小

	for { // 从文件末尾开始读取，直到找到最后一行的起始位置
		b := make([]byte, 1)             // 创建一个字节切片用于读取文件内容
		n, err := file.ReadAt(b, offset) // 从指定位置读取一个字节的内容
		if err != nil {
			break // 如果读取出现错误，则退出循环
		}
		char := string(b[0])
		if char == "\n" {
			break // 如果读取到换行符，则表示找到了最后一行的起始位置，退出循环
		}
		offset--          // 递减偏移量
		lastLineSize += n // 记录最后一行的大小
	}

	lastLine := make([]byte, lastLineSize)   // 创建一个字节切片，用于存储最后一行的内容
	_, err = file.ReadAt(lastLine, offset+1) // 从最后一行的起始位置读取内容
	if err != nil {
		return nil, err // 如果读取出现错误，则返回错误信息
	}
	linesPool := sync.Pool{New: func() interface{} { // 创建一个对象池，用于复用字节切片
		lines := make([]byte, 250*1024) // 每个对象的大小为250KB
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} { // 创建一个对象池，用于复用字符串
		lines := ""
		return lines
	}}

	r := bufio.NewReader(file) // 创建一个读取器
	var wg sync.WaitGroup      // 创建一个等待组，用于等待所有goroutine执行完毕
	for {
		buf := linesPool.Get().([]byte) // 从对象池中获取一个字节切片

		n, err := r.Read(buf) // 从文件中读取内容到字节切片中
		buf = buf[:n]         // 调整字节切片的长度为实际读取的字节数

		if n == 0 { // 如果读取的字节数为0
			if err != nil {
				break // 如果读取出现错误，则退出循环
			}
			if err == io.EOF {
				break // 如果已经读取到文件末尾，则退出循环
			}
			return nil, err // 如果既不是文件末尾也没有出现错误，则返回错误信息
		}
		nextUntilNewline, err := r.ReadBytes('\n') // 读取直到下一个换行符的内容
		if err != io.EOF {
			buf = append(buf, nextUntilNewline...) // 将下一个换行符之前的内容追加到字节切片中
		}

		wg.Add(1) // 增加等待组的计数
		go func() {
			parseSSHLogsChunk(buf, &processedLogs, &linesPool, &stringPool, query, status) // 解析SSH登录日志的片段
			wg.Done()                                                                      // 减少等待组的计数
		}()

	}
	wg.Wait() // 等待所有goroutine执行完毕

	if err != nil {
		return nil, err // 如果解析过程中出现错误，则返回错误信息
	}
	sort.Slice(processedLogs, func(i, j int) bool {
		return processedLogs[i].LoginTime > processedLogs[j].LoginTime // 根据登录时间对日志进行排序
	})
	return processedLogs, nil // 返回处理后的SSH登录日志和nil错误信息
}

// parseSSHLogsChunk 函数用于解析 SSH 日志的分块数据
// chunk 是待处理的日志数据块
// processedLogs 是已处理的 SSH 登录日志的切片指针
// linesPool 是用于重用字符串切片的同步池
// stringPool 是用于重用字符串的同步池
// query 是用于过滤日志的查询字符串
// status 是用于过滤登录状态的标识
func parseSSHLogsChunk(chunk []byte, processedLogs *[]*response.SSHLoginLog, linesPool *sync.Pool, stringPool *sync.Pool, query string, status int) {
	// 创建一个等待组
	var wg2 sync.WaitGroup
	// 将数据块转换为字符串
	logs := stringPool.Get().(string)
	logs = string(chunk)
	linesPool.Put(chunk)
	// 将字符串按行分割为切片
	logsSlice := strings.Split(logs, "\n")
	stringPool.Put(logs)
	// 定义每个线程处理的日志行数
	chunkSize := 300
	// 获取日志行数
	n := len(logsSlice)
	// 计算线程数量
	noOfThread := n / chunkSize
	if n%chunkSize != 0 {
		noOfThread++
	}
	// 遍历每个线程
	for i := 0; i < (noOfThread); i++ {
		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done()
			// 定义正则表达式模式

			// 遍历指定范围内的日志行
			for i := s; i < e; i++ {
				line := logsSlice[i]
				if len(line) == 0 {
					continue
				}
				// 如果查询字符串不匹配直接跳过,兼容query为空的情况
				if !strings.Contains(line, query) {
					continue
				}
				var insertLog *response.SSHLoginLog
				if successPattern.MatchString(line) && status != constant.SSHLoginStatusByFailed { // 匹配成功登录日志
					ipMatch := ipPattern.FindStringSubmatch(line)
					userMatch := userPattern.FindStringSubmatch(line)
					if len(ipMatch) != 3 || len(userMatch) != 2 {
						continue
					}
					insertLog = new(response.SSHLoginLog)
					insertLog.IP = ipMatch[1]
					insertLog.Port = ipMatch[2]
					insertLog.User = userMatch[1]
					insertLog.Status = 1
				} else if failurePattern.MatchString(line) && status != constant.SSHLoginStatusBySuccess { // 匹配失败登录日志
					ipMatch := ipPattern.FindStringSubmatch(line)
					userMatch := invalidUserPattern.FindStringSubmatch(line)
					if len(ipMatch) != 3 || len(userMatch) != 4 {
						continue
					}
					insertLog = new(response.SSHLoginLog)
					insertLog.IP = ipMatch[1]
					insertLog.Port = ipMatch[2]
					insertLog.User = userMatch[3]
					insertLog.Status = 2
				}
				if insertLog != nil { // 将处理完的日志对象添加到已处理日志切片中
					insertLog.LoginTime = util.GetLogTimestampZ(line)
					*processedLogs = append(*processedLogs, insertLog)
				}
			}
		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}
	// 等待所有线程完成
	wg2.Wait()
	// 清空日志切片
	logsSlice = nil
}

func (s *LogAuditService) WriteOperationLog(c *gin.Context, logType int, msg string) {
	data := &model.OperationLog{
		Type: logType,
		Log:  msg,
	}
	if c != nil {
		if u, exist := c.Get("UID"); exist {
			data.Uid = u.(int64)
		}
		if u, exist := c.Get("USERNAME"); exist {
			data.Username = u.(string)
		}
		data.IP = c.ClientIP()
		continent, country := GetIPAttributionByGeoCountry(c.ClientIP())
		data.IPAttribution = fmt.Sprintf("%s-%s", continent, country)
	}

	_, err := data.Create(global.PanelDB)
	if err != nil {
		global.Log.Errorf("写入操作日志失败: %v", err)
		return
	}
	return
}

// GetIPAttributionByGeoCountry 通过GeoIP2-Country获取IP归属地
func GetIPAttributionByGeoCountry(ip string) (continent string, country string) {
	db, err := maxminddb.Open(global.Config.System.PanelPath + "/data/geo/GeoLite2-Country.mmdb")
	if err != nil {
		global.Log.Errorf("open geoip2 database failed, Error:%s", err)
		continent = "Error"
		return
	}
	defer func(db *maxminddb.Reader) {
		_ = db.Close()
	}(db)
	var record struct {
		Continent struct {
			Names struct {
				ZhCn string `maxminddb:"zh-CN"`
			} `maxminddb:"names"`
		} `maxminddb:"continent"`
		Country struct {
			Names struct {
				ZhCn string `maxminddb:"zh-CN"`
			} `maxminddb:"names"`
		} `maxminddb:"country"`
	}
	err = db.Lookup(net.ParseIP(ip), &record)
	if err != nil {
		global.Log.Errorf("获取IP归属地失败, Error:%s", err)
		continent = "Error"
		return
	}
	if util.StrIsEmpty(record.Continent.Names.ZhCn) {
		record.Continent.Names.ZhCn = "未知"
	}
	if util.StrIsEmpty(record.Country.Names.ZhCn) {
		record.Country.Names.ZhCn = "未知"
	}
	return record.Continent.Names.ZhCn, record.Country.Names.ZhCn
}

// GetIPAttributionByCZ 通过纯真IP库获取IP归属地
//func GetIPAttributionByCZ(ip string) (country string, city string) {
//	qqwryPath := global.Config.System.PanelPath + "/data/qqwry/qqwry.dat"
//	q := qqwry.NewQQwry(qqwryPath)
//	q.Find(ip)
//	return q.Country, q.City
//}
