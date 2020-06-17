$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o dmsc_windows_amd64.exe

$env:GOOS="linux"
go build -o dmsc_linux_amd64

$env:GOOS="darwin"
go build -o dmsc_darwin_amd64