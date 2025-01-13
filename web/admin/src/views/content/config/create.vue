<template>
  <div class="container-card" shadow="always">
    <div class="create-header">
      <div class="title">{{ title }}</div>
      <div class="operate-btn">
        <el-button type="primary" @click="$emit('submit')">返回</el-button>
        <el-button type="primary" @click="submit">提交</el-button>
      </div>
    </div>
    <el-form ref="basicForm" :rules="basicFormRules" :model="basicForm" label-width="60px">
      <el-form-item label="名称" prop="title">
        <el-input v-model="basicForm.title" />
      </el-form-item>
      <el-form-item label="别名" prop="alias">
        <el-input v-model="basicForm.alias" :disabled="!!id" />
      </el-form-item>
      <el-form-item label="类型" prop="type">
        <el-select v-model="basicForm.type" :disabled="!!id">
          <el-option label="MD编辑器" :value="1" />
          <el-option label="JSON编辑器" :value="2" />
        </el-select>
      </el-form-item>
      <el-form-item label="项目" prop="projectID">
        <el-select v-model="basicForm.projectID" placeholder="请选择项目">
          <el-option
            v-for="item in optionsMap.projectList"
            :key="item.projectID"
            :label="item.title"
            :value="item.projectID"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input v-model="basicForm.description" type="textarea" />
      </el-form-item>
    </el-form>
    <MdEditor v-if="basicForm.type === 1 && (!id || basicForm.mdContent)" ref="mdEditor" class="config-editor" :md-content="basicForm.mdContent" />
    <vue-json-editor
      v-if="basicForm.type === 2"
      v-model="jsonInfo"
      :show-btns="false"
      mode="code"
      class="json-editor"
      style="height: calc(100vh - 500px);"
    />
  </div>
</template>

<script>
import MdEditor from '@/components/MdEditor'
import vueJsonEditor from 'vue-json-editor-fix-cn'
import { createConfig, updateConfig, getConfigDetail } from '@/api/content/config'
export default {
  name: 'CreateArticle',
  components: {
    MdEditor,
    vueJsonEditor
  },
  props: {
    optionsMap: {
      type: Object,
      default: () => {
        return {
          projectList: []
        }
      }
    },
    id: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      title: '新建配置',
      jsonInfo: {},
      basicForm: {
        title: '',
        description: '',
        type: 1,
        projectID: '',
        alias: '',
        info: ''
      },
      basicFormRules: {
        title: [
          { required: true, message: '请输入名称', trigger: 'blur' },
          { min: 2, message: '名称最少为两个字符', trigger: 'blur' }
        ],
        projectID: [
          { required: true, message: '请选择项目', trigger: 'blur' }
        ],
        alias: [
          { required: true, message: '请填写别名', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入描述', trigger: 'blur' },
          { min: 2, message: '标题最少为两个字符', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getConfigDetail()
  },
  methods: {
    getConfigDetail() {
      if (this.id) {
        this.title = '编辑配置'
        getConfigDetail(this.id).then(res => {
          if (res.code === 200) {
            this.basicForm = res.data.config
            if (this.basicForm.type === 2) {
              this.jsonInfo = JSON.parse(this.basicForm.info)
            }
          }
        })
      }
    },
    submit() {
      const jsonInfo = JSON.stringify(this.jsonInfo)
      if (!this.$refs.mdEditor?.getHtml() && jsonInfo === '{}') {
        this.$message({
          message: '请填写内容',
          type: 'warning'
        })
        return
      }
      this.$refs['basicForm'].validate(valid => {
        if (valid) {
          this.basicForm.info = this.basicForm.type === 1 ? this.$refs.mdEditor?.getHtml() : jsonInfo
          this.basicForm.mdContent = this.$refs.mdEditor?.getMarkdown()
          const configMethod = this.id ? updateConfig : createConfig
          configMethod(this.basicForm).then(res => {
            this.$message({
              message: `${this.id ? '编辑' : '新建'}成功`,
              type: 'success'
            })
            setTimeout(() => {
              this.$emit('submit')
            }, 1000)
          })
        }
      })
    },
    // treeselect
    normalizer(node) {
      return {
        id: node.categoryID,
        label: node.title,
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
    }
  }
}
</script>

<style lang="scss" scoped>
.container-card {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 104px);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 16px 20px 20px 20px;
}
.create-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  .title {
    font-weight: 500;
    font-size: 18px;
  }
}
.drawer {
  padding: 20px;
  overflow: auto;
  &-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    &-title {
      font-weight: 500;
    }
  }
}
</style>
<style lang="scss">
.json-editor {
  height: calc(100vh - 500px) !important;
  min-height: 500px !important;
  .jsoneditor-vue {
    height: 100%;
  }
}

.config-editor {
  height: calc(100vh - 500px) !important;
  min-height: 500px !important;
}
</style>
