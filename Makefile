default: build

all: build build-darwin build-windows build-linux

build:
	mkdir -p bin
	go build -o bin/sproxy-mgmt-mgmt

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o sproxy-mgmt_`tdp tag current`_linux_arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sproxy-mgmt_`tdp tag current`_linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o sproxy-mgmt_`tdp tag current`_linux_386

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o sproxy-mgmt_`tdp tag current`_darwin_amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o sproxy-mgmt_`tdp tag current`_darwin_arm64

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o sproxy-mgmt_`tdp tag current`_windows_amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o sproxy-mgmt_`tdp tag current`_windows_386.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o sproxy-mgmt_`tdp tag current`_windows_arm64.exe

clean:
	rm -rf sproxy-mgmt_*.* bin