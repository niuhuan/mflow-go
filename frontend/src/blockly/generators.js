import { javascriptGenerator, Order } from 'blockly/javascript';

javascriptGenerator.forBlock['start_flow'] = function () {
  // 起始节点本身不产码，运行引擎从它开始沿链路生成后继语句。
  return '';
};

javascriptGenerator.forBlock['define_variable'] = function (block) {
  const name = block.getFieldValue('VAR_NAME') || 'myVar';
  const value = javascriptGenerator.valueToCode(block, 'VALUE', Order.ASSIGNMENT) || 'null';
  return `window.setVar('${name}', ${value});\n`;
};

javascriptGenerator.forBlock['load_account'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "''";
  return `await window.loadAccount(${accountName});\n`;
};

javascriptGenerator.forBlock['save_account'] = function (block) {
  const saveName = javascriptGenerator.valueToCode(block, 'SAVE_NAME', Order.ATOMIC) || "''";
  return `await window.saveAccount(${saveName});\n`;
};

javascriptGenerator.forBlock['wait_seconds'] = function (block) {
  const timeValue = block.getFieldValue('TIME_VALUE');
  const timeUnit = block.getFieldValue('TIME_UNIT');
  return `await window.wait(${timeValue}, '${timeUnit}');\n`;
};

javascriptGenerator.forBlock['full_run'] = function () {
  return 'await window.fullRun();\n';
};

javascriptGenerator.forBlock['daily_mission'] = function () {
  return 'await window.dailyMission();\n';
};

javascriptGenerator.forBlock['refresh_stamina'] = function () {
  return 'await window.refreshStamina();\n';
};

javascriptGenerator.forBlock['simulated_universe'] = function () {
  return 'await window.simulatedUniverse();\n';
};

javascriptGenerator.forBlock['farming'] = function () {
  return 'await window.farming();\n';
};

javascriptGenerator.forBlock['close_game'] = function () {
  return 'await window.closeGame();\n';
};

javascriptGenerator.forBlock['controls_whileUntil'] = function (block) {
  const argument0 = javascriptGenerator.valueToCode(block, 'BOOL', Order.NONE) || 'false';
  const branch = javascriptGenerator.statementToCode(block, 'DO');
  return 'while (' + argument0 + ') {\n' + branch + '}\n';
};

javascriptGenerator.forBlock['controls_if'] = function (block) {
  let n = 0;
  let code = '';
  let branchCode;
  let conditionCode;
  do {
    conditionCode = javascriptGenerator.valueToCode(block, 'IF' + n, Order.NONE) || 'false';
    branchCode = javascriptGenerator.statementToCode(block, 'DO' + n);
    code += (n === 0 ? 'if (' : 'else if (') + conditionCode + ') {\n' + branchCode + '}\n';
    n++;
  } while (block.getInput('IF' + n));
  return code;
};

javascriptGenerator.forBlock['controls_if_else'] = function (block) {
  const conditionCode = javascriptGenerator.valueToCode(block, 'IF0', Order.NONE) || 'false';
  const thenCode = javascriptGenerator.statementToCode(block, 'DO0');
  const elseCode = javascriptGenerator.statementToCode(block, 'ELSE');
  return 'if (' + conditionCode + ') {\n' + thenCode + '} else {\n' + elseCode + '}\n';
};

javascriptGenerator.forBlock['logic_boolean'] = function (block) {
  const code = block.getFieldValue('BOOL') === 'TRUE' ? 'true' : 'false';
  return [code, Order.ATOMIC];
};

javascriptGenerator.forBlock['logic_operation'] = function (block) {
  const operator = block.getFieldValue('OP') === 'AND' ? '&&' : '||';
  const order = operator === '&&' ? Order.LOGICAL_AND : Order.LOGICAL_OR;
  const argument0 = javascriptGenerator.valueToCode(block, 'A', order) || 'false';
  const argument1 = javascriptGenerator.valueToCode(block, 'B', order) || 'false';
  return [argument0 + ' ' + operator + ' ' + argument1, order];
};

javascriptGenerator.forBlock['logic_negate'] = function (block) {
  const argument0 = javascriptGenerator.valueToCode(block, 'BOOL', Order.LOGICAL_NOT) || 'false';
  return ['!' + argument0, Order.LOGICAL_NOT];
};

javascriptGenerator.forBlock['logic_compare'] = function (block) {
  const OPERATORS = { EQ: '==', NEQ: '!=', LT: '<', LTE: '<=', GT: '>', GTE: '>=' };
  const operator = OPERATORS[block.getFieldValue('OP')] || '==';
  const order = Order.RELATIONAL;
  const argument0 = javascriptGenerator.valueToCode(block, 'A', order) || '0';
  const argument1 = javascriptGenerator.valueToCode(block, 'B', order) || '0';
  return [argument0 + ' ' + operator + ' ' + argument1, order];
};

javascriptGenerator.forBlock['custom_function_def'] = function (block) {
  const funcName = block.getFieldValue('FUNC_NAME') || 'myFunc';
  const params = block.getFieldValue('PARAMS') || '';
  const branch = javascriptGenerator.statementToCode(block, 'DO');
  const cleanParams = params
    .split(',')
    .map((p) => p.trim())
    .filter((p) => p.length > 0)
    .join(', ');
  return `async function ${funcName}(${cleanParams}) {\n${branch}}\n`;
};

javascriptGenerator.forBlock['custom_function_call'] = function (block) {
  const funcName = block.getFieldValue('FUNC_NAME') || 'myFunc';
  let args = javascriptGenerator.valueToCode(block, 'ARGS', Order.NONE) || '';
  if (!args || args === "''" || args.trim() === '') {
    return `await ${funcName}();\n`;
  }
  if (args.startsWith("'") && args.endsWith("'")) {
    const argsValue = args.slice(1, -1);
    if (argsValue.indexOf(',') > -1) {
      args = argsValue
        .split(',')
        .map((arg) => `'${arg}'`)
        .join(',');
    }
  }
  return `await ${funcName}(${args});\n`;
};

javascriptGenerator.forBlock['custom_parameter'] = function (block) {
  const paramName = block.getFieldValue('PARAM_NAME') || 'param1';
  return [paramName, Order.ATOMIC];
};

javascriptGenerator.forBlock['run_better_gi'] = function () {
  return 'await window.runBetterGi();\n';
};

javascriptGenerator.forBlock['run_better_gi_by_config'] = function (block) {
  const configName = javascriptGenerator.valueToCode(block, 'CONFIG_NAME', Order.ATOMIC) || "'配置文件'";
  return `await window.runBetterGiByConfig(${configName});\n`;
};

javascriptGenerator.forBlock['close_gi'] = function () {
  return 'await window.closeGi();\n';
};

javascriptGenerator.forBlock['export_gi_account'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "''";
  return `await window.exportGiAccount(${accountName});\n`;
};

javascriptGenerator.forBlock['import_gi_account'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "''";
  return `await window.importGiAccount(${accountName});\n`;
};

javascriptGenerator.forBlock['export_zzz_account'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "''";
  return `await window.exportZzzAccount(${accountName});\n`;
};

javascriptGenerator.forBlock['import_zzz_account'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "''";
  return `await window.importZzzAccount(${accountName});\n`;
};

javascriptGenerator.forBlock['close_zzz'] = function () {
  return 'await window.closeZzz();\n';
};

javascriptGenerator.forBlock['math_number'] = function (block) {
  return [String(block.getFieldValue('NUM')), Order.ATOMIC];
};

javascriptGenerator.forBlock['current_hour_24'] = function () {
  return ['new Date().getHours()', Order.ATOMIC];
};

javascriptGenerator.forBlock['current_time_minus_4h_day_of_week'] = function () {
  return [
    '(function() { var date = new Date(); date.setTime(date.getTime() - 4 * 60 * 60 * 1000); return date.getDay(); })()',
    Order.ATOMIC,
  ];
};

javascriptGenerator.forBlock['wait_until_time'] = function (block) {
  const hour = block.getFieldValue('HOUR');
  const minute = block.getFieldValue('MINUTE');
  return `await window.waitUntilTime(${hour}, ${minute});\n`;
};

javascriptGenerator.forBlock['print_variable'] = function (block) {
  const value = javascriptGenerator.valueToCode(block, 'VALUE', Order.NONE) || "''";
  return `window.log('打印变量: ' + ${value});\n`;
};

javascriptGenerator.forBlock['run_zzzod'] = function () {
  return 'await window.runZzzod();\n';
};

javascriptGenerator.forBlock['start_ok_ww_daily'] = function () {
  return 'await window.startOkWwDaily();\n';
};

javascriptGenerator.forBlock['start_ok_ww_weekly'] = function () {
  return 'await window.startOkWwWeekly();\n';
};

javascriptGenerator.forBlock['kill_ok_ww'] = function () {
  return 'await window.killOkWw();\n';
};

javascriptGenerator.forBlock['run_better_gi_scheduler'] = function (block) {
  const groups = javascriptGenerator.valueToCode(block, 'GROUPS', Order.ATOMIC) || "''";
  return `await window.runBetterGiScheduler(${groups});\n`;
};

javascriptGenerator.forBlock['run_command'] = function (block) {
  const command = javascriptGenerator.valueToCode(block, 'COMMAND', Order.ATOMIC) || "''";
  return `await window.runCommand(${command});\n`;
};

javascriptGenerator.forBlock['run_command_background'] = function (block) {
  const command = javascriptGenerator.valueToCode(block, 'COMMAND', Order.ATOMIC) || "''";
  return `await window.runCommandBackground(${command});\n`;
};

javascriptGenerator.forBlock['genshin_auto_login'] = function (block) {
  const accountName = javascriptGenerator.valueToCode(block, 'ACCOUNT_NAME', Order.ATOMIC) || "'account1'";
  return `await window.genshinAutoLogin(${accountName});\n`;
};

javascriptGenerator.forBlock['clear_hsr_reg'] = function () {
  return 'await window.clearHsrReg();\n';
};

javascriptGenerator.forBlock['clear_gi_reg'] = function () {
  return 'await window.clearGiReg();\n';
};

javascriptGenerator.forBlock['clear_zzz_reg'] = function () {
  return 'await window.clearZzzReg();\n';
};

javascriptGenerator.forBlock['clear_ww_reg'] = function () {
  return 'await window.clearWwReg();\n';
};
