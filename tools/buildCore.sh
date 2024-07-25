#!/bin/bash

if [ $DRONE_BRANCH ];then
release_version=$DRONE_BRANCH
else
release_version=$DRONE_TAG
fi

echo "编译版本号:"${release_version}

########################### MacOs
echo "macos-amd64 编译中"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-macos-amd64/xarr-rss .
echo "√ 编译MacOs-amd64"

echo "macos-arm64 编译中"
CGO_ENABLED=0  GOOS=darwin GOARCH=arm64 GOARM=7 go build -ldflags "-s -w -X 'main.Version=${DRONE_BRANCH}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-macos-arm64/xarr-rss .
echo "√ 编译MacOs-arm64"
########################### MacOs


########################### Linux
echo "编译linux-386 编译中"
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-linux-386/xarr-rss .
echo "√ 编译linux-386"

echo "编译linux-amd64 编译中"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-linux-amd64/xarr-rss .
echo "√ 编译linux-amd64"

echo "编译linux-arm 编译中"
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-linux-arm/xarr-rss .
echo "√ 编译linux-arm"

echo "编译linux-arm64 编译中"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-linux-arm64/xarr-rss .
echo "√ 编译linux-arm64"
########################### Linux


########################### windows
echo "编译Windows-i386 编译中"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-windows-i386/xarr-rss.exe .
echo "√ 编译Windows-i386"

echo "编译Windows-amd64 编译中"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-windows-amd64/xarr-rss.exe .
echo "√ 编译Windows-amd64"

echo "编译Windows-arm64 编译中"
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "-s -w -X 'main.Version=${release_version}'" -gcflags=-trimpath=$(pwd) -asmflags=-trimpath=$(pwd) -o ./release/xarr-rss-windows-arm64/xarr-rss.exe .
echo "√ 编译Windows-arm64"
########################### windows

echo "编译全部完成"