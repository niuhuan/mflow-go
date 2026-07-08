package games

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"mflow-go/internal/sysutil"
)

const (
	giBin     = "BetterGI.exe"
	giGameExe = "YuanShen.exe"
	giRegKey  = `HKEY_CURRENT_USER\SOFTWARE\miHoYo\原神`
)

func (r *Runner) runBetterGi(extraArgs []string, timeoutMinutes uint64) error {
	_ = sysutil.Taskkill(giBin)
	ctx := r.ctx()
	if sysutil.SleepCtx(ctx, 5*time.Second) {
		return nil
	}
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	if cfg.BetterGiPath == "" {
		return fmt.Errorf("BetterGI路径为空")
	}
	r.pm.RegisterExe(giBin)
	r.pm.RegisterExe(giGameExe)

	args := append([]string{}, extraArgs...)
	cmd := exec.Command(filepath.Join(cfg.BetterGiPath, giBin), args...)
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return fmt.Errorf("启动BetterGI失败: %w", err)
	}
	r.watchExeUntilGone(ctx, giBin, time.Duration(timeoutMinutes)*time.Minute)
	_ = sysutil.Taskkill(giBin)
	_ = sysutil.Taskkill(giGameExe)
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	return nil
}

func (r *Runner) RunBetterGi() error {
	cfg, _ := loadCfg()
	return r.runBetterGi([]string{"startOneDragon"}, cfg.BetterGiTimeoutMinutes)
}

func (r *Runner) RunBetterGiByConfig(configName string) error {
	cfg, _ := loadCfg()
	return r.runBetterGi([]string{"startOneDragon", configName}, cfg.BetterGiTimeoutMinutes)
}

func (r *Runner) RunBetterGiScheduler(groups string) error {
	cfg, _ := loadCfg()
	args := []string{"--startGroups"}
	args = append(args, strings.Fields(groups)...)
	return r.runBetterGi(args, cfg.BetterGiSchedulerTimeoutMinutes)
}

func (r *Runner) RunBetterGiGui() error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	cmd := exec.Command(filepath.Join(cfg.BetterGiPath, giBin))
	sysutil.SetupEncodingEnv(cmd)
	return cmd.Start()
}

func (r *Runner) CloseGi() error {
	return sysutil.Taskkill(giGameExe)
}

func (r *Runner) ClearGiReg() error {
	return sysutil.DeleteReg(giRegKey)
}

func (r *Runner) ListGiAccounts() ([]string, error) {
	cfg, err := loadCfg()
	if err != nil {
		return nil, err
	}
	return listSubDirs(filepath.Join(cfg.BetterGiPath, "m7f_accounts"))
}

func (r *Runner) ExportGiAccount(accountName string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.BetterGiPath, "m7f_accounts", accountName)
	if err := os.MkdirAll(accountFolder, 0o755); err != nil {
		return fmt.Errorf("创建账号文件夹失败: %w", err)
	}
	if err := sysutil.ExportReg(giRegKey, filepath.Join(accountFolder, "account.reg")); err != nil {
		return err
	}
	return sysutil.ZipDir(filepath.Join(cfg.BetterGiPath, "User"), filepath.Join(accountFolder, "User.zip"))
}

func (r *Runner) ImportGiAccount(accountName string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.BetterGiPath, "m7f_accounts", accountName)
	if err := sysutil.Unzip(filepath.Join(accountFolder, "User.zip"), cfg.BetterGiPath); err != nil {
		return err
	}
	return sysutil.ImportReg(filepath.Join(accountFolder, "account.reg"))
}

func (r *Runner) GenshinAutoLogin(accountName string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	if cfg.GenshinAutoLoginPath == "" {
		return fmt.Errorf("原神自动登录器文件夹路径未配置")
	}
	exePath := filepath.Join(cfg.GenshinAutoLoginPath, "AutoLoginGenshin.exe")
	if _, err := os.Stat(exePath); err != nil {
		return fmt.Errorf("原神自动登录器文件不存在: %s", exePath)
	}
	cmd := exec.Command(exePath, "--saved", accountName)
	cmd.Dir = cfg.GenshinAutoLoginPath
	sysutil.SetupEncodingEnv(cmd)
	if err := r.startCmd(cmd); err != nil {
		return fmt.Errorf("启动原神自动登录器失败: %w", err)
	}
	_ = cmd.Wait()
	return nil
}
