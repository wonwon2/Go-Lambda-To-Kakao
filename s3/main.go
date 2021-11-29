package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var MYBucket string
var AccessKeyID string
var SecretAccessKey string
var MyRegion string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
	sess := ConnectAws()
	S3upload(sess, "../function.zip")
}

func S3upload(sess *session.Session, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	svc := s3manager.NewUploader(sess)
	fmt.Println("Uploading file to S3...")
	filePath := path.Base(filename)
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MYBucket),
		Key:    aws.String(filePath),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)
	fmt.Println("result : ", result)
	fmt.Println("result Loaction : ", result.Location)
	fmt.Println("result ETag : ", result.ETag)
	fmt.Println("result UploadID: ", result.UploadID)
	fmt.Println("result VersionID: ", result.VersionID)

}

func ConnectAws() *session.Session {
	MYBucket = os.Getenv("BUCKET_NAME")
	AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion = os.Getenv("AWS_REGION")
	sess, err := session.NewSession(
		&aws.Config{
			Region: &MyRegion,
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}
