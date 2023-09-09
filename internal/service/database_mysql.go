package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DatabaseMysqlService struct {
}

var MysqlCoding = map[string]string{
	"utf8":    "utf8_general_ci",
	"utf8mb4": "utf8mb4_general_ci",
	"gbk":     "gbk_chinese_ci",
	"big5":    "big5_chinese_ci",
}

// Create 创建数据库
func (s *DatabaseMysqlService) Create(param *request.CreateMysqlR) error {
	//格式化参数
	param.DatabaseName = util.ClearStr(param.DatabaseName)
	param.UserName = util.ClearStr(param.UserName)
	param.Password = util.ClearStr(param.Password)
	param.AccessPermission = util.ClearStr(param.AccessPermission)
	codingStr := MysqlCoding[param.Coding]
	if !util.IsGeneral(param.DatabaseName) || !util.IsGeneral(param.UserName) {
		return errors.New(helper.Message("database.mysql.NameAndUserNameIsIllegal"))
	}
	if util.StrIsEmpty(param.Coding) || util.StrIsEmpty(codingStr) {
		return errors.New(helper.Message("database.mysql.DatabaseCharsetError"))
	}
	//长度校验
	if len(param.DatabaseName) < 1 || len(param.UserName) < 1 || len(param.Password) < 1 {
		return errors.New("database name, user name, and password is empty")
	}
	//查询数据库是否已存在
	get, err := (&model.Databases{Name: param.DatabaseName}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID > 0 {
		return errors.New(helper.Message("database.mysql.DatabaseAlreadyExists"))
	}

	DBService, err := s.NewMysqlServiceBySid(param.Sid)
	if err != nil {
		return err
	}
	//从MySQL查询数据库是否已存在
	exists, err := s.DatabaseExists(DBService, param.DatabaseName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("not found database")
	}
	//创建数据库
	CResult := DBService.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s;", param.DatabaseName, param.Coding, codingStr))
	if CResult.Error != nil {
		return CResult.Error
	}
	//删除用户
	_ = DBService.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'localhost';", param.UserName))
	//逗号分割权限
	for _, v := range strings.Split(param.AccessPermission, ",") {
		_ = DBService.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s';", param.UserName, v))
	}
	//创建用户
	err = s.CreateDatabaseUser(DBService, param.DatabaseName, param.UserName, param.Password, param.AccessPermission)
	if err != nil {
		return err
	}
	//保存到面板数据库
	_, err = (&model.Databases{
		Name:     param.DatabaseName,
		Username: param.UserName,
		Password: param.Password,
		Accept:   param.AccessPermission,
		Sid:      param.Sid,
		Pid:      param.Pid,
		Type:     constant.DataBaseTypeByMysql,
		Ps:       param.Ps,
	}).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// List 列出数据库
func (s *DatabaseMysqlService) List(param *request.ListMysqlR, offset, limit int) (map[string]interface{}, error) {
	var databases []*model.Databases
	list := make(map[string]interface{})
	var count int64
	var err error
	queryC := model.ConditionsT{"ORDER": "create_time DESC"}

	if param.Sid != 0 {
		queryC["sid"] = param.Sid
		//
	}
	if !util.StrIsEmpty(param.Query) {
		queryC["FIXED"] = "ps LIKE '%" + param.Query + "%' " //备注
		queryC["OR"] = "name LIKE '%" + param.Query + "%' "
	}
	databases, count, err = (&model.Databases{}).List(global.PanelDB, &queryC, offset, limit)
	if err != nil {
		return list, err
	}
	for _, database := range databases {
		database.BackupCount, _ = (&model.Backup{}).Count(global.PanelDB, &model.ConditionsT{"pid": database.ID}, &model.ConditionsT{})
	}
	list["list"] = databases
	list["total_rows"] = count
	return list, nil
}

// ServerList 列出数据库服务
func (s *DatabaseMysqlService) ServerList(param *request.ListMysqlServerR) (map[string]interface{}, error) {
	list := make(map[string]interface{})
	databases, count, err := (&model.DatabaseServer{}).List(global.PanelDB, &model.ConditionsT{"db_type": param.DBType, "ORDER": "create_time DESC"}, 0, 0)
	if err != nil {
		return list, err
	}
	list["list"] = databases
	list["total_rows"] = count
	return list, err
}

// GetRootPwd 获取本地MySQL数据库root密码
func (s *DatabaseMysqlService) GetRootPwd() (string, error) {
	return global.Config.System.MysqlRootPassword, nil
}

// SetRootPwd 设置MySQL数据库root密码
func (s *DatabaseMysqlService) SetRootPwd(param *request.SetRootPwdR) error {
	if util.StrIsEmpty(param.Password) {
		return errors.New("password is empty")
	}

	//修改本地mysql密码
	if param.Sid == 0 {
		err := s.ChangeLocalRootPassword(param.Password)
		if err != nil {
			return err
		}

	} else {
		//修改远程mysql密码
		err := s.ChangeRemoteRootPassword(param.Sid, param.Password)
		if err != nil {
			return err
		}
		//保存到数据库
		get, err := (&model.DatabaseServer{ID: param.Sid}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		get.Password = param.Password
		err = (get).Update(global.PanelDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncGetDB 从服务器同步数据库
func (s *DatabaseMysqlService) SyncGetDB(param *request.SyncGetDBR) []error {
	var dbType int
	var errs []error
	if param.Sid != 0 {
		dbType = 2
	}
	dbServer, err := s.NewMysqlServiceBySid(param.Sid)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	data, err := dbServer.Raw("show databases").Rows()
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	users, err := dbServer.Raw("select User,Host from mysql.user where User!='root' AND Host!='localhost' AND Host!=''").Rows()
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	nameArr := []string{"information_schema", "performance_schema", "mysql", "sys"}
	for data.Next() {
		var value string
		_ = data.Scan(&value)
		b := false
		for _, key := range nameArr {
			if value == key {
				b = true
				break
			}
		}
		if b {
			continue
		}
		get, err := (&model.Databases{Name: value}).Get(global.PanelDB)
		if err != nil {
			continue
		}
		if get.ID != 0 {
			continue
		}

		host := "127.0.0.1"
		for users.Next() {
			var user, h string
			_ = users.Scan(&user, &h)
			if value == user {
				host = h
				break
			}
		}

		ps := "同步数据库"

		if matched, _ := regexp.MatchString(`^[\w+\.-]+$`, value); !matched {
			continue
		}

		if _, err := (&model.Databases{
			Name:     value,
			DbType:   dbType,
			Sid:      param.Sid,
			Username: value,
			Password: "",
			Accept:   host,
			Ps:       ps,
		}).Create(global.PanelDB); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errs
}

// SyncToDB 同步数据库到服务器
func (s *DatabaseMysqlService) SyncToDB(ids []int64) []error {
	var errs []error
	queryC := model.ConditionsT{"type = ?": constant.DataBaseTypeByMysql}
	if len(ids) > 0 {
		queryC["id IN ?"] = ids
	}
	//获取数据库列表
	list, _, err := (&model.Databases{}).List(global.PanelDB, &queryC, 0, 0)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, value := range list {
		dbServer, err := s.NewMysqlServiceBySid(value.Sid)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		_, err = dbServer.Raw("create database if not exists " + value.Name).Rows()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errs
}

// SetAccessPermission 设置数据库访问权限
func (s *DatabaseMysqlService) SetAccessPermission(param *request.SetAccessPermissionR) error {
	databaseInfo, err := (&model.Databases{Username: param.UserName}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	dbServer, err := s.NewMysqlServiceBySid(databaseInfo.Sid)
	if err != nil {
		return err
	}
	param.AccessPermission = strings.TrimSpace(param.AccessPermission)
	if util.StrIsEmpty(param.AccessPermission) {
		return errors.New("AccessPermission is empty")
	}
	_, err = dbServer.Raw("show databases").Rows()
	if err != nil {
		return err
	}
	users, err := dbServer.Raw("select Host from mysql.user where User=? AND Host!='localhost'", databaseInfo.Username).Rows()
	if err != nil {
		return err
	}
	for users.Next() {
		var us string
		_ = users.Scan(&us)
		dbServer.Exec(fmt.Sprintf("drop user '%s'@'%s'", databaseInfo.Username, us))
		dbServer.Exec("FLUSH PRIVILEGES;")
	}
	dbServer.Exec("FLUSH PRIVILEGES;")
	err = s.CreateDatabaseUser(dbServer, databaseInfo.Name, databaseInfo.Username, databaseInfo.Password, param.AccessPermission)
	if err != nil {
		return err
	}
	return nil
}

// GetAccessPermission 获取数据库访问权限
func (s *DatabaseMysqlService) GetAccessPermission(param *request.GetAccessPermissionR) (string, error) {
	databaseInfo, err := (&model.Databases{Username: param.UserName}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	dbServer, err := s.NewMysqlServiceBySid(databaseInfo.Sid)
	if err != nil {
		return "", err
	}
	_, err = dbServer.Raw("show databases").Rows()
	if err != nil {
		return "", err
	}
	users, err := dbServer.Raw("select Host from mysql.user where User=? AND Host!='localhost'", databaseInfo.Username).Rows()
	if err != nil {
		return "", err
	}
	for users.Next() {
		var us string
		_ = users.Scan(&us)
		return us, nil
	}
	return "", nil
}

// SetPwd 设置数据库密码
func (s *DatabaseMysqlService) SetPwd(param *request.SetPwdR) error {
	param.Password = util.ClearStr(param.Password)
	if util.StrIsEmpty(param.Password) {
		return errors.New("password is empty")
	}
	databaseInfo, err := (&model.Databases{Username: param.UserName}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	err = s.ChangeDatabasePassword(databaseInfo.Sid, databaseInfo.Name, databaseInfo.Username, param.Password)
	if err != nil {
		return err
	}

	//保存至面板数据库
	databaseInfo.Password = param.Password
	if err = databaseInfo.Update(global.PanelDB); err != nil {
		return err
	}
	return nil
}

// CheckDeleteDatabase 检查是否可以删除数据库
func (s *DatabaseMysqlService) CheckDeleteDatabase(param *request.CheckDeleteDatabaseR) ([]*response.CheckDeleteDatabaseP, error) {
	var data []*response.CheckDeleteDatabaseP
	for _, v := range param.IDs {
		get, err := (&model.Databases{ID: v}).Get(global.PanelDB)
		if err != nil {
			return nil, err
		}
		dbServer, err := s.NewMysqlServiceBySid(get.Sid)
		if err != nil {
			return nil, err
		}
		size, s2, err := s.GetDatabaseSize(dbServer)
		if err != nil {
			return nil, err
		}
		data = append(data, &response.CheckDeleteDatabaseP{
			Databases: get,
			Size:      size,
			SizeStr:   s2,
		})
	}
	return data, nil

}

// DeleteDatabase 删除数据库
func (s *DatabaseMysqlService) DeleteDatabase(param *request.DeleteDatabaseR) (string, error) {
	get, err := (&model.Databases{ID: param.ID}).Get(global.PanelDB)
	if err != nil {
		return "", err
	}
	dbServer, err := s.NewMysqlServiceBySid(get.Sid)
	if err != nil {
		return "", err
	}
	// 删除数据库
	result := dbServer.Exec("DROP DATABASE IF EXISTS `" + get.Name + "`")
	if result.Error != nil {
		return "", result.Error
	}
	err = s.DeleteDatabaseUser(dbServer, get.Username)
	if err != nil {
		return "", err
	}
	if err = (get).Delete(global.PanelDB, &model.ConditionsT{}); err != nil {
		return "", err
	}
	return get.Name, nil
}

// GetDatabaseSize 获取数据库大小
func (s *DatabaseMysqlService) GetDatabaseSize(dbServer *gorm.DB) (float64, string, error) {

	var result struct {
		Size float64
	}
	dbServer.Raw("SELECT ROUND(SUM(data_length + index_length), 2) AS size FROM information_schema.tables WHERE table_schema = ?", "database_name").Scan(&result)

	return result.Size, util.FormatSize(result.Size), nil
}

// NewMysqlServiceBySid 通过sid生成数据库实例
func (s *DatabaseMysqlService) NewMysqlServiceBySid(sid int64) (*gorm.DB, error) {
	var dsn string
	var err error
	if sid == 0 {
		//使用本地mysql数据库
		myConf, _ := util.ReadFileStringBody("/etc/my.cnf")

		portRe := regexp.MustCompile(`port\s*=\s*([0-9]+)`).FindStringSubmatch(myConf)
		var port string
		if len(portRe) < 2 {
			port = "3306"
		} else {
			port = portRe[1]
		}
		//查询Mysql root密码
		dsn = fmt.Sprintf("root:%s@tcp(localhost:%s)/?charset=utf8&parseTime=True&loc=Local", global.Config.System.MysqlRootPassword, port)
	} else { //使用远程mysql数据库
		//查询数据库服务器信息
		dbServer, err := (&model.DatabaseServer{ID: sid}).Get(global.PanelDB)
		if err != nil || dbServer == nil {
			return nil, err
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local", dbServer.User, dbServer.Password, dbServer.Host, dbServer.Port, dbServer.Charset)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,  // DSN data source name
		DefaultStringSize:        256,  // string 类型字段的默认长度
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	}), &gorm.Config{})
	if err != nil {
		global.Log.Errorf("NewMysqlServiceBySid->gorm.Open Error:%s", err.Error())
		return nil, err
	}

	return db, nil
}

// NewMysqlService 生成数据库实例
//func (s *DatabaseMysqlService) NewMysqlService(dbName string) (*gorm.DB, error) {
//	isCloudDb := false
//	if dbName != "" {
//		get, err := (&database.Databases{Name: dbName}).Get(global.PanelDB)
//		if err != nil {
//			return nil, err
//		}
//		isCloudDb = get.DbType == 1
//	}
//	var dbObj *gorm.DB
//	if isCloudDb {
//		dbObj = db_mysql.PanelMysql()
//		var connConfig struct {
//			DbHost     string
//			DbPort     int
//			DbName     string
//			DbUser     string
//			DbPassword string
//		}
//		if err := json.Unmarshal([]byte(dbFind.ConnCfg), &connConfig); err != nil {
//			return nil, err
//		}
//		if err := dbObj.SetHost(connConfig.DbHost, connConfig.DbPort, connConfig.DbName, connConfig.DbUser, connConfig.DbPassword); err != nil {
//			return nil, err
//		}
//	} else {
//		dbObj = panelMysql.PanelMysql()
//	}
//	return dbObj, nil
//}

// DatabaseExists 使用show databases查询数据库是否存在
func (s *DatabaseMysqlService) DatabaseExists(dbServer *gorm.DB, database string) (bool, error) {
	var databases []string
	result := dbServer.Raw("show databases").Scan(&databases)
	if result.Error != nil {
		return false, result.Error
	}

	// 输出查询结果
	for _, dbName := range databases {
		if dbName == database {
			return true, nil
		}
	}
	return false, nil
}

// CreateDatabaseUser 创建数据库用户
func (s *DatabaseMysqlService) CreateDatabaseUser(dbServer *gorm.DB, dbname, username, password, address string) error {
	dbServer.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY '%s'", username, password))
	result := dbServer.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost'", dbname, username))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "1044") {
			result = dbServer.Exec(fmt.Sprintf("GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,EVENT,TRIGGER ON %s.* TO '%s'@'localhost'", dbname, username))
			if result.Error != nil {
				return result.Error
			}
		} else {
			return result.Error
		}
	}
	for _, a := range strings.Split(address, ",") {
		result = dbServer.Exec(fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, a, password))
		if result.Error != nil {
			return result.Error
		}
		result = dbServer.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%s'", dbname, username, a))
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "1044") {
				result = dbServer.Exec(fmt.Sprintf("GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,EVENT,TRIGGER ON %s.* TO '%s'@'%s'", dbname, username, a))
				if result.Error != nil {
					return result.Error
				}
			} else {
				return result.Error
			}
		}
	}
	result = dbServer.Exec("FLUSH PRIVILEGES;")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteDatabaseUser 删除数据库用户
func (s *DatabaseMysqlService) DeleteDatabaseUser(dbServer *gorm.DB, username string) error {
	// 删除用户
	var users []struct {
		Host string
	}
	dbServer.Raw("SELECT Host FROM mysql.user WHERE User=? AND Host!='localhost'", username).Scan(&users)
	dbServer.Exec("DROP USER IF EXISTS '" + username + "'@'localhost'")
	for _, user := range users {
		dbServer.Exec("DROP USER IF EXISTS '" + username + "'@'" + user.Host + "'")
	}

	// 刷新权限
	dbServer.Exec("FLUSH PRIVILEGES")
	return nil
}

// ChangeLocalRootPassword 修改root密码
func (s *DatabaseMysqlService) ChangeLocalRootPassword(password string) error {
	//读取shell文件
	cmdTemp := `#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

pwd=${password}
service mysqld stop
mysqld_safe --skip-grant-tables&
sleep 3
m_version=$(cat /www/server/mysql/version.pl|grep -E "(5.1.|5.5.|5.6.)")
if [ "$m_version" != "" ];then
    mysql -uroot -e "insert into mysql.user(Select_priv,Insert_priv,Update_priv,Delete_priv,Create_priv,Drop_priv,Reload_priv,Shutdown_priv,Process_priv,File_priv,Grant_priv,References_priv,Index_priv,Alter_priv,Show_db_priv,Super_priv,Create_tmp_table_priv,Lock_tables_priv,Execute_priv,Repl_slave_priv,Repl_client_priv,Create_view_priv,Show_view_priv,Create_routine_priv,Alter_routine_priv,Create_user_priv,Event_priv,Trigger_priv,Create_tablespace_priv,User,Password,host)values('Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','root',password('${pwd}'),'127.0.0.1')"
    mysql -uroot -e "insert into mysql.user(Select_priv,Insert_priv,Update_priv,Delete_priv,Create_priv,Drop_priv,Reload_priv,Shutdown_priv,Process_priv,File_priv,Grant_priv,References_priv,Index_priv,Alter_priv,Show_db_priv,Super_priv,Create_tmp_table_priv,Lock_tables_priv,Execute_priv,Repl_slave_priv,Repl_client_priv,Create_view_priv,Show_view_priv,Create_routine_priv,Alter_routine_priv,Create_user_priv,Event_priv,Trigger_priv,Create_tablespace_priv,User,Password,host)values('Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','Y','root',password('${pwd}'),'localhost')"
    mysql -uroot -e "UPDATE mysql.user SET password=PASSWORD('${pwd}') WHERE user='root'";
else
    mysql -uroot -e "UPDATE mysql.user SET authentication_string='' WHERE user='root'";
    mysql -uroot -e "FLUSH PRIVILEGES";
    mysql -uroot -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '${pwd}';";
fi
mysql -uroot -e "FLUSH PRIVILEGES";
pkill -9 mysqld_safe
pkill -9 mysqld
sleep 2
service mysqld start

echo '==========================================='
echo "root密码成功修改为: ${pwd}"
echo "The root password set ${pwd}  successful"
exit 0
`
	//替换密码
	cmdTemp = strings.Replace(cmdTemp, "${password}", password, -1)
	err := util.WriteFile("/tmp/mysql_root.sh", []byte(cmdTemp), 755)
	defer func() {
		_ = os.Remove("/tmp/mysql_root.sh")
	}()
	if err != nil {
		return err
	}
	fmt.Println("正在修改root密码，请稍后...")
	fmt.Println("The set password...")
	err = util.ExecShellScriptS("bash /tmp/mysql_root.sh")
	if err != nil {
		global.Log.Errorf("修改root密码失败：%s", err.Error())
		return err
	}
	//开始写入配置文件
	global.Config.System.MysqlRootPassword = password
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err = global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}

	return nil
}

// ChangeRemoteRootPassword 修改远程root密码
func (s *DatabaseMysqlService) ChangeRemoteRootPassword(sid int64, password string) error {
	//查询数据库服务器信息
	get, err := (&model.DatabaseServer{ID: sid}).Get(global.PanelDB)
	if err != nil || get.ID == 0 {
		return err
	}
	dbService, err := s.NewMysqlServiceBySid(sid)

	// 获取MySQL版本
	var version string
	_ = dbService.Raw("select version()").Row().Scan(&version)

	// 获取管理员用户
	adminUser := get.User

	// 根据MySQL版本执行不同的操作
	if strings.HasPrefix(version, "5.7") || strings.HasPrefix(version, "8.0") {
		var hosts []string
		dbService.Raw("select Host from mysql.user where User = ?", adminUser).Pluck("Host", &hosts)
		for _, host := range hosts {
			dbService.Exec("UPDATE mysql.user SET authentication_string='' WHERE User = ? and Host = ?", adminUser, host)
			dbService.Exec("ALTER USER ?@? IDENTIFIED BY ?", adminUser, host, password)
		}
	} else if strings.HasPrefix(version, "10.5") || strings.HasPrefix(version, "10.4") {
		var hosts []string
		dbService.Raw("select Host from mysql.user where User = ?", adminUser).Pluck("Host", &hosts)
		for _, host := range hosts {
			dbService.Exec("ALTER USER ?@? IDENTIFIED BY ?", adminUser, host, password)
		}
	} else {
		dbService.Exec("update mysql.user set Password = password(?) where User = ?", password, adminUser)
	}

	dbService.Exec("flush privileges")
	return nil
}

// ChangeDatabasePassword 修改数据库密码
func (s *DatabaseMysqlService) ChangeDatabasePassword(sid int64, dbname string, username string, newPassword string) error {
	var version string
	dbServer, err := s.NewMysqlServiceBySid(sid)

	if err != nil {
		return err
	}
	rows := dbServer.Raw("select version();")
	rows.Scan(&version)
	if strings.HasPrefix(version, "5.7") || strings.HasPrefix(version, "8.0") {
		var accept []string
		dbServer.Raw("select Host from mysql.user where User=? AND Host!='localhost'", dbname).Scan(&accept)
		dbServer.Exec("update mysql.user set authentication_string='' where User=?", username)
		dbServer.Exec(fmt.Sprintf("ALTER USER `%s`@`localhost` IDENTIFIED BY '%s'", username, newPassword))
		for _, myHost := range accept {
			dbServer.Exec(fmt.Sprintf("ALTER USER `%s`@`%s` IDENTIFIED BY '%s'", username, myHost, newPassword))
		}

	} else if strings.Contains(version, "10.5.") || strings.Contains(version, "10.4.") {
		var accept []string
		dbServer.Raw("select Host from mysql.user where User=? AND Host!='localhost'", dbname).Scan(&accept)
		dbServer.Exec(fmt.Sprintf("ALTER USER `%s`@`localhost` IDENTIFIED BY '%s'", username, newPassword))
		for _, myHost := range accept {
			dbServer.Exec(fmt.Sprintf("ALTER USER `%s`@`%s` IDENTIFIED BY '%s'", username, myHost, newPassword))
		}
	} else {
		dbServer.Exec("update mysql.user set Password=password(?) where User=?", newPassword, username)
	}

	dbServer.Exec("flush privileges")
	return nil
}

// GetMysqldumpPath 获取mysqldump路径
func (s *DatabaseMysqlService) GetMysqldumpPath() string {
	binFiles := []string{
		fmt.Sprintf("%s/mysql/bin/mysqldump", global.Config.System.ServerPath),
		"/usr/bin/mysqldump",
		"/usr/local/bin/mysqldump",
		"/usr/sbin/mysqldump",
		"/usr/local/sbin/mysqldump",
	}

	for _, binFile := range binFiles {
		if _, err := os.Stat(binFile); err == nil {
			return binFile
		}
	}

	for _, binFile := range binFiles {
		if _, err := os.Stat(binFile); err == nil {
			return binFile
		}
	}

	return ""
}

func (s *DatabaseMysqlService) GetDatabase(databaseID int64) (databasesInfo *model.Databases, err error) {
	//判断数据库是否存在
	databasesInfo, err = (&model.Databases{ID: databaseID}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if databasesInfo.ID == 0 {
		return nil, errors.New("数据库不存在")
	}
	return
}

// ImportDatabase 导入数据库
func (s *DatabaseMysqlService) ImportDatabase(databasesGet *model.Databases, sqlFilePath string, backupID int64) (logPath string, err error) {
	importFilePath := ""
	deCompressPath := ""
	logPath = fmt.Sprintf("%s/mysql/import_%s_%d.log", global.Config.Logger.RootPath, databasesGet.Name, time.Now().Unix())
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0755)
	if err != nil {
		return "", err
	}
	global.Log.Debugf("sqlFilePath: %s", sqlFilePath)
	go func() {
		defer func(logFile *os.File) {
			_ = logFile.Close()
		}(logFile)
		_, _ = logFile.WriteString(fmt.Sprintf("----------------------------------------\nStart import Databases\n"))
		//判断文件导入还是备份还原
		if !util.StrIsEmpty(sqlFilePath) {
			global.Log.Debugf("sqlFilePath: %s", sqlFilePath)
			if !util.IsFile(sqlFilePath) {
				_, _ = logFile.WriteString(fmt.Sprintf("File %s is not found\n", sqlFilePath))
				return
			}
			//判断文件是否是sql文件
			if strings.HasSuffix(sqlFilePath, ".sql") {
				importFilePath = sqlFilePath
			} else {
				deCompressPath = sqlFilePath
			}
		} else if backupID > 0 {
			backupGet, err := (&model.Backup{ID: backupID}).Get(global.PanelDB)
			if err != nil {
				_, _ = logFile.WriteString(fmt.Sprintf("backup.Get Error:%s\n", err.Error()))
				return
			}
			if backupGet.ID == 0 {
				_, _ = logFile.WriteString(fmt.Sprintf("backupID %d is not found\n", backupID))
				return
			}
			if backupGet.StorageId == 0 { //本地备份
				deCompressPath = backupGet.FilePath
			} else { //远程备份
				storageGet, err := (&model.Storage{ID: backupGet.StorageId}).Get(global.PanelDB)
				if err != nil {
					_, _ = logFile.WriteString(fmt.Sprintf("storage.Get Error:%s\n", err.Error()))
					return
				}
				if storageGet.ID == 0 {
					_, _ = logFile.WriteString(fmt.Sprintf("storageID %d is not found\n", backupGet.StorageId))
					return
				}
				//下载文件
				var config request.StorageConfigR
				err = util.JsonStrToStruct(storageGet.Config, &config)
				if err != nil {
					_, _ = logFile.WriteString(fmt.Sprintf("JsonStrToStruct Error:%s\n", err.Error()))
					return
				}
				deCompressPath = fmt.Sprintf("/tmp/%d/%s", time.Now().Unix(), backupGet.FileName)
				_ = os.MkdirAll(filepath.Dir(deCompressPath), 0777)
				storageCore, err := GroupApp.StorageServiceApp.NewStorageCore(storageGet.Category, &config)
				_, err = storageCore.Download(backupGet.FilePath, deCompressPath)
				if err != nil {
					_, _ = logFile.WriteString(fmt.Sprintf("storageCore.Download Error:%s\n", err.Error()))
					return
				}
				defer func(name string) {
					_ = os.Remove(name)

				}(deCompressPath)
			}
		} else {
			_, _ = logFile.WriteString(fmt.Sprintf("param is Error:%s\n", err.Error()))
			return
		}

		if !util.StrIsEmpty(deCompressPath) {
			_, _ = logFile.WriteString(fmt.Sprintf("Decompress:%s\n", sqlFilePath))
			//解压缩文件
			tempDir := fmt.Sprintf("/tmp/%d", time.Now().Unix())
			_ = os.MkdirAll(tempDir, 0777)
			_, err = (&ExplorerService{}).Decompress(true, sqlFilePath, tempDir, "")
			if err != nil {
				_, _ = logFile.WriteString(fmt.Sprintf("Decompress Error:%s\n", err.Error()))
				return
			}
			files, err := os.ReadDir(tempDir)
			if err != nil {
				_, _ = logFile.WriteString(fmt.Sprintf("os.ReadDir:%s\n", err.Error()))
				return
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				} else {
					if strings.HasSuffix(file.Name(), ".sql") {
						importFilePath = tempDir + "/" + file.Name()
						_, _ = logFile.WriteString(fmt.Sprintf("Decompressed File:%s\n", importFilePath))
						break
					}
				}
			}
		}

		if util.StrIsEmpty(importFilePath) || !util.IsFile(importFilePath) {
			global.Log.Error(fmt.Sprintf("导入文件不存在：%s", importFilePath))
			_, _ = logFile.WriteString(fmt.Sprintf("--Error:file is not found ,importFilePath:%s,\n", importFilePath))
			return
		}
		fileInfo, err := os.Stat(importFilePath)
		if err != nil {
			_, _ = logFile.WriteString(fmt.Sprintf("os.Stat Error:%s\n", err.Error()))
			return
		}
		if fileInfo.Size() < 100 {
			global.Log.Error(fmt.Sprintf("导入文件太小：%s", importFilePath))
			_, _ = logFile.WriteString(fmt.Sprintf("--Error:fileSize too small ,size:%d,\n", fileInfo.Size()))
			return
		}

		cmdStr := "mysql "

		if databasesGet.Sid == 0 { //本地数据库
			rootPwd, _ := s.GetRootPwd()
			rootPwd = util.ShellQuote(rootPwd)
			cmdStr += fmt.Sprintf("-uroot -p%s  --force \"%s\" < \"%s\"", rootPwd, databasesGet.Name, importFilePath)
		} else { //远程数据库
			dbServer, err := (&model.DatabaseServer{ID: databasesGet.Sid}).Get(global.PanelDB)
			if err != nil || dbServer == nil {
				_, _ = logFile.WriteString(fmt.Sprintf("DatabaseServer.Get Error:%s\n", err.Error()))
				return
			}
			cmdStr += fmt.Sprintf("-h %s -P %d -u%s -p%s  --force \"%s\" < \"%s\"", dbServer.Host, dbServer.Port, dbServer.User, dbServer.Password, databasesGet.Name, importFilePath)
		}

		shellResult, err := util.ExecShell(cmdStr)
		if err != nil {
			global.Log.Error(fmt.Sprintf("导入数据库失败：%s ,result:%s", err.Error(), shellResult))
			_, _ = logFile.WriteString(fmt.Sprintf("---Error:%s,result:%s\n", err.Error(), shellResult))
			return
		}
		_, _ = logFile.WriteString(fmt.Sprintf("Result:%s\n----------------------Success---------------------\n", shellResult))
		global.Log.Infof("import database %s successfully!", databasesGet.Name)
		return
	}()
	return logPath, nil
}
