set GOROOT=c:\go
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0
set GOEXPERIMENT=noframepointer4
go build -ldflags "-s -w" -o golang-ots-linux main.go
pause