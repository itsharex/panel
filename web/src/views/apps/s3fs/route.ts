import { $gettext } from '@/utils/gettext'
import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 's3fs',
  path: '/apps/s3fs',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-s3fs-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: $gettext('S3FS'),
        icon: 'logos:aws',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
