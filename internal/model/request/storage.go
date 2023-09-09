package request

import "TTPanel/internal/core/storage"

type AddStorageR struct {
	Category    int    `json:"category" form:"category" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" `
	StorageConfigR
}

type StorageConfigR struct {
	AliOssConfig     storage.AliOssConfig     `json:"aliOss_config"`
	MinioConfig      storage.MinioConfig      `json:"minio_config"`
	QiniuKodoConfig  storage.QiniuKodoConfig  `json:"qiniuKodo_config"`
	S3Config         storage.S3Config         `json:"s3_config"`
	TencentCosConfig storage.TencentCosConfig `json:"tencentCos_config"`
}

//type TencentCOSConfigR struct {
//	TencentCOSSecretId  string `json:"tencentCOS_secret_id"`
//	TencentCOSSecretKey string `json:"tencentCOS_secret_key"`
//	TencentCOSRegion    string `json:"tencentCOS_region"`
//	TencentCOSBucket    string `json:"tencentCOS_bucket"`
//}
//
//type S3ConfigR struct {
//	S3SecretId  string `json:"s3_secret_id"`
//	S3SecretKey string `json:"s3_secret_key"`
//	S3Region    string `json:"s3_region"`
//	S3Endpoint  string `json:"s3_endpoint"`
//	S3Bucket    string `json:"s3_bucket"`
//}
//
//type AliOSSConfigR struct {
//	AliOSSSecretId  string `json:"aliOSS_secret_id"`
//	AliOSSSecretKey string `json:"aliOSS_secret_key"`
//	AliOSSRegion    string `json:"aliOSS_region"`
//	AliOSSEndpoint  string `json:"aliOSS_endpoint"`
//	AliOSSBucket    string `json:"aliOSS_bucket"`
//}
//
//type MinioConfigR struct {
//	MinioSecretId  string `json:"minio_secret_id"`
//	MinioSecretKey string `json:"minio_secret_key"`
//	MinioEndpoint  string `json:"minio_endpoint"`
//	MinioBucket    string `json:"minio_bucket"`
//}
//
//type QiniuKodoConfigR struct {
//	QiniuKodoSecretId  string `json:"qiniuKodo_secret_id"`
//	QiniuKodoSecretKey string `json:"qiniuKodo_secret_key"`
//	QiniuKodoRegion    string `json:"qiniuKodo_region"`
//	QiniuKodoEndpoint  string `json:"qiniuKodo_endpoint"`
//	QiniuKodoBucket    string `json:"qiniuKodo_bucket"`
//}

type StorageBucketListR struct {
	Category int `json:"category" form:"category" binding:"required"`
	StorageConfigR
}

type StorageListR struct {
	Query    string `json:"query" form:"query"`
	Category int    `json:"category" form:"category"`
	Limit    int    `json:"limit" form:"limit" binding:"required"`
	Page     int    `json:"page" form:"page" binding:"required"`
}

type EditStorageR struct {
	ID          int64  `json:"id" form:"id" binding:"required"`
	Category    int    `json:"category" form:"category" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" `
	StorageConfigR
}

type EditLocalStorageR struct {
	Path string `json:"path" form:"path" binding:"required"`
}
