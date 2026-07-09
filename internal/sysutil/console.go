package sysutil

import (
	"os"

	"golang.org/x/sys/windows"
)

// attachParentProcess 是 AttachConsole 的特殊参数，表示附着到父进程的控制台。
const attachParentProcess = ^uintptr(0) // 0xFFFFFFFF (-1)

// AttachParentConsole 尝试附着到父进程（如 cmd.exe / wt.exe 中的 cmd）的控制台，
// 并把 os.Stdout / os.Stderr / os.Stdin 重定向到该控制台，使 GUI 进程的日志可见。
// 若没有父控制台（例如直接双击启动）则静默失败，返回 false。
func AttachParentConsole() bool {
	k32 := windows.NewLazySystemDLL("kernel32.dll")
	attach := k32.NewProc("AttachConsole")
	r1, _, _ := attach.Call(attachParentProcess)
	if r1 == 0 {
		return false
	}
	if out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stdout = out
		os.Stderr = out
	}
	if in, err := os.OpenFile("CONIN$", os.O_RDONLY, 0); err == nil {
		os.Stdin = in
	}
	return true
}
