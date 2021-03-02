package oss

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type FileUploader struct {
	Endpoint string

	session  *session.Session
	uploader *s3manager.Uploader
	s3       *s3.S3
}

func (f *FileUploader) InitDefault() {
	// The session the S3 Uploader will use
	f.session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(f.Endpoint),
	}))
	f.s3 = s3.New(f.session)

	// Create an uploader with the session and default options
	f.uploader = s3manager.NewUploader(f.session)
}

func (f *FileUploader) Upload(content io.Reader, bucket string, key string, cacheControl string) (*s3manager.UploadOutput, error) {
	result, err := f.uploader.Upload(&s3manager.UploadInput{
		Bucket:       &bucket,
		Key:          &key,
		Body:         content,
		CacheControl: &cacheControl,
	})
	return result, err
}

func (f *FileUploader) Remove(bucket string, key string) (*s3.DeleteObjectOutput, error) {
	result, err := f.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return result, err
}
