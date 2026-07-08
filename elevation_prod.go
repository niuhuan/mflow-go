//go:build !dev

package main

import "mflow-go/internal/sysutil"

// maybeElevate 生产构建下若未提权则以管理员重启，返回 true 表示当前进程应退出。
func maybeElevate() bool {
	if !sysutil.IsElevated() {
		_ = sysutil.RunAsAdmin()
		return true
	}
	return false
}
