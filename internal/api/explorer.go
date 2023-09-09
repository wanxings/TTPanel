package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ExplorerApi struct{}

// GetDir
// @Tags      System
// @Summary   获取文件目录
// @Router    /api/explorer/GetNetWork [post]
func (s *ExplorerApi) GetDir(c *gin.Context) {
	response := app.NewResponse(c)
	ResponseData := make(map[string]interface{})
	var err error
	//获取参数
	param := request.GetDirR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetDir.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	offset, limit := app.GetOffsetLimits(param.Page, param.Limit)
	if util.StrIsEmpty(param.Path) {
		param.Path = "/"
	}

	if !util.StrIsEmpty(param.Query) && param.All {
		ResponseData, err = ServiceGroupApp.ExplorerServiceApp.SearchDir(param.Query, param.Path, param.Sort, param.Reverse)
		if err != nil {
			response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		}
	} else {
		ResponseData, err = ServiceGroupApp.ExplorerServiceApp.GetDir(param.Query, param.Path, param.Sort, param.Reverse, offset, limit)
		if err != nil {
			response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		}

	}
	//获取收藏夹
	favorites, err := ServiceGroupApp.ExplorerServiceApp.GetFavorites()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	ResponseData["favorites"] = favorites
	response.ToResponse(&ResponseData)
}

// GetFileBody
// @Tags      System
// @Summary   获取文件内容
// @Router    /api/explorer/GetNetWork [post]
func (s *ExplorerApi) GetFileBody(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetFileBody.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	ResponseData, err := ServiceGroupApp.ExplorerServiceApp.ReadFile(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(&ResponseData)
}

// SaveFileBody
// @Tags      System
// @Summary   保存文件内容
// @Router    /api/explorer/SaveFileBody [post]
func (s *ExplorerApi) SaveFileBody(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SaveFileBodyR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SaveFileBody.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.SaveFileBody(param.Type, param.Path, param.Body)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Save", map[string]any{"Path": param.Path}))
	response.ToResponseMsg(helper.Message("tips.SaveSuccess"))

}

// GetAttribute
// @Tags      System
// @Summary   获取文件（夹）属性
// @Router    /api/explorer/SaveFileBody [post]
func (s *ExplorerApi) GetAttribute(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetAttribute.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	ResponseData, err := ServiceGroupApp.ExplorerServiceApp.GetAttribute(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(&ResponseData)

}

// BatchSetFileAccess
// @Tags      System
// @Summary   批量设置文件（夹）权限
// @Router    /api/explorer/SaveFileBody [post]
func (s *ExplorerApi) BatchSetFileAccess(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchSetFileAccessR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchSetFileAccess.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.BatchChangePermission(param.PathList, param.UserName, param.Permissions)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.SetPermission", map[string]any{"Path": strings.Join(param.PathList, ","), "Permission": param.Permissions, "User": param.UserName}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// BatchDeleteDirFile
// @Tags      System
// @Summary   批量删除文件（夹）
// @Router    /api/explorer/SaveFileBody [post]
func (s *ExplorerApi) BatchDeleteDirFile(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchDeleteDirFileR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchDeleteDirFile.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	err := ServiceGroupApp.ExplorerServiceApp.BatchDelete(param.PathList)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Delete", map[string]any{"Path": strings.Join(param.PathList, ",")}))
	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// Rename
// @Tags      System
// @Summary   重命名文件（夹）
// @Router    /api/explorer/SaveFileBody [post]
func (s *ExplorerApi) Rename(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.RenameR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Rename.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.Rename(param.OldPath, param.NewPath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Rename", map[string]any{"OldPath": param.OldPath, "NewPath": param.NewPath}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// BatchCheckExistsFiles
// @Tags      System
// @Summary   批量检查文件（夹）是否存在
// @Router    /api/explorer/BatchCheckExistsFiles [post]
func (s *ExplorerApi) BatchCheckExistsFiles(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchCheckExistsFilesR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCheckExistsFiles.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.BatchCheckExistsFiles(param.InitPath, param.FileList, param.CheckPath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)

}

// BatchCopy
// @Tags      System
// @Summary   批量复制文件（夹）
// @Router    /api/explorer/BatchCopy [post]
func (s *ExplorerApi) BatchCopy(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchCopyMoveR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchCopy.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.BatchCopy(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Copy", map[string]any{"InitPath": param.InitPath, "FileList": strings.Join(param.FileList, ","), "ToPath": param.ToPath}))
	response.ToResponseMsg(helper.Message("tips.CopySuccess"))
}

// BatchMove
// @Tags      System
// @Summary   批量移动文件（夹）
// @Router    /api/explorer/BatchMove [post]
func (s *ExplorerApi) BatchMove(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchCopyMoveR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchMove.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.BatchMove(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Move", map[string]any{"InitPath": param.InitPath, "FileList": strings.Join(param.FileList, ","), "ToPath": param.ToPath}))
	response.ToResponseMsg(helper.Message("tips.MoveSuccess"))
}

// CreateDir
// @Tags      System
// @Summary   创建文件夹
// @Router    /api/explorer/CreateDir [post]
func (s *ExplorerApi) CreateDir(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("CreateDir.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.CreateDir(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.CreateDir", map[string]any{"Path": param.Path}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// CreateFile
// @Tags      System
// @Summary   创建文件
// @Router    /api/explorer/CreateFile [post]
func (s *ExplorerApi) CreateFile(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("CreateDir.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.CreateFile(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.CreateFile", map[string]any{"Path": param.Path}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// CreateSymlink
// @Tags      System
// @Summary   创建符号链接
// @Router    /api/explorer/CreateSymlink [post]
func (s *ExplorerApi) CreateSymlink(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.RenameR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("CreateSymlink.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.CreateSymlink(param.OldPath, param.NewPath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.CreateSymlink", map[string]any{"NewPath": param.NewPath, "OldPath": param.OldPath}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// CreateDuplicate
// @Tags      System
// @Summary   创建副本
// @Router    /api/explorer/CreateDuplicate [post]
func (s *ExplorerApi) CreateDuplicate(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("CreateDuplicate.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	newPath, err := ServiceGroupApp.ExplorerServiceApp.CreateDuplicate(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.CreateDuplicate", map[string]any{"NewPath": newPath, "OldPath": param.Path}))
	response.ToResponseMsg(helper.Message("tips.CreateSuccess"))
}

// GetPathSize
// @Tags      System
// @Summary   获取目录大小
// @Router    /api/explorer/GetPathSize [post]
func (s *ExplorerApi) GetPathSize(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetPathSize.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.GetPathSize(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(data)
}

// CheckFileExists
// @Tags      System
// @Summary   检查文件（夹）是否存在
// @Router    /api/explorer/CheckFileExists [post]
func (s *ExplorerApi) CheckFileExists(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CommonPathR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("CheckFileExists.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	isExists := ServiceGroupApp.ExplorerServiceApp.BatchCheckFilesExist([]string{param.Path})
	response.ToResponse(isExists)
}

// ChunkUpload
// @Tags      System
// @Summary   分片上传文件
// @Router    /api/explorer/ChunkUpload [post]
func (s *ExplorerApi) ChunkUpload(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.UploadR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ChunkUpload.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//获取目录权限
	dirPerm := util.GetFilePerm(param.Path)
	savePath := param.Path + "/" + param.Name
	//如果不覆盖并且已存在同文件，则直接返回成功
	if !param.Cover && util.PathExists(savePath) {
		response.ToResponse(nil)
	}
	uploadBlob, _ := param.Blob.Open()
	if param.ChunkSize > 0 { //分片上传
		tmpFilePath := param.Path + "/" + param.Name + fmt.Sprintf("%d", param.Size) + ".upload.tmp"

		//可能是文件夹上传，首先尝试创建文件夹
		_ = os.MkdirAll(filepath.Dir(tmpFilePath), dirPerm)
		var dSize int64
		dSize = 0
		if util.IsFile(tmpFilePath) {
			dSize = util.GetFileSize(tmpFilePath)
		}
		if dSize != param.Start {
			response.ToResponse(dSize)
			return
		}
		fileOP, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, dirPerm)
		defer func(fileOP *os.File) {
			_ = fileOP.Close()
		}(fileOP)
		if err != nil {
			response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
			return
		}
		_, err = io.Copy(fileOP, uploadBlob)

		if param.Size < param.Start+param.ChunkSize {
			//文件上传完成，尝试删除重复文件
			_ = os.Remove(savePath)
			//重命名文件
			err = os.Rename(tmpFilePath, savePath)
			if err != nil {
				response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
				return
			}

			go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
				c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Upload", map[string]any{"Path": param.Path + "/" + param.Name}))
		}
		response.ToResponse(err)
	}
}

// Upload
// @Tags      System
// @Summary   上传文件
// @Router    /api/explorer/Upload [post]
func (s *ExplorerApi) Upload(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.UploadR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Upload.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	var err error

	//获取目录权限
	dirPerm := util.GetFilePerm(param.Path)
	savePath := param.Path + "/" + param.Name
	//如果不覆盖并且已存在同文件，则直接返回成功
	if !param.Cover && util.PathExists(savePath) {
		response.ToResponse(nil)
	}
	uploadBlob, _ := param.Blob.Open()
	//可能是文件夹上传，首先尝试创建文件夹
	_ = os.MkdirAll(filepath.Dir(savePath), dirPerm)

	//尝试删除文件
	_ = os.Remove(savePath)
	fileOP, err := os.OpenFile(savePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, dirPerm)
	defer func(fileOP *os.File) {
		_ = fileOP.Close()
	}(fileOP)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	_, err = io.Copy(fileOP, uploadBlob)
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Upload", map[string]any{"Path": param.Path + "/" + param.Name}))
	response.ToResponse(err)

}

// Compress
// @Tags      System
// @Summary   压缩文件
// @Router    /api/explorer/Compress [post]
func (s *ExplorerApi) Compress(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.CompressR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Compress.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	logoPath, err := ServiceGroupApp.ExplorerServiceApp.Compress(false, param.Path, param.FileList, param.CompressPath, param.CompressType, param.Password)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Compress", map[string]any{"Path": param.CompressPath}))

	response.ToResponse(logoPath)
}

// Decompress
// @Tags      System
// @Summary   解压文件
// @Router    /api/explorer/Decompress [post]
func (s *ExplorerApi) Decompress(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DecompressR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("Decompress.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	logPath, err := ServiceGroupApp.ExplorerServiceApp.Decompress(false, param.CompressFilePath, param.DecompressDirPath, param.Password)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.Decompress", map[string]any{"CompressFilePath": param.CompressFilePath, "DecompressDirPath": param.DecompressDirPath}))

	response.ToResponse(logPath)
}

// GetLogContent
// @Tags      System
// @Summary   获取日志
// @Router    /api/explorer/GetLogContent [post]
func (s *ExplorerApi) GetLogContent(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.GetLogContentR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetLog.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	result, err := ServiceGroupApp.ExplorerServiceApp.GetLogContent(param.Path, param.Location, param.Line)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(result)
}

// ClearLogContent
// @Tags      System
// @Summary   清空日志
// @Router    /api/explorer/ClearLogContent [post]
func (s *ExplorerApi) ClearLogContent(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.ClearLogContentR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("ClearLogContent.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.ClearLogContent(param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.ClearLogFile", map[string]any{"Path": param.Path}))

	response.ToResponseMsg(helper.Message("tips.ClearSuccess"))
}

// SearchFileContent
// @Tags      System
// @Summary   搜索文件内容
// @Router    /api/explorer/SearchFileContent [post]
func (s *ExplorerApi) SearchFileContent(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SearchFileContentR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SearchFileContent.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.SearchFileContent(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// RemoteDownload
// @Tags      System
// @Summary   远程下载文件
// @Router    /api/explorer/RemoteDownload [post]
func (s *ExplorerApi) RemoteDownload(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.RemoteDownloadR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RemoteDownload.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	key, err := ServiceGroupApp.ExplorerServiceApp.RemoteDownload(&param)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.RemoteDownload", map[string]any{"Url": param.Url, "Path": param.SavePath}))

	response.ToResponse(key)
}

// RemoteDownloadProcess
// @Tags      System
// @Summary   远程下载文件进度
// @Router    /api/explorer/RemoteDownloadProcess [post]
func (s *ExplorerApi) RemoteDownloadProcess(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.RemoteDownloadProcessR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("RemoteDownloadProcess.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.RemoteDownloadProcess(param.Key)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// Download
// @Tags      System
// @Summary   下载文件
// @Router    /api/explorer/Download [get]
func (s *ExplorerApi) Download(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	filePath := c.Query("file_path")
	if !util.IsFile(filePath) {
		response.ToErrorResponse(errcode.ServerError.WithDetails("not found"))
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	global.Log.Warnf("Download FilePath:%s,IP:%s", filePath, c.ClientIP())
	c.File(filePath)
}

// ExternalDownload
// @Tags      System
// @Summary   外部下载文件
// @Router    /api/explorer/ExternalDownload [post]
func (s *ExplorerApi) ExternalDownload(c *gin.Context) {
	//获取参数
	dToken := c.Query("token")
	//获取文件源路径
	if filePath, ok := global.GoCache.Get("tmp_d_" + dToken); ok {
		if !util.IsFile(filePath.(string)) {
			global.Log.Errorf("ExternalDownload.IsFile errs: file %s is not found", filePath.(string))
			c.Abort()
			return
		}
		global.GoCache.Delete("tmp_d_" + dToken)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath.(string))))
		c.File(filePath.(string))
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.GetExternalDownloadByToken(dToken)
	if err != nil || data == nil {
		global.Log.Errorf("ExternalDownload.IsFile errs: %v", err)
		fmt.Printf("ExternalDownload.IsFile errs:%v", err)
		// 入口不正确响应状态码，如果是200,响应默认提示
		if global.Config.System.EntranceErrorCode == 200 {
			c.File(fmt.Sprintf("%s/data/panel_entry_prompt.html", global.Config.System.PanelPath))

		} else {
			c.AbortWithStatus(global.Config.System.EntranceErrorCode)
		}
		c.Abort()
		return
	}
	if !util.IsFile(data.FilePath) {
		global.Log.Errorf("ExternalDownload.IsFile errs: file %s is not found", data.FilePath)
		c.Abort()
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(data.FilePath)))
	global.Log.Warnf("ExternalDownload FilePath:%s,IP:%s", data.FilePath, c.ClientIP())
	c.File(data.FilePath)
}

// BatchOperateSpecialPermission
// @Tags      System
// @Summary   批量操作特殊权限
// @Router    /api/explorer/BatchOperateSpecialPermission [post]
func (s *ExplorerApi) BatchOperateSpecialPermission(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.BatchOperateSpecialPermissionR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("BatchOperateSpecialPermission.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	failedErrors := ServiceGroupApp.ExplorerServiceApp.BatchSetSpecialPermission(param.Action, param.PathList, param.Permissions, param.All)
	errStr := make([]string, len(failedErrors))
	for i, err := range errs {
		errStr[i] = err.Error()
	}
	if len(errStr) > 0 {
		response.ToErrorResponse(errcode.ServerError.WithDetails(errStr...))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.SetSpecialPermission", map[string]any{"Path": strings.Join(param.PathList, " "), "Permission": param.Permissions}))

	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// SetRemark
// @Tags      System
// @Summary   设置备注
// @Router    /api/explorer/SetRemark [post]
func (s *ExplorerApi) SetRemark(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.SetRemarkR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("SetRemark.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.SetRemark(param.Path, param.Remark)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.SetRemark", map[string]any{"Path": param.Path, "Remark": param.Remark}))
	response.ToResponseMsg(helper.Message("tips.SetSuccess"))
}

// FavoritesList
// @Tags      System
// @Summary   收藏夹列表
// @Router    /api/explorer/FavoritesList [post]
func (s *ExplorerApi) FavoritesList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.ExplorerServiceApp.FavoritesList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseList(data, len(data), 0, 0)
}

// OperateFavorites
// @Tags      System
// @Summary   操作收藏夹
// @Router    /api/explorer/OperateFavorites [post]
func (s *ExplorerApi) OperateFavorites(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.OperateFavoritesR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("OperateFavorites.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.OperateFavorites(param.Action, param.Path)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.OperateFavorites", map[string]any{"Path": param.Path, "Action": param.Action}))
	response.ToResponseMsg(helper.Message("tips.OperateSuccess"))
}

// GenerateDownloadExternalLink
// @Tags      System
// @Summary   生成外部下载链接
// @Router    /api/explorer/GenerateDownloadExternalLink [post]
func (s *ExplorerApi) GenerateDownloadExternalLink(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.GenerateDownloadExternalLinkR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GenerateDownloadExternalLink.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	data, err := ServiceGroupApp.ExplorerServiceApp.GenerateDownloadExternalLink(param.FilePath, param.ExpireTime, param.Description)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.GenerateDownloadExternalLink", map[string]any{"Path": param.FilePath}))
	response.ToResponse(data)
}

// DownloadExternalLinkList
// @Tags      System
// @Summary   外部下载链接列表
// @Router    /api/explorer/DownloadExternalLinkList [post]
func (s *ExplorerApi) DownloadExternalLinkList(c *gin.Context) {
	response := app.NewResponse(c)
	data, err := ServiceGroupApp.ExplorerServiceApp.DownloadExternalLinkList()
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(data)
}

// DeleteDownloadExternalLink
// @Tags      System
// @Summary   删除外部下载链接
// @Router    /api/explorer/DeleteDownloadExternalLink [post]
func (s *ExplorerApi) DeleteDownloadExternalLink(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.DeleteDownloadExternalLinkR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("DeleteDownloadExternalLink.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	err := ServiceGroupApp.ExplorerServiceApp.DeleteDownloadExternalLink(param.ID)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(
		c, constant.OperationLogTypeByExplorer, helper.MessageWithMap("explorer.DeleteDownloadExternalLink", map[string]any{"ID": param.ID}))

	response.ToResponseMsg(helper.Message("tips.DeleteSuccess"))
}

// GetFileTemporaryDownloadLink
// @Tags      System
// @Summary   获取文件临时下载链接
// @Router    /api/explorer/GetFileTemporaryDownloadLink [post]
func (s *ExplorerApi) GetFileTemporaryDownloadLink(c *gin.Context) {
	response := app.NewResponse(c)
	//获取参数
	param := request.GetFileTemporaryDownloadLinkR{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("GetFileTemporaryDownloadLink.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	data, err := ServiceGroupApp.ExplorerServiceApp.GetFileTemporaryDownloadLink(param.FilePath)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	global.Log.Warnf("GetFileTemporaryDownloadLink -> FilePath:%s,IP:%s,Key:%s", param.FilePath, c.ClientIP(), data)

	response.ToResponse(data)
}
