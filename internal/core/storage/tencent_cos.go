package storage

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

type TencentCosConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
}
type TencentCosStorage struct {
	Config TencentCosConfig
	client *cos.Client
}

func (s *TencentCosStorage) GetConfig() interface{} {
	return struct {
		TencentCosConfig TencentCosConfig `json:"tencentCos_config"`
	}{
		TencentCosConfig: s.Config,
	}
}

func NewTencentCos(config TencentCosConfig) (*TencentCosStorage, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	return &TencentCosStorage{Config: config, client: client}, nil
}
func (s *TencentCosStorage) NewBaseClient() *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%s.myqcloud.com", s.Config.Region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.Config.SecretID,
			SecretKey: s.Config.SecretKey,
		},
	})
	return client
}

func (s *TencentCosStorage) Upload(src, target string) error {
	_, _, err := s.client.Object.Upload(
		context.Background(), target, src, nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *TencentCosStorage) GetBucketList() (bucketList []string, err error) {
	client := s.NewBaseClient()
	serviceGetResult, _, err := client.Service.Get(context.Background())
	if err != nil {
		return nil, err
	}
	for _, bucket := range serviceGetResult.Buckets {
		bucketList = append(bucketList, bucket.Name)
	}
	return bucketList, nil
}

func (s *TencentCosStorage) Delete(path string) (bool, error) {
	if _, err := s.client.Object.Delete(context.Background(), path); err != nil {
		return false, err
	}
	return true, nil
}

func (s *TencentCosStorage) ObjectsList(query string) ([]interface{}, error) {
	objectList, _, err := s.client.Bucket.Get(context.Background(), &cos.BucketGetOptions{Prefix: query})
	if err != nil {
		return nil, err
	}

	var result []interface{}
	for _, item := range objectList.Contents {
		result = append(result, item.Key)
	}
	return result, nil
}
func (s *TencentCosStorage) IsExist(path string) (bool, error) {
	_, err := s.client.Object.Head(context.Background(), path, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TencentCosStorage) Download(src, target string) (bool, error) {
	_, err := s.client.Object.GetToFile(context.Background(), src, target, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
