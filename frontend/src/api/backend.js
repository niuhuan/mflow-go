// 封装 Wails 后端调用（wailsjs/go/main/App）。
import * as App from '../../wailsjs/go/main/App';

export const backend = {
  // 版本
  getVersion: () => App.GetVersion(),
  getNewVersion: () => App.GetNewVersion(),
  openReleasePage: () => App.OpenReleasePage(),

  // 文件
  appDataPath: () => App.AppDataPath(),
  exists: (p) => App.Exists(p),
  mkdir: (p) => App.Mkdir(p),
  readTextFile: (p) => App.ReadTextFile(p),
  writeTextFile: (p, c) => App.WriteTextFile(p, c),

  // 配置
  loadBackendConfig: () => App.LoadBackendConfig(),
  saveBackendConfig: (cfg) => App.SaveBackendConfig(cfg),

  // 工程对话框
  openProjectDialog: () => App.OpenProjectDialog(),
  saveProjectDialog: () => App.SaveProjectDialog(),

  // 自动运行
  getAutoRunFile: () => App.GetAutoRunFile(),

  // 运行会话与中断
  startRun: () => App.StartRun(),
  endRun: () => App.EndRun(),
  interrupt: () => App.Interrupt(),
  isRunning: () => App.IsRunning(),

  // 星铁
  loadAccount: (name) => App.LoadAccount(name),
  saveAccount: (name) => App.SaveAccount(name),
  fullRun: () => App.FullRun(),
  dailyMission: () => App.DailyMission(),
  refreshStamina: () => App.RefreshStamina(),
  simulatedUniverse: () => App.SimulatedUniverse(),
  farming: () => App.Farming(),
  closeGame: () => App.CloseGame(),
  clearHsrReg: () => App.ClearHsrReg(),
  runM7Launcher: () => App.RunM7Launcher(),
  listAccounts: () => App.ListAccounts(),
  getAccountUID: () => App.GetAccountUID(),
  exportAccount: (name, u, p) => App.ExportAccount(name, u, p),

  // 原神
  runBetterGi: () => App.RunBetterGi(),
  runBetterGiByConfig: (c) => App.RunBetterGiByConfig(c),
  runBetterGiScheduler: (g) => App.RunBetterGiScheduler(g),
  runBetterGiGui: () => App.RunBetterGiGui(),
  closeGi: () => App.CloseGi(),
  clearGiReg: () => App.ClearGiReg(),
  listGiAccounts: () => App.ListGiAccounts(),
  exportGiAccount: (name) => App.ExportGiAccount(name),
  importGiAccount: (name) => App.ImportGiAccount(name),
  genshinAutoLogin: (name) => App.GenshinAutoLogin(name),

  // 绝区零
  runZzzod: () => App.RunZzzod(),
  runZzzodGui: () => App.RunZzzodGui(),
  closeZzz: () => App.CloseZzz(),
  clearZzzReg: () => App.ClearZzzReg(),
  listZzzAccounts: () => App.ListZzzAccounts(),
  exportZzzAccount: (name) => App.ExportZzzAccount(name),
  importZzzAccount: (name) => App.ImportZzzAccount(name),

  // 鸣潮
  startOkWwDaily: () => App.StartOkWwDaily(),
  startOkWwWeekly: () => App.StartOkWwWeekly(),
  killOkWw: () => App.KillOkWw(),
  clearWwReg: () => App.ClearWwReg(),

  // 通用命令
  runCommand: (c) => App.RunCommand(c),
  runCommandBackground: (c) => App.RunCommandBackground(c),
};
