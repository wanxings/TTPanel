package initialize

import (
	"TTPanel/internal/global"
	"TTPanel/internal/middleware"
	"TTPanel/internal/router"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Routers 初始化面板路由
func Routers() *gin.Engine {
	//gin运行模式
	gin.SetMode(global.Config.System.RunMode)

	Router := gin.New()
	//Router.HandleMethodNotAllowed = true
	//Router.Use(gin.Logger())
	//发生panic时进行恢复
	Router.Use(gin.Recovery())
	// 跨域配置
	Router.Use(middleware.CORS())
	//检查验证类型
	Router.Use(middleware.CheckAuth())
	//gzip中间件
	Router.Use(gzip.Gzip(gzip.DefaultCompression))
	//日志中间件
	Router.Use(middleware.RouterLog())

	//初始化session中间件
	if !gin.IsDebugging() {
		Router.Use(middleware.SessionInit())
	}

	routerGroupApp := router.GroupApp
	//外链下载路由
	ExternalGroup := Router.Group("/")
	{
		routerGroupApp.ExplorerRouterApp.InitExternal(ExternalGroup) //不做入口和登录鉴权，token错误会出现和入口验证失败一样的响应
	}

	if !gin.IsDebugging() {
		//验证系统入口
		Router.Use(middleware.Entrance())
		//BasicAuth认证
		//Router.Use(middleware.BasicAuth())
	}

	//静态文件路由
	routerGroupApp.StaticRouterApp.Init(Router)

	//公共路由（不做登录鉴权）
	PublicGroup := Router.Group("api")
	{
		routerGroupApp.UserBaseRouterApp.InitBase(PublicGroup) //基础路由,不做登录鉴权
	}

	//登录鉴权路由
	PrivateGroup := Router.Group("api")
	//Todo:暂时不使用SessionAuth
	PrivateGroup.Use(middleware.JWT())
	{
		//用户类路由
		routerGroupApp.UserRouterApp.Init(PrivateGroup)
		//系统功能类路由
		routerGroupApp.ExplorerRouterApp.Init(PrivateGroup)   //文件管理
		routerGroupApp.MonitorRouterApp.Init(PrivateGroup)    //监控管理
		routerGroupApp.SettingsRouterApp.Init(PrivateGroup)   //系统设置
		routerGroupApp.PanelRouterApp.Init(PrivateGroup)      //面板管理
		routerGroupApp.RecycleBinRouterApp.Init(PrivateGroup) //回收站管理
		routerGroupApp.WebSSHRouterApp.Init(PrivateGroup)     //webSSH路由
		routerGroupApp.HostRouterApp.Init(PrivateGroup)       //主机管理
		routerGroupApp.NotifyRouterApp.Init(PrivateGroup)     //通知管理
		routerGroupApp.StorageRouterApp.Init(PrivateGroup)    //存储管理
		routerGroupApp.BackupRouterApp.Init(PrivateGroup)     //备份管理
		//任务类路由
		routerGroupApp.CronTaskRouterApp.Init(PrivateGroup)  //计划任务管理
		routerGroupApp.QueueTaskRouterApp.Init(PrivateGroup) //队列任务管理
		//系统类路由
		routerGroupApp.SystemFirewallRouterApp.Init(PrivateGroup) //系统防火墙
		routerGroupApp.TTWafRouterApp.Init(PrivateGroup)          //TTWaf
		routerGroupApp.LinuxToolsRouterApp.Init(PrivateGroup)     //linux工具
		routerGroupApp.SSHManageRouterApp.Init(PrivateGroup)      //SSH管理
		routerGroupApp.TaskManagerRouterApp.Init(PrivateGroup)    //任务管理
		//扩展路由
		routerGroupApp.ExtensionNginxRouterApp.Init(PrivateGroup)      //Nginx
		routerGroupApp.ExtensionPHPRouterApp.Init(PrivateGroup)        //PHP
		routerGroupApp.ExtensionMysqlRouterApp.Init(PrivateGroup)      //Mysql
		routerGroupApp.ExtensionDockerRouterApp.Init(PrivateGroup)     //Docker
		routerGroupApp.ExtensionPhpmyadminRouterApp.Init(PrivateGroup) //phpmyadmin
		routerGroupApp.ExtensionRedisRouterApp.Init(PrivateGroup)      //Redis
		routerGroupApp.ExtensionNodejsRouterApp.Init(PrivateGroup)     //Nodejs

		//Todo:插件路由
		//pluginRouter.InitPluginRouter(PrivateGroup) //插件
		//项目类路由
		routerGroupApp.ProjectRouterApp.Init(PrivateGroup)        //项目通用接口
		routerGroupApp.ProjectGeneralRouterApp.Init(PrivateGroup) //通用项目管理
		routerGroupApp.ProjectPHPRouterApp.Init(PrivateGroup)     //php项目管理
		routerGroupApp.SSLRouterApp.Init(PrivateGroup)            //SSL证书管理
		//数据库类路由
		routerGroupApp.DatabaseMysqlRouterApp.Init(PrivateGroup) //数据库管理

		//日志审计类路由
		routerGroupApp.LogAuditRouterApp.Init(PrivateGroup) //日志审计管理

	}
	// 默认404
	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Not Found",
		})
	})

	// 默认405
	Router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"msg":  "Method Not Allowed",
		})
	})

	return Router
}
