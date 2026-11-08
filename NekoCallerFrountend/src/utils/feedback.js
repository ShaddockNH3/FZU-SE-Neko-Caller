import { createDiscreteApi } from 'naive-ui'
import { themeOverrides } from '@/theme'

/**
 * 组件外面也要弹提示：api 拦截器和 composable 里的异步流程。
 * useMessage 只能在组件 setup 里用，所以挂一份全局实例。
 *
 * 全站提示只从这一处出，业务代码里别再各弹一次，否则失败时会蹦两条。
 */
const { message, dialog } = createDiscreteApi(['message', 'dialog'], {
  configProviderProps: { themeOverrides }
})

export { message, dialog }
