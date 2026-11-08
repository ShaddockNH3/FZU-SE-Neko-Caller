import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import './assets/styles/main.scss'

// Naive UI 按需引入，不全局注册；图标走 phosphor。
const app = createApp(App)

app.use(createPinia())
app.use(router)
app.mount('#app')
