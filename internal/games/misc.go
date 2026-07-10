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

// RunCommandBackground 后台执行命令：接好输出流后立即返回（不阻塞流程），
// 进程 PID 被登记，运行结束或中断时会被杀掉。
func (r *Runner) RunCommandBackground(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", command)
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startBackground(cmd); err != nil {
		return "", fmt.Errorf("后台运行命令失败: %w", err)
	}
	return "", nil
}
