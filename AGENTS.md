# mflow-go 开发指南（AGENTS.md）

本文件是 `mflow-go` 项目的开发大纲与实现路线图。`mflow-go` 是对同仓库内 `mflow/` 目录（原 Tauri + React + Rust 实现）的重写版本，使用 **Wails + Vue3 + Go** 技术栈。

> 面向 AI Agent 与开发者：先读本文件，再动手。所有轮次按顺序推进，每轮完成后自检验收标准。

---

## 1. 项目背景

`mflow`（三月七小助手工作流）是一个游戏一条龙任务编排工具，基于 [Blockly](https://developers.google.com/blockly) 积木编排流程，支持 **崩坏：星穹铁道（三月七小助手）、原神（BetterGI）、绝区零（ZZZOD）、鸣潮（OK-WW）** 的多账号、多游戏自动化任务调度。

原项目通过 Blockly 生成 JavaScript 代码，前端 `eval` 执行，逐条调用后端命令（启动第三方 exe、导入导出账号注册表/配置、清理注册表、等待、命令行等）。

### 1.1 原项目关键实现（迁移参考）

原后端逻辑集中在 `mflow/src-tauri/src/main.rs`，需完整迁移到 Go。核心命令清单（Tauri command → Go 方法）：

| 分类 | 原命令 | 说明 |
|------|--------|------|
| 版本 | `get_version` / `get_new_version` | 读取版本、查询 GitHub 最新 release |
| 文件 | `exists` / `mkdir` / `read_text_file` / `write_text_file` / `app_data_path` | 文件读写、应用数据目录（`%APPDATA%/opensource/mflow`）|
| 配置 | `load_backend_config` / `save_backend_config` | `backend_config.json`（各软件路径 + 各任务超时分钟数）|
| 星铁 | `full_run` / `daily_mission` / `refresh_stamina` / `simulated_universe` / `farming` / `close_game` | 调用 `March7th Assistant.exe <cmd>`，超时后 kill + `taskkill StarRail.exe` |
| 星铁账号 | `load_account` / `save_account` / `export_account` / `list_accounts` / `get_account_uid` / `clear_game_reg` | 复制 `config.yaml` + 注册表导入导出（`崩坏：星穹铁道`）|
| 原神 | `run_better_gi` / `run_better_gi_by_config` / `run_better_gi_scheduler` / `close_gi` / `run_better_gi_gui` | 调用 `BetterGI.exe`（`startOneDragon` / `--startGroups`），轮询进程存在直到结束或超时 |
| 原神账号 | `export_gi_account` / `import_gi_account` / `list_gi_accounts` / `clear_gi_reg` / `genshin_auto_login` | User 目录 zip + 注册表（`原神`）；自动登录调 `AutoLoginGenshin.exe --saved <acc>` |
| 绝区零 | `run_zzzod` / `run_zzzod_gui` / `close_zzz` | 调用 `OneDragon-Launcher.exe -o -c`，轮询 `ZenlessZoneZero.exe` |
| 绝区零账号 | `export_zzz_account` / `import_zzz_account` / `list_zzz_accounts` / `clear_zzz_reg` | config 目录 zip + 注册表（`绝区零`）|
| 鸣潮 | `start_ok_ww_daily` / `start_ok_ww_weekly` / `kill_ok_ww` / `clear_ww_reg` | 调用 `ok-ww.exe -t <1/2> -e`，按路径 kill 进程 |
| 通用 | `run_command` / `run_command_background` | `cmd /c` 执行命令 |
| 其他 | `open_release_page` / `get_auto_run_file` | 打开发布页、命令行 `--auto-run` 自动运行 |
| 管理员 | `mflow/src-tauri/src/win.rs` + `cmd/` | 检测管理员权限并以管理员重启（第三方脚本需要）|

原前端 Blockly 定义参考：
- 积木定义：`mflow/src/blocks/customBlocks.ts`
- 代码生成器：`mflow/src/blocks/generators.js`
- 工具箱与运行引擎：`mflow/src/App.jsx`
- 模板 XML：`mflow/src/assets/single-user.xml`、`mflow/src/assets/mutil-account.xml`

---

## 2. 与原项目的差异（本次重写的核心需求）

1. **前端框架**：Vue3（替代 React）。仍使用 Blockly 官方 npm 包。
2. **主界面即工作空间**：应用启动后直接进入 Blockly 工作空间，而不是先显示「开始/打开工程」页面。
   - 原来「关闭当前工作空间」面板里的大量操作按钮（启动各 GUI、导入/导出各游戏账号等工具菜单）→ **移入模态框**。
   - 「配置/设置」页 → **模态框**。
   - 「打开/新建/模板」等工程操作 → 顶部菜单或模态框。
3. **Blockly 美化**：参考 <https://docs.blockly.com/guides/design/appearance/>，采用指南中**第一种较平铺（flat）的外观**。引入官方/开源 flat 主题（如自定义 `Blockly.Theme`，扁平积木、去重描边、清爽配色、圆角、平铺分类工具箱）。
4. **"开始"积木不可删除**，且运行逻辑改变：
   - 点击运行**只执行「开始」积木链路上的语句**。
   - 以及**孤岛（未连接到开始链路）的「函数定义 / 变量定义」代码**（需被 hoist 以便被调用）。
     - 函数定义 = `custom_function_def`。
     - 变量定义 = **新增的独立「定义变量」积木**（`define_variable`，语句块，形如「定义变量 名称 = 值」）；孤岛存在时同样 hoist 到主链路之前。
   - 其他游离的「运行类」积木（如单独放着的 `full_run` 等）**不执行**。
5. **新增中断功能**：
   - 运行时后端**记录启动的所有子进程**（进程树）以及**指定的 exe**（如 `March7th Assistant.exe`、`BetterGI.exe`、`OneDragon-Launcher.exe`、`ok-ww.exe` 及对应游戏进程）。
   - 点击「中断」后端对**进程树 + 指定 exe 执行 kill**，并让前端运行引擎停止后续步骤。

---

## 3. 技术栈

- **后端**：Go 1.25 + Wails v2 (v2.13.0)。入口 `main.go`，业务 `app.go`（建议拆分为多个文件/包）。
- **前端**：Vue3 + Vite。目录 `frontend/`。
- **积木引擎**：`blockly`（npm，`^12`）+ `blockly/javascript` 生成器。
- **平台**：Windows 优先（进程管理、注册表、`taskkill`/`tasklist` 均为 Windows）。需管理员权限运行。
- **Go 依赖（预计新增）**：
  - `golang.org/x/sys/windows/registry`（注册表导入导出，替代 winreg）
  - `github.com/shirou/gopsutil/v3/process`（进程树遍历，用于中断 kill）
  - 或直接调用系统命令 `reg`、`taskkill`、`tasklist`、`powershell Compress-Archive/Expand-Archive`（与原实现一致，最省事）。

---

## 4. 架构大纲

### 4.1 前端 (`frontend/src/`)

```
frontend/src/
├── main.js                  # Vue 应用入口
├── App.vue                  # 根组件：顶部工具栏 + Blockly 工作空间 + 控制台
├── blockly/
│   ├── customBlocks.js      # 积木定义（迁移自 customBlocks.ts）
│   ├── generators.js        # JS 代码生成器（迁移 + 改造运行链路裁剪）
│   ├── theme.js             # flat 平铺主题 + 工具箱样式
│   ├── toolbox.js           # 工具箱 XML/JSON
│   └── templates.js         # 单账号/多账号模板 XML
├── composables/
│   ├── useBlockly.js         # Blockly 注入、序列化、变更监听、开始积木保护
│   └── useRunner.js          # 运行引擎：裁剪链路、window API 绑定、执行、中断
├── components/
│   ├── BlocklyWorkspace.vue  # Blockly 容器组件
│   ├── ConsolePanel.vue      # 运行日志控制台
│   ├── Toolbar.vue           # 顶部工具栏（运行/中断/保存/工具/设置/工程）
│   └── modals/
│       ├── ToolsModal.vue    # 工具箱（启动各 GUI、导入导出账号）
│       ├── SettingsModal.vue # 设置（各软件路径 + 超时）
│       ├── ProjectModal.vue  # 打开/新建/模板
│       └── SelectOneModal.vue# 通用列表选择框（替代原生 prompt/selectOne）
├── api/
│   └── backend.js            # 封装 wailsjs/go/main/App 调用（替代 fromTauri.ts）
└── style.css
```

### 4.2 后端 (Go, 仓库根)

```
mflow-go/
├── main.go                  # Wails 启动 + --auto-run 参数解析 + 管理员提权
├── app.go                   # App 结构体 + startup + 通用绑定方法
├── internal/
│   ├── config/config.go     # BackendConfig 读写（对齐原 config.rs 字段）
│   ├── games/
│   │   ├── hsr.go           # 星铁命令 + 账号
│   │   ├── genshin.go       # 原神命令 + 账号 + 自动登录
│   │   ├── zzz.go           # 绝区零命令 + 账号
│   │   └── ww.go            # 鸣潮命令
│   ├── procs/
│   │   └── manager.go       # 运行进程注册表 + 进程树 kill（中断核心）
│   ├── sysutil/
│   │   ├── reg.go           # 注册表导入/导出/删除
│   │   ├── zip.go           # 目录压缩/解压
│   │   ├── proc.go          # spawn/taskkill/tasklist/进程存在检测
│   │   └── win_admin.go     # 管理员检测与提权
│   └── version/version.go   # 版本 + GitHub release 查询
└── AGENTS.md
```

> Wails 绑定：所有需被前端调用的方法挂在 `App` 上（或聚合的服务结构体），命名用 Go 导出风格（如 `FullRun`、`ExportGiAccount`），Wails 会生成 `frontend/wailsjs/go/main/App.js`。前端调用它们即可（替代原 `invoke`）。事件推送（运行日志、进程状态）用 Wails 的 `runtime.EventsEmit`。

---

## 5. 技术要点

### 5.1 Blockly 与 Vue 集成
- 在 `onMounted` 中 `Blockly.inject`，`onUnmounted` 中 `dispose`。用 `ref` 持有 workspace，避免响应式代理包裹 Blockly 实例（用 `markRaw` 或普通变量，**不要放进 `reactive`**）。
- 中文语言包：`import * as zhHans from 'blockly/msg/zh-hans'; Blockly.setLocale(zhHans)`。
- 序列化统一使用 `Blockly.serialization.workspaces`（**JSON**）。新存档扩展名为 **`.m7p`**，**不兼容**旧 `.m7f`（XML）文件。

### 5.2 Blockly flat 平铺美化
- 依据官方 appearance 指南创建自定义主题：`Blockly.Theme.defineTheme('mflowFlat', {...})`。
- 要点：扁平积木风格（`blockStyles` 去除高光/渐变、统一 `colourPrimary/Secondary/Tertiary`）、`componentStyles`（工作区背景、工具箱背景、flyout 背景浅色）、`fontStyle`、适度圆角。
- 工具箱使用平铺分类样式（`toolbox` 配置 `renderer: 'thrasos'` 或 `'zelos'` 视效果而定；平铺外观优先考虑非 zelos 的 `geras/thrasos` + flat 配色）。
- 允许引用开源 flat 主题代码，但需与项目分类颜色（星铁/原神/ZZZ/鸣潮/通用/流程/逻辑/数值/变量/函数/集合）协调。

### 5.3 "开始"积木保护（不可删除）
- `start_flow` 积木设置为不可删除：加载后对该 block 调用 `block.setDeletable(false)`；工作区初始必然含且仅含一个 `start_flow`。
- 监听删除事件，若被删除则阻止/自动恢复；工具箱中不再提供「开始」积木（因为已存在且唯一）。

### 5.4 运行链路裁剪（核心逻辑变更）
运行时**不再**对整个 workspace 生成代码，而是：
1. 定位唯一的 `start_flow` 积木，从它开始沿 `nextConnection` 生成主链路代码（`start_flow` 本身不产码，其后继语句依次生成）。
2. 额外收集**孤岛的定义类积木**并前置（hoist）：
   - `custom_function_def`（函数定义）
   - 变量定义 / 赋值类（若引入 `variables_set` 等）
   - 判定"孤岛"：该顶层积木不在 start 链路上，且类型属于「定义类」白名单。
3. **忽略**所有其他不在链路上的顶层运行类积木（不生成、不执行）。

实现方式（二选一）：
- 方案 A（推荐）：用 `workspace.getTopBlocks(true)` 找出所有顶层积木；对 start 链路积木逐个 `javascriptGenerator.blockToCode`；对定义类孤岛积木单独生成并拼到最前面。
- 方案 B：给运行类积木生成器加"仅在链路内产码"的标记，配合 `workspaceToCode` 过滤。

> 注意：函数定义在原实现中生成 `async function ...`，需保证 hoist 到 `eval` 包裹的作用域内、在主链路之前。

### 5.5 运行引擎与后端调用
- 保留原「生成 JS → 前端执行」模型：`useRunner.js` 把各积木对应的 `window.xxx` 函数绑定为调用后端 Wails 方法。
- 用 `new Function`/受控 `eval` 执行拼接后的 `async` 代码；每步 `await`。
- 日志：前端 `log()` 写入 `ConsolePanel`；后端长任务通过 `runtime.EventsEmit` 推送进度，前端订阅。

### 5.6 中断功能（新增，重点）
- **后端进程管理器 `internal/procs/manager.go`**：
  - 维护一个"运行会话"，记录本次运行启动的所有子进程 PID 及其子进程树，以及需要兜底 kill 的**指定 exe 名单**（`StarRail.exe`、`BetterGI.exe`、`YuanShen.exe`、`ZenlessZoneZero.exe`、`ok-ww.exe`、`pythonw.exe`、`Client-Win64-Shipping.exe`、各 Launcher 等）。
  - 每次 spawn 第三方程序时把 `cmd.Process.Pid` 登记进当前会话。
- **`Interrupt()` 后端方法**：遍历会话进程树（`gopsutil` 或 `taskkill /T /PID`）逐个强杀，并对指定 exe 执行 `taskkill /F /IM`。清空会话。
- **前端**：点击「中断」→ 调用 `Interrupt()` → 设置运行中断标志，运行引擎在下一步 `await` 前检查标志并抛出中止错误，停止后续步骤。
- 各长任务的"轮询直到进程消失/超时"循环需响应中断信号（context 取消），及时退出。
- 运行开始时 `StartSession()`，结束（成功/失败/中断）时 `EndSession()` 清理。

### 5.7 界面模态框化
- 顶部工具栏按钮：**运行 / 中断 / 保存 / 工具 / 设置 / 工程（打开·新建·模板·关闭）**。
- `ToolsModal`、`SettingsModal`、`ProjectModal`、`SelectOneModal` 用 Vue 组件实现（替代原生 `prompt`/`alert`/手写 DOM `selectOne`）。
- 未保存关闭时用模态确认（替代 `confirm`）；窗口关闭拦截用 Wails `OnBeforeClose`。

### 5.8 命令行自动运行
- `main.go` 解析 `--auto-run <file.m7f>`，通过 Wails 启动后由前端读取（后端提供 `GetAutoRunFile()`），加载并自动运行。

### 5.9 管理员权限
- 迁移 `win.rs` 逻辑：release 下检测是否已提权，未提权则用 `runas` 重启自身（第三方脚本 kill 进程/改注册表需要）。

---

## 6. 实现轮次（按顺序推进）

> 每轮结束需保证 `wails dev` 可运行、无编译错误、无明显 lint 错误。开发命令：仓库根执行 `wails dev`（前端热重载）；构建 `wails build`。

### 第 0 轮：基线与脚手架确认
- 确认根目录 Wails+Vue 脚手架可 `wails dev` 启动（当前仅 `Greet` demo）。
- 引入前端依赖：`blockly`。清理 demo（`HelloWorld.vue`）。
- 验收：能启动空白 Vue 窗口。

### 第 1 轮：Blockly 主界面骨架（Vue）
- 实现 `BlocklyWorkspace.vue` + `useBlockly.js`：注入 Blockly、中文语言、加载/保存序列化。
- 迁移 `customBlocks.js`、`generators.js`、`toolbox.js`、`templates.js`。
- 启动即进入工作空间，默认含唯一且不可删除的「开始」积木。
- 顶部 `Toolbar.vue`（先放：保存 / 运行占位按钮）+ `ConsolePanel.vue`。
- 验收：可拖拽积木、连接、序列化到内存；开始积木不可删。

### 第 2 轮：Go 后端基础与配置
- `internal/config`：迁移 `BackendConfig`（字段与 `mflow/src-tauri/src/config.rs` 完全对齐）；`app_data_path` = `%APPDATA%/opensource/mflow`。
- 文件 API（`Exists/Mkdir/ReadTextFile/WriteTextFile`）、`GetVersion/GetNewVersion`、`OpenReleasePage`。
- 前端 `api/backend.js` 封装；工程「打开/新建/模板/保存/关闭」走后端文件读写（`.m7f`）。
- 验收：可新建/打开/保存 `.m7f`，配置读写正常。

### 第 3 轮：游戏命令后端迁移
- `internal/sysutil`（proc/reg/zip）+ `internal/games`（hsr/genshin/zzz/ww）逐一迁移原 `main.rs` 命令。
- Wails 绑定全部命令方法；前端 `generators.js` 的 `window.xxx` 绑定到这些方法。
- 长任务的超时/轮询逻辑用 `context.Context` 实现（为中断做准备）。
- 验收：各积木可单独触发对应后端逻辑（可在无第三方软件时打桩/返回明确错误）。

### 第 4 轮：运行引擎 + 链路裁剪
- `useRunner.js`：实现 5.4 的链路裁剪（仅 start 链路 + 孤岛定义类），拼接 async 代码执行。
- 绑定全部 `window.*` 到 `api/backend.js`。
- 运行日志接入 `ConsolePanel`，后端 `EventsEmit` 进度订阅。
- 验收：只运行开始链路；游离运行积木不执行；孤岛函数可被调用。

### 第 5 轮：中断功能
- `internal/procs/manager.go`：运行会话 + 进程登记 + 进程树/指定 exe kill。
- 各 `spawn` 处登记 PID；`Interrupt()` 方法；长任务响应 context 取消。
- 前端「中断」按钮 + 运行引擎中断标志与中止。
- 验收：运行中点击中断，后端强杀相关进程树与游戏进程，前端停止后续步骤且状态复位。

### 第 6 轮：界面模态框化
- `ToolsModal`（启动各 GUI + 导入/导出各游戏账号 + 列表选择 `SelectOneModal`）。
- `SettingsModal`（各软件路径 + 各任务超时）。
- `ProjectModal`（打开/新建/模板/关闭），未保存确认模态，`OnBeforeClose` 拦截。
- 验收：原 React 版所有工具/设置/工程操作在模态框中可用。

### 第 7 轮：Blockly flat 平铺美化
- 实现 `theme.js` flat 主题 + 工具箱平铺样式，接入 workspace。
- 调整分类配色与整体视觉，贴合 appearance 指南第一种平铺外观。
- 验收：外观平铺扁平、清爽，与各分类颜色协调。

### 第 8 轮：命令行自动运行 + 管理员提权 + 打包
- `--auto-run` 解析与自动运行；`win_admin` 提权（release）。
- 版本号来源、GitHub release 检查、更新提示。
- `wails build` 产物验证；`build/` 图标与 NSIS 安装配置。
- 验收：`mflow-go.exe --auto-run x.m7f` 可自动运行；release 自动提权；安装包可构建。

### 第 9 轮：收尾
- 迁移/更新 `README.md`（沿用原教程内容，替换启动命令为 Wails）。
- 存档格式：`.m7p`（JSON）。
- 全流程回归：多账号模板 → 运行 → 中断。

---

## 7. 约定与注意事项

- **响应语言**：与用户沟通使用简体中文。
- **平台**：以 Windows 为主，进程/注册表相关代码用构建标签或运行时判断隔离。
- **不要修改 `mflow/` 目录**：它是只读参考实现，所有新代码写在仓库根的 `mflow-go` 工程内。
- **进程安全**：kill 前务必确认只杀本会话登记的进程树 + 明确的游戏/工具 exe，避免误杀。
- **序列化**：统一 JSON，扩展名 `.m7p`，不兼容旧 `.m7f`（XML）。
- **Blockly 实例**勿被 Vue 响应式代理包裹（用 `markRaw`/普通闭包变量）。
- 提交（commit）仅在用户明确要求时进行。
