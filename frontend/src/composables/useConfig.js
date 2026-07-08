import { backend } from '../api/backend';

const FILE_NAME = 'frontend_config.json';

async function configPath() {
  const dir = await backend.appDataPath();
  return dir + '/' + FILE_NAME;
}

// 读取上次打开的工程路径（用于应用重启后恢复）。
export async function loadLastFile() {
  try {
    const p = await configPath();
    if (await backend.exists(p)) {
      const text = await backend.readTextFile(p);
      const data = JSON.parse(text);
      return data.lastFile || '';
    }
  } catch (_) {
    /* ignore */
  }
  return '';
}

// 保存上次打开的工程路径（空串表示未关联任何文件）。
export async function saveLastFile(lastFile) {
  try {
    const p = await configPath();
    await backend.writeTextFile(p, JSON.stringify({ lastFile: lastFile || '' }));
  } catch (_) {
    /* ignore */
  }
}
