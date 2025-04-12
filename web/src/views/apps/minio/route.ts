import type { RouteType } from '~/types/router'
import { $gettext } from '@/utils/gettext'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'minio',
  path: '/apps/minio',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-minio-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: $gettext('Minio'),
        icon: 'simple-icons:minio',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
