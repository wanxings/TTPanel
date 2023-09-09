package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOssConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
}
type AliOssStorage struct {
	Config AliOssConfig
	client *oss.Client
}

func (a *AliOssStorage) GetConfig() interface{} {
	return struct {
		AliOssStorage AliOssConfig `json:"aliOss_config"`
	}{
		AliOssStorage: a.Config,
	}
}

func NewAliOss(config AliOssConfig) (*AliOssStorage, error) {
	client, err := oss.New(config.Endpoint, config.SecretID, config.SecretKey)
	if err != nil {
		return nil, err
	}
	return &AliOssStorage{Config: config, client: client}, nil
}
func (a *AliOssStorage) Upload(src, target string) error {
	bucket, err := a.GetBucket()
	if err != nil {
		return err
	}
	err = bucket.UploadFile(target, src, 200*1024*1024, oss.Routines(5), oss.Checkpoint(true, ""))
	if err != nil {
		return err
	}
	return nil
}

func (a *AliOssStorage) GetBucketList() (bucketList []string, err error) {
	response, err := a.client.ListBuckets()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, bucket := range response.Buckets {
		result = append(result, bucket.Name)
	}
	return result, err
}

func (a *AliOssStorage) ObjectsList(query string) ([]interface{}, error) {
	bucket, err := a.GetBucket()
	if err != nil {
		return nil, err
	}
	lor, err := bucket.ListObjects(oss.Prefix(query))
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, obj := range lor.Objects {
		result = append(result, obj.Key)
	}
	return result, nil
}

func (a *AliOssStorage) Delete(path string) (bool, error) {
	bucket, err := a.GetBucket()
	if err != nil {
		return false, err
	}
	err = bucket.DeleteObject(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *AliOssStorage) IsExist(path string) (bool, error) {
	bucket, err := a.GetBucket()
	if err != nil {
		return false, err
	}
	return bucket.IsObjectExist(path)
}

func (a *AliOssStorage) Download(src, target string) (bool, error) {
	bucket, err := a.GetBucket()
	if err != nil {
		return false, err
	}
	err = bucket.DownloadFile(src, target, 200*1024*1024, oss.Routines(5), oss.Checkpoint(true, ""))
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *AliOssStorage) GetBucket() (*oss.Bucket, error) {
	bucket, err := a.client.Bucket(a.Config.Bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
