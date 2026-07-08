# mflow-go

三月七小助手、更好的原神、绝区零、鸣潮 一条龙工作流编排工具（Wails + Vue3 + Go 重写版）。

基于 [Blockly](https://developers.google.com/blockly) 积木编排流程，主要解决多账号多游戏任务编排问题。支持 崩坏：星穹铁道、原神、绝区零、鸣潮。

## 相关项目

- [三月七小助手 (March7thAssistant)](https://github.com/moesnow/March7thAssistant)
- [更好的原神 (BetterGI)](https://github.com/babalae/better-genshin-impact)
- [绝区零一条龙 (ZenlessZoneZero-OneDragon)](https://github.com/OneDragon-Anything/ZenlessZoneZero-OneDragon)
- [OK鸣潮](https://github.com/ok-oldking/ok-wuthering-waves)
- [原神自动登录工具](https://github.com/niuhuan/AutoLoginGenshin)

## 功能特性

- 以 Blockly 工作空间为主界面，拖拽积木即可编排一条龙流程。
- 「开始」积木不可删除；点击运行时**只执行「开始」链路上的代码**，以及孤岛的**函数定义 / 变量定义**（会被提升）。其余游离积木不会执行。
- 运行时后端记录启动的进程树与相关游戏/工具 exe，支持一键**中断**（强制结束进程树）。
- 各游戏积木采用统一暖色系配色，Blockly 采用 flat 平铺主题。
- 工具、设置、工程操作均以模态框呈现。
- 工程文件格式为 `.m7p`（JSON）。

## 教程

0. 推荐安装 [Windows Terminal](https://github.com/microsoft/terminal)。本程序原理是检测程序进程是否存在。
1. 下载 三月七小助手完整版、更好的原神、绝区零一条龙完整带运行时版；安装 [PowerShell](https://github.com/powershell/powershell/releases) 7.5 以上。
2. 对三方程序进行配置（设置任务结束时关闭程序和游戏、下载所需插件、设置合理超时等，详见各项目文档）。
3. 打开本程序，在「设置」中填写 三月七小助手、更好的原神、绝区零一条龙、鸣潮、原神自动登录器 的文件夹路径。
4. （可选）多账户时通过「工具」导出/导入账户（会同时备份/恢复对应软件的配置文件与注册表）。
5. 用「工程 → 模板」新建单账号或多账号工程，配置好账户名称与积木，点击「运行」。

## 注意事项

1. 多账号：导入导出账号时会同时备份和恢复各软件的配置文件，请导入修改之后再导出一次。
2. 若使用「更好的原神调度器」，请用空格分隔多个配置组，并在最后加上「退出程序」。
3. 分支宇宙有时会无限循环直至超时，推荐使用模拟宇宙。
4. OK鸣潮出现错误后不会自动退出，会等待到设置的超时后被结束。

## 开发

### 环境

- Go 1.25+
- Node 22、npm
- [Wails v2](https://wails.io/)（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）
- WebView2 运行时（Windows）

### 启动开发环境

```
wails dev
```

> dev 模式不会请求管理员权限；生产构建首次启动会自动请求管理员权限（第三方脚本需要 kill 进程 / 改注册表）。

### 构建

```
wails build
```

产物位于 `build/bin/mflow-go.exe`。

### 命令行自动运行

```
mflow-go.exe --auto-run C:\path\to\workflow.m7p
```

## 目录结构

```
mflow-go/
├── main.go / app.go          # Wails 入口与前端绑定（App 方法）
├── elevation_prod.go/_dev.go # 管理员提权（按 dev 构建标签隔离）
├── internal/
│   ├── config/               # 后端配置（%APPDATA%/opensource/mflow）
│   ├── games/                # 星铁/原神/绝区零/鸣潮 任务
│   ├── procs/                # 运行会话与进程树中断
│   ├── sysutil/              # 进程/注册表/压缩/管理员
│   └── version/              # 版本与更新检查
└── frontend/src/
    ├── blockly/              # 积木定义、生成器、工具箱、主题、模板
    ├── components/           # 工作空间、工具栏、控制台、模态框
    ├── composables/          # 运行引擎、对话框
    └── api/                  # 后端调用封装
```

更详细的架构与实现说明见 [AGENTS.md](./AGENTS.md)。
