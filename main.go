package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-resty/resty/v2"
)

// var AWS_ACCESS_KEY_ID string = "xxxxxxxxxxxxxxxxxxxxxxxx"
// var AWS_SECRET_ACCESS_KEY string = "xxxxxxxxxxxxxxxxxxxxx+xxxxxxxxxxxxxxx"
// 액세스키에 기호가 있을경우 아래와 같은 오류 발생, 기호가 없을때까지 재발급;;
//SignatureDoesNotMatch: The request signature we calculated does not match the signature you provided. Check your key and signing method.

var AWS_REGION string = "AWS_REGION"
var AWS_ACCESS_KEY_ID string = "AWS_ACCESS_KEY_ID"
var AWS_SECRET_ACCESS_KEY string = "AWS_SECRET_ACCESS_KEY"
var BUCKET_NAME string = "bucketName"

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	var filename string = "kakao_code.json"
	alarmMsg := AlarmMsg{}
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		json.Unmarshal([]byte(snsRecord.Message), &alarmMsg)
		fmt.Println("alarmMsg 확인")
		fmt.Println(alarmMsg)
	}
	sess := ConnectAws()
	S3download(sess, filename)
	header := TokenGet()
	fmt.Println("header : ", header)

	uuidList := UuidGet(header)
	fmt.Println(uuidList)
	uuid := ConvertUUID(uuidList)
	fmt.Println(uuid)
	SendMsg(header, uuid, &alarmMsg)
}

func main() {
	lambda.Start(handler)
}

func SendMsg(header string, uuid string, alarmMsg *AlarmMsg) {
	fmt.Println("=========SendMsg() Start==========")

	name := fmt.Sprint("AlarmName:  ", alarmMsg.AlarmName)
	descr := fmt.Sprint("Description:   ", alarmMsg.AlarmDescription)
	region := fmt.Sprint("Region:       ", alarmMsg.Region)
	id := fmt.Sprint("AccountId:     ", alarmMsg.AWSAccountId)
	metric := fmt.Sprint("MetricName:  ", alarmMsg.Trigger.MetricName)
	instance := fmt.Sprint("Instanceid:     ", alarmMsg.Trigger.Dimensions[0].Value)
	fmt.Println(name)
	fmt.Println(descr)
	fmt.Println(region)
	fmt.Println(id)
	fmt.Println(metric)
	fmt.Println(instance)

	textStr := fmt.Sprintf(`template_object={"object_type": "text","text": "%s %s %s %s %s %s","link": {"web_url": "https://developers.kakao.com", "mobile_web_url": "https://developers.kakao.com"},"button_title": "바로 확인"}`, name, descr, region, id, metric, instance)
	send_url := "https://kapi.kakao.com/v1/api/talk/friends/message/default/send"
	fmt.Println("textStr: ", textStr)
	token := fmt.Sprintf("Authorization: %s", header)
	fmt.Println(token)
	c := exec.Command("curl", "-v", "-X", "POST", send_url, "-H", "Content-Type: application/x-www-form-urlencoded", "-H", token, "--data-urlencode", uuid, "--data-urlencode", textStr)
	c.Stdout = os.Stdout
	c.Run()
	fmt.Println()
	fmt.Println("=========SendMsg() END GOOD==========")
}

func ConvertUUID(uuidList []string) string {
	list := []string{}
	for _, item := range uuidList {
		id := fmt.Sprintf("\"%s\"", item)
		list = append(list, id)
		fmt.Println(id)
	}
	uuids := strings.Join(list, ",")
	uuid := fmt.Sprintf("receiver_uuids=[%s]", uuids)
	return uuid
}

func UuidGet(header string) []string {
	fmt.Println("=========UuidGet() Start=========")
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", header).
		Get("https://kapi.kakao.com/v1/api/talk/friends")
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
	}
	friends_list := FriendsList{}
	uuids_list := []string{}
	json.Unmarshal(resp.Body(), &friends_list)
	for _, list := range friends_list.Elements {
		uuids_list = append(uuids_list, list.UUID)
	}
	fmt.Println(uuids_list)
	fmt.Println("=========UuidGet() END===========")

	return uuids_list
}

func TokenGet() string {
	fmt.Println("=========TokenGet() Start=========")
	filepath := "/tmp/kakao_code.json"
	json_data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	kakaoToken := KakaoToken{}
	if err = json.Unmarshal(json_data, &kakaoToken); err != nil {
		log.Fatal(err)
	}

	header := "Bearer " + kakaoToken.Access_token
	fmt.Println("=========TokenGet() END===========")
	return header
}

func S3download(sess *session.Session, filename string) {
	fmt.Println("s3download 시작")
	var filepath string = "/tmp/" + filename
	// s3 다운로드 받을때 에러 발생, read-only file system
	// 람다에서는 오직 /tmp 에서만 파일을 작성 할 수 있다고 한다.

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})
	if err != nil {
		fmt.Println("error ", err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	fmt.Println("s3 다운로드 끝")

}

func ConnectAws() *session.Session {
	fmt.Println("세션 연결 시작")

	sess, err := session.NewSession(
		&aws.Config{
			Region: &AWS_REGION,
			Credentials: credentials.NewStaticCredentials(
				AWS_ACCESS_KEY_ID,
				AWS_SECRET_ACCESS_KEY,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("세션 연결 끝")

	return sess
}

type KakaoToken struct {
	Access_token             string `json:"access_token"`
	Expires_in               int64  `json:"expires_in"`
	Refresh_token            string `json:"refresh_token"`
	Refresh_token_expires_in int64  `json:"refresh_token_expires_in"`
	Scope                    string `json:"scope"`
	Token_type               string `json:"token_type"`
}

type FriendsList struct {
	Elements []struct {
		ProfileNickname       string `json:"profile_nickname,omitempty"`
		ProfileThumbnailImage string `json:"profile_thumbnail_image"`
		ID                    int    `json:"id,omitempty"`
		UUID                  string `json:"uuid,omitempty"`
		Favorite              bool   `json:"favorite,omitempty"`
	} `json:"elements"`
	TotalCount    int         `json:"total_count,omitempty"`
	AfterURL      interface{} `json:"after_url,omitempty"`
	FavoriteCount int         `json:"favorite_count,omitempty"`
}

type AlarmMsg struct {
	AlarmName        string  `json:"AlarmName"`
	AlarmDescription string  `json:"AlarmDescription"`
	AWSAccountId     string  `json:"AWSAccountId"`
	NewStateValue    string  `json:"NewStateValue"`
	NewStateReason   string  `json:"NewStateReason"`
	StateChangeTime  string  `json:"StateChangeTime"`
	Region           string  `json:"Region"`
	AlarmArn         string  `json:"AlarmArn"`
	OldStateValue    string  `json:"OldStateValue"`
	Trigger          Trigger `json:"Trigger"`
}

type Trigger struct {
	MetricName                       string       `json:"MetricName"`
	Namespace                        string       `json:"Namespace"`
	StatisticType                    string       `json:"StatisticType"`
	Statistic                        string       `json:"Statistic"`
	Unit                             string       `json:"Unit,omitempty"`
	Dimensions                       []Dimensions `json:"Dimensions"`
	Period                           int64        `json:"Period"`
	EvaluationPeriods                int64        `json:"EvaluationPeriods"`
	ComparisonOperator               string       `json:"ComparisonOperator"`
	Threshold                        float64      `json:"Threshold"`
	TreatMissingData                 string       `json:"TreatMissingData,omitempty"`
	EvaluateLowSampleCountPercentile string       `json:"EvaluateLowSampleCountPercentile,omitempty"`
}

type Dimensions struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}
