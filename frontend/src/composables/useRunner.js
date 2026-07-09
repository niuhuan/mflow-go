import { ref } from 'vue';
import { javascriptGenerator } from 'blockly/javascript';
import { backend } from '../api/backend';

// 仅这些"定义类"孤岛积木会被 hoist 到主流程之前。
const DEFINITION_TYPES = ['custom_function_def', 'define_variable'];

export function useRunner(getWorkspace, log) {
  const running = ref(false);
  let runId = 0; // 每次运行递增，用于让上一轮的计时器/回调立即失效
  let interrupted = false;
  const policy = { stopScript: false, stopAccount: true };

  // 是否应中止：被中断，或已经不是当前这一轮运行。
  function aborted(myId) {
    return interrupted || myId !== runId;
  }

  // 可中断的睡眠：每 200ms 检查一次，中断时立即 reject。
  function sleep(ms, myId) {
    return new Promise((resolve, reject) => {
      const start = Date.now();
      const iv = setInterval(() => {
        if (aborted(myId)) {
          clearInterval(iv);
          reject(new Error('已中断'));
          return;
        }
        if (Date.now() - start >= ms) {
          clearInterval(iv);
          resolve();
        }
      }, 200);
    });
  }

  function codeOf(block) {
    let c = javascriptGenerator.blockToCode(block);
    if (Array.isArray(c)) c = c[0];
    return c || '';
  }

  // 生成待执行代码：仅「开始」链路 + 孤岛的定义类积木（hoist 到前面）。
  function buildCode(workspace) {
    const starts = workspace.getBlocksByType('start_flow', false);
    if (starts.length === 0) {
      throw new Error('找不到"开始"积木。');
    }
    if (starts.length > 1) {
      throw new Error('只能有一个"开始"积木，请移除多余的。');
    }
    const startBlock = starts[0];

    javascriptGenerator.init(workspace);
    const mainCode = codeOf(startBlock);

    let defs = '';
    for (const b of workspace.getTopBlocks(true)) {
      if (b.id === startBlock.id) continue;
      if (DEFINITION_TYPES.includes(b.type)) {
        defs += codeOf(b);
      }
    }
    return defs + mainCode;
  }

  function bindWindow(myId) {
    const w = window;
    const scope = {};
    w.__mflowVars = scope;
    w.log = (msg) => log(msg);
    w.setVar = (name, value) => {
      scope[name] = value;
      log(`定义变量: ${name} = ${value}`);
    };
    w.getVar = (name) => scope[name];

    // category: 'account' 受"账号错误停止"控制；其余 'script' 受"脚本错误停止"控制。
    const call = (label, fn, category) => async (...args) => {
      if (aborted(myId)) throw new Error('已中断');
      if (label) log(label(...args));
      try {
        return await fn(...args);
      } catch (e) {
        if (aborted(myId)) throw new Error('已中断'); // 中断优先，始终停止
        const stop = category === 'account' ? policy.stopAccount : policy.stopScript;
        const msg = e && e.message ? e.message : String(e);
        log((category === 'account' ? '账号操作失败: ' : '脚本执行失败: ') + msg);
        if (stop) throw e;
        return undefined; // 忽略错误，继续后续流程
      }
    };

    // 账号类（受"备份或恢复账号错误是停止流程"控制，默认停止）
    w.loadAccount = call((n) => '加载星铁账号数据: ' + n, backend.loadAccount, 'account');
    w.saveAccount = call((n) => '保存星铁账号数据: ' + n, backend.saveAccount, 'account');
    w.exportGiAccount = call((n) => '导出原神账号: ' + n, backend.exportGiAccount, 'account');
    w.importGiAccount = call((n) => '导入原神账号: ' + n, backend.importGiAccount, 'account');
    w.exportZzzAccount = call((n) => '导出绝区零账号: ' + n, backend.exportZzzAccount, 'account');
    w.importZzzAccount = call((n) => '导入绝区零账号: ' + n, backend.importZzzAccount, 'account');

    // 脚本/运行类（受"运行脚本错误是停止流程"控制，默认继续）
    w.fullRun = call(() => '执行星铁完整运行...', backend.fullRun, 'script');
    w.dailyMission = call(() => '执行星铁每日任务...', backend.dailyMission, 'script');
    w.refreshStamina = call(() => '执行星铁刷体力...', backend.refreshStamina, 'script');
    w.simulatedUniverse = call(() => '执行星铁模拟宇宙...', backend.simulatedUniverse, 'script');
    w.farming = call(() => '执行星铁锄大地...', backend.farming, 'script');
    w.closeGame = call(() => '关闭星铁...', backend.closeGame, 'script');
    w.clearHsrReg = call(() => '清除星铁注册表...', backend.clearHsrReg, 'script');

    w.runBetterGi = call(() => '运行原神一条龙...', backend.runBetterGi, 'script');
    w.runBetterGiByConfig = call((c) => '运行原神一条龙，配置: ' + c, backend.runBetterGiByConfig, 'script');
    w.runBetterGiScheduler = call((g) => '运行原神调度器: ' + g, backend.runBetterGiScheduler, 'script');
    w.closeGi = call(() => '关闭原神...', backend.closeGi, 'script');
    w.clearGiReg = call(() => '清除原神注册表...', backend.clearGiReg, 'script');
    w.genshinAutoLogin = call((n) => '原神自动登录: ' + n, backend.genshinAutoLogin, 'script');

    w.runZzzod = call(() => '运行绝区零一条龙...', backend.runZzzod, 'script');
    w.closeZzz = call(() => '关闭绝区零...', backend.closeZzz, 'script');
    w.clearZzzReg = call(() => '清除绝区零注册表...', backend.clearZzzReg, 'script');

    w.startOkWwDaily = call(() => '运行鸣潮日常...', backend.startOkWwDaily, 'script');
    w.startOkWwWeekly = call(() => '运行鸣潮周常...', backend.startOkWwWeekly, 'script');
    w.killOkWw = call(() => '关闭鸣潮...', backend.killOkWw, 'script');
    w.clearWwReg = call(() => '清除鸣潮注册表...', backend.clearWwReg, 'script');

    w.runCommand = call((c) => '运行命令: ' + c, backend.runCommand, 'script');
    w.runCommandBackground = call((c) => '后台运行命令: ' + c, backend.runCommandBackground, 'script');

    // 纯前端行为：等待（可快速中断）。
    w.wait = async (timeValue, timeUnit = 'SECONDS') => {
      let seconds = timeValue;
      let unit = '秒';
      if (timeUnit === 'MINUTES') { seconds = timeValue * 60; unit = '分钟'; }
      else if (timeUnit === 'HOURS') { seconds = timeValue * 3600; unit = '小时'; }
      log(`等待 ${timeValue} ${unit}...`);
      await sleep(seconds * 1000, myId);
      log('等待完成。');
    };

    w.waitUntilTime = async (targetHour, targetMinute) => {
      log(`等待到 ${String(targetHour).padStart(2, '0')}:${String(targetMinute).padStart(2, '0')}`);
      // 循环等待，每秒检查一次时间；sleep 会在中断时立即抛出。
      // eslint-disable-next-line no-constant-condition
      while (true) {
        if (aborted(myId)) throw new Error('已中断');
        const now = new Date();
        if (now.getHours() === targetHour && now.getMinutes() === targetMinute) {
          log('已到达指定时间。');
          return;
        }
        await sleep(1000, myId);
      }
    };
  }

  async function run() {
    if (running.value) return;
    const workspace = getWorkspace();
    if (!workspace) return;

    runId += 1;
    const myId = runId;
    interrupted = false;
    running.value = true;

    let code;
    try {
      code = buildCode(workspace);
    } catch (e) {
      log('执行失败: ' + e.message);
      running.value = false;
      return;
    }

    // 读取错误处理策略
    try {
      const cfg = await backend.loadBackendConfig();
      policy.stopScript = !!cfg.stop_on_script_error;
      policy.stopAccount = cfg.stop_on_account_error !== false; // 默认开
    } catch (_) {
      policy.stopScript = false;
      policy.stopAccount = true;
    }

    bindWindow(myId);
    await backend.startRun();
    log('开始运行...');

    try {
      // eslint-disable-next-line no-new-func
      const fn = new Function(`return (async () => { ${code} })();`);
      await fn();
      log(aborted(myId) ? '已中断。' : '运行完成。');
    } catch (e) {
      log('运行结束: ' + (e && e.message ? e.message : e));
    } finally {
      await backend.endRun();
      running.value = false;
      interrupted = false;
    }
  }

  async function interrupt() {
    if (!running.value) return;
    interrupted = true;
    log('正在中断...');
    try {
      await backend.interrupt();
    } catch (e) {
      log('中断出错: ' + e);
    }
  }

  return { running, run, interrupt };
}
