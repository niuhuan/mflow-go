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

// RunAsAdmin 以管理员权限重新启动自身。
// keepConsole 为 true 时，优先通过 wt/cmd 打开可跟随主进程退出的控制台；
// 为 false 时直接提权启动 GUI 程序，不额外创建黑窗口。
func RunAsAdmin(keepConsole bool) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	args := escapeArgs(os.Args[1:])
	if !keepConsole {
		return shellExecuteRunas(exe, args)
	}

	// 内层命令：带引号的可执行文件路径 + 转义后的原命令行参数。
	inner := syscall.EscapeArg(exe)
	if args != "" {
		inner += " " + args
	}
	// `cmd /s /c "<inner>"`：cmd 会等待该命令完成后退出，
	// 因此管理员实例运行期间控制台仍可用于显示日志，应用退出后窗口会自动关闭。
	cmdArgs := `/s /c "` + inner + `"`

	// 优先 wt.exe，失败（如执行别名无法提权）则回退到 cmd.exe。
	if WhereWtExe() {
		if err := shellExecuteRunas("wt.exe", "cmd "+cmdArgs); err == nil {
			return nil
		}
	}
	return shellExecuteRunas("cmd.exe", cmdArgs)
}

func escapeArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	escaped := make([]string, 0, len(args))
	for _, arg := range args {
		escaped = append(escaped, syscall.EscapeArg(arg))
	}
	return strings.Join(escaped, " ")
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
