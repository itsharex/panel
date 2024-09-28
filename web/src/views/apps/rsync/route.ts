import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'rsync',
  path: '/apps/rsync',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-rsync-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Rsync',
        icon: 'mdi:sync',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType