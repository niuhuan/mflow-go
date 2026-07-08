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
	zzzLauncher = "OneDragon-Launcher.exe"
	zzzGameExe  = "ZenlessZoneZero.exe"
	zzzRegKey   = `HKEY_CURRENT_USER\SOFTWARE\miHoYo\绝区零`
)

func (r *Runner) RunZzzod() error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	if cfg.ZzzodPath == "" {
		return fmt.Errorf("绝区零一条龙路径为空")
	}
	ctx := r.ctx()
	r.pm.RegisterExe(zzzGameExe)
	launcher := filepath.Join(cfg.ZzzodPath, zzzLauncher)
	cmd := exec.Command("cmd", "/C", fmt.Sprintf("%s -o -c", launcher))
	cmd.Dir = cfg.ZzzodPath
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return fmt.Errorf("启动绝区零一条龙失败: %w", err)
	}
	r.watchExeUntilGone(ctx, zzzGameExe, time.Duration(cfg.ZzzodTimeoutMinutes)*time.Minute)
	_ = sysutil.Taskkill(zzzGameExe)
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	return nil
}

func (r *Runner) RunZzzodGui() error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	launcher := filepath.Join(cfg.ZzzodPath, zzzLauncher)
	var cmd *exec.Cmd
	if sysutil.WhereWtExe() {
		cmd = exec.Command("wt", "-d", cfg.ZzzodPath, "cmd", "/c", launcher)
	} else {
		cmd = exec.Command("cmd", "/c", launcher)
	}
	cmd.Dir = cfg.ZzzodPath
	sysutil.SetupEncodingEnv(cmd)
	return cmd.Start()
}

func (r *Runner) CloseZzz() error {
	return sysutil.Taskkill(zzzGameExe)
}

func (r *Runner) ClearZzzReg() error {
	return sysutil.DeleteReg(zzzRegKey)
}

func (r *Runner) ListZzzAccounts() ([]string, error) {
	cfg, err := loadCfg()
	if err != nil {
		return nil, err
	}
	return listSubDirs(filepath.Join(cfg.ZzzodPath, "m7f_accounts"))
}

func (r *Runner) ExportZzzAccount(accountName string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.ZzzodPath, "m7f_accounts", accountName)
	if err := os.MkdirAll(accountFolder, 0o755); err != nil {
		return fmt.Errorf("创建账号文件夹失败: %w", err)
	}
	if err := sysutil.ExportReg(zzzRegKey, filepath.Join(accountFolder, "account.reg")); err != nil {
		return err
	}
	return sysutil.ZipDir(filepath.Join(cfg.ZzzodPath, "config"), filepath.Join(accountFolder, "config.zip"))
}

func (r *Runner) ImportZzzAccount(accountName string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.ZzzodPath, "m7f_accounts", accountName)
	if err := sysutil.Unzip(filepath.Join(accountFolder, "config.zip"), cfg.ZzzodPath); err != nil {
		return err
	}
	return sysutil.ImportReg(filepath.Join(accountFolder, "account.reg"))
}
