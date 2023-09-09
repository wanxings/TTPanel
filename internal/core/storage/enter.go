package storage

import (
	"TTPanel/internal/helper/constant"
	"fmt"
)

type Storage interface {
	GetConfig() interface{}
	Upload(src, target string) error
	GetBucketList() (bucketList []string, err error)
	ObjectsList(query string) ([]interface{}, error)
	Delete(path string) (bool, error)
	IsExist(path string) (bool, error)
	Download(src, target string) (bool, error)
}

func New(storageType int, config interface{}) (Storage, error) {
	var err error
	switch storageType {
	case constant.StorageCategoryByS3:
		awsConfig, ok := config.(S3Config)
		if !ok {
			err = fmt.Errorf("invalid aws config")
			break
		}
		return NewS3Storage(awsConfig)
	case constant.StorageCategoryByTencentCOS:
		tencentCosConfig, ok := config.(TencentCosConfig)
		if !ok {
			err = fmt.Errorf("invalid tencentCos config")
			break
		}
		return NewTencentCos(tencentCosConfig)
	case constant.StorageCategoryByAliOSS:
		aliOssConfig, ok := config.(AliOssConfig)
		if !ok {
			err = fmt.Errorf("invalid aliOss config")
			break
		}
		return NewAliOss(aliOssConfig)
	case constant.StorageCategoryByQiniuKodo:
		qiniuKodoConfig, ok := config.(QiniuKodoConfig)
		if !ok {
			err = fmt.Errorf("invalid qiniuKodo config")
			break
		}
		return NewQiniuKodo(qiniuKodoConfig)
	case constant.StorageCategoryByMinio:
		minioConfig, ok := config.(MinioConfig)
		if !ok {
			err = fmt.Errorf("invalid minio config")
			break
		}
		return NewMinio(minioConfig)
	default:
		err = fmt.Errorf("invalid notify type")
	}
	return nil, err
}
