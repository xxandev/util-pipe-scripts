app_name := pipe

default: B

R:
	cd cmd && go run .

B: 
	cd ./cmd && go build -o "../build/bin/$(app_name)"
	cd ./cmd && GOOS=linux GOARCH=arm GOARM=6 go build -o "../build/bin/linux-arm6/$(app_name)"
	cd ./cmd && GOOS=linux GOARCH=arm GOARM=7 go build -o "../build/bin/linux-arm7/$(app_name)"
	cd ./cmd && GOOS=linux GOARCH=arm64 go build -o "../build/bin/linux-arm8/$(app_name)"
	cd ./cmd && GOOS=linux GOARCH=amd64 go build -o "../build/bin/linux-amd64/$(app_name)"
	cd ./cmd && GOOS=linux GOARCH=386 go build -o "../build/bin/linux-x386/$(app_name)"
	cd ./cmd && GOOS=windows GOARCH=amd64 go build -o "../build/bin/windows-x64/$(app_name).exe"
	cd ./cmd && GOOS=windows GOARCH=386 go build -o "../build/bin/windows-x32/$(app_name).exe"
	cd ./cmd && GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o "../build/bin/windows-x64i/$(app_name).exe"
	cd ./cmd && GOOS=windows GOARCH=386 go build -ldflags "-H windowsgui" -o "../build/bin/windows-x32i/$(app_name).exe"
