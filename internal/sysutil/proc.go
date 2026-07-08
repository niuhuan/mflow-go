package sysutil

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// DecodeGBK 将 GBK 字节解码为字符串（Windows 命令行输出常为 GBK）。
func DecodeGBK(b []byte) string {
	dec := simplifiedchinese.GBK.NewDecoder()
	out, err := dec.Bytes(b)
	if err != nil {
		return string(b)
	}
	return string(out)
}

// SetupEncodingEnv 为子进程设置 UTF-8 相关环境变量。
func SetupEncodingEnv(cmd *exec.Cmd) {
	cmd.Env = append(cmd.Environ(),
		"PYTHONIOENCODING=utf-8",
		"LANG=zh_CN.UTF-8",
		"LC_ALL=zh_CN.UTF-8",
	)
}

// Taskkill 按镜像名强制结束进程。
func Taskkill(exe string) error {
	return exec.Command("taskkill", "/f", "/im", exe).Run()
}

// KillTree 按 PID 强制结束进程及其整个进程树。
func KillTree(pid int) error {
	return exec.Command("taskkill", "/f", "/t", "/pid", fmt.Sprintf("%d", pid)).Run()
}

// TaskExists 检查指定镜像名的进程是否存在。
func TaskExists(exe string) (bool, error) {
	out, err := exec.Command("tasklist", "/fi", fmt.Sprintf("imagename eq %s", exe)).Output()
	if err != nil {
		return false, fmt.Errorf("检查任务状态失败: %w", err)
	}
	text := strings.ToLower(DecodeGBK(out))
	return strings.Contains(text, strings.ToLower(exe)), nil
}

func escapePSSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// KillProcessByPath 通过完整路径匹配并结束进程（用于 ok-ww 等）。
func KillProcessByPath(exePath string) error {
	p := escapePSSingleQuote(exePath)
	script := fmt.Sprintf(
		"$p='%s'; Get-Process | Where-Object { $_.Path -eq $p } | ForEach-Object { Stop-Process -Id $_.Id -Force -ErrorAction SilentlyContinue }",
		p,
	)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).CombinedOutput()
	if err != nil {
		return fmt.Errorf("结束进程失败: %s", DecodeGBK(out))
	}
	return nil
}

// TaskExistsByPath 通过完整路径检查进程是否存在。
func TaskExistsByPath(exePath string) (bool, error) {
	p := escapePSSingleQuote(exePath)
	script := fmt.Sprintf(
		"$p='%s'; if (Get-Process | Where-Object { $_.Path -eq $p }) { Write-Output 1 } else { Write-Output 0 }",
		p,
	)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return false, fmt.Errorf("检查任务状态失败: %w", err)
	}
	return strings.TrimSpace(string(out)) == "1", nil
}

// SleepCtx 在可被 context 取消的前提下睡眠，返回 true 表示被取消。
func SleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return true
	case <-t.C:
		return false
	}
}
