package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// BackendConfig 对齐原 Tauri 版本的后端配置字段。
type BackendConfig struct {
	M7Path              string `json:"m7_path"`
	BetterGiPath        string `json:"better_gi_path"`
	ZzzodPath           string `json:"zzzod_path"`
	GenshinAutoLoginPath string `json:"genshin_auto_login_path"`
	OkWwPath            string `json:"ok_ww_path"`
	KeepConsoleWindow   bool   `json:"keep_console_window"`
	ScriptLogToStdout   bool   `json:"script_log_to_stdout"`
	ScriptLogToAppConsole bool `json:"script_log_to_app_console"`

	// 错误处理策略
	StopOnScriptError  bool `json:"stop_on_script_error"`  // 运行脚本错误时停止流程（默认关）
	StopOnAccountError bool `json:"stop_on_account_error"` // 备份/恢复账号错误时停止流程（默认开）

	FullRunTimeoutMinutes           uint64 `json:"full_run_timeout_minutes"`
	DailyMissionTimeoutMinutes      uint64 `json:"daily_mission_timeout_minutes"`
	RefreshStaminaTimeoutMinutes    uint64 `json:"refresh_stamina_timeout_minutes"`
	SimulatedUniverseTimeoutMinutes uint64 `json:"simulated_universe_timeout_minutes"`
	FarmingTimeoutMinutes           uint64 `json:"farming_timeout_minutes"`
	BetterGiTimeoutMinutes          uint64 `json:"better_gi_timeout_minutes"`
	BetterGiSchedulerTimeoutMinutes uint64 `json:"better_gi_scheduler_timeout_minutes"`
	ZzzodTimeoutMinutes             uint64 `json:"zzzod_timeout_minutes"`
	OkWwTimeoutMinutes              uint64 `json:"ok_ww_timeout_minutes"`
}

// Default 返回带默认超时的配置。
func Default() BackendConfig {
	return BackendConfig{
		KeepConsoleWindow:               true,
		ScriptLogToStdout:               true,
		ScriptLogToAppConsole:           false,
		StopOnScriptError:               false,
		StopOnAccountError:              true,
		FullRunTimeoutMinutes:           60,
		DailyMissionTimeoutMinutes:      60,
		RefreshStaminaTimeoutMinutes:    60,
		SimulatedUniverseTimeoutMinutes: 120,
		FarmingTimeoutMinutes:           360,
		BetterGiTimeoutMinutes:          60,
		BetterGiSchedulerTimeoutMinutes: 300,
		ZzzodTimeoutMinutes:             60,
		OkWwTimeoutMinutes:              60,
	}
}

// AppDataPath 返回应用数据目录：%APPDATA%/opensource/mflow
func AppDataPath() string {
	return filepath.Join(os.Getenv("APPDATA"), "opensource", "mflow")
}

func configFile() string {
	return filepath.Join(AppDataPath(), "backend_config.json")
}

// Load 读取后端配置，不存在则返回默认值。
func Load() (BackendConfig, error) {
	cfg := Default()
	path := configFile()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save 保存后端配置。
func Save(cfg BackendConfig) error {
	if err := os.MkdirAll(AppDataPath(), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile(), data, 0o644)
}
