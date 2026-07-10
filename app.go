package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"mflow-go/internal/config"
	"mflow-go/internal/games"
	"mflow-go/internal/procs"
	"mflow-go/internal/sysutil"
	"mflow-go/internal/version"
)

// App 是绑定到前端的主结构体。嵌入 *games.Runner 使其游戏任务方法被 Wails 自动导出。
type App struct {
	*games.Runner
	ctx context.Context
	pm  *procs.Manager
}

// NewApp 创建 App。
func NewApp() *App {
	pm := procs.New()
	return &App{
		Runner: games.NewRunner(pm),
		pm:     pm,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 将控制台代码页设为 UTF-8，使我们与子进程的 UTF-8 输出不乱码。
	sysutil.SetConsoleUTF8()
	// 将子进程输出转发到前端控制台。
	a.Runner.SetLogger(func(line string) {
		runtime.EventsEmit(ctx, "console:output", line)
	})
}

// ---------- 版本 ----------

func (a *App) GetVersion() string {
	return version.Current
}

func (a *App) GetNewVersion() string {
	return version.NewVersion()
}

func (a *App) OpenReleasePage() {
	url := version.ReleaseURL()
	if url != "" {
		runtime.BrowserOpenURL(a.ctx, url)
	}
}

// ---------- 文件 ----------

func (a *App) AppDataPath() string {
	return config.AppDataPath()
}

func (a *App) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (a *App) Mkdir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func (a *App) ReadTextFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *App) WriteTextFile(path, content string) error {
	if dir := filepath.Dir(path); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// ---------- 配置 ----------

func (a *App) LoadBackendConfig() (config.BackendConfig, error) {
	return config.Load()
}

func (a *App) SaveBackendConfig(cfg config.BackendConfig) error {
	return config.Save(cfg)
}

// OpenProjectDialog 弹出打开文件对话框，返回所选 .m7p 路径（取消返回空串）。
func (a *App) OpenProjectDialog() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "打开工程",
		Filters: []runtime.FileFilter{
			{DisplayName: "mflow 工程 (*.m7p, *.m7f)", Pattern: "*.m7p;*.m7f"},
			{DisplayName: "所有文件 (*.*)", Pattern: "*.*"},
		},
	})
}

// SaveProjectDialog 弹出另存为对话框，返回目标 .m7p 路径（取消返回空串）。
func (a *App) SaveProjectDialog() (string, error) {
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存工程",
		DefaultFilename: "workflow.m7p",
		Filters: []runtime.FileFilter{
			{DisplayName: "mflow 工程 (*.m7p)", Pattern: "*.m7p"},
		},
	})
}

// ---------- 自动运行 ----------

func (a *App) GetAutoRunFile() string {
	return autoRunFile
}

// ---------- 运行会话与中断 ----------

// StartRun 开始运行会话（供前端在执行流程前调用）。
func (a *App) StartRun() {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}
	a.Runner.StartRunSession(cfg)
}

// EndRun 结束运行会话（成功/失败/中断后调用）。
func (a *App) EndRun() {
	a.Runner.EndRunSession()
}

// Interrupt 中断：结束本次会话登记的进程树与相关 exe。
func (a *App) Interrupt() {
	a.pm.Interrupt()
}

// IsRunning 是否处于运行会话中。
func (a *App) IsRunning() bool {
	return a.pm.IsRunning()
}
