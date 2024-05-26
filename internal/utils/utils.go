package utils

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func IsFile(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return !fi.IsDir()
	}
	return false
}

func IsDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

func GenCMD(script string, params ...string) *exec.Cmd {
	OSpEXT := strings.ToLower(runtime.GOOS + filepath.Ext(script))
	sCMD := regexp.MustCompile("[;|&$]").ReplaceAllString(strings.Join(params, " "), "")
	switch OSpEXT {
	case "windows.bat":
		return exec.Command("cmd", "/c", fmt.Sprintf("%s %s", script, sCMD))
	case "windows.ps1", "windows":
		return exec.Command("powershell", fmt.Sprintf("%s %s", script, sCMD))
	case "windows.py":
		return exec.Command("powershell", fmt.Sprintf("py %s %s", script, sCMD))
	case "linux.sh", "linux":
		return exec.Command("bash", "-c", fmt.Sprintf("%s %s", script, sCMD))
	case "linux.py":
		return exec.Command("bash", "-c", fmt.Sprintf("python3 %s %s", script, sCMD))
	case "linux.pl":
		return exec.Command("bash", "-c", fmt.Sprintf("perl %s %s", script, sCMD))
	}
	return exec.Command(script, sCMD)
}

func ExecCommand(cmd *exec.Cmd) (string, error) {
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8").Run()
	}
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(pipe)
	var response string
	for scanner.Scan() {
		response += scanner.Text() + "\n"
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return response, nil
}

func OnTermination(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if fn != nil {
			fn()
		}
	}()
}

func GenPass(length int) string {
	password := make([]byte, length)
	for n := range password {
		if index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset)))); err == nil {
			password[n] = charset[index.Int64()]
		} else {
			password[n] = 0x0
		}
	}
	return string(password)
}

func GetFreePort(b, e uint16) int {
	for port := int(b); port < int(e); port++ {
		if conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			defer conn.Close()
			if err := conn.Close(); err == nil {
				return port
			}
		}
	}
	return 8888
}
