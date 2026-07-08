import { ref } from 'vue';
import { javascriptGenerator } from 'blockly/javascript';
import { backend } from '../api/backend';

// 仅这些"定义类"孤岛积木会被 hoist 到主流程之前。
const DEFINITION_TYPES = ['custom_function_def', 'define_variable'];

export function useRunner(getWorkspace, log) {
  const running = ref(false);
  let interrupted = false;

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

  // 运行时可被中断的保护点：每次调用后端前检查。
  function guard() {
    if (interrupted) {
      throw new Error('已中断');
    }
  }

  function bindWindow() {
    const w = window;
    const scope = {};
    w.__mflowVars = scope;
    w.log = (msg) => log(msg);
    w.setVar = (name, value) => {
      scope[name] = value;
      log(`定义变量: ${name} = ${value}`);
    };
    w.getVar = (name) => scope[name];

    const call = (label, fn) => async (...args) => {
      guard();
      if (label) log(label(...args));
      return fn(...args);
    };

    w.loadAccount = call((n) => '加载星铁账号数据: ' + n, backend.loadAccount);
    w.saveAccount = call((n) => '保存星铁账号数据: ' + n, backend.saveAccount);
    w.fullRun = call(() => '执行星铁完整运行...', backend.fullRun);
    w.dailyMission = call(() => '执行星铁每日任务...', backend.dailyMission);
    w.refreshStamina = call(() => '执行星铁刷体力...', backend.refreshStamina);
    w.simulatedUniverse = call(() => '执行星铁模拟宇宙...', backend.simulatedUniverse);
    w.farming = call(() => '执行星铁锄大地...', backend.farming);
    w.closeGame = call(() => '关闭星铁...', backend.closeGame);
    w.clearHsrReg = call(() => '清除星铁注册表...', backend.clearHsrReg);

    w.runBetterGi = call(() => '运行原神一条龙...', backend.runBetterGi);
    w.runBetterGiByConfig = call((c) => '运行原神一条龙，配置: ' + c, backend.runBetterGiByConfig);
    w.runBetterGiScheduler = call((g) => '运行原神调度器: ' + g, backend.runBetterGiScheduler);
    w.closeGi = call(() => '关闭原神...', backend.closeGi);
    w.clearGiReg = call(() => '清除原神注册表...', backend.clearGiReg);
    w.exportGiAccount = call((n) => '导出原神账号: ' + n, backend.exportGiAccount);
    w.importGiAccount = call((n) => '导入原神账号: ' + n, backend.importGiAccount);
    w.genshinAutoLogin = call((n) => '原神自动登录: ' + n, backend.genshinAutoLogin);

    w.runZzzod = call(() => '运行绝区零一条龙...', backend.runZzzod);
    w.closeZzz = call(() => '关闭绝区零...', backend.closeZzz);
    w.clearZzzReg = call(() => '清除绝区零注册表...', backend.clearZzzReg);
    w.exportZzzAccount = call((n) => '导出绝区零账号: ' + n, backend.exportZzzAccount);
    w.importZzzAccount = call((n) => '导入绝区零账号: ' + n, backend.importZzzAccount);

    w.startOkWwDaily = call(() => '运行鸣潮日常...', backend.startOkWwDaily);
    w.startOkWwWeekly = call(() => '运行鸣潮周常...', backend.startOkWwWeekly);
    w.killOkWw = call(() => '关闭鸣潮...', backend.killOkWw);
    w.clearWwReg = call(() => '清除鸣潮注册表...', backend.clearWwReg);

    w.runCommand = call((c) => '运行命令: ' + c, backend.runCommand);
    w.runCommandBackground = call((c) => '后台运行命令: ' + c, backend.runCommandBackground);

    // 纯前端行为：等待，等待期间也响应中断。
    w.wait = (timeValue, timeUnit = 'SECONDS') =>
      new Promise((resolve, reject) => {
        let seconds = timeValue;
        let unit = '秒';
        if (timeUnit === 'MINUTES') { seconds = timeValue * 60; unit = '分钟'; }
        else if (timeUnit === 'HOURS') { seconds = timeValue * 3600; unit = '小时'; }
        log(`等待 ${timeValue} ${unit}...`);
        const timer = setInterval(() => {
          if (interrupted) {
            clearInterval(timer);
            reject(new Error('已中断'));
          }
        }, 500);
        setTimeout(() => {
          clearInterval(timer);
          if (interrupted) { reject(new Error('已中断')); return; }
          log('等待完成。');
          resolve();
        }, seconds * 1000);
      });

    w.waitUntilTime = (targetHour, targetMinute) =>
      new Promise((resolve, reject) => {
        log(`等待到 ${String(targetHour).padStart(2, '0')}:${String(targetMinute).padStart(2, '0')}`);
        const check = () => {
          if (interrupted) { reject(new Error('已中断')); return; }
          const now = new Date();
          if (now.getHours() === targetHour && now.getMinutes() === targetMinute) {
            log('已到达指定时间。');
            resolve();
          } else {
            setTimeout(check, 30000);
          }
        };
        check();
      });
  }

  async function run() {
    if (running.value) return;
    const workspace = getWorkspace();
    if (!workspace) return;

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

    bindWindow();
    await backend.startRun();
    log('开始运行...');

    try {
      // eslint-disable-next-line no-new-func
      const fn = new Function(`return (async () => { ${code} })();`);
      await fn();
      log('运行完成。');
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
