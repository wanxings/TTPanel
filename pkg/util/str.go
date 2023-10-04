package util

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/idna"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type StrType int

const (
	NUM   StrType = iota // 数字
	LOWER                // 小写字母
	UPPER                // 大写字母
	ALL                  // 全部
	CLEAR                // 去除部分易混淆的字符
)

var fontKinds = [][]int{{10, 48}, {26, 97}, {26, 65}}
var letters = []byte("123457890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr 生成随机字符串
// size 个数 kind 模式
func RandStr(size int, kind StrType) []byte {
	iKind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			iKind = StrType(rand.Intn(int(ALL)))
		}
		scope, base := fontKinds[iKind][0], fontKinds[iKind][1]
		result[i] = uint8(base + rand.Intn(scope))
		// 不易混淆字符模式：重新生成字符
		if kind == 4 {
			result[i] = letters[rand.Intn(len(letters))]
		}
	}
	return result
}

// UniqueStrSlice 字符串切片去重
func UniqueStrSlice(slice []string) []string {
	set := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if _, ok := set[v]; !ok {
			set[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// StrIsEmpty 判断字符串是否为空,为空返回true
func StrIsEmpty(str string) bool {
	return len(str) == 0
}

// ClearStr 清理字符串中的空格、换行符、制表符
func ClearStr(str string) string {
	return string(ClearBytes([]byte(str)))
}

func ClearBytes(str []byte) []byte {
	return bytes.ReplaceAll(bytes.ReplaceAll(bytes.ReplaceAll(bytes.ReplaceAll(str, []byte(" "), []byte("")), []byte("\r"), []byte("")), []byte("\t"), []byte("")), []byte("\n"), []byte(""))
}

// IsGeneral 正则校验字符串是否由字母、数字、_、-、小数点组成
func IsGeneral(str string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	return reg.MatchString(str)
}

// IsPHPVersion 校验PHP版本号格式是否正确
func IsPHPVersion(str string) bool {
	reg := regexp.MustCompile(`^[0-9]+$`)
	return reg.MatchString(str)
}

// IsVersion 校验Nginx版本号格式是否正确
func IsVersion(str string) bool {
	reg := regexp.MustCompile(`^[0-9.]+$`)
	return reg.MatchString(str)
}

// IsMysqlVersion 校验Mysql版本号格式是否正确
func IsMysqlVersion(str string) bool {
	reg := regexp.MustCompile(`^[0-9.]+$`)
	return reg.MatchString(str)
}

// ToPunycode 将域名转换为punycode编码
func ToPunycode(domain string) string {
	tmp := strings.Split(domain, ".")
	var newDomain string
	for _, dKey := range tmp {
		if dKey == "*" {
			continue
		}
		match, _ := regexp.MatchString(`[\x80-\xff]+`, dKey)
		if !match {
			match, _ = regexp.MatchString(`[\u4e00-\u9fa5]+`, dKey)
		}
		if !match {
			newDomain += dKey + "."
		} else {
			punycode, _ := idna.ToASCII(dKey)
			if len(punycode) > 0 {
				newDomain += "xn--" + punycode + "."
			} else {
				newDomain += dKey + "."
			}
		}
	}
	if tmp[0] == "*" {
		newDomain = "*." + newDomain
	}
	return newDomain[:len(newDomain)-1]
}

// TrimPath 去除路径的空格和结尾的斜杠
func TrimPath(path string) string {
	path = strings.ReplaceAll(path, " ", "")
	if strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}
	return path
}

// TrimStr 去除字符串中的空格、换行、制表符和特殊符号
func TrimStr(str string) string {
	reg, _ := regexp.Compile(`[\s\\\/:\*\?"<>\|]+`) // 匹配一个或多个空白字符和特殊符号
	str = reg.ReplaceAllString(str, "")             // 去除空格、换行、制表符和特殊符号
	return str
}

// IsValidProjectName 校验项目名称 去除字符串中的空格和换行符, 不能含有/!@#$%^&*()+|-'"特殊字符,不能以小数点开头或者结尾,文件名长度不能超过200个字符
func IsValidProjectName(name string) (string, error) {
	// 去除字符串中的空格和换行符
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "", errors.New("project name cannot be empty")
	}
	// 文件名长度不能超过255个字符
	if len(name) > 200 {
		return "", errors.New("project name length cannot exceed 200 characters")
	}
	// 不能以小数点开头或者结尾
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return "", errors.New("project name cannot start or end with '.'")
	}
	//不能含有/!@#$%^&*()+|-'"特殊字符
	if strings.ContainsAny(name, "/!@#$%^&*()+|-'\"") {
		return "", errors.New("project name cannot contain special characters")
	}
	return name, nil
}

func GetCmdDelimiter() string {
	return "=========================================================================================="
}
func GetCmdDelimiter2() string {
	return "----------------------------------------\n"
}

func ConvertSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	index := 0
	fSize := float64(size)
	for fSize >= 1024 && index < len(units)-1 {
		fSize /= 1024
		index++
	}
	return fmt.Sprintf("%.2f %s", fSize, units[index])
}

func FormatDuration(duration int64) string {
	seconds := duration % 60
	minutes := (duration / 60) % 60
	hours := duration / 3600

	if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	} else {
		return fmt.Sprintf("%d秒", seconds)
	}
}

// StrInArray 判断一个字符串是否在一个字符串数组中
func StrInArray(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
