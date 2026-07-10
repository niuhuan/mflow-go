package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultWindowWidth  = 1200
	defaultWindowHeight = 800
)

// WindowState 保存主窗口的普通态尺寸/位置，以及是否以上次全屏状态启动。
type WindowState struct {
	X          int  `json:"x"`
	Y          int  `json:"y"`
	Width      int  `json:"width"`
	Height     int  `json:"height"`
	Fullscreen bool `json:"fullscreen"`
}

// DefaultWindowState 返回默认窗口状态。
func DefaultWindowState() WindowState {
	return WindowState{
		Width:  defaultWindowWidth,
		Height: defaultWindowHeight,
	}
}

func windowStateFile() string {
	return filepath.Join(AppDataPath(), "window_state.json")
}

// LoadWindowState 读取窗口状态。未保存过时返回 ok=false。
func LoadWindowState() (state WindowState, ok bool, err error) {
	state = DefaultWindowState()
	data, err := os.ReadFile(windowStateFile())
	if err != nil {
		if os.IsNotExist(err) {
			return state, false, nil
		}
		return state, false, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return DefaultWindowState(), false, err
	}
	if state.Width <= 0 {
		state.Width = defaultWindowWidth
	}
	if state.Height <= 0 {
		state.Height = defaultWindowHeight
	}
	return state, true, nil
}

// SaveWindowState 保存窗口状态。
func SaveWindowState(state WindowState) error {
	if err := os.MkdirAll(AppDataPath(), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(windowStateFile(), data, 0o644)
}
