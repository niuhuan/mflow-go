package games

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"mflow-go/internal/sysutil"
)

const (
	wwExe     = "ok-ww.exe"
	wwGameExe = "Client-Win64-Shipping.exe"
	wwRegKey  = `HKEY_CURRENT_USER\SOFTWARE\Kuro Games\鸣潮`
)

func (r *Runner) wwPaths() (dir, exe, pythonw string, err error) {
	cfg, e := loadCfg()
	if e != nil {
		return "", "", "", e
	}
	if cfg.OkWwPath == "" {
		return "", "", "", fmt.Errorf("ok-ww路径为空")
	}
	dir = cfg.OkWwPath
	exe = filepath.Join(dir, wwExe)
	pythonw = filepath.Join(dir, "data", "apps", "ok-ww", "python", "pythonw.exe")
	return dir, exe, pythonw, nil
}

func (r *Runner) KillOkWw() error {
	_, exe, pythonw, err := r.wwPaths()
	if err != nil {
		return err
	}
	_ = sysutil.KillProcessByPath(exe)
	_ = sysutil.KillProcessByPath(pythonw)
	_ = sysutil.Taskkill(wwGameExe)
	return nil
}

func (r *Runner) startOkWw(taskType string) error {
	dir, exe, pythonw, err := r.wwPaths()
	if err != nil {
		return err
	}
	if _, err := os.Stat(exe); err != nil {
		return fmt.Errorf("ok-ww.exe文件不存在: %s", exe)
	}
	ctx := r.ctx()
	r.pm.RegisterExe(wwGameExe)

	cmd := exec.Command(exe, "-t", taskType, "-e")
	cmd.Dir = dir
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return fmt.Errorf("启动 ok-ww 失败: %w", err)
	}
	if sysutil.SleepCtx(ctx, 100*time.Second) {
		return nil
	}
	deadline := time.Now().Add(time.Hour)
	for time.Now().Before(deadline) {
		if sysutil.SleepCtx(ctx, 100*time.Second) {
			return nil
		}
		exists, err := sysutil.TaskExistsByPath(pythonw)
		if err != nil || !exists {
			break
		}
	}
	_ = r.KillOkWw()
	sysutil.SleepCtx(ctx, 10*time.Second)
	return nil
}

func (r *Runner) StartOkWwDaily() error {
	return r.startOkWw("1")
}

func (r *Runner) StartOkWwWeekly() error {
	return r.startOkWw("2")
}

func (r *Runner) ClearWwReg() error {
	return sysutil.DeleteReg(wwRegKey)
}
