set GOOS=linux
set GOARCH=amd64
go build -buildmode=plugin -o=describe/describe.so  describe/describe.go