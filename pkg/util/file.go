package util

import (
	"TTPanel/internal/global"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var SensitivePath = map[string]bool{
	"":                 true,
	"/":                true,
	"/*":               true,
	"/www":             true,
	"/root":            true,
	"/boot":            true,
	"/bin":             true,
	"/etc":             true,
	"/home":            true,
	"/dev":             true,
	"/sbin":            true,
	"/var":             true,
	"/usr":             true,
	"/tmp":             true,
	"/sys":             true,
	"/proc":            true,
	"/media":           true,
	"/mnt":             true,
	"/opt":             true,
	"/lib":             true,
	"/srv":             true,
	"/selinux":         true,
	"/www/panel":       true,
	"/www/server":      true,
	"/www/server/data": true,
	"/www/recycle_bin": true,
}

type DirInfo struct {
	DirName string `json:"dir_name"` //文件夹名称
	//DirPath    string `json:"dir_path"`    //文件夹路径
	IsDir      bool   `json:"is_dir"`      //是否是文件夹
	Perm       string `json:"perm"`        //权限
	PermString string `json:"perm_string"` //权限字符串
	Size       int64  `json:"size"`        //文件大小
	ModTime    int64  `json:"mod_time"`    //修改时间
	Link       string `json:"link"`        //软连接目标路径
	Remark     string `json:"remark"`      //备注
	Owner      string `json:"owner"`       //拥有者
	OwnerId    string `json:"owner_id"`    //拥有者id
}

type FileHistoryInfo struct {
	Path    string `json:"path"`     //所处路径
	MD5     string `json:"md5"`      //文件md5
	Size    int64  `json:"size"`     //文件大小
	ModTime int64  `json:"mod_time"` //修改时间
}

// CheckSensitivePath 检查是否是敏感路径
func CheckSensitivePath(path string) bool {
	path = strings.Replace(path, "//", "/", -1)
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	if _, ok := SensitivePath[path]; ok {
		return true
	}
	return false
}

// PathExists 判断所给路径文件/文件夹是否存在,true 存在，false 不存在
func PathExists(path string) bool {
	_, err := os.Lstat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// FileSize 获取文件大小
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// IsTextFile 是否是文本文件
func IsTextFile(path string) bool {
	fileTypeStr, err := ExecShell(fmt.Sprintf("file -bk %s", path))
	if err != nil {
		return false
	}
	return strings.Contains(fileTypeStr, "text") || strings.Contains(fileTypeStr, "empty") || strings.Contains(fileTypeStr, "data")
}

// RenameFileIfExists 如果文件路径存在，返回一个不存在的带序列的新文件路径
func RenameFileIfExists(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return path, nil
	}

	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]
	for i := 1; ; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
		_, err := os.Stat(newPath)
		if os.IsNotExist(err) {
			err = os.Rename(path, newPath)
			if err != nil {
				return "", err
			}
			return newPath, nil
		}
	}
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// GetFileSize 获取文件大小
func GetFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fi.Size()
}

// GetFilePerm 获取文件权限
func GetFilePerm(path string) os.FileMode {
	file, err := os.Open(path)
	if err != nil {
		return 0755
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	fi, err := file.Stat()
	if err != nil {
		return 0755
	}
	return fi.Mode().Perm()
}

// ChangePermission 修改文件（夹）权限和所有者
func ChangePermission(path string, userName string, permissions string) error {
	permI, _ := strconv.ParseUint(permissions, 8, 10)
	userInfo, _ := user.Lookup(userName) //获取用户信息
	err := os.Chmod(path, os.FileMode(permI))
	if err != nil {
		return err
	}
	uid, _ := strconv.Atoi(userInfo.Uid)
	gid, _ := strconv.Atoi(userInfo.Gid)
	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}
	return nil
}

// FormatFileContent 格式化内容
func FormatFileContent(content string) []string {
	//根据换行符分割
	contentList := strings.Split(content, "\n")
	for i, v := range contentList {
		//如果开头是#，则跳过
		if strings.HasPrefix(v, "#") {
			continue
		}
		//去除空格
		contentList[i] = strings.TrimSpace(v)
	}
	return contentList
}

// FileModePermToString 将文件权限转换为字符串格式“755”
func FileModePermToString(mode os.FileMode) string {
	return strconv.FormatInt(int64(mode), 8)
}

// GetFileDetails 获取文件（夹）详细信息
func GetFileDetails(path string) *DirInfo {
	fileInfo, err := os.Stat(path)
	if err != nil {
		global.Log.Debugf("GetFileDetails->os.Stat-err: %s,path:%s", err.Error(), path)
		global.Log.Error(err)
		fileInfo, err = os.Lstat(path)
		if err != nil {
			global.Log.Debugf("GetFileDetails->os.Lstat-err: %s,path:%s", err.Error(), path)
			global.Log.Error(err)
			return &DirInfo{}
		}
	}
	userId := fmt.Sprint(fileInfo.Sys().(*syscall.Stat_t).Uid)
	owner := "root"
	userInfo, err := user.LookupId(userId)
	if err != nil {
		owner = err.Error()
	} else {
		owner = userInfo.Username
	}
	Ln, _ := os.Readlink(path) //尝试获取软链接
	isDir := fileInfo.IsDir()
	size := fileInfo.Size()
	if isDir {
		size = -1
	}
	return &DirInfo{
		DirName:    fileInfo.Name(),
		IsDir:      fileInfo.IsDir(),
		Perm:       FileModePermToString(fileInfo.Mode().Perm()),
		PermString: fileInfo.Mode().String(),
		Size:       size,
		ModTime:    fileInfo.ModTime().Unix(),
		Link:       Ln,
		Owner:      owner,
		OwnerId:    userId,
	}
}

// GetFileContentType 获取文件类型
func GetFileContentType(file *os.File) (string, error) {
	// 只读取前512字节，判断文件类型
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// GetGroupNameByGID 根据GID获取组名
func GetGroupNameByGID(gid int) string {
	group, _ := user.LookupGroupId(fmt.Sprint(gid))
	return group.Name
}

// IsUserExist 判断用户是否存在
func IsUserExist(username string) bool {
	_, err := user.Lookup(username)
	if err != nil {
		return false
	}
	return true
}

// DeleteUser 删除用户
func DeleteUser(username string) error {
	return exec.Command("userdel", username).Run()
}

// WriteFile 写入文件
func WriteFile(path string, content []byte, perm fs.FileMode) error {
	_ = os.MkdirAll(filepath.Dir(path), perm) //创建文件夹
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		fmt.Println("open file failed,err:", err)
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("close file failed,err:", err)
		}
	}(file)
	_, err = file.Write(content) //写入字节切片数据
	//_, err = file.WriteString(content) //直接写入字符串数据
	if err != nil {
		fmt.Println("WriteString failed,err:", err)
		return err
	}
	return nil
}

// ReadFileStringBody 读取文件
func ReadFileStringBody(path string) (string, error) {
	//检查path是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	//检查path是否是文件
	if fileInfo.IsDir() {
		return "", errors.New(path + " is not a file")
	}
	//判断文件类型是否是text
	if !IsTextFile(path) {
		return "", errors.New("file is not of type text")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("ReadFileStringBody->File reading error", err)
		return "", err
	}
	return string(data), nil
}

// CreateDir 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// CreateFile 创建文件
func CreateFile(path string) error {
	_, err := os.Create(path)
	return err
}

// CheckDirName 检查文件夹名称是否合法
func CheckDirName(name string) error {
	var sensitiveNames = []string{`\\`, `&`, `*`, `|`, `;`, `"`, `'`, `<`, `>`}
	for _, p := range sensitiveNames {
		if strings.Contains(name, p) {
			return errors.New(fmt.Sprintf("Folder name cannot contain '%s'", p))
		}
	}
	//名称结尾是否是.
	if strings.HasSuffix(name, ".") {
		return errors.New("Folder ending cannot be '.' ")
	}
	return nil
}

// TailFile 取文件指定尾行数
func TailFile(path string, line int) (string, error) {
	body, err := ExecShell(fmt.Sprintf("tail -n %d %s", line, path))
	if err != nil {
		return "", err
	}
	return body, nil
}

func IsValidFileName(fileName string) bool {
	if len(fileName) > 255 {
		return false
	}
	// 定义文件名的正则表达式
	pattern := "^[^./][^/\\\\]*$"
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(fileName) {
		return false
	}
	reservedWords := []string{"/dev/null", "/usr/bin"}
	for _, word := range reservedWords {
		if fileName == word {
			return false
		}
	}
	return true
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir 复制文件夹
func CopyDir(src string, dst string) error {
	var err error

	// 获取源文件夹信息
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 如果目标文件夹不存在，则创建
	if !fi.IsDir() {
		return fmt.Errorf("src is not a directory")
	}

	_, err = os.Open(dst)
	if err != nil {
		_ = os.MkdirAll(dst, fi.Mode())
	}

	// 遍历源文件夹中的所有文件和目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcFile := filepath.Join(src, entry.Name())
		dstFile := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcFile, dstFile)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = CopyFile(srcFile, dstFile)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

// DownloadFile 下载文件保存到指定的路径（包含文件名）如果存在并且要求不覆盖文件则追加时间戳
func DownloadFile(filePath string, url string, overwrite bool) (string, error) {
	// 如果文件存在且不覆盖，则追加当前时间戳
	if !overwrite && PathExists(filePath) {
		ext := filepath.Ext(filePath)
		filename := strings.TrimSuffix(filePath, ext)
		now := time.Now().UnixNano()
		filePath = fmt.Sprintf("%s-%d%s", filename, now, ext)
	}

	//尝试创建文件夹
	fileDir := filepath.Dir(filePath)
	err := os.MkdirAll(fileDir, 0644)
	if err != nil {
		return "", err
	}

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	// 发起 GET 请求获取文件
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 将响应内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// TarGzFiles 将指定的文件或目录压缩到一个tar.gz文件

// MoveFiles 移动文件或文件夹
func MoveFiles(files []string, destDir string) error {
	for _, file := range files {
		// 获取文件信息
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		// 目标路径
		destPath := filepath.Join(destDir, fileInfo.Name())

		if fileInfo.IsDir() {
			// 移动文件夹
			err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// 目标路径
				dest := filepath.Join(destDir, path[len(file):])

				if info.IsDir() {
					// 创建目录
					err = os.MkdirAll(dest, info.Mode())
					if err != nil {
						return err
					}
				} else {
					// 移动文件
					err = os.Rename(path, dest)
					if err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				return err
			}

			// 删除源文件夹
			err = os.RemoveAll(file)
			if err != nil {
				return err
			}
		} else {
			// 移动文件
			err = os.Rename(file, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// IsLink 文件是否是软件链接
func IsLink(path string) bool {
	// 获取文件信息
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	// 判断文件类型
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}

// GetFileMD5 取文件的MD5
func GetFileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(fmt.Errorf(err.Error()))
		}
	}(file)

	// 创建一个MD5哈希对象
	hash := md5.New()

	// 将文件内容拷贝到哈希对象中
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 计算MD5值
	md5Hash := hash.Sum(nil)

	// 将MD5值转换为16进制字符串
	md5String := hex.EncodeToString(md5Hash)

	return md5String, nil
}

// GetFileSHA1 获取文件的sha1
func GetFileSHA1(path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// 创建一个SHA1哈希对象
	hash := sha1.New()

	// 将文件内容拷贝到哈希对象中
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 计算哈希值
	hashValue := hash.Sum(nil)

	// 将哈希值转换为字符串
	hashString := fmt.Sprintf("%x", hashValue)

	return hashString, nil
}

// GetFolderStats 获取文件夹的文件数量和总大小
func GetFolderStats(folderPath string) (fileCount int, totalSize int64) {
	// 检查文件夹是否存在
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		return
	}

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
			totalSize += info.Size()
		}
		return nil
	})
	return
}
