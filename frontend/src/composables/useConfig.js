import { backend } from '../api/backend';

const FILE_NAME = 'frontend_config.json';

async function configPath() {
  const dir = await backend.appDataPath();
  return dir + '/' + FILE_NAME;
}

// 读取整个前端配置对象。
export async function loadConfig() {
  try {
    const p = await configPath();
    if (await backend.exists(p)) {
      const text = await backend.readTextFile(p);
      return JSON.parse(text) || {};
    }
  } catch (_) {
    /* ignore */
  }
  return {};
}

// 合并写入前端配置（读取现有内容后合并，避免覆盖其它字段）。
export async function updateConfig(patch) {
  const cur = await loadConfig();
  const next = { ...cur, ...patch };
  try {
    const p = await configPath();
    await backend.writeTextFile(p, JSON.stringify(next));
  } catch (_) {
    /* ignore */
  }
  return next;
}

// 上次打开的工程路径
export async function loadLastFile() {
  return (await loadConfig()).lastFile || '';
}

export async function saveLastFile(lastFile) {
  await updateConfig({ lastFile: lastFile || '' });
}

// 控制台位置：'right' | 'bottom'，默认右侧
export async function loadConsolePosition() {
  const pos = (await loadConfig()).consolePosition;
  return pos === 'bottom' ? 'bottom' : 'right';
}

export async function saveConsolePosition(position) {
  await updateConfig({ consolePosition: position });
}

// 控制台尺寸：右侧时为宽度、下方时为高度（像素）。
export async function loadConsoleSize() {
  const cfg = await loadConfig();
  const width = Number(cfg.consoleWidth) || 380;
  const height = Number(cfg.consoleHeight) || 180;
  return { width, height };
}

export async function saveConsoleWidth(width) {
  await updateConfig({ consoleWidth: Math.round(width) });
}

export async function saveConsoleHeight(height) {
  await updateConfig({ consoleHeight: Math.round(height) });
}
