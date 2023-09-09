package util

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ExecShell 执行单个shell命令，返回执行结果（无法执行多个命令）
func ExecShell(cmdStr string) (string, error) {
	if !strings.HasPrefix(cmdStr, "sudo") {
		cmdStr = "sudo " + cmdStr
	}
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errMsg := "Shell: " + cmdStr + ";"
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		fmt.Printf("ExecShell-errMsg-> %s\n", errMsg)
		return errMsg, fmt.Errorf(errMsg)
	}
	return stdout.String(), nil
}
func ExecShellScript(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errMsg := "Shell: " + cmdStr + ";"
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		err = errors.New(fmt.Sprintf("%s,%s", err.Error(), errMsg))
		fmt.Printf("ExecShellScript-errMsg-> %s\n", errMsg)
		return errMsg, err
	}
	return stdout.String(), nil
}

// ExecShellScriptS 无返回值
func ExecShellScriptS(cmdStr string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	return cmd.Run()
}
func ExecShellAsUser(cmdStr string, user string) (string, error) {
	cmd := exec.Command("su", user, "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found in") {
			return false
		}
		panic(err)
	}

	return true
}
func ShellQuote(cmd string) string {
	var buf bytes.Buffer
	buf.WriteByte('\'')
	for _, r := range cmd {
		switch r {
		case '\'', '\\':
			buf.WriteByte('\\')
			buf.WriteRune(r)
		default:
			if strconv.IsPrint(r) {
				buf.WriteRune(r)
			} else {
				buf.WriteString(`\x`)
				buf.WriteString(strconv.FormatInt(int64(r), 16))
			}
		}
	}
	buf.WriteByte('\'')
	return buf.String()
}
