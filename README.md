# **UTIL Pipe**

Utility pipe - the utility was created to run bash scripts from root for any user

## **Installation**
- install [golang](https://go.dev/) 1.18+
- go get github.com/xxandev/util-pipe
- cd ..../util-pipe
- make B or go build .

## **Run**
- run pipe from root
```
.../pipe -host ':8091' -scripts-path '/root/.scripts' -wiki-path '/root/.wiki'
```

- call script
```
curl http://localhost:8091/script-exec?script=test.sh
curl http://localhost:8091/script-exec?script=test.sh&params=tester,text-text
curl http://localhost:8091/script-exec?script=test.sh&params=$USER,text-text
curl --connect-timeout 30 http://localhost:8091/script-exec?script=test.sh
curl --max-time 30 curl http://localhost:8091/script-exec?script=test.sh
```