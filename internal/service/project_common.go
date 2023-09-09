package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	projectR "TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type NginxConfig struct {
	ProjectName string
	Domain      []*projectR.DomainItem
	ProjectPath string
	SSL         bool
	PHPVersion  string
}

// GetJsonFileData 传入文件路径和结构体，读取文件内容并转换为结构体，如果文件不存在，则创建一个内容为”{}“的json文件
func GetJsonFileData(path string, data any) error {
	if !util.PathExists(path) {
		err := util.WriteFile(path, []byte("{}"), 0644)
		if err != nil {
			return err
		}
	}
	listBody, err := util.ReadFileStringBody(path)
	if err != nil {
		return err
	}

	err = util.JsonStrToStruct(listBody, data)
	if err != nil {
		return err
	}
	return nil
}

func CheckConfigPath() {
	nginxPath := GetExtensionsPath(constant.ExtensionNginxName)
	mainPath := fmt.Sprintf("%s/vhost/main", nginxPath)
	projectPath := fmt.Sprintf("%s/vhost/project", nginxPath)
	if !util.PathExists(mainPath) {
		_ = os.MkdirAll(mainPath, 0755)
	}
	if !util.PathExists(projectPath) {
		_ = os.MkdirAll(projectPath, 0755)
	}
}

// ReloadNginx 重载nginx配置
func ReloadNginx() (err error) {
	err = GroupApp.ExtensionNginxServiceApp.SetStatus("reload")
	if err != nil {
		return
	}
	return
}

func CheckDomainItem(domainItem *projectR.DomainItem) (err error) {
	if util.StrIsEmpty(domainItem.Name) {
		return errors.New("domain is empty")
	}
	//检查域名是否合法
	if !util.CheckDomain(domainItem.Name) {
		return errors.New(helper.MessageWithMap("project.DomainNameIsIllegal", map[string]any{"Domain": domainItem.Name}))
	}
	//转换域名为punycode
	domainItem.Name = util.ToPunycode(domainItem.Name)
	//检查端口是否合法
	if !util.CheckProjectPort(strconv.Itoa(domainItem.Port)) {
		return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": domainItem.Port}))
	}
	//检查域名是否存在
	_, err = DomainExists(domainItem.Name)
	if err != nil {
		return err
	}
	return nil
}

func InsertDomain(domain *projectR.DomainItem, pid int64) error {
	//开放系统防火墙端口
	firewall, _ := GroupApp.SystemFirewallServiceApp.New()
	_ = firewall.BatchCreatePortRule([]*projectR.CreatePortRuleR{{Port: domain.Port, Strategy: constant.SystemFirewallStrategyAllow, Protocol: "tcp", Ps: "项目域名端口"}})

	//添加至数据库
	_, err := (&model.ProjectDomain{Domain: domain.Name, Port: domain.Port, ProjectId: pid}).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// DomainExists 检查项目域名是否已经存在(添加域名检查)
func DomainExists(domain string) (bool, error) {
	domainGet, total, err := (&model.ProjectDomain{}).List(global.PanelDB, &model.ConditionsT{"domain": domain, "ORDER": "create_time DESC"}, 0, 0)
	if err != nil {
		return true, err
	}
	if total > 0 {
		return true, errors.New(helper.MessageWithMap("project.DomainNameHasExisted", map[string]any{"Domain": domainGet[0].Domain}))
	}

	return false, nil
}

// GetProjectPHPVersion 获取项目的php版本
func GetProjectPHPVersion(projectName string) (version string, err error) {
	confBody, err := util.ReadFileStringBody(ProjectMainConfFilePath(projectName))
	if err != nil {
		global.Log.Errorf("GetProjectPHPVersion->ReadFileStringBody Error:%s\n", err.Error())
		return "00", err
	}
	re := regexp.MustCompile(`enable-php-(\d+)\.conf`)
	match := re.FindStringSubmatch(confBody)
	if len(match) > 1 {
		version = match[1]
	} else {
		version = "00"
	}
	return
}

// GetProjectSSLInfo 获取项目的ssl信息
func GetProjectSSLInfo(projectName string) (sslDetails *response.SSLDetails, err error) {
	sslConfPath := fmt.Sprintf("%s/ssl.conf", ProjectConfDirPath(projectName))
	confBody, err := util.ReadFileStringBody(sslConfPath)
	if err != nil {
		return
	}
	if util.StrIsEmpty(confBody) {
		return
	}

	certPath := fmt.Sprintf("%s/cert", ProjectConfDirPath(projectName))
	privateBody, err := util.ReadFileStringBody(fmt.Sprintf("%s/private.pem", certPath))
	if err != nil {
		return
	}
	fullchainBody, err := util.ReadFileStringBody(fmt.Sprintf("%s/fullchain.pem", certPath))
	if err != nil {
		return
	}
	cert, err := util.ParseCert([]byte(fullchainBody))
	if err != nil {
		return
	}
	sslDetails = &response.SSLDetails{
		Domains:      cert.DNSNames,
		ExpireDate:   cert.NotAfter,
		StartDate:    cert.NotBefore,
		Type:         cert.Issuer.CommonName,
		Organization: cert.Issuer.Organization[0],
		Key:          privateBody,
		Csr:          fullchainBody,
	}
	start := strings.Index(confBody, "#AlwaysUseHttps_Start")
	end := strings.Index(confBody, "#AlwaysUseHttps_End")
	if start != -1 && end != -1 {
		sslDetails.AlwaysUseHttps = true
	}
	return
}

// GetProjectRootPath 从nginx配置文件取项目Root路径
func GetProjectRootPath(projectName string) (rootPath string, err error) {
	confBody, err := util.ReadFileStringBody(ProjectMainConfFilePath(projectName))
	if err != nil {
		global.Log.Errorf("GetProjectRootPath->ReadFileStringBody Error:%s\n", err.Error())
		return "00", err
	}
	re := regexp.MustCompile(`\s*root\s+(.+);`)
	match := re.FindStringSubmatch(confBody)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", errors.New("can not find root path")
}

// GenerateNginxConfig 生成nginx配置并写入文件
func GenerateNginxConfig(config *NginxConfig) (err error) {
	if util.StrIsEmpty(config.PHPVersion) {
		config.PHPVersion = "00"
	}
	mainConfigPath := ProjectMainConfFilePath(config.ProjectName)

	// 读取模板项目配置文件
	tmp, err := util.ReadFileStringBody(fmt.Sprintf("%s/template/main.conf", GetExtensionsPath(constant.ExtensionNginxName)))
	if err != nil {
		return err
	}
	projectRootDir := ProjectConfDirPath(config.ProjectName)
	sslFilePath := fmt.Sprintf("%s/ssl.conf", projectRootDir)
	if !util.PathExists(sslFilePath) {
		_ = os.MkdirAll(path.Dir(sslFilePath), 0755)
		_ = util.WriteFile(sslFilePath, []byte(""), 0644)

	}
	errorPageFilePath := fmt.Sprintf("%s/error_page.conf", projectRootDir)
	if !util.PathExists(errorPageFilePath) {
		_ = os.MkdirAll(path.Dir(errorPageFilePath), 0755)
		_ = util.WriteFile(errorPageFilePath, []byte(""), 0644)
	}
	antiLeechFilePath := fmt.Sprintf("%s/anti_leech.conf", projectRootDir)
	if !util.PathExists(antiLeechFilePath) {
		_ = os.MkdirAll(path.Dir(antiLeechFilePath), 0755)
		_ = util.WriteFile(antiLeechFilePath, []byte(""), 0644)
	}
	rewriteFilePath := fmt.Sprintf("%s/rewrite.conf", projectRootDir)
	if !util.PathExists(rewriteFilePath) {
		_ = os.MkdirAll(path.Dir(rewriteFilePath), 0755)
		_ = util.WriteFile(rewriteFilePath, []byte(""), 0644)

	}
	AccessRuleDirPath := fmt.Sprintf("%s/access_rule", projectRootDir)
	if !util.PathExists(AccessRuleDirPath) {
		_ = os.MkdirAll(AccessRuleDirPath, 0755)
	}
	ProxyDirPath := fmt.Sprintf("%s/proxy", projectRootDir)
	if !util.PathExists(ProxyDirPath) {
		_ = os.MkdirAll(ProxyDirPath, 0755)
	}
	redirectDirPath := fmt.Sprintf("%s/redirect", projectRootDir)
	if !util.PathExists(redirectDirPath) {
		_ = os.MkdirAll(redirectDirPath, 0755)
	}

	listenIpv4BlockStr := ""
	portMap := make(map[int]bool)
	var serverNames []string
	for _, item := range config.Domain {
		serverNames = append(serverNames, item.Name)
		if _, ok := portMap[item.Port]; !ok {
			portMap[item.Port] = true
			listenIpv4BlockStr += fmt.Sprintf("listen %d;\n", item.Port)
		}
	}
	//替换配置文件中的变量

	tmp = strings.Replace(tmp, "{{listen_ipv4_block}}", listenIpv4BlockStr, -1)
	tmp = strings.Replace(tmp, "{{listen_ipv6_block}}", "", -1) //如果是ipv6则替换，暂时不写，直接填空
	tmp = strings.Replace(tmp, "{{server_name_block}}", strings.Join(serverNames, " "), -1)
	tmp = strings.Replace(tmp, "{{root_block}}", config.ProjectPath, -1)
	tmp = strings.Replace(tmp, "{{include_ssl_block}}", sslFilePath, -1)
	tmp = strings.Replace(tmp, "{{include_error_page_block}}", errorPageFilePath, -1)
	tmp = strings.Replace(tmp, "{{include_access_rule_block}}", fmt.Sprintf("%s/*.conf", AccessRuleDirPath), -1)
	tmp = strings.Replace(tmp, "{{include_proxy_block}}", fmt.Sprintf("%s/*.conf", ProxyDirPath), -1)
	tmp = strings.Replace(tmp, "{{include_anti_leech_block}}", antiLeechFilePath, -1)
	tmp = strings.Replace(tmp, "{{include_redirect_block}}", fmt.Sprintf("%s/*.conf", redirectDirPath), -1)
	tmp = strings.Replace(tmp, "{{include_php_block}}", fmt.Sprintf("enable-php-%s.conf", config.PHPVersion), -1)
	tmp = strings.Replace(tmp, "{{include_rewrite_block}}", rewriteFilePath, -1)
	tmp = strings.Replace(tmp, "{{access_log_block}}", fmt.Sprintf("%s/%s.log", global.Config.System.WwwLogPath, config.ProjectName), -1)
	tmp = strings.Replace(tmp, "{{error_log_block}}", fmt.Sprintf("%s/%s.error.log", global.Config.System.WwwLogPath, config.ProjectName), -1)
	//include {{include_anti_leech_block}};
	//格式化配置文件
	tmp, err = FormatNginxConf(tmp)
	if err != nil {
		return err
	}
	// 写入配置文件
	err = util.WriteFile(mainConfigPath, []byte(tmp), 0644)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNginxConfig 删除nginx配置文件
func DeleteNginxConfig(projectName string) {
	mainConfigPath := ProjectMainConfFilePath(projectName)
	confDirPath := ProjectConfDirPath(projectName)
	_ = os.Remove(mainConfigPath)
	_ = os.RemoveAll(confDirPath)
	//重启nginx
	err := ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
}

// GenerateTtwafConfig 生成ttwaf配置并写入文件
//func GenerateTtwafConfig(projectName string) (err error) {
//	var projectConfigs map[string]safeResp.TTWafProjectConfig
//	var globalConfigs safeResp.TTWafConfig
//	//读取项目配置文件
//	projectConfigs, err = GetTTWafProjectConfig()
//	if err != nil {
//		return err
//	}
//	//读取全局配置文件
//	globalConfigs, err = GetTTWafConfig()
//	if err != nil {
//		return err
//	}
//	//生成项目配置文件
//	projectConfigs[projectName] = safeResp.TTWafProjectConfig{
//		Status: true,
//		DisablePhpPath: []string{
//			"^/cache/",
//			"^/config/",
//			"^/runtime/",
//			"^/application/",
//			"^/temp/",
//			"^/logs/",
//			"^/log/",
//		},
//		DisablePath: nil,
//		DisableExt: []string{
//			"sql",
//			"bak",
//			"swp",
//		},
//		DisableUploadExt: []string{
//			"sh",
//			"php",
//			"jsp",
//		},
//		Cc:               globalConfigs.Cc,
//		AttackTolerance:  globalConfigs.AttackTolerance,
//		SemanticAnalysis: globalConfigs.SemanticAnalysis,
//		BlockCountry:     globalConfigs.BlockCountry,
//		Get:              true,
//		Post:             true,
//		Cookie:           true,
//		UserAgent:        true,
//	}
//	//写入项目配置文件
//	err = WriteTTWafProjectConfig(projectConfigs)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// GenerateSslConfig 生成ssl配置文件并写入文件
func GenerateSslConfig(projectName string, private, fullchain []byte) (err error) {
	sslConfPath := fmt.Sprintf("%s/ssl.conf", ProjectConfDirPath(projectName))
	//生成private和fullchain文件
	certPath := fmt.Sprintf("%s/cert", ProjectConfDirPath(projectName))
	_ = os.MkdirAll(certPath, 0755)
	privatePath := fmt.Sprintf("%s/private.pem", certPath)
	fullchainPath := fmt.Sprintf("%s/fullchain.pem", certPath)
	if err = util.WriteFile(certPath+"/private.pem", private, 0644); err != nil {
		return
	}
	if err = util.WriteFile(certPath+"/fullchain.pem", fullchain, 0644); err != nil {
		return
	}
	// 读取模板SSL配置文件
	tmp, err := util.ReadFileStringBody(fmt.Sprintf("%s/template/ssl.conf", GetExtensionsPath(constant.ExtensionNginxName)))
	if err != nil {
		return err
	}
	tmp = strings.Replace(tmp, "{{certificate_path}}", fullchainPath, -1)
	tmp = strings.Replace(tmp, "{{certificate_key_path}}", privatePath, -1)
	//格式化配置文件
	tmp, err = FormatNginxConf(tmp)
	if err != nil {
		return
	}
	if err = util.WriteFile(sslConfPath, []byte(tmp), 0644); err != nil {
		return
	}

	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}

	return
}

// GenerateRedirectConfig 生成redirect配置文件并写入文件
func GenerateRedirectConfig(projectName string, ID string, param *projectR.CreateRedirectR) (err error) {
	redirectDirPath := fmt.Sprintf("%s/redirect", ProjectConfDirPath(projectName))
	_ = os.MkdirAll(redirectDirPath, 0755)
	redirectFilePath := fmt.Sprintf("%s/%s.conf", redirectDirPath, ID)

	var configTmp string
	//$request_uri

	if param.Type == constant.ProjectRedirectTypeDomain {
		domainTmp := `if ($host ~ '^{{domain}}'){
            return {{code}} {{target_url}}{{preserve_uri_args}};
    	}`
		for _, domain := range param.Domains {
			domainTmp = strings.Replace(domainTmp, "{{domain}}", domain, -1)
			domainTmp = strings.Replace(domainTmp, "{{code}}", strconv.Itoa(param.Code), -1)
			domainTmp = strings.Replace(domainTmp, "{{target_url}}", param.TargetUrl, -1)
			puaStr := ""
			if param.PreserveUriArgs {
				puaStr = "$request_uri"
			} else {
				puaStr = ""
			}
			domainTmp = strings.Replace(domainTmp, "{{preserve_uri_args}}", puaStr, -1)
			configTmp += domainTmp
		}

	}
	if param.Type == constant.ProjectRedirectTypePath {
		pathTmp := ` rewrite ^{{path}}(.*) {{target_url}}{{preserve_uri_args}} {{code}};`
		codeMap := map[int]string{
			301: "permanent",
			302: "redirect",
		}
		pathTmp = strings.Replace(pathTmp, "{{path}}", param.Path, -1)
		pathTmp = strings.Replace(pathTmp, "{{target_url}}", param.TargetUrl, -1)
		puaStr := ""
		if param.PreserveUriArgs {
			puaStr = "$1"
		} else {
			puaStr = ""
		}
		pathTmp = strings.Replace(pathTmp, "{{preserve_uri_args}}", puaStr, -1)
		pathTmp = strings.Replace(pathTmp, "{{code}}", codeMap[param.Code], -1)
		configTmp = pathTmp
	}
	//格式化配置文件
	configTmp, err = FormatNginxConf(configTmp)
	if err != nil {
		return
	}

	if err = util.WriteFile(redirectFilePath, []byte(configTmp), 0644); err != nil {
		return
	}

	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return
}

// GenerateAccessRuleConfig 生成访问规则配置并写入文件
func GenerateAccessRuleConfig(projectName string, ID string, param *projectR.CreateAccessRuleConfigR) (err error) {
	accessRuleDirPath := fmt.Sprintf("%s/access_rule", ProjectConfDirPath(projectName))
	_ = os.MkdirAll(accessRuleDirPath, 0755)
	accessRuleFilePath := fmt.Sprintf("%s/%s.conf", accessRuleDirPath, ID)

	accessRuleConf := ""
	//读取模板配置文件
	if param.RuleType == constant.ProjectAccessRuleTypeBasicAuth {
		if util.StrIsEmpty(param.BasicAuthConfig.User) || util.StrIsEmpty(param.BasicAuthConfig.Password) {
			return errors.New("又要设置密码访问，又特么留空，离谱他妈给离谱开门")
		}
		verifyFilePath := fmt.Sprintf("%s/%s.pass", accessRuleDirPath, ID)
		verifyStr := fmt.Sprintf("%s:%s", param.BasicAuthConfig.User, param.BasicAuthConfig.Password)
		if err = util.WriteFile(verifyFilePath, []byte(verifyStr), 0644); err != nil {
			return err
		}
		accessRuleConf = `location ~* ^{{dir}}* {
    		auth_basic "Authorization";
    		auth_basic_user_file {{verify_file_path}};
		}`
		accessRuleConf = strings.Replace(accessRuleConf, "{{dir}}", param.Dir, -1)
		accessRuleConf = strings.Replace(accessRuleConf, "{{verify_file_path}}", verifyFilePath, -1)
	} else if param.RuleType == constant.ProjectAccessRuleTypeNoAccess {
		accessRuleConf = `location ~* ^{{dir}}.*{{suffix_list}}$ {
        	deny all;
    	}`
		suffixTmp := ""
		if len(param.NoAccessConfig.SuffixList) > 0 {
			suffixTmp = ".(" + strings.Join(param.NoAccessConfig.SuffixList, "|") + ")"
		} else {
			suffixTmp = ""
		}
		accessRuleConf = strings.Replace(accessRuleConf, "{{dir}}", param.Dir, -1)
		accessRuleConf = strings.Replace(accessRuleConf, "{{suffix_list}}", suffixTmp, -1)
	} else {
		return errors.New("unknown access_rule type")
	}
	//格式化配置文件
	accessRuleConf, err = FormatNginxConf(accessRuleConf)
	if err != nil {
		return
	}
	if err = util.WriteFile(accessRuleFilePath, []byte(accessRuleConf), 0644); err != nil {
		return
	}
	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return
}

// GenerateReverseProxyConfig 生成反向代理配置并且写入文件
func GenerateReverseProxyConfig(projectName string, ID string, param *projectR.CreateReverseProxyConfigR) (err error) {
	reverseProxyDirPath := fmt.Sprintf("%s/proxy", ProjectConfDirPath(projectName))
	_ = os.MkdirAll(reverseProxyDirPath, 0755) //$host:$server_port
	reverseProxyFilePath := fmt.Sprintf("%s/%s.conf", reverseProxyDirPath, ID)
	confTmp := `location {{exact_match}} {{proxy_dir}} {
    proxy_pass {{target_url}};
    proxy_set_header Host {{host}};
	proxy_set_header X-Scheme $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header REMOTE-HOST $remote_addr;
    add_header X-Cache $upstream_cache_status;
    proxy_set_header Connection "upgrade";
    proxy_connect_timeout {{proxy_connect_timeout}}s;
   	proxy_read_timeout {{proxy_read_timeout}}s;
    proxy_send_timeout {{proxy_send_timeout}}s;
    proxy_http_version {{proxy_http_version}};
    proxy_set_header Upgrade $http_upgrade;
    {{sub_filters}}
	{{cache}}
    }
	`
	confTmp = strings.Replace(confTmp, "{{proxy_dir}}", param.ProxyDir, -1)
	confTmp = strings.Replace(confTmp, "{{target_url}}", param.TargetUrl, -1)
	confTmp = strings.Replace(confTmp, "{{host}}", param.Host, -1)
	confTmp = strings.Replace(confTmp, "{{proxy_connect_timeout}}", strconv.Itoa(param.ProxyConnectTimeout), -1)
	confTmp = strings.Replace(confTmp, "{{proxy_read_timeout}}", strconv.Itoa(param.ProxyReadTimeout), -1)
	confTmp = strings.Replace(confTmp, "{{proxy_send_timeout}}", strconv.Itoa(param.ProxySendTimeout), -1)
	confTmp = strings.Replace(confTmp, "{{proxy_http_version}}", param.ProxyHttpVersion, -1)

	//是否精准匹配
	if param.ExactMatch {
		confTmp = strings.Replace(confTmp, "{{exact_match}}", "^~", -1)
	} else {
		confTmp = strings.Replace(confTmp, "{{exact_match}}", "", -1)
	}

	//处理sub_filters
	subFiltersTmp := ""
	if len(param.SubFilters) > 0 {
		subFiltersTmp = `proxy_set_header Accept-Encoding "";`
		subFiltersTmp += fmt.Sprintf("\nsub_filter_types %s;", strings.Join(param.SubFilterTypes, " "))
		for _, v := range param.SubFilters {
			subFiltersTmp += fmt.Sprintf("\nsub_filter \"%s\" \"%s\";", v.Old, v.New)
		}
		subFiltersTmp += fmt.Sprintf("\nsub_filter_once %s;", param.SubFilterOnce)
		confTmp = strings.Replace(confTmp, "{{sub_filters}}", subFiltersTmp, -1)
	}
	confTmp = strings.Replace(confTmp, "{{sub_filters}}", subFiltersTmp, -1)
	//处理缓存
	cacheTmp := ""
	if param.CacheTime > 0 {
		cacheTmp = `if ( $uri ~* "\.(gif|png|jpg|css|js|woff|woff2)$" )
						{
        					expires 12h;
    					}
    					proxy_ignore_headers Set-Cookie Cache-Control expires;
    					proxy_cache cache_one;
    					proxy_cache_key $host$uri$is_args$args;
    					proxy_cache_valid 200 304 301 302 {{cache_time}}m;`
	} else {
		cacheTmp = `add_header Cache-Control no-cache;`
	}
	cacheTmp = strings.Replace(cacheTmp, "{{cache_time}}", strconv.Itoa(param.CacheTime), -1)
	confTmp = strings.Replace(confTmp, "{{cache}}", cacheTmp, -1)

	//格式化配置文件
	confTmp, err = FormatNginxConf(confTmp)
	if err != nil {
		return
	}

	if err = util.WriteFile(reverseProxyFilePath, []byte(confTmp), 0644); err != nil {
		return
	}
	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return
}

// GenerateAntiLeechConfig 生成防盗链配置并且写入文件
func GenerateAntiLeechConfig(projectName string, param *projectR.CreateAntiLeechConfigR) (err error) {
	antiLeechConfPath := fmt.Sprintf("%s/anti_leech.conf", ProjectConfDirPath(projectName))
	confTmp := `location ~ .*\.({{suffix_list}})$
    {
        expires      30d;
        access_log /dev/null;
        valid_referers {{referer_none}} {{pass_domains}};
        if ($invalid_referer){
           return {{response_status_code}};
        }
    }`
	confTmp = strings.Replace(confTmp, "{{suffix_list}}", strings.Join(param.SuffixList, "|"), -1)
	RefererNoneStr := ""
	if param.RefererNone {
		RefererNoneStr = "none blocked "
	}
	confTmp = strings.Replace(confTmp, "{{referer_none}}", RefererNoneStr, -1)
	confTmp = strings.Replace(confTmp, "{{pass_domains}}", strings.Join(param.PassDomains, " "), -1)
	confTmp = strings.Replace(confTmp, "{{response_status_code}}", strconv.Itoa(param.ResponseStatusCode), -1)
	//格式化配置文件
	confTmp, err = FormatNginxConf(confTmp)
	if err != nil {
		return
	}

	if err = util.WriteFile(antiLeechConfPath, []byte(confTmp), 0644); err != nil {
		return
	}
	//重启nginx
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		return
	}
	return
}

// GetTTWafProjectConfig 获取ttwaf项目配置
//func GetTTWafProjectConfig() (ttWafProjectConfig map[string]safeResp.TTWafProjectConfig, err error) {
//	projectConfigPath := fmt.Sprintf("%s/ttwaf/project.json", global.Config.System.ServerPath)
//	configBody, err := util.ReadFileStringBody(projectConfigPath)
//	if err != nil {
//		fmt.Println("GetTTWafProjectConfig.util.ReadFileStringBody.ERROR：", err)
//		return
//	}
//	err = util.JsonStrToStruct(configBody, &ttWafProjectConfig)
//	if err != nil {
//		return
//	}
//	return
//}

// WriteTTWafProjectConfig 写入ttwaf项目配置
//func WriteTTWafProjectConfig(projectConfigs map[string]safeResp.TTWafProjectConfig) (err error) {
//	projectConfigPath := fmt.Sprintf("%s/ttwaf/project.json", global.Config.System.ServerPath)
//	projectConfigsStr, err := util.StructToJsonStr(projectConfigs)
//	if err != nil {
//		return err
//	}
//	err = util.WriteFile(projectConfigPath, []byte(projectConfigsStr), 0644)
//	if err != nil {
//		return err
//	}
//	return
//}

// GetTTWafConfig 获取ttwaf全局配置
//func GetTTWafConfig() (ttWafConfig safeResp.TTWafConfig, err error) {
//	configPath := fmt.Sprintf("%s/ttwaf/config.json", global.Config.System.ServerPath)
//	configBody, err := util.ReadFileStringBody(configPath)
//	if err != nil {
//		fmt.Println("GetTTWafConfig.util.ReadFileStringBody.ERROR：", err)
//		return
//	}
//	err = util.JsonStrToStruct(configBody, &ttWafConfig)
//	if err != nil {
//		return
//	}
//	return
//}

// GetTTWafDomainConfig 获取ttwaf域名配置
//func GetTTWafDomainConfig() (ttWafDomainConfig map[string]*safeModel.TTWafDomainSetting, err error) {
//	//读取域名配置文件
//	configPath := fmt.Sprintf("%s/ttwaf/domain.json", global.Config.System.ServerPath)
//	projectBody, err := util.ReadFileStringBody(configPath)
//	if err != nil {
//		fmt.Println("OperateTTWafDomainConfig.util.ReadFileStringBody.ERROR：", err)
//		return
//	}
//	err = util.JsonStrToStruct(projectBody, &ttWafDomainConfig)
//	if err != nil {
//		return
//	}
//	return
//}
//
//// WriteTTWafDomainConfig 写入ttwaf域名配置
//func WriteTTWafDomainConfig(domainConfigs map[string]*safeModel.TTWafDomainSetting) (err error) {
//	configPath := fmt.Sprintf("%s/ttwaf/domain.json", global.Config.System.ServerPath)
//	//写入域名配置文件
//	domainConfigsStr, err := util.StructToJsonStr(&domainConfigs)
//	if err != nil {
//		return err
//	}
//	err = util.WriteFile(configPath, []byte(domainConfigsStr), 0644)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// OperateTTWafDomainConfig 操作ttwaf域名配置
//func OperateTTWafDomainConfig(projectID int64, projectName string, path string, action string) (err error) {
//	domainConfigs, err := GetTTWafDomainConfig()
//	if err != nil {
//		return
//	}
//	switch action {
//	case constant.OperateTTWafDomainConfigByUpdate:
//		domainList, total, errs := (&project_manager.ProjectDomain{}).List(global.PanelDB, &common.ConditionsT{"project_id": projectID}, 0, 0)
//		if errs != nil {
//			err = errs
//			return
//		}
//		if total == 0 {
//			delete(domainConfigs, projectName)
//			break
//		}
//		var domainListStr []string
//		for _, domainInfo := range domainList {
//			domainListStr = append(domainListStr, domainInfo.Domain)
//		}
//		domainConfigs[projectName] = &safeModel.TTWafDomainSetting{
//			Path:       path,
//			DomainList: domainListStr,
//		}
//	case constant.OperateTTWafDomainConfigByDelete:
//		delete(domainConfigs, projectName)
//	}
//	err = WriteTTWafDomainConfig(domainConfigs)
//	if err != nil {
//		return
//	}
//	return
//}

// AddDomainToNginxConfig nginx配置文件添加域名
func AddDomainToNginxConfig(projectName string, domainMap, portMap map[string]bool) (err error) {
	configPath := ProjectMainConfFilePath(projectName)
	configBody, err := util.ReadFileStringBody(configPath)
	if err != nil {
		return err
	}
	// 添加域名
	rep := regexp.MustCompile(`server_name\s*(.*);`)
	tmp := rep.FindStringSubmatch(configBody)[1]
	domains := strings.Split(tmp, " ")
	//合并域名
	for _, v := range domains {
		domainMap[v] = true
	}
	var newServerName string
	for s := range domainMap {
		newServerName += s + " "
	}
	configBody = strings.Replace(configBody, tmp, newServerName, 1)

	// 添加端口
	rep = regexp.MustCompile(`listen\s+[\[\]:]*([0-9]+).*;`)

	tmp2 := rep.FindAllStringSubmatch(configBody, -1)
	for _, tmp3 := range tmp2 {
		for _, tmp4 := range tmp3 {
			delete(portMap, tmp4)
		}
	}
	if len(portMap) > 0 {
		listen := rep.FindString(configBody)
		newListen := ""
		for k, _ := range portMap {
			listenIPv6 := "" //if self.isIPv6 {
			//	listenIPv6 = "listen [::]:" + port + ";"
			//}
			newListen += "listen " + k + ";" + listenIPv6
		}
		global.Log.Debugf("AddDomainToNginxConfig->listen.rep.FindString:%s\n", listen)
		global.Log.Debugf("AddDomainToNginxConfig->newListen:%s\n", newListen)

		configBody = strings.Replace(configBody, listen, listen+newListen, 1)
	}
	configBody, err = FormatNginxConf(configBody)
	if err != nil {
		return err
	}
	// 写入配置文件
	err = util.WriteFile(configPath, []byte(configBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

// SetRootPath 设置运行目录
func SetRootPath(projectName string, rootPath string) (err error) {
	configPath := ProjectMainConfFilePath(projectName)
	configBody, err := util.ReadFileStringBody(configPath)
	if err != nil {
		return err
	}
	backupConfig := configBody
	re := regexp.MustCompile(`root\s+([^\s;]+);`)
	match := re.FindStringSubmatch(configBody)
	if len(match) > 1 {
		configBody = strings.Replace(configBody, match[1], rootPath, 1)
	} else {
		return errors.New("not found root path")
	}
	err = util.WriteFile(configPath, []byte(configBody), 0644)
	if err != nil {
		return err
	}
	err = ReloadNginx()
	if err != nil {
		global.Log.Error(err)
		_ = util.WriteFile(configPath, []byte(backupConfig), 0644)
		return
	}
	return nil
}

// DelDomainToNginxConfig nginx配置文件删除域名
func DelDomainToNginxConfig(projectName string, delDomain []string, delPort []string) (err error) {
	configPath := ProjectMainConfFilePath(projectName)
	configBody, err := util.ReadFileStringBody(configPath)
	if err != nil {
		return err
	}
	// 删除域名
	if len(delDomain) > 0 {
		rep := regexp.MustCompile(`server_name\s*(.*);`)
		tmp := rep.FindStringSubmatch(configBody)[1]
		nowDomains := strings.Split(tmp, " ")
		newDomainMap := make(map[string]bool)
		for _, v := range nowDomains {
			newDomainMap[v] = true
		}

		//删除键
		for _, v := range delDomain {
			delete(newDomainMap, v)
		}

		//拼接
		newServerName := ""
		for k, _ := range newDomainMap {
			newServerName += k + " "
		}
		configBody = strings.Replace(configBody, tmp, newServerName, 1)
	}

	// 删除端口
	for _, v := range delPort {
		rep4 := regexp.MustCompile(`listen\s+` + v + `\s*;`)
		rep6 := regexp.MustCompile(`listen\s+\[::\]:` + v + `\s*;`)
		configBody = rep4.ReplaceAllString(configBody, "")
		configBody = rep6.ReplaceAllString(configBody, "")
	}

	configBody, err = FormatNginxConf(configBody)
	if err != nil {
		return err
	}
	// 写入配置文件
	err = util.WriteFile(configPath, []byte(configBody), 0644)
	if err != nil {
		return err
	}

	return nil
}

// FormatNginxConf 格式化nginx配置文件
func FormatNginxConf(conf string) (confN string, err error) {
	fmt.Println(conf)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("nginx Config Format Error")
		}
	}()
	p := parser.NewStringParser(conf)
	configP := p.Parse()
	confN = gonginx.DumpConfig(configP, gonginx.IndentedStyle)
	return
}

// ProjectMainConfFilePath 项目主配置文件路径
func ProjectMainConfFilePath(projectName string) (filePath string) {
	filePath = GetExtensionsPath(constant.ExtensionNginxName) + "/vhost/main/" + projectName + ".conf"
	return
}

// ProjectConfDirPath 项目配置文件目录路径
func ProjectConfDirPath(projectName string) (filePath string) {
	filePath = GetExtensionsPath(constant.ExtensionNginxName) + "/vhost/project/" + projectName
	return
}
