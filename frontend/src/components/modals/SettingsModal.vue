<script setup>
import { ref, onMounted } from 'vue';
import BaseModal from './BaseModal.vue';
import { backend } from '../../api/backend';
import { dialogs } from '../../composables/useDialogs';

const emit = defineEmits(['close']);

const cfg = ref(null);

const pathFields = [
  { key: 'm7_path', label: '三月七小助手 文件夹' },
  { key: 'better_gi_path', label: '更好的原神(BetterGI) 文件夹' },
  { key: 'zzzod_path', label: '绝区零一条龙 文件夹' },
  { key: 'genshin_auto_login_path', label: '原神自动登录器 文件夹' },
  { key: 'ok_ww_path', label: 'OK鸣潮 文件夹' },
];

const timeoutFields = [
  { key: 'full_run_timeout_minutes', label: '星铁完整运行' },
  { key: 'daily_mission_timeout_minutes', label: '星铁每日任务' },
  { key: 'refresh_stamina_timeout_minutes', label: '星铁刷体力' },
  { key: 'simulated_universe_timeout_minutes', label: '星铁模拟宇宙' },
  { key: 'farming_timeout_minutes', label: '星铁锄大地' },
  { key: 'better_gi_timeout_minutes', label: '原神一条龙' },
  { key: 'better_gi_scheduler_timeout_minutes', label: '原神调度器' },
  { key: 'zzzod_timeout_minutes', label: '绝区零一条龙' },
  { key: 'ok_ww_timeout_minutes', label: '鸣潮' },
];

onMounted(async () => {
  try {
    cfg.value = await backend.loadBackendConfig();
  } catch (e) {
    await dialogs.alert('加载配置失败: ' + e);
    emit('close');
  }
});

async function save() {
  try {
    // 确保超时为数字
    for (const f of timeoutFields) {
      cfg.value[f.key] = Number(cfg.value[f.key]) || 0;
    }
    await backend.saveBackendConfig(cfg.value);
    await dialogs.alert('设置已保存。');
    emit('close');
  } catch (e) {
    await dialogs.alert('保存失败: ' + e);
  }
}
</script>

<template>
  <BaseModal title="设置" width="560px" @close="emit('close')">
    <div v-if="cfg">
      <h4 class="section">软件路径</h4>
      <div v-for="f in pathFields" :key="f.key" class="field">
        <label>{{ f.label }}</label>
        <input v-model="cfg[f.key]" class="input" placeholder="请输入文件夹路径" />
      </div>

      <h4 class="section">错误处理</h4>
      <label class="check">
        <input type="checkbox" v-model="cfg.stop_on_script_error" />
        运行脚本错误时停止流程
      </label>
      <label class="check">
        <input type="checkbox" v-model="cfg.stop_on_account_error" />
        备份或恢复账号错误时停止流程
      </label>

      <h4 class="section">启动与日志</h4>
      <label class="check">
        <input type="checkbox" v-model="cfg.keep_console_window" />
        启动时保留黑窗口（重启生效）
      </label>
      <label class="check">
        <input type="checkbox" v-model="cfg.script_log_to_stdout" />
        将脚本日志打印到标准输出（重新运行生效）
      </label>
      <label class="check">
        <input type="checkbox" v-model="cfg.script_log_to_app_console" />
        将脚本日志打印到应用程序控制台（重新运行生效）
      </label>

      <h4 class="section">任务超时（分钟）</h4>
      <div class="grid">
        <div v-for="f in timeoutFields" :key="f.key" class="field small">
          <label>{{ f.label }}</label>
          <input v-model.number="cfg[f.key]" type="number" min="0" class="input" />
        </div>
      </div>
    </div>
    <div v-else class="loading">加载中...</div>

    <template #footer>
      <button class="btn" @click="emit('close')">取消</button>
      <button class="btn primary" :disabled="!cfg" @click="save">保存</button>
    </template>
  </BaseModal>
</template>

<style scoped>
.section {
  margin: 4px 0 10px;
  font-size: 13px;
  color: #4a5568;
  border-left: 3px solid #2f855a;
  padding-left: 8px;
}
.field {
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.field label {
  font-size: 12px;
  color: #4a5568;
}
.input {
  width: 100%;
  box-sizing: border-box;
  padding: 7px 9px;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  font-size: 13px;
}
.check {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
  color: #2d3748;
  cursor: pointer;
}
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.field.small label {
  font-size: 12px;
}
.loading {
  padding: 20px;
  text-align: center;
  color: #a0aec0;
}
.btn {
  padding: 6px 14px;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  font-size: 13px;
}
.btn.primary {
  background: #2f855a;
  color: #fff;
  border-color: #2f855a;
}
.btn.primary:disabled {
  opacity: 0.5;
}
</style>
