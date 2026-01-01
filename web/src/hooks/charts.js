// 本组件参考 arco-pro 的实现
// https://github.com/arco-design/arco-design-pro-vue/blob/main/arco-design-pro-vite/src/hooks/chart-option.ts

import { computed } from 'vue'
import { useAppStore } from '@/pinia'

export default function useChartOption(sourceOption) {
  const appStore = useAppStore()
  const isDark = computed(() => {
    return appStore.isDark
  })
  // 在computed内部调用sourceOption，这样能追踪sourceOption函数内部访问的所有响应式变量
  const chartOption = computed(() => {
    // 调用sourceOption时传入isDark.value，函数内部访问的其他响应式变量也会被追踪
    return sourceOption(isDark.value)
  })
  return {
    chartOption
  }
}
