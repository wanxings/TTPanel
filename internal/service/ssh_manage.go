package service

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SSHManageService struct{}

var SSHConfigPath = "/etc/ssh/sshd_config"

func (s *SSHManageService) GetSSHInfo() (map[string]any, error) {
	var data = make(map[string]any)
	//获取ssh状态
	data["status"] = s.GetSSHStatus()

	//获取ssh端口号
	data["port"] = 22
	sshConfigBody, err := util.ReadFileStringBody(SSHConfigPath)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(?m)^\s*Port\s+([0-9]+)`)
	portFind := re.FindStringSubmatch(sshConfigBody)
	if len(portFind) > 1 {
		data["port"] = portFind[1]
	}
	//获取其他配置
	data["rsa_authentication"] = "no"
	r := regexp.MustCompile(`(?m)^\s*RSAAuthentication\s*(.*)`)
	rsaFind := r.FindStringSubmatch(sshConfigBody)
	fmt.Printf("%v", rsaFind)
	if len(rsaFind) > 1 {
		data["rsa_authentication"] = rsaFind[1]
	}

	data["pubkey_authentication"] = "no"
	if true {
		//r = regexp.MustCompile(`^\s*PubkeyAuthentication\s*(yes|no)`)
		r = regexp.MustCompile(`(?m)^\s*PubkeyAuthentication\s*(.*)`)
		pubkeyFind := r.FindStringSubmatch(sshConfigBody)
		fmt.Printf("%v", pubkeyFind)
		if len(pubkeyFind) > 1 {
			data["pubkey_authentication"] = pubkeyFind[1]
		}
	}

	data["password_authentication"] = "no"
	r = regexp.MustCompile(`(?m)^\s*PasswordAuthentication\s*(.*)`)
	sshPasswordFind := r.FindStringSubmatch(sshConfigBody)
	if len(sshPasswordFind) > 1 {
		data["password_authentication"] = sshPasswordFind[1]
	}

	// 是否允许root登录，默认允许
	//yes = 允许
	//no = 不允许
	//without-password = 允许，但不允许使用密码登录
	//forced-commands-only = 允许，但只允许执行命令，不能使用终端
	data["root_is_login"] = "yes"
	data["root_login_type"] = "yes"
	r = regexp.MustCompile(`(?m)^\s*PermitRootLogin\s*(.*)`)
	rootIsLoginFind := r.FindStringSubmatch(sshConfigBody)
	if len(rootIsLoginFind) > 1 {
		data["root_is_login"] = rootIsLoginFind[1]
		if data["root_is_login"] != "yes" {
			data["root_login_type"] = data["root_is_login"]
		}
	}
	return data, nil
}

func (s *SSHManageService) GetSSHStatus() bool {
	shell, err := util.ExecShellScript("if pgrep sshd &> /dev/null; then echo \"have\"; fi")
	if err == nil && strings.Contains(shell, "have") {
		return true
	} else {
		return false
	}
}

func (s *SSHManageService) SetSSHStatus(action string) error {
	shell, _ := util.ExecShell(fmt.Sprintf("bash %s/data/shell/set_ssh_status.sh %s", global.Config.System.PanelPath, action))
	if strings.Contains(shell, "successfully") {
		return nil
	}
	return fmt.Errorf(" Error: %s", shell)
}

func (s *SSHManageService) OperateSSHKeyLogin(action, keyType string) (string, error) {
	key := ""
	if action == "on" {
		authorizedKeys := "/root/.ssh/authorized_keys"
		files := []string{
			fmt.Sprintf("/root/.ssh/id_%s.pub", keyType),
			fmt.Sprintf("/root/.ssh/id_%s", keyType),
		}

		//尝试删除旧密钥
		for _, file := range files {
			if util.PathExists(file) {
				if err := os.Remove(file); err != nil {
					return key, err
				}
			}
		}
		if _, err := util.ExecShellScript(fmt.Sprintf("ssh-keygen -t %s -P '' -f /root/.ssh/id_%s |echo y", keyType, keyType)); err != nil {
			return key, err
		}

		if !util.PathExists(files[0]) {
			return key, errors.New("failed to generate key")
		}

		_, err := util.ExecShellScript(fmt.Sprintf("cat %s > %s && chmod 600 %s", files[0], authorizedKeys, authorizedKeys))
		if err != nil {
			return key, err
		}

		sshConfigBody, err := util.ReadFileStringBody(SSHConfigPath)
		if err != nil {
			return key, err
		}

		rec := regexp.MustCompile(`\n#?RSAAuthentication\s\w+`)
		rec2 := regexp.MustCompile(`\n#?PubkeyAuthentication\s\w+`)

		if len(rec.FindStringIndex(sshConfigBody)) == 0 {
			sshConfigBody += "\nRSAAuthentication yes"
		}

		if len(rec2.FindStringIndex(sshConfigBody)) == 0 {
			sshConfigBody += "\nPubkeyAuthentication yes"
		}

		fileSSH := rec.ReplaceAllString(sshConfigBody, "\nRSAAuthentication yes")
		fileResult := rec2.ReplaceAllString(fileSSH, "\nPubkeyAuthentication yes")

		//fmt.Printf("SSHConfig:%v", fileResult)
		err = util.WriteFile(SSHConfigPath, []byte(fileResult), 0644)
		if err != nil {
			return key, err
		}

		err = restartSSH()
		if err != nil {
			return key, err
		}
		key, err = util.ReadFileStringBody(files[1])
		if err != nil {
			return key, err
		}
	} else {
		sshStatus := s.GetSSHStatus()

		rec := `\n\s*#?\s*RSAAuthentication\s+\w+`
		rec2 := `\n\s*#?\s*PubkeyAuthentication\s+\w+`

		SSHConfigBody, err := util.ReadFileStringBody(SSHConfigPath)
		if err != nil {
			return key, err
		}
		fileSSH := regexp.MustCompile(rec).ReplaceAllString(SSHConfigBody, "\nRSAAuthentication no")
		fileResult := regexp.MustCompile(rec2).ReplaceAllString(fileSSH, "\nPubkeyAuthentication no")

		err = util.WriteFile(SSHConfigPath, []byte(fileResult), 0644)
		if err != nil {
			return key, err
		}
		//fmt.Printf("SSHConfig:%v", fileResult)

		if sshStatus {
			err = restartSSH()
			if err != nil {
				return key, err
			}

		}
	}
	return key, nil
}

// OperatePasswordLogin 操作密码登录
func (s *SSHManageService) OperatePasswordLogin(action string) error {
	//读取配置文件
	sshConfigBody, err := util.ReadFileStringBody(SSHConfigPath)
	if err != nil {
		return err
	}
	newSSHConfigBody := sshConfigBody
	if action == "on" {
		sshPassword := `\n#?PasswordAuthentication\s\w+`
		if len(regexp.MustCompile(sshPassword).FindAll([]byte(sshConfigBody), -1)) == 0 {
			newSSHConfigBody = sshConfigBody + "\nPasswordAuthentication yes"
		} else {
			newSSHConfigBody = regexp.MustCompile(sshPassword).ReplaceAllString(sshConfigBody, "\nPasswordAuthentication yes")
		}
	} else {
		sshPasswordPattern := `(?m)^\s*PasswordAuthentication\s+\w+`
		replacement := "PasswordAuthentication no"
		newSSHConfigBody = regexp.MustCompile(sshPasswordPattern).ReplaceAllString(sshConfigBody, replacement)
	}
	err = util.WriteFile(SSHConfigPath, []byte(newSSHConfigBody), 0644)
	if err != nil {
		return err
	}
	return restartSSH()
}

// GetSSHLoginStatistics 获取SSH登录统计
func (s *SSHManageService) GetSSHLoginStatistics(refresh bool) (map[string]any, error) {
	//先判断有无缓存
	failCount, ok := global.GoCache.Get("ssh_login_fail_count")
	if !ok {
		refresh = true
	}
	trueCount, ok := global.GoCache.Get("ssh_login_true_count")
	if !ok {
		//无缓存
		refresh = true
	}

	if refresh {
		dirList, err := os.ReadDir("/var/log")
		if err != nil {
			return nil, err
		}
		var checkList []string
		for _, v := range dirList {
			if v.IsDir() {
				continue
			}
			if v.Name() == "secure" || strings.Contains(v.Name(), "auth.log") {
				fmt.Printf("%v", v.Name())
				checkList = append(checkList, "/var/log/"+v.Name())
			}
		}
		var failCountTotal, trueCountTotal int
		for _, logPath := range checkList {
			failTotalStr, err := util.ExecShellScript(fmt.Sprintf("cat %s |grep 'Failed password' |wc -l", logPath))
			if err != nil {
				return nil, err
			}
			//转成int
			total, _ := strconv.Atoi(util.ClearStr(failTotalStr))
			failCountTotal += total

			trueTotalStr, err := util.ExecShellScript(fmt.Sprintf("cat %s |grep 'Accepted password' |wc -l", logPath))
			if err != nil {
				return nil, err
			}
			//转成int
			total, _ = strconv.Atoi(util.ClearStr(trueTotalStr))
			trueCountTotal += total
		}
		//设置缓存
		global.GoCache.Set("ssh_login_fail_count", failCountTotal, -1)
		global.GoCache.Set("ssh_login_true_count", trueCountTotal, -1)
		failCount = failCountTotal
		trueCount = trueCountTotal
	} else {
		failCount = failCount.(int)
		trueCount = trueCount.(int)
	}
	return map[string]any{
		"fail_count": failCount,
		"true_count": trueCount,
	}, nil
}

func restartSSH() error {
	if util.PathExists("/etc/redhat-release") {
		if _, err := util.ExecShell("systemctl restart sshd"); err != nil {
			_, _ = util.ExecShell("/etc/init.d/sshd restart")
		}
	} else {
		_, _ = util.ExecShell("service ssh restart")
	}
	return nil
}
