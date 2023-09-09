package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type TTWafService struct {
}

// Config ttwaf配置
func (t *TTWafService) Config() (*response.TTWafConfig, error) {
	fileBody, err := util.ReadFileStringBody(global.Config.System.ServerPath + "/ttwaf/config/config.json")
	if err != nil {
		global.Log.Errorf("Config->ReadFileStringBody  Error:%s", err.Error())
		return nil, err
	}
	var config response.TTWafConfig
	err = json.Unmarshal([]byte(fileBody), &config)
	if err != nil {
		global.Log.Errorf("Config->json.Unmarshal  Error:%s", err.Error())
		return nil, err
	}
	return &config, nil
}

// ProjectConfig ttwaf项目配置
func (t *TTWafService) ProjectConfig(projectID int64) (config response.TTWafProjectConfig, err error) {
	projectInfo, err := (&model.Project{ID: projectID}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if projectInfo.ID > 0 {
		configList, errs := t.GetProjectConfig()
		if errs != nil {
			err = errs
			return
		}
		return configList[projectInfo.Name], nil
	}

	return
}

// SaveConfig 保存ttwaf配置
func (t *TTWafService) SaveConfig(param *response.TTWafConfig) error {
	// 将结构体转换为 JSON 字符串
	configJsonByte, err := json.Marshal(param)
	if err != nil {
		global.Log.Errorf("SaveConfig->json.Marshal  Error:%s", err.Error())
		return err
	}
	// 写入文件
	err = util.WriteFile(global.Config.System.ServerPath+"/ttwaf/config/config.json", configJsonByte, 0666)
	if err != nil {
		global.Log.Errorf("SaveConfig->WriteFile  Error:%s", err.Error())
		return err
	}
	return nil
}

// SaveProjectConfig 保存ttwaf项目配置
func (t *TTWafService) SaveProjectConfig(projectID int64, config response.TTWafProjectConfig) (err error) {
	projectInfo, err := (&model.Project{ID: projectID}).Get(global.PanelDB)
	if err != nil {
		return
	}
	if projectInfo.ID == 0 {
		return errors.New("project is not found")
	}
	configList, errs := t.GetProjectConfig()
	if errs != nil {
		err = errs
		return
	}
	configList[projectInfo.Name] = config
	err = t.WriteProjectConfig(configList)
	if err != nil {
		return err
	}
	return nil
}

// GlobalSet ttwaf全局设置
func (t *TTWafService) GlobalSet(param *request.TTWafGlobalSetR) error {
	//读取tt_waf的项目配置文件
	fileBody, err := util.ReadFileStringBody(global.Config.System.ServerPath + "/ttwaf/config/project.json")
	if err != nil {
		global.Log.Errorf("GlobalSet->ReadFileStringBody  Error:%s", err.Error())
		return err
	}
	var projectConfigMap map[string]response.TTWafProjectConfig
	err = json.Unmarshal([]byte(fileBody), &projectConfigMap)
	if err != nil {
		global.Log.Errorf("GlobalSet->json.Unmarshal  Error:%s", err.Error())
		return err
	}
	for _, projectConfig := range projectConfigMap {
		if param.Cc != nil {
			projectConfig.Cc = *param.Cc
		}
		if param.AttackTolerance != nil {
			projectConfig.AttackTolerance = *param.AttackTolerance
		}
	}
	// 将结构体转换为 JSON 字符串
	projectConfigMapJsonByte, err := json.Marshal(projectConfigMap)
	if err != nil {
		global.Log.Errorf("GlobalSet->json.Marshal  Error:%s", err.Error())
		return err
	}
	// 写入文件
	err = util.WriteFile(global.Config.System.ServerPath+"/ttwaf/config/project.json", projectConfigMapJsonByte, 0666)
	if err != nil {
		global.Log.Errorf("GlobalSet->WriteFile  Error:%s", err.Error())
		return err
	}

	//读取tt_waf的全局配置文件
	config, err := t.Config()
	if err != nil {
		return err
	}
	if param.Cc != nil {
		config.Cc = *param.Cc
	}
	if param.AttackTolerance != nil {
		config.AttackTolerance = *param.AttackTolerance
	}
	// 将结构体转换为 JSON 字符串
	configMapJsonByte, err := json.Marshal(config)
	if err != nil {
		global.Log.Errorf("GlobalSet->json.Marshal  Error:%s", err.Error())
		return err
	}
	// 写入文件
	err = util.WriteFile(global.Config.System.ServerPath+"/ttwaf/config/config.json", configMapJsonByte, 0666)
	if err != nil {
		global.Log.Errorf("GlobalSet->json.WriteFile  Error:%s", err.Error())
		return err
	}
	return nil
}

// BlockList ttwaf封锁列表
func (t *TTWafService) BlockList(param *request.TTWafBlockListR, offset, limit int) ([]*model.TTWafBlockIpLog, int64, error) {
	//构造条件
	whereT := model.ConditionsT{
		"ORDER": "create_time DESC",
	}
	whereOrT := model.ConditionsT{}
	//分割ip
	if !util.StrIsEmpty(param.Ip) {
		ipList := strings.Split(param.Ip, ",")
		if len(ipList) > 1 {
			whereT["ip IN ?"] = ipList
		} else {
			if strings.HasSuffix(param.Ip, "*") {
				whereT["ip LIKE ?"] = param.Ip + "%"
			} else {
				whereT["ip"] = param.Ip
			}
		}
	}
	//分割url
	if !util.StrIsEmpty(param.Url) {
		urlList := strings.Split(param.Url, ",")
		if len(urlList) > 1 {
			whereT["uri IN ?"] = urlList
		} else {
			whereT["uri"] = param.Ip
		}
	}
	if !util.StrIsEmpty(param.Domain) {
		//查询域名对应的server_name
		whereT["server_name"] = param.Domain
	}
	if !util.StrIsEmpty(param.UserAgent) {
		whereT["user_agent"] = param.UserAgent
	}
	if param.StartTime != 0 && param.EndTime != 0 {
		whereT["create_time > ?"] = param.StartTime
		whereT["create_time < ?"] = param.EndTime
	}

	blockIPs, total, err := (&model.TTWafBlockIpLog{}).List(global.TTWafDB, &whereT, &whereOrT, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return blockIPs, total, err
}

// BanList ttwaf封禁列表
func (t *TTWafService) BanList(param *request.TTWafBanListR, offset, limit int) ([]*model.TTWafBanIpLog, int64, error) {
	//构造条件
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	//分割ip
	if !util.StrIsEmpty(param.Ip) {
		ipList := strings.Split(param.Ip, ",")
		if len(ipList) > 1 {
			whereT["ip IN ?"] = ipList
		} else {
			if strings.HasSuffix(param.Ip, "*") {
				whereT["ip LIKE ?"] = param.Ip + "%"
			} else {
				whereT["ip"] = param.Ip
			}
		}
	}
	//分割url
	if !util.StrIsEmpty(param.Url) {
		urlList := strings.Split(param.Url, ",")
		if len(urlList) > 1 {
			whereT["uri IN ?"] = urlList
		} else {
			whereT["uri"] = param.Ip
		}
	}
	if !util.StrIsEmpty(param.Domain) {
		//查询域名对应的server_name
		whereT["server_name"] = param.Domain
	}
	if !util.StrIsEmpty(param.UserAgent) {
		whereT["user_agent"] = param.UserAgent
	}
	if param.StartTime != 0 && param.EndTime != 0 {
		whereT["create_time > ?"] = param.StartTime
		whereT["create_time < ?"] = param.EndTime
	}
	banIPs, total, err := (&model.TTWafBanIpLog{}).List(global.TTWafDB, &whereT, &whereOrT, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return banIPs, total, err
}

// GetRegRule ttwaf获取规则
func (t *TTWafService) GetRegRule(ruleName string) ([]response.TTWafRegRule, error) {
	fileBody, err := util.ReadFileStringBody(global.Config.System.ServerPath + "/ttwaf/rule/" + ruleName + ".json")
	if err != nil {
		global.Log.Errorf("GetRegRule->ReadFileStringBody  Error:%s", err.Error())
		return nil, err
	}
	var list []response.TTWafRegRule
	err = json.Unmarshal([]byte(fileBody), &list)
	if err != nil {
		global.Log.Errorf("GetRegRule->json.Unmarshal  Error:%s", err.Error())
		return nil, err
	}
	return list, nil
}

// SaveRegRule ttwaf保存规则
func (t *TTWafService) SaveRegRule(ruleName string, ruleList []response.TTWafRegRule) error {
	ruleStr, err := util.StructToJsonStr(ruleList)
	if err != nil {
		return err
	}
	err = util.WriteFile(global.Config.System.ServerPath+"/ttwaf/rule/"+ruleName+".json", []byte(ruleStr), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Overview ttwaf概览
func (t *TTWafService) Overview() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	//查询block_ip总数量
	blockIpCountAll, err := (&model.TTWafBlockIpLog{}).Count(global.TTWafDB, &model.ConditionsT{})
	if err != nil {
		return nil, err
	}
	result["block_ip_count_all"] = blockIpCountAll
	//查询24小时block_ip数量
	BlockIpCount24, err := (&model.TTWafBlockIpLog{}).Count(global.TTWafDB, &model.ConditionsT{
		"create_time > ?": time.Now().Add(-24 * time.Hour).Unix(),
		"create_time < ?": time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	result["block_ip_count_24h"] = BlockIpCount24
	//查询正在封锁的IP数量
	banIpCountIng, err := (&model.TTWafBanIpLog{}).Count(global.TTWafDB, &model.ConditionsT{
		"status": 1,
	})
	if err != nil {
		return nil, err
	}
	result["ban_ip_count_ing"] = banIpCountIng
	//查询24小时ban_ip数量
	banIpCount24h, err := (&model.TTWafBanIpLog{}).Count(global.TTWafDB, &model.ConditionsT{
		"create_time > ?": time.Now().Add(-24 * time.Hour).Unix(),
		"create_time < ?": time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	result["ban_ip_count_24h"] = banIpCount24h
	//查询近七天block_ip数量
	blockIpCount7d, err := (&model.TTWafBlockIpLog{}).CountByDay(global.TTWafDB, 7)
	if err != nil {
		return nil, err
	}
	if len(blockIpCount7d) < 7 {
		// 存储结果到map中
		resultMap := make(map[string]int)
		for _, r := range blockIpCount7d {
			resultMap[r.Date] = r.Count
		}
		blockIpCount7d = nil
		// 遍历近七天的日期范围，如果map中不存在该日期，则将数量设置为0
		for i := 0; i < 7; i++ {
			date := time.Now().AddDate(0, 0, -6+i).Format("2006-01-02")
			if _, ok := resultMap[date]; !ok {
				resultMap[date] = 0
			}
			blockIpCount7d = append(blockIpCount7d, model.CountByDayData{Date: date, Count: resultMap[date]})
		}

	}
	result["block_ip_count_7d"] = blockIpCount7d
	//查询今天block_ip前五的项目
	blockIpTopProject, err := (&model.TTWafBlockIpLog{}).TopServerNames(global.TTWafDB, "server_name", 5)
	if err != nil {
		return nil, err
	}
	result["block_ip_top_project"] = blockIpTopProject
	//查询正在封锁的ip列表
	banIPList, _, err := (&model.TTWafBanIpLog{}).List(global.TTWafDB, &model.ConditionsT{"ORDER": "create_time DESC"}, &model.ConditionsT{}, 0, 20)
	if err != nil {
		return nil, err
	}
	result["ban_ip_list"] = banIPList
	//查询拦截列表
	blockIPList, _, err := (&model.TTWafBlockIpLog{}).List(global.TTWafDB, &model.ConditionsT{"ORDER": "create_time DESC"}, &model.ConditionsT{}, 0, 20)
	if err != nil {
		return nil, err
	}
	result["block_ip_list"] = blockIPList
	return result, nil
}

// AllowIP 解封IP
func (t *TTWafService) AllowIP(ip string) error {
	err := (&model.TTWafBanIpLog{
		Status: 0,
	}).BatchUpdates(global.TTWafDB, "status", &model.ConditionsT{
		"ip = ?": ip,
	})
	if err != nil {
		return err
	}

	_, err = t.RequestTTWaf(fmt.Sprintf("http://127.0.0.1/ttwaf/remove_drop_ip?ip=%s", ip))
	if err != nil {
		return err
	}

	return nil
}

// AddIpBlackList 添加ip黑名单
func (t *TTWafService) AddIpBlackList(ipItem *request.TTWafIPBlackWhiteR) error {
	config, err := t.Config()
	if err != nil {
		return err
	}

	if len(ipItem.IPV4) > 0 {
		for _, ipv4 := range ipItem.IPV4 {
			if len(ipv4) == 2 {
				config.IpBlackList.IPV4 = append(config.IpBlackList.IPV4, ipv4)
			}
		}
	}
	if len(ipItem.IPV6) > 0 {
		for _, ipv6 := range ipItem.IPV6 {
			config.IpBlackList.IPV6 = append(config.IpBlackList.IPV6, ipv6)
		}
	}

	config.IpBlackList.IPV6 = util.UniqueStrSlice(config.IpBlackList.IPV6)
	err = t.SaveConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTWafService) RequestTTWaf(u string) (body []byte, err error) {
	body = nil
	//获取防火墙访问key
	key, err := util.ReadFileStringBody(fmt.Sprintf("%s/ttwaf/config/access_key", global.Config.System.ServerPath))
	if err != nil {
		return
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return
	}
	params := parsedURL.Query()
	params.Set("access_key", key)
	parsedURL.RawQuery = params.Encode()
	u = parsedURL.String()

	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建一个HTTP请求
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// 关闭响应体
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 读取响应体
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(string(body))
		return
	}
	return
}

// GetProjectConfig 获取ttwaf项目配置
func (t *TTWafService) GetProjectConfig() (ttWafProjectConfig map[string]response.TTWafProjectConfig, err error) {
	projectConfigPath := fmt.Sprintf("%s/ttwaf/config/project.json", global.Config.System.ServerPath)
	configBody, err := util.ReadFileStringBody(projectConfigPath)
	if err != nil {
		fmt.Println("GetTTWafProjectConfig.util.ReadFileStringBody.ERROR：", err)
		return
	}
	err = util.JsonStrToStruct(configBody, &ttWafProjectConfig)
	if err != nil {
		return
	}
	return
}

// WriteProjectConfig 写入ttwaf项目配置
func (t *TTWafService) WriteProjectConfig(projectConfigs map[string]response.TTWafProjectConfig) (err error) {
	projectConfigPath := fmt.Sprintf("%s/ttwaf/config/project.json", global.Config.System.ServerPath)
	projectConfigsStr, err := util.StructToJsonStr(projectConfigs)
	if err != nil {
		return err
	}
	err = util.WriteFile(projectConfigPath, []byte(projectConfigsStr), 0644)
	if err != nil {
		return err
	}
	return
}

// GenerateProjectConfig 生成ttwaf项目配置并写入文件
func (t *TTWafService) GenerateProjectConfig(projectName string) (err error) {
	var projectConfigs map[string]response.TTWafProjectConfig
	var globalConfigs *response.TTWafConfig
	//读取项目配置文件
	projectConfigs, err = t.GetProjectConfig()
	if err != nil {
		return err
	}
	//读取全局配置文件
	globalConfigs, err = t.Config()
	if err != nil {
		return err
	}
	//生成项目配置文件
	projectConfigs[projectName] = response.TTWafProjectConfig{
		Status: true,
		DisablePhpPath: []string{
			"^/cache/",
			"^/config/",
			"^/runtime/",
			"^/application/",
			"^/temp/",
			"^/logs/",
			"^/log/",
		},
		DisablePath: nil,
		DisableExt: []string{
			"sql",
			"bak",
			"swp",
		},
		DisableUploadExt: []string{
			"sh",
			"php",
			"jsp",
		},
		Cc:               globalConfigs.Cc,
		AttackTolerance:  globalConfigs.AttackTolerance,
		SemanticAnalysis: globalConfigs.SemanticAnalysis,
		BlockCountry:     globalConfigs.BlockCountry,
		Get:              true,
		Post:             true,
		Cookie:           true,
		UserAgent:        true,
	}
	//写入项目配置文件
	err = t.WriteProjectConfig(projectConfigs)
	if err != nil {
		return err
	}
	return nil
}

// GetDomainConfig 获取ttwaf域名配置
func (t *TTWafService) GetDomainConfig() (ttWafDomainConfig map[string]*response.TTWafDomainSetting, err error) {
	//读取域名配置文件
	configPath := fmt.Sprintf("%s/ttwaf/config/domain.json", global.Config.System.ServerPath)
	projectBody, err := util.ReadFileStringBody(configPath)
	if err != nil {
		fmt.Println("OperateTTWafDomainConfig.util.ReadFileStringBody.ERROR：", err)
		return
	}
	err = util.JsonStrToStruct(projectBody, &ttWafDomainConfig)
	if err != nil {
		return
	}
	return
}

// WriteDomainConfig 写入ttwaf域名配置
func (t *TTWafService) WriteDomainConfig(domainConfigs map[string]*response.TTWafDomainSetting) (err error) {
	configPath := fmt.Sprintf("%s/ttwaf/config/domain.json", global.Config.System.ServerPath)
	//写入域名配置文件
	domainConfigsStr, err := util.StructToJsonStr(&domainConfigs)
	if err != nil {
		return err
	}
	err = util.WriteFile(configPath, []byte(domainConfigsStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

// OperateDomainConfig 操作ttwaf域名配置
func (t *TTWafService) OperateDomainConfig(projectID int64, projectName string, path string, action string) (err error) {
	domainConfigs, err := t.GetDomainConfig()
	if err != nil {
		return
	}
	switch action {
	case constant.OperateTTWafDomainConfigByUpdate:
		domainList, total, errs := (&model.ProjectDomain{}).List(global.PanelDB, &model.ConditionsT{"project_id": projectID, "ORDER": "create_time DESC"}, 0, 0)
		if errs != nil {
			err = errs
			return
		}
		if total == 0 {
			delete(domainConfigs, projectName)
			break
		}
		var domainListStr []string
		for _, domainInfo := range domainList {
			domainListStr = append(domainListStr, domainInfo.Domain)
		}
		domainConfigs[projectName] = &response.TTWafDomainSetting{
			Path:       path,
			DomainList: domainListStr,
		}
	case constant.OperateTTWafDomainConfigByDelete:
		delete(domainConfigs, projectName)
	}
	err = t.WriteDomainConfig(domainConfigs)
	if err != nil {
		return
	}
	return
}

// CountryList ttwaf国家列表
func (t *TTWafService) CountryList() []string {
	var countryList []string
	countryList = []string{"中国大陆以外的地区(包括[中国特别行政区:港,澳,台])", "中国大陆(不包括[中国特别行政区:港,澳,台])", "中国香港", "中国澳门", "中国台湾",
		"美国", "日本", "英国", "德国", "韩国", "法国", "巴西", "加拿大", "意大利", "澳大利亚", "荷兰", "俄罗斯", "印度", "瑞典", "西班牙", "墨西哥",
		"比利时", "南非", "波兰", "瑞士", "阿根廷", "印度尼西亚", "埃及", "哥伦比亚", "土耳其", "越南", "挪威", "芬兰", "丹麦", "乌克兰", "奥地利",
		"伊朗", "智利", "罗马尼亚", "捷克", "泰国", "沙特阿拉伯", "以色列", "新西兰", "委内瑞拉", "摩洛哥", "马来西亚", "葡萄牙", "爱尔兰", "新加坡",
		"欧洲联盟", "匈牙利", "希腊", "菲律宾", "巴基斯坦", "保加利亚", "肯尼亚", "阿拉伯联合酋长国", "阿尔及利亚", "塞舌尔", "突尼斯", "秘鲁", "哈萨克斯坦",
		"斯洛伐克", "斯洛文尼亚", "厄瓜多尔", "哥斯达黎加", "乌拉圭", "立陶宛", "塞尔维亚", "尼日利亚", "克罗地亚", "科威特", "巴拿马", "毛里求斯", "白俄罗斯",
		"拉脱维亚", "多米尼加", "卢森堡", "爱沙尼亚", "苏丹", "格鲁吉亚", "安哥拉", "玻利维亚", "赞比亚", "孟加拉国", "巴拉圭", "波多黎各", "坦桑尼亚",
		"塞浦路斯", "摩尔多瓦", "阿曼", "冰岛", "叙利亚", "卡塔尔", "波黑", "加纳", "阿塞拜疆", "马其顿", "约旦", "萨尔瓦多", "伊拉克", "亚美尼亚", "马耳他",
		"危地马拉", "巴勒斯坦", "斯里兰卡", "特立尼达和多巴哥", "黎巴嫩", "尼泊尔", "纳米比亚", "巴林", "洪都拉斯", "莫桑比克", "尼加拉瓜", "卢旺达", "加蓬",
		"阿尔巴尼亚", "利比亚", "吉尔吉斯坦", "柬埔寨", "古巴", "喀麦隆", "乌干达", "塞内加尔", "乌兹别克斯坦", "黑山", "关岛", "牙买加", "蒙古", "文莱",
		"英属维尔京群岛", "留尼旺", "库拉索岛", "科特迪瓦", "开曼群岛", "巴巴多斯", "马达加斯加", "伯利兹", "新喀里多尼亚", "海地", "马拉维", "斐济", "巴哈马",
		"博茨瓦纳", "扎伊尔", "阿富汗", "莱索托", "百慕大", "埃塞俄比亚", "美属维尔京群岛", "列支敦士登", "津巴布韦", "直布罗陀", "苏里南", "马里", "也门",
		"老挝", "塔吉克斯坦", "安提瓜和巴布达", "贝宁", "法属玻利尼西亚", "圣基茨和尼维斯", "圭亚那", "布基纳法索", "马尔代夫", "泽西岛", "摩纳哥", "巴布亚新几内亚",
		"刚果", "塞拉利昂", "吉布提", "斯威士兰", "缅甸", "毛里塔尼亚", "法罗群岛", "尼日尔", "安道尔", "阿鲁巴", "布隆迪", "圣马力诺", "利比里亚",
		"冈比亚", "不丹", "几内亚", "圣文森特岛", "荷兰加勒比区", "圣马丁", "多哥", "格陵兰", "佛得角", "马恩岛", "索马里", "法属圭亚那", "西萨摩亚",
		"土库曼斯坦", "瓜德罗普", "马里亚那群岛", "瓦努阿图", "马提尼克", "赤道几内亚", "南苏丹", "梵蒂冈", "格林纳达", "所罗门群岛", "特克斯和凯科斯群岛", "多米尼克",
		"乍得", "汤加", "瑙鲁", "圣多美和普林西比", "安圭拉岛", "法属圣马丁", "图瓦卢", "库克群岛", "密克罗尼西亚联邦", "根西岛", "东帝汶", "中非",
		"几内亚比绍", "帕劳", "美属萨摩亚", "厄立特里亚", "科摩罗", "圣皮埃尔和密克隆", "瓦利斯和富图纳", "英属印度洋领地", "托克劳", "马绍尔群岛", "基里巴斯",
		"纽埃", "诺福克岛", "蒙特塞拉特岛", "朝鲜", "马约特", "圣卢西亚", "圣巴泰勒米岛"}
	return countryList
}
