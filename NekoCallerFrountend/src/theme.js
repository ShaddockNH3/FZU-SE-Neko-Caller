/**
 * 主题：浅色底 + 纯白卡片 + 深蓝主色 + 金色强调。
 *
 * 和港区那个项目用同一套色值，两个后台看起来是一家的。
 * 颜色只在这里改，组件里不写死。
 */
export const palette = {
  bg: '#f4f6fa',
  surface: '#ffffff',
  primary: '#1e40af',
  primaryHover: '#1d4ed8',
  primaryPressed: '#1e3a8a',
  gold: '#c8963e',
  success: '#2f9e6f',
  warning: '#c8963e',
  danger: '#c2453d',
  border: '#e5e9f2',
  text: '#1f2937',
  textSecondary: '#4b5563',
  textMuted: '#8b95a7'
}

export const themeOverrides = {
  common: {
    primaryColor: palette.primary,
    primaryColorHover: palette.primaryHover,
    primaryColorPressed: palette.primaryPressed,
    primaryColorSuppl: palette.primaryHover,
    infoColor: palette.primary,
    infoColorHover: palette.primaryHover,
    successColor: palette.success,
    warningColor: palette.warning,
    errorColor: palette.danger,
    bodyColor: palette.bg,
    cardColor: palette.surface,
    modalColor: palette.surface,
    popoverColor: palette.surface,
    tableColor: palette.surface,
    borderColor: palette.border,
    dividerColor: '#eef1f7',
    textColorBase: palette.text,
    textColor1: palette.text,
    textColor2: palette.textSecondary,
    textColor3: palette.textMuted,
    borderRadius: '10px',
    fontWeightStrong: '600'
  },
  Card: {
    borderRadius: '14px',
    color: palette.surface,
    borderColor: palette.border
  },
  DataTable: {
    thColor: '#f7f9fc',
    thTextColor: '#5b6472',
    thFontWeight: '600',
    tdColorHover: '#f7f9ff',
    borderRadius: '10px'
  },
  Menu: {
    itemColorActive: '#eef3ff',
    itemColorActiveHover: '#e6edff',
    itemTextColorActive: palette.primary,
    itemTextColorActiveHover: palette.primary,
    itemIconColorActive: palette.primary,
    itemTextColor: palette.textSecondary,
    itemIconColor: palette.textSecondary,
    borderRadius: '8px'
  },
  Layout: {
    siderColor: palette.surface,
    headerColor: palette.surface
  },
  Tag: {
    borderRadius: '6px'
  },
  Button: {
    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px'
  }
}
