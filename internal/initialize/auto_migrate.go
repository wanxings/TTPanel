package initialize

import (
	"TTPanel/internal/global"
	"TTPanel/internal/model"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func InitMigrate() {
	var err error
	//备份表
	var backupTable = &gormigrate.Migration{
		ID: "backupTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Backup{})
		},
	}
	//计划任务表
	var cronTaskTable = &gormigrate.Migration{
		ID: "cronTaskTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.CronTask{})
		},
	}
	//数据库服务表
	var databaseServerTable = &gormigrate.Migration{
		ID: "databaseServerTable-202307030703",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.DatabaseServer{})
		},
	}
	//数据库表
	var databasesTable = &gormigrate.Migration{
		ID: "databasesTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Databases{})
		},
	}
	//docker...表
	var dockersTable = &gormigrate.Migration{
		ID: "dockersTable-202308190703",
		Migrate: func(tx *gorm.DB) error {
			err = tx.AutoMigrate(&model.DockerCompose{})
			if err != nil {
				return err
			}
			//docker-compose模板表
			err = tx.AutoMigrate(&model.DockerComposeTemplate{})
			if err != nil {
				return err
			}
			//docker仓库表
			err = tx.AutoMigrate(&model.DockerRepository{})
			if err != nil {
				return err
			}
			_, _ = (&model.DockerRepository{
				ID:     1,
				Name:   "官方仓库",
				Url:    "https://docker.io",
				Remark: "Docker官方仓库",
			}).Create(global.PanelDB)
			return nil
		},
	}
	//外链下载表
	var externalDownloadTable = &gormigrate.Migration{
		ID: "externalDownloadTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.ExternalDownload{})
		},
	}
	//系统防火墙规则..表
	var firewallsRuleTable = &gormigrate.Migration{
		ID: "firewallsRuleTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			err = tx.AutoMigrate(&model.FirewallRuleForward{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.FirewallRuleIp{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.FirewallRulePort{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	//主机...表
	var hostsTable = &gormigrate.Migration{
		ID: "hostsTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			err = tx.AutoMigrate(&model.Host{})
			if err != nil {
				return err
			}
			//主机分类表
			err = tx.AutoMigrate(&model.HostCategory{})
			if err != nil {
				return err
			}
			_, _ = (&model.HostCategory{
				ID:     1,
				Name:   "默认分类",
				Remark: "默认分类",
			}).Create(global.PanelDB)
			if err != nil {
				return err
			}
			//主机快捷命令表
			err = tx.AutoMigrate(&model.HostShortcutCommand{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	//监控...表
	var monitorsTable = &gormigrate.Migration{
		ID: "monitorsTable-202307030703",
		Migrate: func(tx *gorm.DB) error {
			err = tx.AutoMigrate(&model.MonitorAppCrash{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.MonitorNetwork{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.MonitorHighResourceProcesses{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.MonitorMetrics{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.MonitorIo{})
			if err != nil {
				return err
			}
			err = tx.AutoMigrate(&model.MonitorEvent{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	//通知通道表
	var notifyChannelTable = &gormigrate.Migration{
		ID: "notifyChannelTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.NotifyChannel{})
		},
	}
	//面板操作日志表
	var operationLogTable = &gormigrate.Migration{
		ID: "operationLogTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.OperationLog{})
		},
	}
	//项目...表
	var projectsTable = &gormigrate.Migration{
		ID: "projectsTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			err = tx.AutoMigrate(&model.Project{})
			if err != nil {
				return err
			}
			//项目分类表
			err = tx.AutoMigrate(&model.ProjectCategory{})
			if err != nil {
				return err
			}
			_, _ = (&model.ProjectCategory{
				ID:   1,
				Name: "默认分类",
				Ps:   "默认分类",
			}).Create(global.PanelDB)
			//项目域名表
			err = tx.AutoMigrate(&model.ProjectDomain{})
			if err != nil {
				return err
			}
			//项目绑定表
			err = tx.AutoMigrate(&model.ProjectBinding{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	//队列任务表
	var queueTable = &gormigrate.Migration{
		ID: "queueTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.QueueTask{})
		},
	}
	//存储表
	var storageTable = &gormigrate.Migration{
		ID: "storageTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Storage{})
		},
	}
	//用户表
	var userTable = &gormigrate.Migration{
		ID: "userTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.User{})
		},
	}
	//面板数据库迁移
	panel := gormigrate.New(global.PanelDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		backupTable,
		cronTaskTable,
		databaseServerTable,
		databasesTable,
		dockersTable,
		externalDownloadTable,
		firewallsRuleTable,
		hostsTable,
		monitorsTable,
		notifyChannelTable,
		operationLogTable,
		projectsTable,
		queueTable,
		storageTable,
		userTable,
	})
	if err = panel.Migrate(); err != nil {
		global.Log.Error(fmt.Sprintf("Migrate Panel Error: %v", err))
		panic(err)
	}
	//ttwaf拦截ip日志表
	var blockIPTable = &gormigrate.Migration{
		ID: "blockIPTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.TTWafBlockIpLog{})
		},
	}
	//ttwaf屏蔽ip日志表
	var banIPTable = &gormigrate.Migration{
		ID: "banIPTable-202307030702",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.TTWafBanIpLog{})
		},
	}
	ttwaf := gormigrate.New(global.TTWafDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		blockIPTable,
		banIPTable,
	})
	if err = ttwaf.Migrate(); err != nil {
		global.Log.Error(fmt.Sprintf("Migrate TTwaf Error: %v", err))
		panic(err)
	}

	//analytics表处理 Todo:好像有点问题
	//serverNameS := []string{}
	//// 遍历指定目录下的文件
	//err = filepath.Walk(global.Config.System.PanelPath+"/extensions/nginx/vhost/main", func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	// 检查文件是否以 .conf 结尾
	//	if strings.HasSuffix(info.Name(), ".conf") {
	//		// 检查文件名称是否包含 1_1_default 或 nginx_firewall
	//		if !strings.Contains(info.Name(), "1_1_default") && !strings.Contains(info.Name(), "nginx_firewall") {
	//			// 提取文件名称
	//			name := strings.TrimSuffix(info.Name(), ".conf")
	//			serverNameS = append(serverNameS, name)
	//		}
	//	}
	//
	//	return nil
	//})
	//dbDir := global.Config.System.ServerPath + "/ttwaf/analytics"
	//if len(serverNameS) > 0 {
	//	for _, serverName := range serverNameS {
	//		dbPath := fmt.Sprintf("%s/%s/data.db", dbDir, serverName)
	//		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
	//			_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	//			_, err = os.Create(dbDir)
	//			if err != nil {
	//				global.Log.Errorf("%s Create Error:%s", dbDir, err.Error())
	//			}
	//		}
	//
	//		// 连接数据库
	//		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// 自动迁移表结构
	//		err = db.AutoMigrate(models...)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//}

	global.Log.Info("Migration done")
	_, _ = fmt.Fprintf(color.Output, "Migration....   %s\n",
		color.GreenString("done"),
	)
}
