import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import Overview from './views/Overview.vue'
import FilesRoot from './views/FilesRoot.vue'
import BrowseCollection from './views/BrowseCollection.vue'
import Settings from './views/Settings.vue'

const routes = [
  { path: '/overview', component: Overview },
  { path: '/', component: FilesRoot },
  { path: '/browse/:id', component: BrowseCollection },
  { path: '/settings', component: Settings },
]

const router = createRouter({ history: createWebHashHistory(), routes })
const app = createApp(App)
app.use(router)
app.mount('#app')
