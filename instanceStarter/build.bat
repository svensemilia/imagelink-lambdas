
set GOOS=linux
go build -o output/main main.go
build-lambda-zip -output output/main.zip output/main
cmd /k