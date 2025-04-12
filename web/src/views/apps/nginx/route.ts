import type { RouteType } from '~/types/router'
import { $gettext } from '@/utils/gettext'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'nginx',
  path: '/apps/nginx',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-nginx-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: $gettext('OpenResty (Nginx)'),
        icon: 'logos:nginx',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
