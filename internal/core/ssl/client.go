package ssl

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"crypto"
	"errors"
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/acme/api"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"github.com/go-acme/lego/v4/providers/http/webroot"
	"github.com/go-acme/lego/v4/registration"
	"time"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}

func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

type AcmeClient struct {
	Config *lego.Config
	Client *lego.Client
	User   *AcmeUser
}

func NewAcmeClient(email, privateKey string) (*AcmeClient, error) {
	if email == "" {
		return nil, errors.New("email can not blank")
	}
	if privateKey == "" {
		client, err := NewRegisterClient(email)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else {
		client, err := NewPrivateKeyClient(email, privateKey)
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}

func (c *AcmeClient) UseDns(dnsType string, Authorization map[string]string) error {
	var p challenge.Provider
	var err error
	if dnsType == constant.DnsTypeTencentCloud {
		tencentcloudConfig := tencentcloud.NewDefaultConfig()
		tencentcloudConfig.SecretID = Authorization["secret_id"]
		tencentcloudConfig.SecretKey = Authorization["secret_key"]
		tencentcloudConfig.Region = Authorization["region"]
		tencentcloudConfig.SessionToken = Authorization["session_token"]
		p, err = tencentcloud.NewDNSProviderConfig(tencentcloudConfig)
		if err != nil {
			return err
		}

	}
	if dnsType == constant.DnsTypeAliyun {
		alidnsConfig := alidns.NewDefaultConfig()
		alidnsConfig.SecretKey = Authorization["secret_key"]
		alidnsConfig.APIKey = Authorization["api_key"]
		p, err = alidns.NewDNSProviderConfig(alidnsConfig)
		if err != nil {
			return err
		}

	}
	if dnsType == constant.DnsTypeCloudflare {
		cloudflareConfig := cloudflare.NewDefaultConfig()
		global.Log.Debugf("UseDns->Authorization:%v \n", Authorization)
		cloudflareConfig.AuthToken = Authorization["auth_token"]
		cloudflareConfig.ZoneToken = Authorization["auth_token"]
		global.Log.Debugf("UseDns->cloudflareConfig.AuthToken:%v \n", cloudflareConfig.AuthToken)
		p, err = cloudflare.NewDNSProviderConfig(cloudflareConfig)
		if err != nil {
			return err
		}
	}

	return c.Client.Challenge.SetDNS01Provider(p, dns01.AddDNSTimeout(1*time.Minute))
}

func (c *AcmeClient) UseManualDns() error {
	p := &manualDnsProvider{}
	if err := c.Client.Challenge.SetDNS01Provider(p, dns01.AddDNSTimeout(3*time.Minute)); err != nil {
		return err
	}
	return nil
}

func (c *AcmeClient) UseHTTP(path string) error {
	httpProvider, err := webroot.NewHTTPProvider(path)
	if err != nil {
		return err
	}

	err = c.Client.Challenge.SetHTTP01Provider(httpProvider)
	if err != nil {
		return err
	}
	return nil
}

func (c *AcmeClient) ObtainSSL(domains []string) (certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := c.Client.Certificate.Obtain(request)
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

func (c *AcmeClient) RenewSSL(certUrl string) (certificate.Resource, error) {
	certificates, err := c.Client.Certificate.Get(certUrl, true)
	if err != nil {
		return certificate.Resource{}, err
	}
	certificates, err = c.Client.Certificate.Renew(*certificates, true, true, "")
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

type Resolve struct {
	Key   string
	Value string
	Err   string
}

type manualDnsProvider struct {
	Resolve *Resolve
}

func (p *manualDnsProvider) Present(domain, token, keyAuth string) error {
	return nil
}

func (p *manualDnsProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}

func (c *AcmeClient) GetDNSResolve(domains []string) (map[string]Resolve, error) {
	core, err := api.New(c.Config.HTTPClient, c.Config.UserAgent, c.Config.CADirURL, c.User.Registration.URI, c.User.Key)
	if err != nil {
		return nil, err
	}
	order, err := core.Orders.New(domains)
	if err != nil {
		return nil, err
	}
	resolves := make(map[string]Resolve)
	resC, errC := make(chan acme.Authorization), make(chan domainError)
	for _, authorization := range order.Authorizations {
		go func(authorization string) {
			authZ, err := core.Authorizations.Get(authorization)
			if err != nil {
				errC <- domainError{Domain: authZ.Identifier.Value, Error: err}
				return
			}
			resC <- authZ
		}(authorization)
	}

	var responses []acme.Authorization
	for i := 0; i < len(order.Authorizations); i++ {
		select {
		case res := <-resC:
			responses = append(responses, res)
		case err := <-errC:
			resolves[err.Domain] = Resolve{Err: err.Error.Error()}
		}
	}
	close(resC)
	close(errC)

	for _, auth := range responses {
		domain := challenge.GetTargetedDomain(auth)
		ChallengeN, err := challenge.FindChallenge(challenge.DNS01, auth)
		if err != nil {
			resolves[domain] = Resolve{Err: err.Error()}
			continue
		}
		keyAuth, err := core.GetKeyAuthorization(ChallengeN.Token)
		if err != nil {
			resolves[domain] = Resolve{Err: err.Error()}
			continue
		}
		challengeInfo := dns01.GetChallengeInfo(domain, keyAuth)
		resolves[domain] = Resolve{
			Key:   challengeInfo.FQDN,
			Value: challengeInfo.Value,
		}
	}

	return resolves, nil
}
