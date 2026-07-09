package games

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
	"unicode/utf8"

	"mflow-go/internal/config"
	"mflow-go/internal/procs"
	"mflow-go/internal/sysutil"
)

// Runner 承载所有游戏任务，持有进程会话管理器以支持中断。
type Runner struct {
	pm     *procs.Manager
	logMu  sync.RWMutex
	logger func(string)
}

// NewRunner 创建 Runner。
func NewRunner(pm *procs.Manager) *Runner {
	return &Runner{pm: pm}
}

func (r *Runner) ctx() context.Context {
	return r.pm.Context()
}

// SetLogger 设置输出回调（每行脚本输出会调用一次），用于把子进程输出转发到前端控制台。
func (r *Runner) SetLogger(fn func(string)) {
	r.logMu.Lock()
	r.logger = fn
	r.logMu.Unlock()
}

func (r *Runner) emit(line string) {
	r.logMu.RLock()
	fn := r.logger
	r.logMu.RUnlock()
	if fn != nil {
		fn(line)
	}
}

// lineWriter 将子进程输出按行拆分后转发给 emit。自动处理 UTF-8/GBK 编码。
type lineWriter struct {
	emit func(string)
	mu   sync.Mutex
	buf  []byte
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := bytes.TrimRight(w.buf[:i], "\r")
		w.flushLine(line)
		w.buf = w.buf[i+1:]
	}
	return len(p), nil
}

func (w *lineWriter) flushLine(line []byte) {
	if w.emit == nil {
		return
	}
	var s string
	if utf8.Valid(line) {
		s = string(line)
	} else {
		s = sysutil.DecodeGBK(line)
	}
	w.emit(s)
}

// newWriter 返回一个把子进程输出转发到前端控制台的 writer。
func (r *Runner) newWriter() *lineWriter {
	return &lineWriter{emit: r.emit}
}

func loadCfg() (config.BackendConfig, error) {
	return config.Load()
}

// startCmd 启动进程并登记 PID 到当前会话。
// 子进程输出同时写入：1) 本进程 stdout/stderr（在 wt/cmd 控制台窗口可见）；2) 前端页面控制台。
func (r *Runner) startCmd(cmd *exec.Cmd) error {
	if cmd.Stdout == nil {
		cmd.Stdout = io.MultiWriter(os.Stdout, r.newWriter())
	}
	if cmd.Stderr == nil {
		cmd.Stderr = io.MultiWriter(os.Stderr, r.newWriter())
	}
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
