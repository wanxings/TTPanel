package initialize

import (
	"TTPanel/internal/conf"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitPanelDB(m conf.Sqlite) *gorm.DB {
	return NewSqlDB(m.PanelDsn(), m)
}
func InitTTWafDB(m conf.Sqlite) *gorm.DB {
	return NewSqlDB(m.TTWafDsn(), m)
}

func NewSqlDB(dsn string, m conf.Sqlite) *gorm.DB {
	newLogger := logger.New(
		logrus.StandardLogger(), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second,     // 慢 SQL 阈值
			LogLevel:                  m.GetLogLevel(), // 日志级别
			IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,           // 禁用彩色打印
		},
	)

	config := &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	plugin := dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour)

	var (
		db  *gorm.DB
		err error
	)

	logrus.Debugln("use SQLite DB")
	if db, err = gorm.Open(sqlite.Open(dsn), config); err == nil {
		err = db.Use(plugin)
		if err != nil {
			panic(fmt.Sprintf("SQLite failed to db.Use(plugin): %v", err))
		}
	} else {
		panic(fmt.Sprintf("SQLite failed to connect to database: %v", err))
	}

	return db

}
