import { computed, onMounted, onUnmounted, ref } from 'vue'

/**
 * 视口宽度断点。表格列、侧栏折叠都要按宽度切换，纯 CSS 做不到（列是 JS 定义的）。
 */
export function useBreakpoints() {
  const width = ref(typeof window === 'undefined' ? 1280 : window.innerWidth)

  const update = () => {
    width.value = window.innerWidth
  }

  onMounted(() => window.addEventListener('resize', update))
  onUnmounted(() => window.removeEventListener('resize', update))

  return {
    width,
    /** < 1100：表格藏次要列 */
    isCompact: computed(() => width.value < 1100),
    /** < 720：侧栏收成图标、页面留白减小 */
    isNarrow: computed(() => width.value < 720)
  }
}
