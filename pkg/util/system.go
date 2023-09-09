package util

import (
	"fmt"
	"strings"
	"syscall"
)

// GetDiskSpaceAndInode 获取指定路径所在分区的磁盘空间和Inode 信息
func GetDiskSpaceAndInode(path string) (totalSize, freeSize, availSize, totalInodes, freeInodes int64, err error) {
	var stat syscall.Statfs_t
	if err = syscall.Statfs(path, &stat); err != nil {
		fmt.Println(err)
		return
	}
	// 获取磁盘空间信息
	blockSize := stat.Bsize
	totalSize = int64(stat.Blocks) * blockSize
	freeSize = int64(stat.Bfree) * blockSize
	availSize = int64(stat.Bavail) * blockSize
	// 获取 Inode 信息
	totalInodes = int64(stat.Files)
	freeInodes = int64(stat.Ffree)
	return
}

func IsArch64() bool {
	shell, err := ExecShell("uname -m")
	if err != nil {
		return false
	}
	if strings.Contains(shell, "aarch64") {
		return true
	} else {
		return false
	}
}
