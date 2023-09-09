package response

import (
	"TTPanel/pkg/util"
	"syscall"
)

type DirFileAttribute struct {
	Name              string          `json:"name"`               //文件(夹)名
	Path              string          `json:"path"`               //文件(夹)路径
	Size              int64           `json:"size"`               //文件大小
	IsDir             bool            `json:"is_dir"`             //是否是文件夹
	IsLink            bool            `json:"is_link"`            //是否是软连接
	Owner             string          `json:"owner"`              //拥有者
	OwnerId           uint32          `json:"owner_id"`           //拥有者id
	Group             string          `json:"group"`              //所属组
	GroupId           uint32          `json:"group_id"`           //所属组id
	SpecialPermission string          `json:"special_permission"` //特殊权限
	Perm              string          `json:"perm"`               //权限
	PermString        string          `json:"perm_string"`        //权限字符串
	ModTime           int64           `json:"mod_time"`           //修改时间
	StatT             *syscall.Stat_t `json:"stat_t"`             //文件信息
	FileHistoryList   any             `json:"file_history_list"`  //文件历史版本
}

type BatchCheckExistsFilesP struct {
	RepeatFiles               []*ComparedFiles //重复文件
	OldExistFiles             []string         //旧文件不存在
	ExistSameNameFiles        []*ComparedFiles //存在同名文件
	TargetIsSubDirectoryFiles []*ComparedFiles //目标是子目录
}

type ComparedFiles struct {
	OldFile *util.DirInfo `json:"old_file"` //旧文件
	NewFile *util.DirInfo `json:"new_file"` //新文件
}
