package oss

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type FileUploader struct {
	Endpoint string

	session  *session.Session
	uploader *s3manager.Uploader
}

func (f *FileUploader) InitDefault() {
	// The session the S3 Uploader will use
	f.session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(f.Endpoint),
	}))

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
