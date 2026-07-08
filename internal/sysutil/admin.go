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

// RunAsAdmin 以管理员权限重新启动自身（保留命令行参数）。
func RunAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	var argPtr *uint16
	if args != "" {
		argPtr, _ = syscall.UTF16PtrFromString(args)
	}
	cwd, _ := os.Getwd()
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, windows.SW_NORMAL)
}

// WhereWtExe 检测系统是否存在 Windows Terminal (wt.exe)。
func WhereWtExe() bool {
	return exec.Command("where", "wt.exe").Run() == nil
}
