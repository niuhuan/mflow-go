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
	hsrBin      = "March7th Assistant.exe"
	hsrLauncher = "March7th Launcher.exe"
	hsrGameExe  = "StarRail.exe"
	hsrRegKey   = `HKEY_CURRENT_USER\SOFTWARE\miHoYo\崩坏：星穹铁道`
)

// runM7fCommand 运行三月七小助手命令并在超时后强制结束。
func (r *Runner) runM7fCommand(command string, timeoutMinutes uint64) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	if cfg.M7Path == "" {
		return fmt.Errorf("三月七小助手路径为空")
	}
	ctx := r.ctx()
	r.pm.RegisterExe(hsrBin)
	r.pm.RegisterExe(hsrGameExe)

	cmd := exec.Command(filepath.Join(cfg.M7Path, hsrBin), command)
	cmd.Dir = cfg.M7Path
	sysutil.SetupEncodingEnv(cmd)
	stdin, _ := cmd.StdinPipe()
	if err := r.startCmd(cmd); err != nil {
		return fmt.Errorf("运行命令失败: %w", err)
	}
	if stdin != nil {
		go func() {
			buf := make([]byte, 4096)
			for i := range buf {
				buf[i] = '\n'
			}
			_, _ = stdin.Write(buf)
			_ = stdin.Close()
		}()
	}
	timeout := time.Duration(timeoutMinutes) * time.Minute
	r.waitOrTimeout(ctx, cmd, timeout)
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	_ = sysutil.Taskkill(hsrGameExe)
	return nil
}

func (r *Runner) FullRun() error {
	cfg, _ := loadCfg()
	return r.runM7fCommand("main", cfg.FullRunTimeoutMinutes)
}

func (r *Runner) DailyMission() error {
	return r.runM7fCommand("daily", 60)
}

func (r *Runner) RefreshStamina() error {
	return r.runM7fCommand("power", 60)
}

func (r *Runner) SimulatedUniverse() error {
	return r.runM7fCommand("universe", 120)
}

func (r *Runner) Farming() error {
	return r.runM7fCommand("fight", 360)
}

func (r *Runner) CloseGame() error {
	return sysutil.Taskkill(hsrGameExe)
}

func (r *Runner) ClearHsrReg() error {
	return sysutil.DeleteReg(hsrRegKey)
}

func (r *Runner) RunM7Launcher() error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	cmd := exec.Command(filepath.Join(cfg.M7Path, hsrLauncher))
	cmd.Dir = cfg.M7Path
	sysutil.SetupEncodingEnv(cmd)
	return cmd.Start()
}

// ListAccounts 列出星铁已保存账号。
func (r *Runner) ListAccounts() ([]string, error) {
	cfg, err := loadCfg()
	if err != nil {
		return nil, err
	}
	return listSubDirs(filepath.Join(cfg.M7Path, "m7f_accounts"))
}

func listSubDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

// GetAccountUID 读取当前星铁账号 UID。
func (r *Runner) GetAccountUID() (int64, error) {
	return sysutil.ReadHsrUID()
}

// LoadAccount 载入账号：复制 config.yaml 并导入注册表。
func (r *Runner) LoadAccount(name string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.M7Path, "m7f_accounts", name)
	if err := copyFile(filepath.Join(accountFolder, "config.yaml"), filepath.Join(cfg.M7Path, "config.yaml")); err != nil {
		return fmt.Errorf("复制配置文件失败: %w", err)
	}
	return sysutil.ImportReg(filepath.Join(accountFolder, "account.reg"))
}

// SaveAccount 保存当前账号。
func (r *Runner) SaveAccount(name string) error {
	return r.ExportAccount(name, "", "")
}

// ExportAccount 导出星铁账号（config.yaml + 注册表）。
func (r *Runner) ExportAccount(accountName, username, password string) error {
	cfg, err := loadCfg()
	if err != nil {
		return err
	}
	accountFolder := filepath.Join(cfg.M7Path, "m7f_accounts", accountName)
	if err := os.MkdirAll(accountFolder, 0o755); err != nil {
		return fmt.Errorf("创建账号文件夹失败: %w", err)
	}
	regPath := filepath.Join(accountFolder, "account.reg")
	if err := sysutil.ExportReg(hsrRegKey, regPath); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(cfg.M7Path, "config.yaml"), filepath.Join(accountFolder, "config.yaml")); err != nil {
		return fmt.Errorf("复制配置文件失败: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
