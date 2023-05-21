package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"util-pipe-scripts/internal/utils"
)

var config Config

func init() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("[PIPE-SCRIPTS] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	mute := flag.Bool("mute", false, "mute log")
	flag.StringVar(&config.Host, "host", ":8090", "server host, example: 127.0.0.1:8090, :8090...")
	flag.StringVar(&config.Path, "scripts-path", "", "scripts path")
	flag.Parse()
	if *mute {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		arg := r.URL.Query().Get("script")
		script, err := filepath.Abs(filepath.Join(config.Path, arg))
		if err != nil {
			io.WriteString(w, fmt.Sprintf("%v\n", err))
			return
		}
		log.Printf("call script: %s\n", script)
		if !utils.IsFile(script) {
			io.WriteString(w, "the script does not exist or not a script\n")
			return
		}
		resp, err := utils.ExecCommand(script)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("%v\n", err))
			return
		}
		io.WriteString(w, resp+"\n")
	})
	log.Fatalf("server run: %v\n", http.ListenAndServe(config.Host, nil))
}
