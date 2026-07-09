package sysutil

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// IsElevated 判断当前进程是否具有管理员权限。
func IsElevated() bool {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token); err != nil {
		return false
	}
	defer token.Close()
	return token.IsElevated()
}

// RunAsAdmin 以管理员权限重新启动自身，并在一个控制台窗口中运行以便查看日志：
// 优先使用 Windows Terminal (wt.exe)，否则回退到 cmd.exe，均以 `cmd /s /k` 保留窗口。
func RunAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// 内层命令：带引号的可执行文件路径 + 原有命令行参数。
	inner := `"` + exe + `"`
	if extra := strings.Join(os.Args[1:], " "); extra != "" {
		inner += " " + extra
	}
	// `cmd /s /k "<inner>"`：/k 保留窗口以便查看脚本输出；/s 保证外层引号被正确剥离。
	cmdArgs := `/s /k "` + inner + `"`

	// 优先 wt.exe，失败（如执行别名无法提权）则回退到 cmd.exe。
	if WhereWtExe() {
		if err := shellExecuteRunas("wt.exe", "cmd "+cmdArgs); err == nil {
			return nil
		}
	}
	return shellExecuteRunas("cmd.exe", cmdArgs)
}

// shellExecuteRunas 以管理员权限（runas）通过 ShellExecute 启动指定程序。
func shellExecuteRunas(file, params string) error {
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	filePtr, _ := syscall.UTF16PtrFromString(file)
	paramPtr, _ := syscall.UTF16PtrFromString(params)
	cwd, _ := os.Getwd()
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	return windows.ShellExecute(0, verbPtr, filePtr, paramPtr, cwdPtr, windows.SW_NORMAL)
}

// WhereWtExe 检测系统是否存在 Windows Terminal (wt.exe)。
func WhereWtExe() bool {
	return exec.Command("where", "wt.exe").Run() == nil
}
