package games

import (
	"fmt"
	"os/exec"

	"mflow-go/internal/sysutil"
)

// RunCommand 前台执行命令并等待结束。
func (r *Runner) RunCommand(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", command)
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return "", fmt.Errorf("运行命令失败: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("运行命令失败: %w", err)
	}
	return "", nil
}

// RunCommandBackground 新开窗口后台执行命令。
func (r *Runner) RunCommandBackground(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", "start", "", "cmd", "/c", command)
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return "", fmt.Errorf("后台运行命令失败: %w", err)
	}
	_ = cmd.Wait()
	return "", nil
}
