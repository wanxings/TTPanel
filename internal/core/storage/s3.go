package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3Config struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
}
type S3Storage struct {
	Config S3Config
	Sess   session.Session
}

func (s *S3Storage) GetConfig() interface{} {
	return struct {
		S3Config S3Config `json:"s3_config"`
	}{
		S3Config: s.Config,
	}
}

func NewS3Storage(config S3Config) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.SecretID, config.SecretKey, ""),
		Endpoint:         aws.String(config.Endpoint),
		Region:           aws.String(config.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	return &S3Storage{Config: config, Sess: *sess}, nil
}

func (s *S3Storage) Upload(src, target string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	uploader := s3manager.NewUploader(&s.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(target),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) GetBucketList() (bucketList []string, err error) {
	var result []string
	svc := s3.New(&s.Sess)
	res, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}
	for _, b := range res.Buckets {
		result = append(result, *b.Name)
	}
	return result, nil
}

func (s *S3Storage) Delete(path string) (bool, error) {
	svc := s3.New(&s.Sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.Config.Bucket), Key: aws.String(path)})
	if err != nil {
		return false, err
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *S3Storage) ObjectsList(query string) ([]interface{}, error) {
	svc := s3.New(&s.Sess)
	var result []interface{}
	if err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &s.Config.Bucket,
		Prefix: &query,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			result = append(result, *obj.Key)
		}
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *S3Storage) IsExist(path string) (bool, error) {
	svc := s3.New(&s.Sess)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *S3Storage) Download(src, target string) (bool, error) {
	file, err := os.Create(target)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	downloader := s3manager.NewDownloader(&s.Sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(src),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
