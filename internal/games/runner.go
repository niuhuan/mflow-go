package games

import (
	"bytes"
	"context"
	"fmt"
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
	outMu  sync.Mutex // 串行化写入本进程控制台，避免多流交错
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

// lineWriter 将子进程输出按行拆分、统一解码为 UTF-8 后交给 sink 处理。
type lineWriter struct {
	sink func(string)
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
	if w.sink == nil {
		return
	}
	// 子进程输出可能是 UTF-8 或 GBK，统一解码成 UTF-8 再输出。
	var s string
	if utf8.Valid(line) {
		s = string(line)
	} else {
		s = sysutil.DecodeGBK(line)
	}
	w.sink(s)
}

// writeConsole 把一行（UTF-8）输出到本进程 stdout（wt/cmd 控制台窗口），串行化避免交错。
func (r *Runner) writeConsole(s string) {
	r.outMu.Lock()
	fmt.Fprintln(os.Stdout, s)
	r.outMu.Unlock()
}

// newWriter 返回一个把子进程输出（解码为 UTF-8）同时写到本进程控制台与前端页面控制台的 writer。
func (r *Runner) newWriter() *lineWriter {
	return &lineWriter{sink: func(s string) {
		r.writeConsole(s)
		r.emit(s)
	}}
}

func loadCfg() (config.BackendConfig, error) {
	return config.Load()
}

// attachOutput 把子进程 stdout/stderr 接到统一的输出通道（控制台 + 页面控制台）。
func (r *Runner) attachOutput(cmd *exec.Cmd) {
	if cmd.Stdout == nil {
		cmd.Stdout = r.newWriter()
	}
	if cmd.Stderr == nil {
		cmd.Stderr = r.newWriter()
	}
}

// startCmd 启动进程并登记 PID 到当前会话，转发其输出。
func (r *Runner) startCmd(cmd *exec.Cmd) error {
	r.attachOutput(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	if cmd.Process != nil {
		r.pm.RegisterPID(cmd.Process.Pid)
	}
	return nil
}

// startBackground 后台启动进程：接好输出流后立即返回（不阻塞），并把 PID 登记为后台进程，
// 运行结束或中断时统一杀掉。用于「后台运行命令」。
func (r *Runner) startBackground(cmd *exec.Cmd) error {
	r.attachOutput(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	if cmd.Process != nil {
		r.pm.RegisterBackgroundPID(cmd.Process.Pid)
	}
	// 在后台回收进程资源，不阻塞调用方。
	go func() { _ = cmd.Wait() }()
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
