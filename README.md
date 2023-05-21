# **UTIL Pipe-Scripts**

Utility pipe-script - the utility was created to run bash scripts from root for any user

## **Installation**
- install [golang](https://go.dev/) 1.18+
- go get github.com/xxandev/util-pipe-script
- cd ..../util-pipe-script
- make [ build | arm6 | arm7 | arm8 | linux64 | linux32 | win64 | win32 | win64i | win32i ] or go build .

## **Run**
- run pipe-scripts from root
```
.../pipe-scripts -host ':10080' -scripts-path '/root/scripts'
```

- call script
```
curl http://localhost:10080?script=test.sh
```
OR
```
curl --connect-timeout 30 http://localhost:10080?script=test.sh
```
OR
```
curl --max-time 30 curl http://localhost:10080?script=test.sh
```