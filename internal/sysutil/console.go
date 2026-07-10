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

// cpUTF8 是 UTF-8 代码页编号。
const cpUTF8 = 65001

// SetConsoleUTF8 将当前控制台的输入/输出代码页设为 UTF-8，
// 使本程序及子进程的 UTF-8 输出在控制台中正确显示（不再乱码）。
func SetConsoleUTF8() {
	k32 := windows.NewLazySystemDLL("kernel32.dll")
	_, _, _ = k32.NewProc("SetConsoleOutputCP").Call(uintptr(cpUTF8))
	_, _, _ = k32.NewProc("SetConsoleCP").Call(uintptr(cpUTF8))
}
