package main

import (
	"embed"
	"flag"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const (
	appTitle            = "mflow-go"
	defaultWindowWidth  = 1200
	defaultWindowHeight = 800
)

//go:embed all:frontend/dist
var assets embed.FS

// autoRunFile 保存命令行 --auto-run 指定的工程文件路径。
var autoRunFile string

func main() {
	// 生产构建需要管理员权限（第三方脚本会 kill 进程 / 改注册表）；dev 模式跳过。
	if maybeElevate() {
		return
	}

	var autoRun string
	flag.StringVar(&autoRun, "auto-run", "", "自动运行指定的工程文件(.m7p)")
	flag.Parse()
	autoRunFile = autoRun

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  appTitle,
		Width:  defaultWindowWidth,
		Height: defaultWindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
