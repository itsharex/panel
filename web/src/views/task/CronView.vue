<script setup lang="ts">
import cronstrue from 'cronstrue'
import 'cronstrue/locales/zh_CN'

import Editor from '@guolao/vue-monaco-editor'
import { NButton, NDataTable, NInput, NPopconfirm, NSwitch, NTag } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import cron from '@/api/panel/cron'
import file from '@/api/panel/file'
import { decodeBase64, formatDateTime, renderIcon } from '@/utils'
import { CronNaive } from '@vue-js-cron/naive-ui'

const { $gettext } = useGettext()
const logPath = ref('')
const logModal = ref(false)
const editModal = ref(false)

const editTask = ref({
  id: 0,
  name: '',
  time: '',
  script: ''
})

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: $gettext('Task Name'),
    key: 'name',
    minWidth: 150,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: $gettext('Task Type'),
    key: 'type',
    width: 100,
    resizable: true,
    render(row: any) {
      return h(
        NTag,
        {
          type: row.type === 'shell' ? 'warning' : row.type === 'backup' ? 'success' : 'info'
        },
        {
          default: () => {
            return row.type === 'shell'
              ? $gettext('Run Script')
              : row.type === 'backup'
                ? $gettext('Backup Data')
                : $gettext('Log Rotation')
          }
        }
      )
    }
  },
  {
    title: $gettext('Enabled'),
    key: 'status',
    width: 120,
    resizable: true,
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.status,
        onUpdateValue: () => handleStatusChange(row)
      })
    }
  },
  {
    title: $gettext('Task Schedule'),
    key: 'time',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return cronstrue.toString(row.time, { locale: 'zh_CN' })
    }
  },
  {
    title: $gettext('Creation Time'),
    key: 'created_at',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatDateTime(row.created_at)
    }
  },
  {
    title: $gettext('Last Update Time'),
    key: 'updated_at',
    width: 200,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatDateTime(row.updated_at)
    }
  },
  {
    title: $gettext('Actions'),
    key: 'actions',
    width: 280,
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'warning',
            secondary: true,
            onClick: () => {
              logPath.value = row.log
              logModal.value = true
            }
          },
          {
            default: () => $gettext('Logs'),
            icon: renderIcon('majesticons:eye-line', { size: 14 })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => handleEdit(row)
          },
          {
            default: () => $gettext('Edit'),
            icon: renderIcon('material-symbols:edit-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.id)
          },
          {
            default: () => {
              return $gettext('Are you sure you want to delete this task?')
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px;'
                },
                {
                  default: () => $gettext('Delete'),
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                }
              )
            }
          }
        )
      ]
    }
  }
]

const { loading, data, page, total, pageSize, pageCount, refresh } = usePagination(
  (page, pageSize) => cron.list(page, pageSize),
  {
    initialData: { total: 0, list: [] },
    initialPageSize: 20,
    total: (res: any) => res.total,
    data: (res: any) => res.items
  }
)

const handleStatusChange = (row: any) => {
  useRequest(cron.status(row.id, !row.status)).onSuccess(() => {
    row.status = !row.status
    window.$message.success($gettext('Modified successfully'))
  })
}

const handleEdit = (row: any) => {
  useRequest(cron.get(row.id)).onSuccess(({ data }) => {
    useRequest(file.content(data.shell)).onSuccess(({ data }) => {
      editTask.value.id = row.id
      editTask.value.name = row.name
      editTask.value.time = row.time
      editTask.value.script = decodeBase64(data.content)
      editModal.value = true
    })
  })
}

const handleDelete = async (id: number) => {
  useRequest(cron.delete(id)).onSuccess(() => {
    window.$message.success($gettext('Deleted successfully'))
    window.$bus.emit('task:refresh-cron')
  })
}

const saveTaskEdit = async () => {
  useRequest(
    cron.update(editTask.value.id, editTask.value.name, editTask.value.time, editTask.value.script)
  ).onSuccess(() => {
    window.$message.success($gettext('Modified successfully'))
    window.$bus.emit('task:refresh-cron')
  })
}

onMounted(() => {
  refresh()
  window.$bus.on('task:refresh-cron', () => {
    refresh()
  })
})

onUnmounted(() => {
  window.$bus.off('task:refresh-cron')
})
</script>

<template>
  <n-data-table
    striped
    remote
    :scroll-x="1300"
    :loading="loading"
    :columns="columns"
    :data="data"
    :row-key="(row: any) => row.id"
    v-model:page="page"
    v-model:pageSize="pageSize"
    :pagination="{
      page: page,
      pageCount: pageCount,
      pageSize: pageSize,
      itemCount: total,
      showQuickJumper: true,
      showSizePicker: true,
      pageSizes: [20, 50, 100, 200]
    }"
  />
  <realtime-log-modal v-model:show="logModal" :path="logPath" />
  <n-modal
    v-model:show="editModal"
    preset="card"
    :title="$gettext('Edit Task')"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="saveTaskEdit"
  >
    <n-form inline>
      <n-form-item :label="$gettext('Task Name')">
        <n-input v-model:value="editTask.name" :placeholder="$gettext('Task Name')" />
      </n-form-item>
      <n-form-item :label="$gettext('Task Schedule')">
        <cron-naive v-model="editTask.time" locale="zh-cn"></cron-naive>
      </n-form-item>
    </n-form>
    <Editor
      v-model:value="editTask.script"
      language="shell"
      theme="vs-dark"
      height="60vh"
      mt-8
      :options="{
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true
      }"
    />
  </n-modal>
</template>
