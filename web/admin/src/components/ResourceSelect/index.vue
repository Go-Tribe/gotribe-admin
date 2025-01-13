<template>
  <div class="image-select">
    <el-input :placeholder="placeholder" :value="value" @input="inputChange">
      <template slot="append"><i class="el-icon-folder" @click="showDialog" /></template>
    </el-input>
    <el-dialog
      title="选择资源"
      :visible.sync="dialogVisible"
      width="1200px"
      :modal="modal"
      top="10vh"
    >
      <div>
        <div class="top-content">
          <el-select v-model="params.type" @change="initData">
            <el-option
              v-for="item in resourceType"
              :key="item.id"
              :label="item.type"
              :value="item.id"
            />
          </el-select>
          <el-upload
            action=""
            :before-upload="uploadResource"
            style="margin-left: 20px;"
          >
            <el-button icon="el-icon-upload" type="primary">上传</el-button>
          </el-upload>
        </div>
        <div class="resource-list">
          <div
            v-for="item in resourceList"
            :key="item.resourceID"
            class="resource-item"
            @click="selectResource(item.url+item.path)"
          >
            <el-image
              :src="item.url+item.path"
              style="width: 100%; height: 160px;vertical-align: top;background: #f3f4f6;"
              fit="contain"
            />
            <div class="resource-title">{{ item.title }}</div>
          </div>
        </div>
        <el-pagination
          :current-page="params.pageNum"
          :page-size="params.pageSize"
          :total="total"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, prev, pager, next, sizes"
          background
          style="margin-top: 20px;text-align: center;"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getResourceList, uploadResource } from '@/api/content/resource'
import { resourceType } from '@/constant'
export default {
  name: 'ResourceSelect',
  props: {
    value: {
      type: String,
      default: ''
    },
    modal: {
      type: Boolean,
      default: true
    },
    placeholder: {
      type: String,
      default: '请输入内容'
    }
  },
  data() {
    return {
      params: {
        title: '',
        description: '',
        pageNum: 1,
        pageSize: 10,
        type: 0
      },
      resourceList: [],
      total: 0,
      dialogVisible: false
    }
  },
  computed: {
    resourceType() {
      return resourceType
    }
  },
  methods: {
    async getResourceList() {
      const { data } = await getResourceList(this.params)
      this.resourceList = data.resources
      this.total = data.total
    },
    showDialog() {
      this.dialogVisible = true
      this.initData()
    },
    initData() {
      this.params.pageNum = 1
      this.resourceList = []
      this.getResourceList()
    },
    uploadResource(file) {
      const formData = new FormData()
      formData.append('file', file)
      uploadResource(formData).then(res => {
        this.params.type = 0
        this.initData()
      })
      return false
    },
    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.initData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.initData()
    },
    inputChange(value) {
      this.$emit('input', value)
    },
    selectResource(value) {
      this.$emit('input', value)
      this.dialogVisible = false
    }
  }
}
</script>

<style lang="scss">
.image-select {
  .el-input-group__append {
    background: white;
    cursor: pointer;
    &:hover {
      background: #F5F7FA;
    }
  }
  .el-dialog__body {
    border-top: 1px solid rgb(232, 230, 230);
    padding: 20px;
    overflow: scroll;
    height: calc(80vh - 76px);
  }
  .el-dialog__header {
    padding-bottom: 20px;
  }
  .el-dialog__footer {
    display: none;
  }
  .el-dialog {
    height: 80vh;
  }
  .top-content {
    display: flex;
    margin-bottom: 20px;
  }
  .resource-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    .resource-item {
      width: calc((100% - 50px)/6);
      position: relative;
      border-radius: 4px;
      border: 1px solid rgb(232, 228, 228);
      cursor: pointer;
      .resource-title {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 12px;
        font-weight: 500;
        padding: 4px 8px;
        text-align: center;
      }
    }
  }
}
</style>
