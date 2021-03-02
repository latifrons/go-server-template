package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func Download() error {
	//url := "https://atom8.s3.us-east-2.amazonaws.com/s1"
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast2RegionID),
	}))
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	filename := "down.mod"
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String("atom8"),
		Key:    aws.String("s1"),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
	return nil
}

func Upload() error {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast2RegionID),
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	filename := "go.mod"

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	cache := "public, max-age=86400"

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String("file.894569.site"),
		Key:          aws.String("s2"),
		Body:         f,
		CacheControl: &cache,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(BreakerWrapper)

	v1 := r.Group("/v1")
	{
		v1.POST("/panic", Panic)
	}

	return r
}

func Panic(c *gin.Context) {
	time.Sleep(10 * time.Second) // 模拟接口超时
	c.JSON(200, gin.H{"message": "this is panic"})
}

func BreakerWrapper(c *gin.Context) {
	name := c.Request.Method + "-" + c.Request.RequestURI
	hystrix.Do(name, func() error {
		c.Next()

		statusCode := c.Writer.Status()

		if statusCode >= http.StatusInternalServerError {
			str := fmt.Sprintf("status code %d", statusCode)
			return errors.New(str)
		}

		return nil
	}, func(e error) error {
		if e == hystrix.ErrCircuitOpen {
			c.String(http.StatusAccepted, "请稍后重试")
		}

		return e
	})
}

func main() {

	//r := router()
	//r.Run()

	err := Upload()
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
	}
	//err = Download()
	//if err != nil {
	//	fmt.Printf("failed to download file, %v", err)
	//}
}
