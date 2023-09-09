package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	response2 "TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type ExtensionNodejsService struct {
}

var (
	envNodeBin = "/usr/bin/node"
	envNpmBin  = "/usr/bin/npm"
	envNpxBin  = "/usr/bin/npx"
	envPm2Bin  = "/usr/bin/pm2"
	envYarnBin = "/usr/bin/yarn"
)
var (
	InstallNodejsShellName      = "install_nodejs.sh"
	InstallNodeModulesShellName = "install_modules.sh"
)

func (s *ExtensionNodejsService) GetNodejsVersionListJsonPath() string {
	return fmt.Sprintf("%s/data/extensions/%s/version.json", global.Config.System.PanelPath, constant.ExtensionNodejsName)
}
func (s *ExtensionNodejsService) GetNodejsConfigJsonPath() string {
	return fmt.Sprintf("%s/data/extensions/%s/config.json", global.Config.System.PanelPath, constant.ExtensionNodejsName)
}

// Info 获取Nodejs信息
func (s *ExtensionNodejsService) Info() (*response2.ExtensionsInfoResponse, error) {
	var nodejsInfo response2.ExtensionsInfoResponse
	err := ReadExtensionsInfo(constant.ExtensionNodejsName, &nodejsInfo)
	if err != nil {
		return nil, err
	}
	//获取版本号和安装状态
	nodejsInfo.Description.Version = "0.0.1"
	nodejsInfo.Description.Install = true
	//获取运行状态
	nodejsInfo.Description.Status = true
	return &nodejsInfo, nil
}

// Config 获取Nodejs配置
func (s *ExtensionNodejsService) Config() (*response2.NodejsConfigResponse, error) {
	nodejsConfig, err := s.GetNodejsConfig()
	if err != nil {
		return nil, err
	}
	//获取命令行版本
	nodejsConfig.CliVersion, _ = s.GetNodejsVersion()
	//获取版本列表上次更新时间
	if util.PathExists(s.GetNodejsVersionListJsonPath()) {
		file, err := os.Stat(s.GetNodejsVersionListJsonPath())
		if err != nil {
			return nil, err
		}
		nodejsConfig.VersionUrl.LastUpdateTime = file.ModTime().Unix()
	}
	return nodejsConfig, nil
}

func (s *ExtensionNodejsService) GetNodejsConfig() (*response2.NodejsConfigResponse, error) {
	configBody, err := util.ReadFileStringBody(s.GetNodejsConfigJsonPath())
	if err != nil {
		return nil, err
	}
	var nodejsConfig response2.NodejsConfigResponse
	err = util.JsonStrToStruct(configBody, &nodejsConfig)
	if err != nil {
		return nil, err
	}
	return &nodejsConfig, nil
}

// SetRegistrySources 设置Nodejs镜像源
func (s *ExtensionNodejsService) SetRegistrySources(name string) error {
	config, err := s.Config()
	if err != nil {
		return err
	}
	if !util.StrIsEmpty(config.RegistrySources.List[name]) {
		config.RegistrySources.Use = name
	} else {
		return errors.New("what are you doing? ")
	}

	configStr, err := util.StructToJsonStr(config)
	if err != nil {
		return err
	}
	err = util.WriteFile(s.GetNodejsConfigJsonPath(), []byte(configStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

// SetVersionUrl 设置Nodejs版本源
func (s *ExtensionNodejsService) SetVersionUrl(url string) error {
	config, err := s.Config()
	if err != nil {
		return err
	}
	config.VersionUrl.Use = url

	configStr, err := util.StructToJsonStr(config)
	if err != nil {
		return err
	}
	err = util.WriteFile(s.GetNodejsConfigJsonPath(), []byte(configStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

// VersionList 获取Nodejs版本列表
func (s *ExtensionNodejsService) VersionList() (versionList []*response2.NodejsVersion, err error) {
	versionBody := "{}"
	//判断是否有保存的本地版本列表
	if !util.PathExists(s.GetNodejsVersionListJsonPath()) {
		err = s.UpdateVersionList()
		if err != nil {
			return nil, err
		}
	}
	versionBody, err = util.ReadFileStringBody(s.GetNodejsVersionListJsonPath())
	if err != nil {
		return nil, err
	}
	var versionListTmp []*response2.NodejsVersion
	uName := s.GetVersionName()
	err = util.JsonStrToStruct(versionBody, &versionListTmp)
	if err != nil {
		return nil, err
	}
	glibcVersion := s.GetGlibcVersion()
	for _, version := range versionListTmp {
		if version.Files == nil {
			continue
		}
		for _, v := range version.Files {
			if uName == v {
				if glibcVersion <= 2.17 {
					nodeVersion, err := strconv.Atoi(strings.ReplaceAll(strings.Split(version.Version, ".")[0], "v", ""))
					if err != nil {
						global.Log.Debugf("nodejs version atoi err: %s,version:%s", err.Error(), v)
						continue
					}
					if nodeVersion >= 18 {
						continue
					}
				}
				//判断是否安装
				if util.PathExists(fmt.Sprintf("%s/nodejs/%s/bin/node", global.Config.System.ServerPath, version.Version)) {
					version.Install = true
				}
				if len(versionList) > 0 {
					if versionList[len(versionList)-1].Version == version.Version {
						continue
					}
				}

				//npm_version = self.get_module_version(v['version'],'npm')
				//if npm_version: v['npm'] = npm_version

				if version.Install {
					nodejsPrefix := fmt.Sprintf("%s/nodejs/%s/", global.Config.System.ServerPath, version.Version)
					etcPath := fmt.Sprintf("%s/etc", nodejsPrefix)
					if !util.PathExists(etcPath) {
						_ = os.MkdirAll(etcPath, 0755)
					}
					npmrcFile := fmt.Sprintf("%s/npmrc", etcPath)
					if !util.PathExists(npmrcFile) {
						wBody, err := s.GetNpmrcInfo(version.Version)
						if err != nil {
							return nil, err
						}
						err = util.WriteFile(npmrcFile, []byte(wBody), 0644)
						if err != nil {
							return nil, err
						}
					}
				}
				versionList = append(versionList, version)
			}
		}
	}
	// 排序
	sort.SliceStable(versionList, func(i, j int) bool {
		return versionList[i].Install && !versionList[j].Install
	})
	return
}

// UpdateVersionList 更新Nodejs版本列表
func (s *ExtensionNodejsService) UpdateVersionList() error {
	config, err := s.Config()
	if err != nil {
		return err
	}
	getBody, err := util.RequestGet(config.VersionUrl.Use)
	if err != nil {
		return err
	}
	err = util.WriteFile(s.GetNodejsVersionListJsonPath(), []byte(getBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Install 安装Nodejs
func (s *ExtensionNodejsService) Install(version string) error {
	config, err := s.Config()
	if err != nil {
		return err
	}
	if config.DownloadUrl != nil && util.StrIsEmpty(config.DownloadUrl.List[config.DownloadUrl.Use]) {
		return errors.New("nodejs download url is empty")
	}

	//检查是否在等待或者进行队列中
	taskName := "安装[Nodejs-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	downloadUrl := fmt.Sprintf("%s/%s/node-%s-%s.tar.gz", config.DownloadUrl.List[config.DownloadUrl.Use], version, version, s.GetVersionName())

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash install_lib.sh && cd %s && /bin/bash %s install %s "%s"`,
		global.Config.System.PanelPath+"/data/shell", s.GetShellPath(), InstallNodejsShellName, version, downloadUrl)
	err = AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

// Uninstall 卸载Nodejs
func (s *ExtensionNodejsService) Uninstall(version string) error {
	//检查是否在等待或者进行队列中
	taskName := "卸载[Nodejs-" + version + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": taskName}))
	}

	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash %s uninstall %s`, s.GetShellPath(), InstallNodejsShellName, version)
	err := AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultEnv 设置默认环境
func (s *ExtensionNodejsService) SetDefaultEnv(version string) error {
	//清空默认环境
	s.DeleteEnv()
	if !util.StrIsEmpty(version) {
		//设置默认环境
		nodejsPath := fmt.Sprintf("%s/nodejs/%s", global.Config.System.ServerPath, version)
		nodejsBinPath := fmt.Sprintf("%s/bin/node", nodejsPath)
		if !util.PathExists(nodejsBinPath) {
			global.Log.Errorf("SetDefaultEnv->nodejsBinPath version:%v,nodejsBinPath:%v", version, nodejsBinPath)
			return errors.New(helper.Message("extension.NotInstalled"))
		}
		srcNpmBin := fmt.Sprintf("%s/lib/node_modules/npm/bin/npm-cli.js", nodejsPath)
		srcNpxBin := fmt.Sprintf("%s/lib/node_modules/npm/bin/npx-cli.js", nodejsPath)
		srcPm2Bin := fmt.Sprintf("%s/lib/node_modules/pm2/bin/pm2", nodejsPath)
		srcYarnBin := fmt.Sprintf("%s/lib/node_modules/yarn/bin/yarn.js", nodejsPath)
		if util.PathExists(srcNpmBin) {
			_ = os.Symlink(srcNpmBin, envNpmBin)
		}
		if util.PathExists(srcNpxBin) {
			_ = os.Symlink(srcNpxBin, envNpxBin)
		}
		if util.PathExists(srcPm2Bin) {
			_ = os.Symlink(srcPm2Bin, envPm2Bin)
		}
		if util.PathExists(srcYarnBin) {
			_ = os.Symlink(srcYarnBin, envYarnBin)
		}
		_ = os.Symlink(nodejsBinPath, envNodeBin)
	}
	return nil
}

// NodeModulesList 获取Nodejs模块列表
func (s *ExtensionNodejsService) NodeModulesList(version string) (list []*response2.NodejsNodeModulesInfo, err error) {
	NodeModulesPath := fmt.Sprintf("%s/nodejs/%s/lib/node_modules", global.Config.System.ServerPath, version)
	if !util.PathExists(NodeModulesPath) {
		return nil, errors.New(fmt.Sprintf("not found node_modules path:%v", NodeModulesPath))
	}
	nodeModulesList, err := os.ReadDir(NodeModulesPath)
	if err != nil {
		return nil, err
	}
	for _, nodeModule := range nodeModulesList {
		packagePath := fmt.Sprintf("%s/%s/package.json", NodeModulesPath, nodeModule.Name())
		if nodeModule.IsDir() && util.PathExists(packagePath) {
			packageBody, err := util.ReadFileStringBody(packagePath)
			if err != nil {
				return nil, err
			}
			var nodeModulesInfo response2.NodejsNodeModulesInfo
			err = util.JsonStrToStruct(packageBody, &nodeModulesInfo)
			if err != nil {
				return nil, err
			}
			list = append(list, &nodeModulesInfo)
		}
	}
	return
}

// OperationNodeModules 操作Nodejs模块
func (s *ExtensionNodejsService) OperationNodeModules(version string, modules string, operation string) error {
	errorMsg := ""
	config, err := s.Config()
	if err != nil {
		return err
	}
	if config.RegistrySources != nil && util.StrIsEmpty(config.RegistrySources.List[config.RegistrySources.Use]) {
		return errors.New("registry sources is empty")
	}
	switch operation {
	case "install":
		errorMsg = helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": "node_modules"})
	case "uninstall":
		errorMsg = helper.MessageWithMap("queuetask.UninstallingOrWaiting", map[string]any{"Name": "node_modules"})
	case "upgrade":
		errorMsg = helper.MessageWithMap("queuetask.InstallingOrWaiting", map[string]any{"Name": "node_modules"})
	default:
		return errors.New("operation error")
	}
	//检查是否在等待或者进行队列中
	taskName := operation + "[node_modules-" + modules + "]"
	if exists, err := CheckTaskQueueExists(taskName); err != nil {
		return err
	} else if exists {
		return errors.New(errorMsg)
	}
	//添加面板队列任务
	execStr := fmt.Sprintf(`cd %s && /bin/bash %s %s %s "%s" "%s"`,
		s.GetShellPath(),
		InstallNodeModulesShellName,
		operation,
		version,
		modules,
		config.RegistrySources.List[config.RegistrySources.Use],
	)
	err = AddTaskQueue(taskName, execStr)
	if err != nil {
		return err
	}
	return nil
}

func (s *ExtensionNodejsService) DeleteEnv() {
	if util.PathExists(envNodeBin) {
		_ = os.Remove(envNodeBin)
	}
	if util.PathExists(envNpmBin) {
		_ = os.Remove(envNpmBin)
	}
	if util.PathExists(envNpxBin) {
		_ = os.Remove(envNpxBin)
	}
	if util.PathExists(envPm2Bin) {
		_ = os.Remove(envPm2Bin)
	}
	if util.PathExists(envYarnBin) {
		_ = os.Remove(envYarnBin)
	}
}

// GetNodejsVersion 获取Nodejs版本
func (s *ExtensionNodejsService) GetNodejsVersion() (version string, err error) {
	shell, err := util.ExecShell("node -v")
	if err != nil {
		return "", err
	}
	version = util.ClearStr(shell)
	return
}
func (s *ExtensionNodejsService) GetVersionName() string {
	uname := syscall.Utsname{}
	if err := syscall.Uname(&uname); err != nil {
		log.Println(err)
		return ""
	}

	machine, err := util.ExecShell("uname -m")
	if err != nil {
		return ""
	}

	machine = util.ClearStr(machine)

	switch machine {
	case "x86_64":
		return fmt.Sprintf("linux-x64")
	case "i686":
		return fmt.Sprintf("linux-x86")
	case "aarch64":
		return fmt.Sprintf("linux-arm64")
	case "armv7l":
		return fmt.Sprintf("linux-armv71")
	case "armv6l":
		return fmt.Sprintf("linux-armv61")
	case "armv5l":
		return fmt.Sprintf("linux-armv51")
	case "armv4l":
		return fmt.Sprintf("linux-armv41")
	case "armv3l":
		return fmt.Sprintf("linux-armv31")
	case "armv2l":
		return fmt.Sprintf("linux-armv21")
	case "mips":
		return fmt.Sprintf("linux-mips")
	case "mips64":
		return fmt.Sprintf("linux-mips64")
	case "ppc64":
		return fmt.Sprintf("linux-ppc64")
	case "ppc64le":
		return fmt.Sprintf("linux-ppc64le")
	case "s390x":
		return fmt.Sprintf("linux-s390x")
	case "sparc64":
		return fmt.Sprintf("linux-sparc64")
	case "sparc":
		return fmt.Sprintf("linux-sparc")
	default:
		return ""
	}
}

// GetShellPath 获取shell安装文件路径
func (s *ExtensionNodejsService) GetShellPath() string {
	return global.Config.System.PanelPath + "/data/extensions/nodejs/install"
}

// GetNpmrcInfo 获取指定版本的npmrc默认信息
func (s *ExtensionNodejsService) GetNpmrcInfo(version string) (string, error) {
	prefix := fmt.Sprintf("%s/nodejs/%s", global.Config.System.ServerPath, version)
	nodejsConfig, err := s.GetNodejsConfig()
	if err != nil {
		return "", err
	}
	registry := nodejsConfig.RegistrySources.List[nodejsConfig.RegistrySources.Use]
	cachePath := fmt.Sprintf("%s/nodejs/cache/", global.Config.System.ServerPath)
	if !util.PathExists(cachePath) {
		_ = os.MkdirAll(cachePath, 0755)
	}
	initModule := fmt.Sprintf("%s/etc/init-module.js", prefix)
	npmrcBody := fmt.Sprintf("prefix = %s \nregistry = %s \ncache = %s \ninit.module = %s \n", prefix, registry, cachePath, initModule)
	return npmrcBody, nil
}

// GetGlibcVersion 获取glibc版本号
func (s *ExtensionNodejsService) GetGlibcVersion() float64 {
	shell, err := util.ExecShell("ldd --version | awk 'NR==1{print $NF}'")
	if err != nil {
		global.Log.Debugf("获取glibc版本号失败: %s", err.Error())
		return 2.17
	}
	version := util.ClearStr(shell)
	f, err := strconv.ParseFloat(version, 64)
	if err != nil {
		global.Log.Debugf("转换glibc版本号失败:version-%s err-%s", version, err.Error())
		return 2.17
	}
	return f
}
