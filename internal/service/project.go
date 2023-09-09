package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	request2 "TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ProjectService struct {
}

func (s *ProjectService) GetProjectInfoByID(projectID int64) (projectGet *model.Project, err error) {
	projectGet, err = (&model.Project{ID: projectID}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if projectGet.ID == 0 {
		err = errors.New("project does not exist")
		return
	}
	//vhostPath := extensions.GetExtensionsPath(extensions.NginxName) + "/vhost"
	//s.ProjectData = projectGet
	//s.ProjectMainConfFilePath = vhostPath + "/main/" + projectGet.Name + ".conf"
	//s.ProjectConfDirPath = vhostPath + "/project/" + projectGet.Name
	return
}

func (s *ProjectService) AddDomains(projectData *model.Project, param *request2.AddDomainsR) (errs []error) {
	//检查nginx配置
	nginxService := ExtensionNginxService{}
	err := nginxService.CheckConfig()
	if err != nil {
		errs = append(errs, err)
		return
	}

	domainList := make(map[string]bool)
	portList := make(map[string]bool)
	var firewallRList []*request2.CreatePortRuleR
	for _, domainItem := range param.Domains {
		//检查域名和端口
		err = CheckDomainItem(domainItem)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		//添加域名至数据库
		_, err := (&model.ProjectDomain{Domain: domainItem.Name, Port: domainItem.Port, ProjectId: param.ProjectID}).Create(global.PanelDB)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		domainList[domainItem.Name] = true
		portList[strconv.Itoa(domainItem.Port)] = true
	}
	//开放系统防火墙端口
	for k, _ := range portList {
		port, _ := strconv.Atoi(k)
		firewallRList = append(firewallRList, &request2.CreatePortRuleR{
			Port:     port,
			Strategy: constant.SystemFirewallStrategyAllow,
			Protocol: "tcp",
			Ps:       "项目域名端口",
		})
	}
	firewall, _ := GroupApp.SystemFirewallServiceApp.New()
	_ = firewall.BatchCreatePortRule(firewallRList)
	//添加域名至nginx配置文件
	err = AddDomainToNginxConfig(projectData.Name, domainList, portList)
	if err != nil {
		errs = append(errs, err)
		return
	}
	//更新ttwaf域名配置文件
	err = GroupApp.TTWafServiceApp.OperateDomainConfig(projectData.ID, projectData.Name, projectData.Path, constant.OperateTTWafDomainConfigByUpdate)
	if err != nil {
		errs = append(errs, err)
		return
	}
	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return errs
}

// DomainList 域名列表
func (s *ProjectService) DomainList(projectID int64) (domainList []*model.ProjectDomain, total int64, err error) {
	return (&model.ProjectDomain{}).List(global.PanelDB, &model.ConditionsT{"project_id": projectID, "ORDER": "create_time DESC"}, 0, 0)
}

// BatchDeleteDomain 删除域名
func (s *ProjectService) BatchDeleteDomain(projectData *model.Project, ids []int64) (errs []error) {
	var delDomains []string
	var delPorts []string
	for _, id := range ids {
		//查询需要删除的域名信息
		delDomainInfo, err := (&model.ProjectDomain{ID: id}).Get(global.PanelDB)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		delDomain := delDomainInfo.Domain
		delPort := strconv.Itoa(delDomainInfo.Port)
		if (&model.ProjectDomain{}).Count(global.PanelDB, &model.ConditionsT{"project_id": projectData.ID, "port": delDomainInfo.Port}) > 1 {
			delPort = ""
		}

		err = delDomainInfo.Delete(global.PanelDB, &model.ConditionsT{})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		delDomains = append(delDomains, delDomain)
		delPorts = append(delPorts, delPort)
		global.Log.Warnf("BatchDeleteDomain->domain:%s,port:%s", delDomain, delPort)
	}
	//从配置文件删除域名
	err := DelDomainToNginxConfig(projectData.Name, delDomains, delPorts)
	if err != nil {
		errs = append(errs, err)
		return
	}
	//更新ttwaf域名配置文件
	err = GroupApp.TTWafServiceApp.OperateDomainConfig(projectData.ID, projectData.Name, projectData.Path, constant.OperateTTWafDomainConfigByUpdate)
	if err != nil {
		errs = append(errs, err)
		return
	}

	//重启nginx
	defer func() {
		err = ReloadNginx()
		if err != nil {
			global.Log.Error(err)
			return
		}
	}()
	return
}

// RewriteTemplateList 伪静态模板列表
func (s *ProjectService) RewriteTemplateList() (map[string]interface{}, error) {
	templateList := make([]string, 0)
	path := GetExtensionsPath(constant.ExtensionNginxName) + "/template/rewrite"
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			// 去掉后缀
			templateList = append(templateList, strings.TrimSuffix(file.Name(), ".conf"))
		}
	}
	var data = make(map[string]interface{})
	data["list"] = templateList
	data["path"] = path
	return data, nil
}

// DefaultIndex 默认首页
func (s *ProjectService) DefaultIndex(projectName string) (string, error) {
	confPath := ProjectMainConfFilePath(projectName)
	confBody, err := util.ReadFileStringBody(confPath)
	if err != nil {
		return "", err
	}
	rep := `\s+index\s+(.+);`
	re := regexp.MustCompile(rep)
	if re.MatchString(confBody) {
		tmp := re.FindStringSubmatch(confBody)[1]
		return strings.ReplaceAll(tmp, " ", ","), nil
	}
	return "", nil
}

// SaveDefaultIndex 保存默认首页
func (s *ProjectService) SaveDefaultIndex(projectName string, index string) error {
	index = strings.ReplaceAll(index, "  ", "")
	index = strings.ReplaceAll(index, ",,", ",")
	indexItem := strings.Split(index, ",")
	for _, v := range indexItem {
		if !strings.Contains(v, ".") {
			return errors.New(fmt.Sprintf("Default index %v is illegal", v))
		}
	}
	index = strings.Join(indexItem, " ")

	confPath := ProjectMainConfFilePath(projectName)
	confBody, err := util.ReadFileStringBody(confPath)
	if err != nil {
		return err
	}
	rep := `\s+index\s+.+;`
	re := regexp.MustCompile(rep)
	confBody = re.ReplaceAllString(confBody, "\n\tindex "+index+";")

	err = util.WriteFile(confPath, []byte(confBody), 0644)
	if err != nil {
		return err
	}

	err = (&ExtensionNginxService{}).SetStatus("reload")
	if err != nil {
		return err
	}
	return nil
}

// CategoryList 分类列表
func (s *ProjectService) CategoryList() (list []*model.ProjectCategory, total int64, err error) {
	list, total, err = (&model.ProjectCategory{}).List(global.PanelDB, &model.ConditionsT{}, 0, 0)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// CreateCategory 创建分类
func (s *ProjectService) CreateCategory(param *request2.CreateCategoryR) error {
	//检查分类是否存在
	Get, err := (&model.ProjectCategory{Name: param.Name}).GetByName(global.PanelDB)
	if err != nil {
		return err
	}
	if Get.ID > 0 {
		return errors.New(helper.MessageWithMap("project.CategoryExists", map[string]any{"Name": param.Name}))
	}
	_, err = (&model.ProjectCategory{Name: param.Name, Ps: param.Ps}).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// EditCategory 编辑分类
func (s *ProjectService) EditCategory(param *request2.EditCategoryR) error {
	//检查分类是否存在
	Get, err := (&model.ProjectCategory{Name: param.Name}).GetByName(global.PanelDB)
	if err != nil {
		return err
	}
	if Get.ID > 0 && Get.ID != param.ID {
		return errors.New(helper.MessageWithMap("project.CategoryExists", map[string]any{"Name": param.Name}))
	}
	err = (&model.ProjectCategory{ID: param.ID, Name: param.Name, Ps: param.Ps}).Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// SetSSL 设置SSL
func (s *ProjectService) SetSSL(projectName string, param *request2.SetSslR) (err error) {
	var private, fullchain string
	if !util.StrIsEmpty(param.SSLKey) {
		sslPath := fmt.Sprintf("%s/data/ssl/%s", global.Config.System.PanelPath, param.SSLKey)
		if !util.PathExists(sslPath) {
			return errors.New(helper.MessageWithMap("project.SSLIsNotExists", map[string]any{"Name": param.SSLKey}))
		}
		private, err = util.ReadFileStringBody(sslPath + "/fullchain.pem")
		if err != nil {
			return err
		}
		fullchain, err = util.ReadFileStringBody(sslPath + "/private.pem")
		if err != nil {
			return err
		}

	} else if !util.StrIsEmpty(param.Csr) && !util.StrIsEmpty(param.Key) {
		private = param.Key
		fullchain = param.Csr
		cert, err := util.ParseCert([]byte(fullchain))
		if err != nil {
			return err
		}
		sslDetails := response.SSLDetails{
			AutoRenew:    false,
			DNSAccount:   "",
			AcmeAccount:  "",
			Domains:      cert.DNSNames,
			CertURL:      "",
			ExpireDate:   cert.NotAfter,
			StartDate:    cert.NotBefore,
			Type:         cert.Issuer.CommonName,
			Organization: cert.Issuer.Organization[0],
		}
		//保存证书和证书信息
		savePath := fmt.Sprintf("%s/data/ssl/%s", global.Config.System.PanelPath, cert.Subject.CommonName)
		_ = os.MkdirAll(savePath, 0755)
		err = util.WriteFile(savePath+"/fullchain.pem", []byte(fullchain), 0644)
		if err != nil {
			return errors.New("ERROR：" + err.Error())
		}
		err = util.WriteFile(savePath+"/private.pem", []byte(private), 0644)
		if err != nil {
			return errors.New("ERROR：" + err.Error())
		}
		sslDetailsStr, err := util.StructToJsonStr(sslDetails)
		if err != nil {
			return errors.New("ERROR：" + err.Error())
		}
		err = util.WriteFile(savePath+"/info.json", []byte(sslDetailsStr), 0644)
		if err != nil {
			return errors.New("ERROR：" + err.Error())
		}
	} else {
		return errors.New("empty ssl_key or csr and key")
	}
	if err = GenerateSslConfig(projectName, []byte(private), []byte(fullchain)); err != nil {
		return err
	}
	return nil
}

// CloseSSL 关闭SSL
func (s *ProjectService) CloseSSL(projectName string) (err error) {
	sslConfPath := fmt.Sprintf("%s/ssl.conf", ProjectConfDirPath(projectName))
	certPath := fmt.Sprintf("%s/cert", ProjectConfDirPath(projectName))
	_ = util.WriteFile(sslConfPath, []byte(""), 0644)
	_ = os.Remove(certPath + "/private.pem")
	_ = os.Remove(certPath + "/fullchain.pem")
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return
}

// AlwaysUseHttps 设置强制HTTPS
func (s *ProjectService) AlwaysUseHttps(projectName string, action bool) (err error) {
	var newConfig string
	var isAlwaysUseHttps bool
	//判断是否开启了https
	fileBody, err := util.ReadFileStringBody(ProjectConfDirPath(projectName) + "/ssl.conf")
	if err != nil {
		return err
	}
	if util.StrIsEmpty(fileBody) {
		return errors.New("please set ssl first")
	}
	start := strings.Index(fileBody, "#AlwaysUseHttps_Start")
	end := strings.Index(fileBody, "#AlwaysUseHttps_End")

	if start == -1 && end == -1 {
		isAlwaysUseHttps = false
	} else if start != -1 && end != -1 {
		isAlwaysUseHttps = true
	} else {
		return errors.New("config error")
	}
	if action && !isAlwaysUseHttps {
		tmp := `#AlwaysUseHttps_Start
				if ($server_port !~ 443){
					rewrite ^(/.*)$ https://$host$1 permanent;
				}
				#AlwaysUseHttps_End` + "\n"
		newConfig = tmp + fileBody
	}
	if !action && isAlwaysUseHttps {
		newConfig = fileBody[:start] + fileBody[end+len("#AlwaysUseHttps_End"):]
	}

	//格式化配置文件
	newConfig, err = FormatNginxConf(newConfig)
	if err != nil {
		return
	}
	if err = util.WriteFile(ProjectConfDirPath(projectName)+"/ssl.conf", []byte(newConfig), 0644); err != nil {
		return
	}
	return nil
}

// CreateRedirect 创建重定向
func (s *ProjectService) CreateRedirect(projectName string, param *request2.CreateRedirectR) (err error) {

	//检查参数
	if param.Type == constant.ProjectRedirectTypeDomain {
		if len(param.Domains) == 0 {
			err = errors.New("empty domains")
			return
		}
	}
	if param.Type == constant.ProjectRedirectTypePath {
		if util.StrIsEmpty(param.Path) {
			err = errors.New("empty path")
			return
		}
	}

	//生成配置文件
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	configID := fmt.Sprintf("%d", time.Now().UnixNano())
	if param.Status {
		err = GenerateRedirectConfig(projectName, configID, param)
		if err != nil {
			return err
		}
	}

	listFilePath := ProjectConfDirPath(projectName) + "/redirect/list.json"
	if !util.PathExists(listFilePath) {
		err = util.WriteFile(listFilePath, []byte("{}"), 0644)
		if err != nil {
			return
		}
	}
	listBody, err := util.ReadFileStringBody(listFilePath)
	if err != nil {
		return
	}

	var list map[string]request2.CreateRedirectR
	err = util.JsonStrToStruct(listBody, &list)
	if err != nil {
		return
	}

	list[configID] = *param
	listBody, err = util.StructToJsonStr(list)
	if err != nil {
		return
	}
	err = util.WriteFile(listFilePath, []byte(listBody), 0644)
	if err != nil {
		return
	}

	return nil
}

// RedirectList 重定向列表
func (s *ProjectService) RedirectList(projectName string) (list map[string]request2.CreateRedirectR, err error) {
	filePath := ProjectConfDirPath(projectName) + "/redirect/list.json"
	err = GetJsonFileData(filePath, &list)
	if err != nil {
		return nil, err
	}
	return
}

// BatchEditRedirect 批量编辑重定向
func (s *ProjectService) BatchEditRedirect(projectName string, newList map[string]request2.CreateRedirectR) (err error) {
	oldList, err := s.RedirectList(projectName)
	if err != nil {
		return
	}

	for key, config := range newList {
		//如果旧列表中存在则进行操作
		if _, ok := oldList[key]; ok {
			//如果存在该conf文件则删除
			confPath := ProjectConfDirPath(projectName) + "/redirect/" + key + ".conf"
			if util.PathExists(confPath) {
				_ = os.Remove(confPath)
			}
			//根据状态生成配置文件
			if config.Status {
				err = GenerateRedirectConfig(projectName, key, &config)
				if err != nil {
					return err
				}
			}
			//更新旧列表中的数据
			oldList[key] = config
		}
	}

	//更新列表文件
	saveBody, err := util.StructToJsonStr(oldList)
	if err != nil {
		return
	}
	listFilePath := ProjectConfDirPath(projectName) + "/redirect/list.json"
	err = util.WriteFile(listFilePath, []byte(saveBody), 0644)
	if err != nil {
		return
	}
	return
}

// BatchDeleteRedirect 批量删除重定向
func (s *ProjectService) BatchDeleteRedirect(projectName string, keys []string) (err error) {
	list, err := s.RedirectList(projectName)
	if err != nil {
		return err
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	for _, key := range keys {
		//如果存在该conf文件则删除
		confPath := ProjectConfDirPath(projectName) + "/redirect/" + key + ".conf"
		if util.PathExists(confPath) {
			_ = os.Remove(confPath)
		}
		delete(list, key)
	}
	//更新列表文件
	saveBody, err := util.StructToJsonStr(list)
	if err != nil {
		return
	}
	listFilePath := ProjectConfDirPath(projectName) + "/redirect/list.json"
	err = util.WriteFile(listFilePath, []byte(saveBody), 0644)
	if err != nil {
		return
	}
	return
}

// GetAntiLeechConfig 获取防盗链配置
func (s *ProjectService) GetAntiLeechConfig(projectData *model.Project) (*request2.CreateAntiLeechConfigR, error) {

	antiLeechConfPath := ProjectConfDirPath(projectData.Name) + "/anti_leech.conf"
	if !util.PathExists(antiLeechConfPath) {
		return nil, errors.New("not found file:" + antiLeechConfPath)
	}
	antiLeechConfBody, err := util.ReadFileStringBody(antiLeechConfPath)
	if err != nil {
		return nil, err
	}
	if util.StrIsEmpty(antiLeechConfBody) {
		return nil, nil
	}
	config := &request2.CreateAntiLeechConfigR{}
	config.ProjectID = projectData.ID
	//提取后缀
	suffixItem := regexp.MustCompile(`\.\((.+)\)\$`).FindStringSubmatch(antiLeechConfBody)
	if len(suffixItem) > 1 {
		config.SuffixList = strings.Split(suffixItem[1], "|")
	}
	//提取域名和是否允许空请求
	ReferersItem := regexp.MustCompile(`valid_referers\s+(none\s+blocked)?\s+(.+);`).FindStringSubmatch(antiLeechConfBody)
	if len(ReferersItem) > 2 {
		config.RefererNone = !(ReferersItem[1] == "")
		config.PassDomains = strings.Split(ReferersItem[2], " ")
	}
	//提取响应状态码
	statusItem := regexp.MustCompile(`return\s+(.+);`).FindStringSubmatch(antiLeechConfBody)
	if len(statusItem) > 1 {
		config.ResponseStatusCode, _ = strconv.Atoi(statusItem[1])
	}
	return config, nil
}

// CreateAntiLeechConfig 创建防盗链配置
func (s *ProjectService) CreateAntiLeechConfig(projectName string, param *request2.CreateAntiLeechConfigR) (err error) {
	err = GenerateAntiLeechConfig(projectName, param)
	if err != nil {
		return
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	return
}

// CloseAntiLeechConfig 关闭防盗链配置
func (s *ProjectService) CloseAntiLeechConfig(projectName string) (projectGet *model.Project, err error) {

	antiLeechConfPath := ProjectConfDirPath(projectName) + "/anti_leech.conf"
	if util.PathExists(antiLeechConfPath) {
		_ = util.WriteFile(antiLeechConfPath, []byte(""), 0644)
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	return
}

// CreateReverseProxyConfig 创建反向代理配置
func (s *ProjectService) CreateReverseProxyConfig(projectName string, param *request2.CreateReverseProxyConfigR) (err error) {
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}

	configID := fmt.Sprintf("%d", time.Now().UnixNano())
	if param.Status {
		err = GenerateReverseProxyConfig(projectName, configID, param)
		if err != nil {
			return
		}
	}

	listFilePath := ProjectConfDirPath(projectName) + "/proxy/list.json"
	if !util.PathExists(listFilePath) {
		err = util.WriteFile(listFilePath, []byte("{}"), 0644)
		if err != nil {
			return
		}
	}
	listBody, err := util.ReadFileStringBody(listFilePath)
	if err != nil {
		return
	}

	var list map[string]request2.CreateReverseProxyConfigR
	err = util.JsonStrToStruct(listBody, &list)
	if err != nil {
		return
	}

	//检查代理目录是否重复
	for _, v := range list {
		if v.ProxyDir == param.ProxyDir {
			return errors.New("proxy path is exist")
		}
	}

	list[configID] = *param
	listBody, err = util.StructToJsonStr(list)
	if err != nil {
		return
	}
	err = util.WriteFile(listFilePath, []byte(listBody), 0644)
	if err != nil {
		return
	}

	return
}

// ReverseProxyList 反向代理列表
func (s *ProjectService) ReverseProxyList(projectName string) (list map[string]request2.CreateReverseProxyConfigR, err error) {
	listFilePath := ProjectConfDirPath(projectName) + "/proxy/list.json"
	err = GetJsonFileData(listFilePath, &list)
	if err != nil {
		return nil, err
	}
	return
}

// BatchEditReverseProxyConfig 批量编辑反向代理配置
func (s *ProjectService) BatchEditReverseProxyConfig(projectName string, newList map[string]request2.CreateReverseProxyConfigR) (err error) {

	var oldList map[string]request2.CreateReverseProxyConfigR
	oldList, err = s.ReverseProxyList(projectName)
	if err != nil {
		return
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
			}
		}()
	}
	for key, config := range newList {
		//如果旧列表中存在则进行操作
		if _, ok := oldList[key]; ok {
			//如果存在该conf文件则删除
			confPath := ProjectConfDirPath(projectName) + "/proxy/" + key + ".conf"
			if util.PathExists(confPath) {
				_ = os.Remove(confPath)
			}
			//根据状态生成配置文件
			if config.Status {
				err = GenerateReverseProxyConfig(projectName, key, &config)
				if err != nil {
					return
				}
			}
			//更新旧列表中的数据
			oldList[key] = config
		}
	}
	oldListBody, err := util.StructToJsonStr(oldList)
	if err != nil {
		return
	}
	err = util.WriteFile(ProjectConfDirPath(projectName)+"/proxy/list.json", []byte(oldListBody), 0644)
	if err != nil {
		return
	}
	defer func() {
		err = ReloadNginx()
		if err != nil {
			global.Log.Error(err)
			return
		}
	}()
	return

}

// BatchDeleteReverseProxyConfig 批量删除反向代理配置
func (s *ProjectService) BatchDeleteReverseProxyConfig(projectName string, keys []string) (err error) {
	list, err := s.ReverseProxyList(projectName)
	if err != nil {
		return
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	for _, key := range keys {
		//如果存在该conf文件则删除
		confPath := ProjectConfDirPath(projectName) + "/proxy/" + key + ".conf"
		if util.PathExists(confPath) {
			_ = os.Remove(confPath)
		}
		delete(list, key)
	}
	//更新列表文件
	saveBody, err := util.StructToJsonStr(list)
	if err != nil {
		return
	}
	listFilePath := ProjectConfDirPath(projectName) + "/proxy/list.json"
	err = util.WriteFile(listFilePath, []byte(saveBody), 0644)
	if err != nil {
		return
	}
	return
}

// CreateAccessRuleConfig 创建访问规则配置
func (s *ProjectService) CreateAccessRuleConfig(projectName string, param *request2.CreateAccessRuleConfigR) (err error) {

	configID := fmt.Sprintf("%d", time.Now().UnixNano())
	err = GenerateAccessRuleConfig(projectName, configID, param)
	if err != nil {
		return
	}

	listFilePath := ProjectConfDirPath(projectName) + "/access_rule/list.json"
	if !util.PathExists(listFilePath) {
		err = util.WriteFile(listFilePath, []byte("{}"), 0644)
		if err != nil {
			return
		}
	}
	listBody, err := util.ReadFileStringBody(listFilePath)
	if err != nil {
		return
	}

	var list map[string]request2.CreateAccessRuleConfigR
	err = util.JsonStrToStruct(listBody, &list)
	if err != nil {
		return
	}

	list[configID] = *param
	listBody, err = util.StructToJsonStr(list)
	if err != nil {
		return
	}
	err = util.WriteFile(listFilePath, []byte(listBody), 0644)
	if err != nil {
		return
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	return
}

// AccessRuleConfigList 访问规则列表
func (s *ProjectService) AccessRuleConfigList(projectName string) (list map[string]request2.CreateAccessRuleConfigR, err error) {
	listFilePath := ProjectConfDirPath(projectName) + "/access_rule/list.json"
	err = GetJsonFileData(listFilePath, &list)
	if err != nil {
		return nil, err
	}
	return
}

// EditAccessRuleConfig 编辑访问规则配置
func (s *ProjectService) EditAccessRuleConfig(projectName string, configID string, param request2.CreateAccessRuleConfigR) (err error) {
	configList, err := s.AccessRuleConfigList(projectName)
	if err != nil {
		return err
	}
	if _, ok := configList[configID]; !ok {
		return errors.New("not found AccessRule")
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}

	//删除对应的配置文件
	_ = os.Remove(ProjectConfDirPath(projectName) + "/access_rule/" + configID + ".conf")
	if configList[configID].RuleType == constant.ProjectAccessRuleTypeBasicAuth {
		_ = os.Remove(ProjectConfDirPath(projectName) + "/access_rule/" + configID + ".pass")
	}

	//更新配置文件
	newConf := request2.CreateAccessRuleConfigR{
		ProjectId:       configList[configID].ProjectId,
		Name:            param.Name,
		Dir:             param.Dir,
		RuleType:        configList[configID].RuleType,
		BasicAuthConfig: param.BasicAuthConfig,
		NoAccessConfig:  param.NoAccessConfig,
	}
	configList[configID] = newConf

	err = GenerateAccessRuleConfig(projectName, configID, &newConf)
	if err != nil {
		return
	}

	//更新列表文件
	configListBody, err := util.StructToJsonStr(configList)
	if err != nil {
		return
	}
	err = util.WriteFile(ProjectConfDirPath(projectName)+"/access_rule/list.json", []byte(configListBody), 0644)
	if err != nil {
		return
	}
	return
}

// BatchDeleteAccessRuleConfig 批量删除访问规则配置
func (s *ProjectService) BatchDeleteAccessRuleConfig(projectName string, keys []string) (err error) {
	list, err := s.AccessRuleConfigList(projectName)
	if err != nil {
		return
	}
	if err == nil {
		defer func() {
			err = ReloadNginx()
			if err != nil {
				global.Log.Error(err)
				return
			}
		}()
	}
	for _, key := range keys {
		//如果存在该conf文件则删除
		confPath := ProjectConfDirPath(projectName) + "/access_rule/" + key + ".conf"
		if util.PathExists(confPath) {
			_ = os.Remove(confPath)
			if list[key].RuleType == constant.ProjectAccessRuleTypeBasicAuth {
				_ = os.Remove(ProjectConfDirPath(projectName) + "/access_rule/" + key + ".pass")
			}
		}
		delete(list, key)
	}
	//更新列表文件
	saveBody, err := util.StructToJsonStr(list)
	if err != nil {
		return
	}
	listFilePath := ProjectConfDirPath(projectName) + "/access_rule/list.json"
	err = util.WriteFile(listFilePath, []byte(saveBody), 0644)
	if err != nil {
		return
	}
	return
}
