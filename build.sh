GOARCH=386 CGO_ENABLED=1 go build -o release/secExt.dll -buildmode=c-shared .
GOARCH=amd64 go build -o release/secExt_x64.dll -buildmode=c-shared .
