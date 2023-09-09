package response

import "time"

type DNSResolves struct {
	Key    string `json:"resolve"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Err    string `json:"err"`
}
type DnsTypeConfig struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Help        string `json:"help"`
	Form        []struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"form"`
}

type DnsAccount struct {
	Type          string            `json:"type"`
	Authorization map[string]string `json:"authorization"`
}

type AcmeAccount struct {
	Url        string `json:"url"`
	PrivateKey string `json:"private_key"`
}

type DNSTencentCloudConfig struct {
	SecretID     string `json:"secret_id"`
	SecretKey    string `json:"secret_key"`
	Region       string `json:"region"`
	SessionToken string `json:"session_token"`
}

type SSLDetails struct {
	DNSAccount     string    `json:"dns_account"`
	AcmeAccount    string    `json:"acme_account"`
	VerifyMode     int       `json:"verify_mode"`
	Domains        []string  `json:"domains"`
	CertURL        string    `json:"cert_url"`
	ExpireDate     time.Time `json:"expire_date"`
	StartDate      time.Time `json:"start_date"`
	Type           string    `json:"type"`
	Organization   string    `json:"organization"`
	AutoRenew      bool      `json:"auto_renew"`
	AlwaysUseHttps bool      `json:"always_use_https"`
	Key            string    `json:"key"`
	Csr            string    `json:"csr"`
}
