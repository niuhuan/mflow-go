package games

import (
	"context"
	"os/exec"
	"time"

	"mflow-go/internal/config"
	"mflow-go/internal/procs"
	"mflow-go/internal/sysutil"
)

// Runner 承载所有游戏任务，持有进程会话管理器以支持中断。
type Runner struct {
	pm *procs.Manager
}

// NewRunner 创建 Runner。
func NewRunner(pm *procs.Manager) *Runner {
	return &Runner{pm: pm}
}

func (r *Runner) ctx() context.Context {
	return r.pm.Context()
}

func loadCfg() (config.BackendConfig, error) {
	return config.Load()
}

// startCmd 启动进程并登记 PID 到当前会话。
func (r *Runner) startCmd(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if cmd.Process != nil {
		r.pm.RegisterPID(cmd.Process.Pid)
	}
	return nil
}

// waitOrTimeout 等待子进程退出、超时或被中断，返回后不再等待。
func (r *Runner) waitOrTimeout(ctx context.Context, cmd *exec.Cmd, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	case <-time.After(timeout):
	}
}

// watchExeUntilGone 轮询指定 exe 是否仍存在，直到消失、超时或被中断。
func (r *Runner) watchExeUntilGone(ctx context.Context, exe string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if sysutil.SleepCtx(ctx, 100*time.Second) {
			return // 被中断
		}
		exists, err := sysutil.TaskExists(exe)
		if err != nil || !exists {
			return
		}
	}
}
