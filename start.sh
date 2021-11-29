rm function.zip
GOOS=linux CGO_ENABLED=0 go build main.go
zip function.zip main
cd s3
./s3up