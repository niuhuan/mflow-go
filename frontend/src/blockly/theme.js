import * as Blockly from 'blockly/core';

// flat 平铺风格主题：浅色工作区、扁平积木、清爽工具箱。
export const mflowFlatTheme = Blockly.Theme.defineTheme('mflowFlat', {
  name: 'mflowFlat',
  base: Blockly.Themes.Classic,
  componentStyles: {
    workspaceBackgroundColour: '#f6f8fa',
    toolboxBackgroundColour: '#ffffff',
    toolboxForegroundColour: '#2d3748',
    flyoutBackgroundColour: '#eef2f7',
    flyoutForegroundColour: '#4a5568',
    flyoutOpacity: 1,
    scrollbarColour: '#cbd5e0',
    scrollbarOpacity: 0.6,
    insertionMarkerColour: '#4a5568',
    insertionMarkerOpacity: 0.4,
    markerColour: '#dd6b20',
    cursorColour: '#dd6b20',
    selectedGlowColour: '#dd6b20',
    selectedGlowOpacity: 0.5,
  },
  fontStyle: {
    family: "-apple-system, 'Segoe UI', 'Microsoft YaHei', 'PingFang SC', sans-serif",
    weight: '500',
    size: 12,
  },
  // 顶部无前置连接的积木绘制"帽子"，更贴合平铺外观（开始、函数定义等）。
  startHats: true,
});
