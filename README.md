# Go-Lambda-To-Kakao
CloudWatch Alarm 발생을, Lambda로 이용하여 카카오톡으로 발생한 알람 전송 (Go 언어로 작성)

(만든지 얼마 안되서 미흡한 부분이 많으므로 양해 바랍니다. 그래도 갓 때린 팽이마냥 일단 잘 돌아갑니다.)


AWS EC2(Amazon Linux)에서 람다 함수를 압축한 후에 S3에 전송 후 람다에 올립니다.

![구성도](https://user-images.githubusercontent.com/60952823/143803082-7c68a8c6-2539-429c-8ed7-9461a13ec39c.png)
![image](https://user-images.githubusercontent.com/60952823/143803003-e17c340a-7850-4086-86ae-2b6798fed6c2.png)

사용 순서 Kakao_auth  -> start.sh
