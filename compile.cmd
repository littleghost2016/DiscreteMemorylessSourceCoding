set GOARCH=amd64
set GOOS=windows
go build main.go
move main.exe dmsc_windows_amd64.exe

set GOOS=linux
go build main.go
move main dmsc_linux_amd64

set GOOS=darwin
go build main.go
move main dmsc_darwin_amd64