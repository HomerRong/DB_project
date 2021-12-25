import { createApp } from 'vue'
import App from './App.vue'
import installElementPlus from './plugins/element'
import { createRouter, createWebHistory,createWebHashHistory } from 'vue-router'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import Resetpassword from './views/Resetpassword.vue'
import Category from './views/Category.vue'
import Collection from './views/Collection.vue'
import UploadAvatar from './views/UploadAvatar.vue'

// 1. 定义路由组件.
// 也可以从其他文件导入

// 2. 定义一些路由
// 每个路由都需要映射到一个组件。
// 我们后面再讨论嵌套路由。
const routes = [
  { path: '/', name: 'home',component: Home },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/resetpassword', component: Resetpassword },
  { path: '/category', component: Category },
  { path: '/collection', component: Collection },
  { path: '/uploadavatar', component: UploadAvatar },
]

// 3. 创建路由实例并传递 `routes` 配置
// 你可以在这里输入更多的配置，但我们在这里
// 暂时保持简单

// 还有 createWebHashHistory 和 createMemoryHistory
 
const router = createRouter({
  history: createWebHashHistory(),
  base: __dirname,
  routes,
})


const app = createApp(App)
installElementPlus(app)
app.use(router)
app.mount('#app')












