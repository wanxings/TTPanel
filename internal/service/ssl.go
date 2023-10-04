package service

import (
	"TTPanel/internal/core/ssl"
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"time"
)

type SSLService struct {
}

func (s *SSLService) DnsTypeList() ([]*response.DnsTypeConfig, error) {
	dnsInitPath := global.Config.System.PanelPath + "/data/lego/dns_init.json"
	jsonBody, err := util.ReadFileStringBody(dnsInitPath)
	if err != nil {
		return nil, err
	}
	var list []*response.DnsTypeConfig
	err = util.JsonStrToStruct(jsonBody, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *SSLService) DnsAccountList() (list map[string]*response.DnsAccount, err error) {
	dnsAccountPath := global.Config.System.PanelPath + "/data/lego/dns_account.json"
	if !util.PathExists(dnsAccountPath) {
		_ = util.WriteFile(dnsAccountPath, []byte("{}"), 644)
	}
	jsonBody, err := util.ReadFileStringBody(dnsAccountPath)
	if err != nil {
		return nil, err
	}
	err = util.JsonStrToStruct(jsonBody, &list)
	if err != nil {
		return nil, err
	}
	return
}

func (s *SSLService) CreateDnsAccount(name string, dnsType string, authorization map[string]string) error {
	dnsAccountPath := global.Config.System.PanelPath + "/data/lego/dns_account.json"
	jsonBody, err := util.ReadFileStringBody(dnsAccountPath)
	if err != nil {
		return err
	}
	var accountList map[string]*response.DnsAccount
	err = util.JsonStrToStruct(jsonBody, &accountList)
	if err != nil {
		return err
	}
	if _, ok := accountList[name]; ok {
		return errors.New("account already exists")
	}
	var account response.DnsAccount
	account.Type = dnsType
	account.Authorization = authorization
	accountList[name] = &account
	newJsonBody, err := util.StructToJsonStr(accountList)
	if err != nil {
		return err
	}
	err = util.WriteFile(dnsAccountPath, []byte(newJsonBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *SSLService) EditDnsAccount(name string, authorization map[string]string) error {
	dnsAccountPath := global.Config.System.PanelPath + "/data/lego/dns_account.json"
	jsonBody, err := util.ReadFileStringBody(dnsAccountPath)
	if err != nil {
		return err
	}
	var accountList map[string]*response.DnsAccount
	err = util.JsonStrToStruct(jsonBody, &accountList)
	if err != nil {
		return err
	}
	if _, ok := accountList[name]; !ok {
		return errors.New("account does not exist")
	}
	accountList[name].Authorization = authorization
	newJsonBody, err := util.StructToJsonStr(accountList)
	if err != nil {
		return err
	}
	err = util.WriteFile(dnsAccountPath, []byte(newJsonBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *SSLService) CreateAcmeAccount(email string) error {
	accountList, err := s.AcmeAccountList()
	if err != nil {
		return err
	}

	if _, ok := accountList[email]; ok {
		return errors.New("account already exists")
	}

	client, err := ssl.NewAcmeClient(email, "")
	if err != nil {
		return err
	}

	acmeAccount := response.AcmeAccount{
		Url:        client.User.Registration.URI,
		PrivateKey: string(ssl.GetPrivateKey(client.User.GetPrivateKey())),
	}

	accountList[email] = &acmeAccount

	return s.SaveAcmeAccountList(accountList)
}

func (s *SSLService) SaveAcmeAccountList(accountList map[string]*response.AcmeAccount) error {
	acmeAccountPath := global.Config.System.PanelPath + "/data/lego/acme_account.json"
	newJsonBody, err := util.StructToJsonStr(accountList)
	if err != nil {
		return err
	}
	err = util.WriteFile(acmeAccountPath, []byte(newJsonBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *SSLService) AcmeAccountList() (accountList map[string]*response.AcmeAccount, err error) {
	acmeAccountPath := global.Config.System.PanelPath + "/data/lego/acme_account.json"
	if !util.PathExists(acmeAccountPath) {
		_ = util.WriteFile(acmeAccountPath, []byte("{}"), 644)
	}
	jsonBody, err := util.ReadFileStringBody(acmeAccountPath)
	if err != nil {
		return
	}
	err = util.JsonStrToStruct(jsonBody, &accountList)
	if err != nil {
		return
	}
	return
}

func (s *SSLService) DeleteAcmeAccount(email string) (err error) {
	accountList, err := s.AcmeAccountList()
	if err != nil {
		return err
	}
	delete(accountList, email)
	return s.SaveAcmeAccountList(accountList)
}

func (s *SSLService) ApplyCertificate(params *request.ApplyCertificateR) (logPath string, err error) {
	logPath = fmt.Sprintf("%s/ssl/%d.log", global.Config.Logger.RootPath, time.Now().Unix())
	_ = os.MkdirAll(path.Dir(logPath), 0755)
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	acmeAccountList, err := s.AcmeAccountList()
	if err != nil {
		return
	}
	if _, ok := acmeAccountList[params.AcmeAccount]; !ok {
		err = errors.New("acme account does not exist")
		return
	}
	client, err := ssl.NewPrivateKeyClient(params.AcmeAccount, acmeAccountList[params.AcmeAccount].PrivateKey)
	if err != nil {
		return
	}

	var projectData *model.Project

	if params.ProjectId > 0 {
		projectData, err = (&model.Project{ID: params.ProjectId}).Get(global.PanelDB)
		if err != nil {
			return
		}
		if projectData.ID == 0 {
			err = errors.New("project not found")
			return
		}
	}

	_, _ = logFile.WriteString("|-Use AcmeAccount: " + params.AcmeAccount + "\n")
	switch params.VerifyMode {
	case constant.SSLVerifyModeByDNSAccount:
		dnsAccountList, err := s.DnsAccountList()
		if err != nil {
			return logPath, err
		}
		if _, ok := dnsAccountList[params.DnsAccount]; !ok {
			err = errors.New("dns account does not exist")
			return logPath, err
		}
		if err := client.UseDns(dnsAccountList[params.DnsAccount].Type, dnsAccountList[params.DnsAccount].Authorization); err != nil {
			return logPath, err
		}
	case constant.SSLVerifyModeByManual:
		if err := client.UseManualDns(); err != nil {
			return logPath, err
		}
	case constant.SslVerifyModeByFile:
		if projectData == nil {
			return logPath, errors.New("must specify a project")
		}
		rootPath, err := GetProjectRootPath(projectData.Name)
		if err != nil {
			return logPath, err
		}
		if err = client.UseHTTP(rootPath); err != nil {
			return logPath, err
		}
	default:
		return logPath, errors.New("verification mode does not exist")
	}

	go func() {
		defer func(logFile *os.File) {
			_ = logFile.Close()

		}(logFile)

		_, _ = logFile.WriteString("|-Start applying for a certificate,may take a few minutes...\n")
		resource, err := client.ObtainSSL(params.Domains)
		if err != nil {
			_, _ = logFile.WriteString("\nERROR：" + err.Error() + "\n")
			return
		}

		cert, err := util.ParseCert(resource.Certificate)
		if err != nil {
			return
		}
		sslDetails := response.SSLDetails{
			AutoRenew:    true,
			DNSAccount:   params.DnsAccount,
			AcmeAccount:  params.AcmeAccount,
			Domains:      params.Domains,
			CertURL:      resource.CertURL,
			ExpireDate:   cert.NotAfter,
			StartDate:    cert.NotBefore,
			Type:         cert.Issuer.CommonName,
			Organization: cert.Issuer.Organization[0],
		}
		//申请成功，保存证书和证书信息
		savePath := fmt.Sprintf("%s/data/ssl/%s", global.Config.System.PanelPath, cert.Subject.CommonName)
		_ = os.MkdirAll(savePath, 0755)
		err = util.WriteFile(savePath+"/fullchain.pem", resource.Certificate, 0644)
		if err != nil {
			_, _ = logFile.WriteString("\nERROR：" + err.Error() + "\n")
			return
		}
		err = util.WriteFile(savePath+"/private.pem", resource.PrivateKey, 0644)
		if err != nil {
			_, _ = logFile.WriteString("\nERROR：" + err.Error() + "\n")
			return
		}
		sslDetailsStr, err := util.StructToJsonStr(sslDetails)
		if err != nil {
			_, _ = logFile.WriteString("\nERROR：" + err.Error() + "\n")
			return
		}
		err = util.WriteFile(savePath+"/info.json", []byte(sslDetailsStr), 0644)
		if err != nil {
			_, _ = logFile.WriteString("\nERROR：" + err.Error() + "\n")
			return
		}
		_, _ = logFile.WriteString("\nSuccessful application\n")
		//设置项目
		if params.SetUp {
			if projectData == nil {
				err = errors.New("must specify a project")
				return
			}
			err = GenerateSslConfig(projectData.Name, resource.PrivateKey, resource.Certificate)
			if err != nil {
				_, _ = logFile.WriteString("\n|--ERROR:" + err.Error() + "\n")
				return
			}
			_, _ = logFile.WriteString("\n|--successful deployment to " + projectData.Name + "\n")
		}
		_, _ = logFile.WriteString("\n---------------Over------------------------\n")
	}()
	return
}

// GetResolve 获取域名解析
func (s *SSLService) GetResolve(acmeEmail string, domains []string) (resolve []response.DNSResolves, err error) {
	acmeAccountList, err := s.AcmeAccountList()
	if err != nil {
		return
	}
	if _, ok := acmeAccountList[acmeEmail]; !ok {
		return nil, errors.New("acme account does not exist")
	}
	client, err := ssl.NewPrivateKeyClient(acmeEmail, acmeAccountList[acmeEmail].PrivateKey)
	if err != nil {
		return
	}
	resolves, err := client.GetDNSResolve(domains)
	if err != nil {
		return
	}
	for k, v := range resolves {
		resolve = append(resolve, response.DNSResolves{
			Domain: k,
			Key:    v.Key,
			Value:  v.Value,
			Err:    v.Err,
		})
	}
	return
}

// CheckDnsRecords 检查解析记录
func (s *SSLService) CheckDnsRecords(nameList []string) (checkList map[string]any) {
	checkList = make(map[string]any)
	for _, v := range nameList {
		records, err := net.LookupTXT(v)
		fmt.Printf("name:%s,records:%v", v, records)
		if err != nil {
			checkList[v] = err.Error()
		} else {
			checkList[v] = records
		}
	}
	return
}

// CertList 证书列表
func (s *SSLService) CertList() (certList []response.SSLDetails, err error) {
	certRootPath := fmt.Sprintf("%s/data/ssl", global.Config.System.PanelPath)
	certDirs, err := os.ReadDir(certRootPath)
	if err != nil {
		return nil, err
	}
	for _, certDir := range certDirs {
		if certDir.IsDir() {
			if util.PathExists(certRootPath + "/" + certDir.Name() + "/info.json") {
				certInfo, err := util.ReadFileStringBody(certRootPath + "/" + certDir.Name() + "/info.json")
				if err != nil {
					return nil, err
				}
				var cert response.SSLDetails
				err = util.JsonStrToStruct(certInfo, &cert)
				if err != nil {
					return nil, err
				}
				certList = append(certList, cert)
			}
		}
	}
	return
}
