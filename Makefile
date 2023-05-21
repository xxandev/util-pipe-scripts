app_name := pipe-scripts

default: run

# run for current os
run:
	cd ./cmd/$(app_name) && go run .

# build for current os
build: 
	cd ./cmd/$(app_name) && go build -o "../../build/bin/$(app_name)"

# build for linux debian x86, arch arm6
arm6:
	cd ./cmd/$(app_name) && GOOS=linux GOARCH=arm GOARM=6 go build -o "../../build/bin/linux-arm6/$(app_name)"

# build for linux debian x32, arch arm7
arm7:
	cd ./cmd/$(app_name) && GOOS=linux GOARCH=arm GOARM=7 go build -o "../../build/bin/linux-arm7/$(app_name)"

# build for linux debian x64, arch arm8(arm64)
arm8:
	cd ./cmd/$(app_name) && GOOS=linux GOARCH=arm64 go build -o "../../build/bin/linux-arm8/$(app_name)"

# build for linux debian x64, arch amd64
linux64:
	cd ./cmd/$(app_name) && GOOS=linux GOARCH=amd64 go build -o "../../build/bin/linux-amd64/$(app_name)"

# build for linux debian x32, arch 386
linux32:
	cd ./cmd/$(app_name) && GOOS=linux GOARCH=386 go build -o "../../build/bin/linux-x386/$(app_name)"

# build for windows x64, arch amd64
win64:
	cd ./cmd/$(app_name) && GOOS=windows GOARCH=amd64 go build -o "../../build/bin/windows-x64/$(app_name).exe"

# build for windows x32, arch 386
win32:
	cd ./cmd/$(app_name) && GOOS=windows GOARCH=386 go build -o "../../build/bin/windows-x32/$(app_name).exe"

# build for windows x64, arch amd64, invisible start
win64i:
	cd ./cmd/$(app_name) && GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o "../../build/bin/windows-x64i/$(app_name).exe"

# build for windows x32, arch 386, invisible start
win32i:
	cd ./cmd/$(app_name) && GOOS=windows GOARCH=386 go build -ldflags "-H windowsgui" -o "../../build/bin/windows-x32i/$(app_name).exe"

all: arm6 arm7 arm8 linux64 linux32 win64 win32 win64i win32i