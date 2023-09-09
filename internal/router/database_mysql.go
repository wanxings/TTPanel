package router

import (
	"TTPanel/internal/api"
	"github.com/gin-gonic/gin"
)

type DatabaseMysqlRouter struct{}

func (s *DatabaseMysqlRouter) Init(Router *gin.RouterGroup) {
	mysqlRouter := Router.Group("database_mysql")
	mysqlApi := api.GroupApp.DatabaseMysqlApiApp
	{
		mysqlRouter.POST("Create", mysqlApi.Create) // 创建数据库
		mysqlRouter.POST("List", mysqlApi.List)     // 数据库列表

		mysqlRouter.POST("GetRootPwd", mysqlApi.GetRootPwd)                   // 获取数据库root密码
		mysqlRouter.POST("SetRootPwd", mysqlApi.SetRootPwd)                   // 设置数据库root密码
		mysqlRouter.POST("SetAccessPermission", mysqlApi.SetAccessPermission) // 设置数据库访问权限
		mysqlRouter.POST("GetAccessPermission", mysqlApi.GetAccessPermission) // 获取数据库访问权限
		mysqlRouter.POST("SetPwd", mysqlApi.SetPwd)                           // 设置数据库密码
		mysqlRouter.POST("CheckDeleteDatabase", mysqlApi.CheckDeleteDatabase) // 检查是否可以删除数据库
		mysqlRouter.POST("DeleteDatabase", mysqlApi.DeleteDatabase)           // 删除数据库
		mysqlRouter.POST("ImportDatabase", mysqlApi.ImportDatabase)           // 导入数据库
	}
	{
		mysqlRouter.POST("ServerList", mysqlApi.ServerList) // 数据库服务列表
		//mysqlRouter.POST("CreateServer", mysqlApi.CreateServer) // 创建数据库服务
		//mysqlRouter.POST("SaveServer", mysqlApi.SaveServer)     // 保存数据库服务
		//mysqlRouter.POST("DeleteServer", mysqlApi.DeleteServer) // 删除数据库服务
	}
	{
		mysqlRouter.POST("SyncToDB", mysqlApi.SyncToDB)   // 同步数据库到服务器
		mysqlRouter.POST("SyncGetDB", mysqlApi.SyncGetDB) // 从服务器同步数据库

	}
}
