// 工具箱定义。注意：不包含「开始」积木（start_flow 已存在且唯一、不可删除）。
export const toolboxXml = `
  <xml>
    <category name="星铁" colour="25">
      <block type="load_account">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="save_account">
        <value name="SAVE_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="full_run"></block>
      <block type="daily_mission"></block>
      <block type="refresh_stamina"></block>
      <block type="simulated_universe"></block>
      <block type="farming"></block>
      <block type="close_game"></block>
      <block type="clear_hsr_reg"></block>
    </category>
    <category name="原神" colour="40">
      <block type="run_better_gi"></block>
      <block type="run_better_gi_by_config">
        <value name="CONFIG_NAME">
          <shadow type="text"><field name="TEXT">配置文件</field></shadow>
        </value>
      </block>
      <block type="run_better_gi_scheduler">
        <value name="GROUPS">
          <shadow type="text"><field name="TEXT">配置组名称1 配置组名称2 退出程序</field></shadow>
        </value>
      </block>
      <block type="close_gi"></block>
      <block type="export_gi_account">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="import_gi_account">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="genshin_auto_login">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">account1</field></shadow>
        </value>
      </block>
      <block type="clear_gi_reg"></block>
    </category>
    <category name="ZZZ" colour="10">
      <block type="run_zzzod"></block>
      <block type="close_zzz"></block>
      <block type="export_zzz_account">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="import_zzz_account">
        <value name="ACCOUNT_NAME">
          <shadow type="text"><field name="TEXT">默认账号</field></shadow>
        </value>
      </block>
      <block type="clear_zzz_reg"></block>
    </category>
    <category name="鸣潮" colour="55">
      <block type="start_ok_ww_daily"></block>
      <block type="start_ok_ww_weekly"></block>
      <block type="kill_ok_ww"></block>
      <block type="clear_ww_reg"></block>
    </category>
    <category name="通用" colour="80">
      <block type="wait_seconds">
        <field name="TIME_VALUE">10</field>
        <field name="TIME_UNIT">SECONDS</field>
      </block>
      <block type="wait_until_time">
        <field name="HOUR">4</field>
        <field name="MINUTE">10</field>
      </block>
      <block type="print_variable">
        <value name="VALUE">
          <shadow type="text"><field name="TEXT">变量名</field></shadow>
        </value>
      </block>
      <block type="run_command">
        <value name="COMMAND">
          <shadow type="text"><field name="TEXT">echo Hello World</field></shadow>
        </value>
      </block>
      <block type="run_command_background">
        <value name="COMMAND">
          <shadow type="text"><field name="TEXT">echo Hello World</field></shadow>
        </value>
      </block>
    </category>
    <category name="流程" colour="120">
      <block type="controls_whileUntil"></block>
      <block type="controls_repeat_ext">
        <value name="TIMES">
          <shadow type="math_number"><field name="NUM">10</field></shadow>
        </value>
      </block>
      <block type="controls_if"></block>
      <block type="controls_if_else"></block>
      <block type="controls_forEach"></block>
    </category>
    <category name="逻辑" colour="150">
      <block type="logic_boolean"><field name="BOOL">TRUE</field></block>
      <block type="logic_operation">
        <value name="A"><shadow type="logic_boolean"><field name="BOOL">TRUE</field></shadow></value>
        <value name="B"><shadow type="logic_boolean"><field name="BOOL">FALSE</field></shadow></value>
      </block>
      <block type="logic_negate">
        <value name="BOOL"><shadow type="logic_boolean"><field name="BOOL">TRUE</field></shadow></value>
      </block>
      <block type="logic_compare">
        <value name="A"><shadow type="math_number"><field name="NUM">0</field></shadow></value>
        <value name="B"><shadow type="math_number"><field name="NUM">0</field></shadow></value>
      </block>
    </category>
    <category name="数值" colour="160">
      <block type="math_number"><field name="NUM">0</field></block>
      <block type="current_hour_24"></block>
      <block type="current_time_minus_4h_day_of_week"></block>
      <block type="text"></block>
    </category>
    <sep></sep>
    <category name="变量" colour="%{BKY_VARIABLES_HUE}">
      <block type="define_variable">
        <value name="VALUE"><shadow type="text"><field name="TEXT"></field></shadow></value>
      </block>
      <block type="variables_get"></block>
      <block type="variables_set"></block>
    </category>
    <category name="函数" colour="330">
      <block type="custom_function_def"></block>
      <block type="custom_function_call">
        <value name="ARGS"><shadow type="text"><field name="TEXT"></field></shadow></value>
      </block>
      <block type="custom_parameter"><field name="PARAM_NAME">param1</field></block>
    </category>
    <category name="集合" colour="260">
      <block type="lists_create_with"></block>
    </category>
  </xml>
`;
