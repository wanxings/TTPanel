package storage

import (
	"context"
	"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/http"
	"os"
	"strings"
)

type MinioConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
}
type MinioStorage struct {
	Config MinioConfig
	client *minio.Client
}

func (m MinioStorage) GetConfig() interface{} {
	return struct {
		MinioConfig MinioConfig `json:"minio_config"`
	}{
		MinioConfig: m.Config,
	}
}

func NewMinio(config MinioConfig) (*MinioStorage, error) {
	secure := false
	tlsConfig := &tls.Config{}
	if strings.HasPrefix(config.Endpoint, "https") {
		secure = true
		tlsConfig.InsecureSkipVerify = true
	}
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(config.SecretID, config.SecretKey, ""),
		Secure:    secure,
		Transport: transport,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{
		Config: config,
		client: client,
	}, nil
}
func (m MinioStorage) Upload(src, target string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = m.client.PutObject(context.Background(), m.Config.Bucket, target, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}

func (m MinioStorage) GetBucketList() (bucketList []string, err error) {
	buckets, err := m.client.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}
	var result []string
	for _, bucket := range buckets {
		result = append(result, bucket.Name)
	}
	return result, err
}

func (m MinioStorage) ObjectsList(query string) ([]interface{}, error) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    query,
	}

	var result []interface{}
	for object := range m.client.ListObjects(context.Background(), m.Config.Bucket, opts) {
		if object.Err != nil {
			continue
		}
		result = append(result, object.Key)
	}
	return result, nil
}

func (m MinioStorage) Delete(path string) (bool, error) {
	object, err := m.client.GetObject(context.Background(), m.Config.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return false, err
	}
	info, err := object.Stat()
	if err != nil {
		return false, err
	}
	err = m.client.RemoveObject(context.Background(), m.Config.Bucket, path, minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        info.VersionID,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m MinioStorage) IsExist(path string) (bool, error) {
	_, err := m.client.GetObject(context.Background(), m.Config.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return true, err
	}
	return false, nil
}

func (m MinioStorage) Download(src, target string) (bool, error) {
	object, err := m.client.GetObject(context.Background(), m.Config.Bucket, src, minio.GetObjectOptions{})
	if err != nil {
		return false, err
	}
	localFile, err := os.Create(target)
	if err != nil {
		return false, err
	}
	if _, err = io.Copy(localFile, object); err != nil {
		return false, err
	}
	return true, nil
}
