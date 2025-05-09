<script setup lang="ts">
import type { UploadCustomRequestOptions } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import api from '@/api/panel/backup'

const { $gettext } = useGettext()
const show = defineModel<boolean>('show', { type: Boolean, required: true })
const type = defineModel<string>('type', { type: String, required: true })
const upload = ref<any>(null)

const uploadRequest = ({ file, onFinish, onError, onProgress }: UploadCustomRequestOptions) => {
  const formData = new FormData()
  formData.append('file', file.file as File)
  const { uploading } = useRequest(api.upload(type.value, formData))
    .onSuccess(() => {
      onFinish()
      window.$bus.emit('backup:refresh')
      window.$message.success(
        $gettext('Upload %{ filename } successfully', { filename: file.name })
      )
    })
    .onError(() => {
      onError()
    })
    .onComplete(() => {
      stopWatch()
    })
  const stopWatch = watch(uploading, (progress) => {
    onProgress({ percent: Math.ceil((progress.loaded / progress.total) * 100) })
  })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="$gettext('Upload Backup')"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-upload ref="upload" multiple directory-dnd :custom-request="uploadRequest">
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <the-icon :size="48" icon="bi:arrow-up-square" />
          </div>
          <NText text-18>{{ $gettext('Click or drag files to this area to upload') }}</NText>
          <NP depth="3" m-10>{{
            $gettext('For large files, it is recommended to use SFTP or other methods to upload')
          }}</NP>
        </n-upload-dragger>
      </n-upload>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
