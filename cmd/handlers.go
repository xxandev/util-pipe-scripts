package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"util-pipe/internal/dbg"
	"util-pipe/internal/utils"
	"util-pipe/internal/xj"
)

func hGetStr2MD5(w http.ResponseWriter, r *http.Request) {
	type tmp struct {
		Hex string `json:"md5_hash_hex,omitempty"`
		B64 string `json:"md5_hash_base64,omitempty"`
	}
	w.Header().Add("Content-Type", "application/json")
	if value := r.URL.Query().Get("value"); len(value) > 0 {
		hash := md5.Sum([]byte(value))
		if err := xj.Parser.NewEncoder(w).Encode(tmp{Hex: hex.EncodeToString(hash[:]), B64: base64.StdEncoding.EncodeToString(hash[:])}); err != nil {
			w.Write(xj.Errf("encode json - %v", err).Bts())
		}
		return
	}
	w.Write(xj.Err("invalid value").Bts())
}

func hGetStatistic(w http.ResponseWriter, r *http.Request) {
	statistic.Upgrade()
	w.Header().Add("Content-Type", "application/json")
	if err := xj.Parser.NewEncoder(w).Encode(statistic); err != nil {
		w.Write(xj.Errf("encode json - %v", err).Bts())
	}
}

func hGetScriptExec(w http.ResponseWriter, r *http.Request) {
	script := filepath.Join(config.ScriptsPath, r.URL.Query().Get("script"))
	w.Header().Add("Content-Type", "application/json")
	if !utils.IsFile(script) {
		w.Write(xj.Err("script not found or not file").Bts())
		return
	}
	cmd := utils.GenCMD(script, strings.Split(r.URL.Query().Get("params"), ",")...)
	dbg.Log.Println("call script:", cmd.String())
	resp, err := utils.ExecCommand(cmd)
	if err != nil {
		w.Write(xj.Errf("exec command - %v", err).Bts())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resp))
}

func hPostScriptExec(w http.ResponseWriter, r *http.Request) {
	type ScrExTmp struct {
		Script string   `json:"script,omitempty" yaml:"script,omitempty"`
		Params []string `json:"params,omitempty" yaml:"params,omitempty"`
	}
	w.Header().Add("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(xj.Errf("body read - %v", err).Bts())
		return
	}
	var tmp ScrExTmp
	if err := xj.Parser.Unmarshal(body, &tmp); err != nil {
		w.Write(xj.Errf("unmarshal json - %v", err).Bts())
		return
	}
	script := filepath.Join(config.ScriptsPath, tmp.Script)
	if !utils.IsFile(script) {
		w.Write(xj.Err("script not found or not file").Bts())
		return
	}
	cmd := utils.GenCMD(script, tmp.Params...)
	dbg.Log.Println("call script:", cmd.String())
	resp, err := utils.ExecCommand(cmd)
	if err != nil {
		w.Write(xj.Errf("exec command - %v", err).Bts())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resp))
}

func hGetJsonRead(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("filepath")
	w.Header().Add("Content-Type", "application/json")
	content, err := os.ReadFile(path)
	if err != nil {
		w.Write(xj.Errf("read file json - %v", err).Bts())
		return
	}
	var value any
	if err := xj.Parser.Unmarshal(content, &value); err != nil {
		w.Write(xj.Errf("unmarshal json - %v", err).Bts())
		return
	}
	if err := xj.Parser.NewEncoder(w).Encode(value); err != nil {
		w.Write(xj.Errf("encode json - %v", err).Bts())
	}
}

func hPostJsonCreate(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("filepath")
	w.Header().Add("Content-Type", "application/json")
	if strings.ToLower(filepath.Ext(path)) != ".json" {
		w.Write(xj.Err("filepath - file is not json").Bts())
		return
	}
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			w.Write(xj.Errf("remove file - %v", err).Bts())
			return
		}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(xj.Errf("body read - %v", err).Bts())
		return
	}
	var value any
	if err := xj.Parser.Unmarshal(body, &value); err != nil {
		w.Write(xj.Errf("unmarshal json - %v", err).Bts())
		return
	}
	content, err := xj.Parser.MarshalIndent(value, "", "    ")
	if err != nil {
		w.Write(xj.Errf("marshal json - %v", err).Bts())
		return
	}
	file, err := os.Create(path)
	if err != nil {
		w.Write(xj.Errf("create file - %v", err).Bts())
		return
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		w.Write(xj.Errf("write file - %v", err).Bts())
		return
	}
	if err := file.Close(); err != nil {
		w.Write(xj.Errf("close file - %v", err).Bts())
		return
	}
	w.Write(xj.Succes("create json end").Bts())
}

func hGetCheckAuth(w http.ResponseWriter, r *http.Request) {
	// ldapwhoami -x -w "******" -D "uid=tester,ou=users,dc=ldap,dc=example,dc=com" -H "ldap://ldap.example.com" > /dev/null 2>&1;result=$?;echo $result
	w.Header().Add("Content-Type", "application/json")
	w.Write(xj.Succes("auth successful").Bts())
}

func hGetGottyKill(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if runtime.GOOS != "linux" {
		w.Write(xj.Err("this function is intended for linux os only").Bts())
		return
	}
	w.Write(xj.Info("func no").Bts())
}

func hGetGottyRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if runtime.GOOS != "linux" {
		w.Write(xj.Err("this function is intended for linux os only").Bts())
		return
	}
	expath, err := os.Executable()
	if err != nil {
		w.Write(xj.Errf("get path executable file - %v", err).Bts())
		return
	}
	path := filepath.Join(filepath.Dir(expath), "gotty")
	if !utils.IsFile(path) {
		w.Write(xj.Errf("get path executable file - %v", err).Bts())
		return
	}
	if err := exec.Command("bash", "-c", "id \"gotty\" &>/dev/null || { useradd gotty -m -s /bin/bash; }").Run(); err != nil {
		w.Write(xj.Errf("create user gotty - %v", err).Bts())
		return
	}
	if err := exec.Command("bash", "-c", "[ -e \"/home/gotty/gotty\" ] || { mkdir -p /home/gotty && cp "+path+" /home/gotty/gotty && chmod 700 /home/gotty/gotty && chown gotty:gotty /home/gotty/gotty; }").Run(); err != nil {
		w.Write(xj.Errf("preparation gotty - %v", err).Bts())
		return
	}
	os.Remove("/home/gotty/.bash_history")
	if !utils.IsFile("/home/gotty/gotty") {
		w.Write(xj.Err("gotty not found").Bts())
		return
	}
	type temp struct {
		User string `json:"user,omitempty" yaml:"user,omitempty"`
		Pass string `json:"pass,omitempty" yaml:"pass,omitempty"`
		Port int    `json:"port,omitempty" yaml:"port,omitempty"`
	}
	tmp := temp{User: utils.GenPass(8), Pass: utils.GenPass(16), Port: utils.GetFreePort(8100, 8800)}
	starter := exec.Command("bash", "-c", fmt.Sprintf("su - gotty -c '( /home/gotty/gotty --once --permit-write --timeout 30 --port %v --credential %v:%v bash ) &'", tmp.Port, tmp.User, tmp.Pass))
	if err := starter.Start(); err != nil {
		w.Write(xj.Errf("start gotty - %v", err).Bts())
		return
	}
	if err := xj.Parser.NewEncoder(w).Encode(tmp); err != nil {
		w.Write(xj.Errf("encode json - %v", err).Bts())
	}
}
