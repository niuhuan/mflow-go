package procs

import (
	"context"
	"sync"

	"mflow-go/internal/sysutil"
)

// Manager 记录本次运行会话启动的进程与需兜底结束的 exe，用于「中断」。
type Manager struct {
	mu      sync.Mutex
	running bool
	pids    map[int]struct{}
	bgPids  map[int]struct{} // 后台进程 PID：运行结束(EndSession)或中断时都会被杀掉
	exes    map[string]struct{}
	ctx     context.Context
	cancel  context.CancelFunc
}

// New 创建进程会话管理器。
func New() *Manager {
	return &Manager{
		pids:   make(map[int]struct{}),
		bgPids: make(map[int]struct{}),
		exes:   make(map[string]struct{}),
	}
}

// StartSession 开始一个新的运行会话，返回可用于取消的 context。
func (m *Manager) StartSession() context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cancel != nil {
		m.cancel()
	}
	m.pids = make(map[int]struct{})
	m.bgPids = make(map[int]struct{})
	m.exes = make(map[string]struct{})
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.running = true
	return m.ctx
}

// Context 返回当前会话 context（无会话时返回 Background）。
func (m *Manager) Context() context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ctx == nil {
		return context.Background()
	}
	return m.ctx
}

// IsRunning 是否处于运行会话中。
func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// RegisterPID 登记本次会话启动的进程 PID。
func (m *Manager) RegisterPID(pid int) {
	if pid <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pids[pid] = struct{}{}
}

// RegisterBackgroundPID 登记后台进程 PID（运行结束或中断时都会被杀掉）。
func (m *Manager) RegisterBackgroundPID(pid int) {
	if pid <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bgPids[pid] = struct{}{}
}

// RegisterExe 登记需在中断时兜底 taskkill 的镜像名。
func (m *Manager) RegisterExe(exe string) {
	if exe == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exes[exe] = struct{}{}
}

// EndSession 结束会话（成功/失败均调用），取消 context，并杀掉遗留的后台进程。
func (m *Manager) EndSession() {
	m.mu.Lock()
	bgPids := make([]int, 0, len(m.bgPids))
	for pid := range m.bgPids {
		bgPids = append(bgPids, pid)
	}
	m.bgPids = make(map[int]struct{})
	m.running = false
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	m.mu.Unlock()

	for _, pid := range bgPids {
		_ = sysutil.KillTree(pid)
	}
}

// Interrupt 中断：取消 context，杀死所有登记进程树（含后台进程）+ 兜底 exe。
func (m *Manager) Interrupt() {
	m.mu.Lock()
	pids := make([]int, 0, len(m.pids)+len(m.bgPids))
	for pid := range m.pids {
		pids = append(pids, pid)
	}
	for pid := range m.bgPids {
		pids = append(pids, pid)
	}
	exes := make([]string, 0, len(m.exes))
	for exe := range m.exes {
		exes = append(exes, exe)
	}
	if m.cancel != nil {
		m.cancel()
	}
	m.bgPids = make(map[int]struct{})
	m.running = false
	m.mu.Unlock()

	for _, pid := range pids {
		_ = sysutil.KillTree(pid)
	}
	for _, exe := range exes {
		_ = sysutil.Taskkill(exe)
	}
}
