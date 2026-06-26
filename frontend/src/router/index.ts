import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('@/pages/Dashboard/DashboardPage.vue'),
        },
        {
          path: 'files',
          name: 'Files',
          component: () => import('@/pages/Files/FilesPage.vue'),
        },
        {
          path: 'tags',
          name: 'Tags',
          component: () => import('@/pages/Tags/TagsPage.vue'),
        },
        {
          path: 'trash',
          name: 'Trash',
          component: () => import('@/pages/Files/FilesPage.vue'),
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('@/pages/Settings/SettingsPage.vue'),
        },
      ],
    },
  ],
})

export default router
