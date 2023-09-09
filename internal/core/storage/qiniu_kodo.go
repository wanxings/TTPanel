package storage

import (
	"TTPanel/pkg/util"
	"context"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	kodoStorage "github.com/qiniu/go-sdk/v7/storage"
	"time"
)

type QiniuKodoConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}
type QiniuKodoStorage struct {
	Config QiniuKodoConfig
	client *kodoStorage.BucketManager
}

func (q *QiniuKodoStorage) GetConfig() interface{} {
	return struct {
		QiniuKodoConfig QiniuKodoConfig `json:"qiniuKodo_config"`
	}{
		QiniuKodoConfig: q.Config,
	}
}

func NewQiniuKodo(config QiniuKodoConfig) (*QiniuKodoStorage, error) {
	conn := qbox.NewMac(config.SecretID, config.SecretKey)
	cfg := kodoStorage.Config{
		UseHTTPS: false,
	}
	bucketManager := kodoStorage.NewBucketManager(conn, &cfg)
	return &QiniuKodoStorage{Config: config, client: bucketManager}, nil
}
func (q *QiniuKodoStorage) Upload(src, target string) error {
	putPolicy := kodoStorage.PutPolicy{
		Scope: q.Config.Bucket,
	}
	mac := qbox.NewMac(q.Config.SecretID, q.Config.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := kodoStorage.Config{UseHTTPS: true, UseCdnDomains: false}
	resumeUploader := kodoStorage.NewResumeUploaderV2(&cfg)
	ret := kodoStorage.PutRet{}
	putExtra := kodoStorage.RputV2Extra{}
	if err := resumeUploader.PutFile(context.Background(), &ret, upToken, target, src, &putExtra); err != nil {
		return err
	}
	return nil
}

func (q *QiniuKodoStorage) GetBucketList() (bucketList []string, err error) {
	buckets, err := q.client.Buckets(true)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, bucket := range buckets {
		list = append(list, bucket)
	}
	return list, nil
}

func (q *QiniuKodoStorage) ObjectsList(query string) ([]interface{}, error) {
	var result []interface{}
	marker := ""
	for {
		entries, _, nextMarker, hashNext, err := q.client.ListFiles(q.Config.Bucket, query, "", marker, 1000)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			result = append(result, entry.Key)
		}
		if hashNext {
			marker = nextMarker
		} else {
			break
		}
	}
	return result, nil
}

func (q *QiniuKodoStorage) Delete(path string) (bool, error) {
	if err := q.client.Delete(q.Config.Bucket, path); err != nil {
		return false, err
	}
	return true, nil
}

func (q *QiniuKodoStorage) IsExist(path string) (bool, error) {
	if _, err := q.client.Stat(q.Config.Bucket, path); err != nil {
		return true, err
	}
	return false, nil
}

func (q *QiniuKodoStorage) Download(src, target string) (bool, error) {
	mac := auth.New(q.Config.SecretID, q.Config.SecretKey)
	deadline := time.Now().Add(time.Second * 3600).Unix()
	privateAccessURL := kodoStorage.MakePrivateURL(mac, q.Config.Domain, src, deadline)
	_, err := util.DownloadFile(target, privateAccessURL, false)
	if err != nil {
		return false, err
	}
	return true, nil
}
