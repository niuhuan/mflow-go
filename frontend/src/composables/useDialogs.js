import { reactive } from 'vue';

// 全局对话框状态，由 DialogHost 渲染，业务侧通过 Promise 使用。
export const dialogState = reactive({
  visible: false,
  type: 'alert', // alert | confirm | prompt | select
  title: '',
  message: '',
  value: '',
  items: [],
  resolver: null,
});

function open(opts) {
  return new Promise((resolve) => {
    dialogState.type = opts.type;
    dialogState.title = opts.title || '';
    dialogState.message = opts.message || '';
    dialogState.value = opts.value || '';
    dialogState.items = opts.items || [];
    dialogState.resolver = resolve;
    dialogState.visible = true;
  });
}

function finish(result) {
  const resolve = dialogState.resolver;
  dialogState.visible = false;
  dialogState.resolver = null;
  dialogState.items = [];
  if (resolve) resolve(result);
}

export const dialogs = {
  alert: (message, title = '提示') => open({ type: 'alert', title, message }),
  confirm: (message, title = '确认') => open({ type: 'confirm', title, message }),
  prompt: (title = '请输入', value = '') => open({ type: 'prompt', title, value }),
  select: (items, title = '请选择') => open({ type: 'select', title, items }),
  _finish: finish,
};
