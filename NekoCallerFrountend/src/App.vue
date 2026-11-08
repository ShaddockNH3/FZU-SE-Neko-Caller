<script setup>
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NConfigProvider,
  NDialogProvider,
  NIcon,
  NLayout,
  NLayoutContent,
  NLayoutSider,
  NMenu,
  NMessageProvider,
  dateZhCN,
  zhCN
} from 'naive-ui'
import { PhCat } from '@phosphor-icons/vue'
import { themeOverrides } from '@/theme'
import { useBreakpoints } from '@/composables/useBreakpoints'

const router = useRouter()
const route = useRoute()
const { isNarrow } = useBreakpoints()

// 花名册是从班级点进来的，侧栏还该高亮「班级管理」
const activeKey = computed(() => route.meta.parent || route.path)

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
  router.push(key)
}
</script>

<template>
  <n-config-provider :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-dialog-provider>
        <n-layout has-sider class="shell">
          <n-layout-sider
            bordered
            :width="220"
            :collapsed-width="64"
            :collapsed="isNarrow"
            :show-trigger="false"
            class="sider"
          >
            <div class="brand" :class="{ 'is-collapsed': isNarrow }">
              <n-icon :size="26" color="#1e40af"><PhCat weight="duotone" /></n-icon>
              <div v-if="!isNarrow" class="brand-text">
                <h1>Neko Caller</h1>
                <span>课堂点名系统</span>
              </div>
            </div>

            <n-menu
              :value="activeKey"
              :options="menuOptions"
              :collapsed="isNarrow"
              :collapsed-width="64"
              :collapsed-icon-size="20"
              :root-indent="18"
              @update:value="go"
            />
          </n-layout-sider>

          <n-layout-content class="content">
            <div class="content-inner">
              <router-view />
            </div>
          </n-layout-content>
        </n-layout>
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

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 64px;
  padding: 0 18px;
  border-bottom: 1px solid #eef1f7;

  &.is-collapsed {
    justify-content: center;
    padding: 0;
  }
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;

  h1 {
    font-size: 15px;
    font-weight: 600;
    color: #1f2937;
  }

  span {
    font-size: 11px;
    color: #8b95a7;
    letter-spacing: 0.6px;
  }
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
