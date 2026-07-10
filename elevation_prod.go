//go:build !dev

package main

import (
	"os"
	"path/filepath"
	"strings"

	"mflow-go/internal/config"
	"mflow-go/internal/sysutil"
)

// maybeElevate 生产构建下若未提权则以管理员重启，返回 true 表示当前进程应退出。
func maybeElevate() bool {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	// 绑定生成阶段：wails 会编译并运行临时的 wailsbindings.exe（执行 main 直到 wails.Run）
	// 来收集绑定信息，此时不能提权/重启，否则会破坏该过程。
	exe, _ := os.Executable()
	if strings.Contains(strings.ToLower(filepath.Base(exe)), "wailsbindings") {
		return false
	}

	if !sysutil.IsElevated() {
		_ = sysutil.RunAsAdmin(cfg.KeepConsoleWindow)
		return true
	}
	// 已提权：附着到父进程（wt/cmd）的控制台，使脚本日志可见。
	if cfg.KeepConsoleWindow {
		sysutil.AttachParentConsole()
	}
	return false
}
