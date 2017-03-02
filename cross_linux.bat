set GOROOT=c:\go
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0
set GOEXPERIMENT=noframepointer4
go build -ldflags "-s -w" -o golang-ots-linux main.go
copy golang-ots-linux C:\Works\goWorks\src\github.com\blemobi\go-platform-manager\views
copy run_linux.sh C:\Works\goWorks\src\github.com\blemobi\go-platform-manager\views
pause