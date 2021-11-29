package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MYBucket string

func OauthGet() {
	fmt.Println("=========OauthGet() Start=========")
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {"XXXXXXXXXXXXXXXXXXXXXXX"},
		"redirect_url": {"your_redirect_url"},
		"code":         {"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"},
	}

	resp, err := http.PostForm("https://kauth.kakao.com/oauth/token", data)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	b, _ := json.MarshalIndent(res, "", "	")
	fmt.Println(string(b))
	err = ioutil.WriteFile("kakao_code.json", b, 0644)
	// JSON 파일로 저장
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=========OauthGet() END===========")
}

func main() {
	var filename string = "kakao_code.json"

	LoadEnv()

	sess := ConnectAws()
	OauthGet()
	S3upload(sess, filename)
}

func ConnectAws() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")
	MYBucket = GetEnvWithKey("BUCKET_NAME")
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
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func S3upload(sess *session.Session, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	svc := s3manager.NewUploader(sess)
	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MYBucket),
		Key:    aws.String(filename),
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
