package utils

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"runtime"
)

func IsFile(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		return !fi.IsDir()
	}
	return false
}

func ExecCommand(c string) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		exec.Command("powershell", "$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8").Run()
		cmd = exec.Command("powershell", c)
	case "linux":
		cmd = exec.Command("bash", "-c", c)
	default:
		return "", errors.New("unknown os")
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
