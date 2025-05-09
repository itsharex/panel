<script setup lang="ts">
import setting from '@/api/panel/setting'
import { useThemeStore } from '@/store'
import { locales as availableLocales } from '@/utils'
import { useGettext } from 'vue3-gettext'

const { $gettext } = useGettext()
const themeStore = useThemeStore()

const { data: model } = useRequest(setting.list, {
  initialData: {
    name: '',
    locale: '',
    username: '',
    password: '',
    email: '',
    port: 8888,
    entrance: '',
    offline_mode: false,
    website_path: '',
    backup_path: '',
    https: false,
    cert: '',
    key: ''
  }
})

const locales = computed(() => {
  return Object.entries(availableLocales).map(([code, name]: [string, string]) => {
    return {
      label: name,
      value: code
    }
  })
})

const handleSave = () => {
  useRequest(setting.update(model.value)).onSuccess(() => {
    window.$message.success($gettext('Saved successfully'))
    if (model.value.locale !== themeStore.locale) {
      themeStore.setLocale(model.value.locale)
      window.$message.info($gettext('Panel is restarting, page will refresh in 3 seconds'))
      setTimeout(() => {
        window.location.reload()
      }, 3000)
    }
  })
}
</script>

<template>
  <n-space vertical>
    <n-alert type="info">
      {{
        $gettext(
          'Modifying panel port/entrance requires corresponding changes in the browser address bar to access the panel!'
        )
      }}
    </n-alert>
    <n-form>
      <n-form-item :label="$gettext('Panel Name')">
        <n-input v-model:value="model.name" :placeholder="$gettext('Panel Name')" />
      </n-form-item>
      <n-form-item :label="$gettext('Language')">
        <n-select v-model:value="model.locale" :options="locales"> </n-select>
      </n-form-item>
      <n-form-item :label="$gettext('Username')">
        <n-input v-model:value="model.username" :placeholder="$gettext('admin')" />
      </n-form-item>
      <n-form-item :label="$gettext('Password')">
        <n-input v-model:value="model.password" :placeholder="$gettext('admin')" />
      </n-form-item>
      <n-form-item :label="$gettext('Certificate Default Email')">
        <n-input v-model:value="model.email" :placeholder="$gettext('admin@yourdomain.com')" />
      </n-form-item>
      <n-form-item :label="$gettext('Port')">
        <n-input-number v-model:value="model.port" :placeholder="$gettext('8888')" />
      </n-form-item>
      <n-form-item :label="$gettext('Access Entrance')">
        <n-input v-model:value="model.entrance" :placeholder="$gettext('admin')" />
      </n-form-item>
      <n-form-item :label="$gettext('Offline Mode')">
        <n-switch v-model:value="model.offline_mode" />
      </n-form-item>
      <n-form-item :label="$gettext('Auto Update')">
        <n-switch v-model:value="model.auto_update" />
      </n-form-item>
      <n-form-item :label="$gettext('Default Website Directory')">
        <n-input v-model:value="model.website_path" :placeholder="$gettext('/www/wwwroot')" />
      </n-form-item>
      <n-form-item :label="$gettext('Default Backup Directory')">
        <n-input v-model:value="model.backup_path" :placeholder="$gettext('/www/backup')" />
      </n-form-item>
    </n-form>
  </n-space>
  <n-button type="primary" @click="handleSave">
    {{ $gettext('Save') }}
  </n-button>
</template>

<style scoped lang="scss"></style>
