<template>
  <div class="resource-select">
    <div class="resource-select-img-list">
      <img v-if="value && !multi" class="reource" :src="value" @click="showDialog">
      <template v-if="value && multi">
        <img
          v-for="item in value"
          :key="item"
          class="reource"
          :src="item"
        >
      </template>
      <div v-show="multi || !value" class="resource-select-img-list-add" @click="showDialog">
        <i class="el-icon-plus" />
      </div>
    </div>
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
            class="ml-20"
          >
            <el-button icon="el-icon-upload" type="primary">上传</el-button>
          </el-upload>
          <el-button
            v-if="multi"
            class="ml-20"
            type="primary"
            @click="finishSelect"
          >使用选中图片</el-button>
        </div>
        <div class="resource-list">
          <div
            v-for="item in resourceList"
            :key="item.resourceID"
            :class="['resource-item', item.selected ? 'selected' : '']"
            @click="selectResource(item)"
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
  name: 'ResourceSelectV2',
  props: {
    value: {
      type: [String, Array],
      default: ''
    },
    modal: {
      type: Boolean,
      default: true
    },
    multi: {
      type: Boolean,
      default: false
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
    selectResource(item) {
      if (this.multi) {
        this.$set(item, 'selected', !item.selected)
      } else {
        this.$emit('input', item.url + item.path)
        this.dialogVisible = false
      }
    },
    finishSelect() {
      const resourceList = this.resourceList
        .filter(item => item.selected)
        .map(item => item.url + item.path)
      if (!resourceList.length) {
        this.$message({
          message: '请选择图片',
          type: 'warning'
        })
        return
      }
      this.$emit('input', resourceList)
      this.dialogVisible = false
    }
  }
}
</script>

<style lang="scss">
.resource-select {
  &-img-list {
    display: inline-flex;
    gap: 8px;
    img {
      height: 80px;
      width: 80px;
      border: 1px dashed #d9d9d9;
      border-radius: 2px;
    }
    &-add {
      border: 1px dashed #d9d9d9;
      border-radius: 2px;
      height: 80px;
      width: 80px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      &:hover {
        border-color: #409eff;
      }
    }
  }
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
      overflow: hidden;
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
    .selected {
      border-color: #689cf9;
    }
  }
}
</style>
