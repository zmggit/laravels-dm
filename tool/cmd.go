package tool

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"os/exec"
)

// ExecCommand 执行linux 命令
func ExecCommand(command string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("错误信息---recover---->", err)
		}
	}()
	cmd := exec.Command("bash", "-c", command)

	var output bytes.Buffer
	//cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	if err != nil {
		return output.Bytes(), err
	}
	//logs.Debug(cmd.Process.Pid, "调试cmd.Process")
	//logs.Debug(string(output.Bytes()), "output.Bytes()")
	return output.Bytes(), nil
}
