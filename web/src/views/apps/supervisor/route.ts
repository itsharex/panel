import { $gettext } from '@/utils/gettext'
import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'supervisor',
  path: '/apps/supervisor',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-supervisor-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: $gettext('Supervisor'),
        icon: 'mdi:monitor-dashboard',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
