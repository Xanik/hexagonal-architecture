package files

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"study/config"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFile(key string, req *http.Request) string {
	_, handler, err := req.FormFile(key)

	if err != nil {
		fmt.Println(err)
	}

	go func(head multipart.FileHeader, key string) {
		file, err := head.Open()
		if err != nil {
			log.Println("Error while opening file: ", err)
			return
		}
		size := head.Size
		buffer := make([]byte, size)
		file.Read(buffer)

		s, err := session.NewSession(&aws.Config{
			Region: aws.String(config.Env.S3REGION),
		})
		if err != nil {
			log.Println("Error while logging into s3: ", err)
			return
		}
		_, err = s3.New(s, &aws.Config{MaxRetries: aws.Int(3)}).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(config.Env.S3BUCKET),
			Key:                  aws.String(handler.Filename),
			ACL:                  aws.String("public-read"),
			Body:                 bytes.NewReader(buffer),
			ContentLength:        aws.Int64(int64(size)),
			ContentType:          aws.String(http.DetectContentType(buffer)),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
		})
		log.Println(err)
	}(*handler, key)

	// Create Url To Send Back
	url := fmt.Sprint("https://", config.Env.S3BUCKET, ".s3.amazonaws.com/", handler.Filename)

	return url
}

//PresignedURL generates a temporary url thats accessible to the user
func PresignedURL(url string) string {
	s3REGION := config.Env.S3REGION
	s3BUCKET := config.Env.S3BUCKET

	//new s3 session
	s, err := session.NewSession(&aws.Config{Region: aws.String(s3REGION)})
	// Create S3 service client
	svc := s3.New(s)
	// Config settings: this is where you choose the bucket, filename, content-type etc.
	request, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3BUCKET),
		Key:    aws.String(url),
	})
	urlStr, err := request.Presign(time.Hour * 168)

	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println("The URL is", urlStr)

	return urlStr
}
