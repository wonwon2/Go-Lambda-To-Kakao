# Go-Lambda-To-Kakao
CloudWatch Alarm 발생을, Lambda로 거쳐 카카오톡으로 전송 

(만든지 얼마 안되서 미흡한 부분이 많으므로 양해 바랍니다. 그래도 갓 때린 팽이마냥 잘 돌아갑니다.)
(현재 액세스키 문제로 사용불가 ㅠㅠ)

AWS EC2(Amazon Linux)에서 람다 함수를 압축한 후에 S3에 전송 후 람다에 올립니다.

![구성도](https://user-images.githubusercontent.com/60952823/143803082-7c68a8c6-2539-429c-8ed7-9461a13ec39c.png)
![image](https://user-images.githubusercontent.com/60952823/143803003-e17c340a-7850-4086-86ae-2b6798fed6c2.png)


사용 순서
1. https://kauth.kakao.com/oauth/authorize?client_id={REST_API_KEY}&redirect_uri={REDIRECT_URI}&response_type=code&scope=talk_message,friends 
 : 위의 양식에 맞춰 값을 넣고 주소창에 넣고 엔터 후 ex)xxxcode=abcedf 와 같이 값이 나옴, code= 이부분 이후부터 복사
2. Kakao_auth에 코드안에 code라는 변수안에 값을 넣고 실행 ( 사용자 토큰 가져오기 그리고 S3에 저장 )
3. start.sh를 실행하고 콘솔에 s3주소를 복사 후 -> 람다 업로드 버튼 클릭 -> Amazon S3 위치 클릭 -> 복사한 주소를 아래에 붙여넣기 
 ![image](https://user-images.githubusercontent.com/60952823/143810535-30066ac4-61c9-4d89-ba59-42949bf08014.png)

