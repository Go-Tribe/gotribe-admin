<template>
  <el-card class="m-10" shadow="always">
    <el-form ref="form" size="small" :model="formData" :rules="formRules" label-width="120px">
      <el-form-item label="网站名称" prop="title">
        <el-input v-model.trim="formData.title" placeholder="网站名称" />
      </el-form-item>
      <el-form-item label="网站logo" prop="logo">
        <ResourceSelectV2 v-model="formData.logo" />
      </el-form-item>
      <el-form-item label="网站图标" prop="icon">
        <ResourceSelectV2 v-model="formData.icon" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">更新</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import ResourceSelectV2 from '@/components/ResourceSelectV2'
import { getConfig, updateConfig } from '@/api/system/config'
import store from '@/store'

export default {
  name: 'Config',
  components: {
    ResourceSelectV2
  },
  data() {
    return {
      formData: {
        title: '',
        logo: '',
        icon: ''
      },
      formRules: {
        title: [
          { required: true, message: '请输入网站名称', trigger: 'blur' }
        ],
        logo: [
          { required: true, message: '请上传网站logo', trigger: 'blur' }
        ],
        icon: [
          { required: true, message: '请上传网站图标', trigger: 'blur' }
        ]
      }
    }
  },
  mounted() {
    this.getConfigData()
  },
  methods: {
    getConfigData() {
      getConfig().then(res => {
        this.formData = res.data.systemConfig || {}
      })
    },
    onSubmit() {
      this.$refs['form'].validate(async valid => {
        if (valid) {
          const { message } = await updateConfig(this.formData)
          store.dispatch('app/getSystemConfig')
          this.$message({
            showClose: true,
            message: message,
            type: 'success'
          })
        } else {
          return false
        }
      })
    }
  }
}
</script>
