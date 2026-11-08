<script setup>
import { computed, h, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton,
  NConfigProvider,
  NDialogProvider,
  NDrawer,
  NDrawerContent,
  NIcon,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
  NMenu,
  NMessageProvider,
  dateZhCN,
  zhCN
} from 'naive-ui'
import { PhList } from '@phosphor-icons/vue'
import { themeOverrides } from '@/theme'
import { useBreakpoints } from '@/composables/useBreakpoints'
import BrandMark from '@/components/BrandMark.vue'

const router = useRouter()
const route = useRoute()
const { isNarrow } = useBreakpoints()

const drawerOpen = ref(false)

// 花名册是从班级点进来的，侧栏还该高亮「班级管理」
const activeKey = computed(() => route.meta.parent || route.path)
const pageTitle = computed(() => route.meta.title || 'Neko Caller')

const menuOptions = computed(() =>
  router.options.routes
    .filter((item) => item.meta?.title && item.meta?.icon)
    .map((item) => ({
      key: item.path,
      label: item.meta.title,
      icon: () => h(NIcon, null, { default: () => h(item.meta.icon) })
    }))
)

function go(key) {
  drawerOpen.value = false
  router.push(key)
}

// 窗口拉宽时把抽屉收掉，不然会留个空抽屉挂在左边
watch(isNarrow, (narrow) => {
  if (!narrow) drawerOpen.value = false
})
</script>

<template>
  <n-config-provider :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-dialog-provider>
        <n-layout has-sider class="shell">
          <!-- 窄屏不再用「折叠侧栏」：Naive 折叠后内容仍按 220px 排版，会溢到正文上 -->
          <n-layout-sider v-if="!isNarrow" bordered :width="220" :show-trigger="false" class="sider">
            <brand-mark />

            <n-menu :value="activeKey" :options="menuOptions" :root-indent="18" @update:value="go" />
          </n-layout-sider>

          <n-layout class="body">
            <n-layout-header v-if="isNarrow" bordered class="mobile-bar">
              <n-button quaternary circle @click="drawerOpen = true">
                <template #icon>
                  <n-icon :size="20"><PhList /></n-icon>
                </template>
              </n-button>
              <span class="mobile-title">{{ pageTitle }}</span>
            </n-layout-header>

            <n-layout-content class="content">
              <div class="content-inner">
                <router-view />
              </div>
            </n-layout-content>
          </n-layout>
        </n-layout>

        <n-drawer v-model:show="drawerOpen" placement="left" :width="230">
          <n-drawer-content :native-scrollbar="false" body-content-style="padding:0">
            <brand-mark compact />

            <n-menu :value="activeKey" :options="menuOptions" :root-indent="18" @update:value="go" />
          </n-drawer-content>
        </n-drawer>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style lang="scss" scoped>
.shell {
  // 撑满一屏：不然内容短的时候侧栏只有半截
  height: 100vh;
  background: #f4f6fa;
}

.sider {
  background: #fff;
}

.body {
  background: #f4f6fa;
}

.mobile-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 56px;
  padding: 0 12px;
  background: #fff;
}

.mobile-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
}

.content {
  background: #f4f6fa;
}

.content-inner {
  max-width: 1240px;
  margin: 0 auto;
  padding: 24px 28px 40px;
}

@media (max-width: 760px) {
  .content-inner {
    padding: 16px 12px 32px;
  }
}
</style>
