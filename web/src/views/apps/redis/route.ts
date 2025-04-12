import type { RouteType } from '~/types/router'
import { $gettext } from '@/utils/gettext'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'redis',
  path: '/apps/redis',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-redis-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: $gettext('Redis'),
        icon: 'logos:redis',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
