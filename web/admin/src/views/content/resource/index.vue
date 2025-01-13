<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="资源类型">
          <el-select v-model="params.type">
            <el-option
              v-for="item in resourceType"
              :key="item.id"
              :label="item.type"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-upload
            action=""
            :before-upload="uploadResource"
          >
            <el-button icon="el-icon-upload" type="primary">上传</el-button>
          </el-upload>
        </el-form-item>
      </el-form>

      <div class="resource-list">
        <div v-for="item in tableData" :key="item.resourceID" class="resource-item">
          <el-image
            :src="item.fileType === 1 ? item.url+item.path : getResourceIcon(item.fileType)"
            style="width: 100%; height: 160px;vertical-align: top;background: #f3f4f6;"
            fit="contain"
          />
          <div class="resource-title">{{ item.title }}</div>
          <div class="hover-content">
            <i class="el-icon-edit" @click="update(item)" />
            <i class="el-icon-delete" style="margin-left: 12px;" @click="singleDelete(item.resourceID)" />
          </div>
        </div>
      </div>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[12, 24, 36, 48]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <el-dialog title="修改资源" :visible.sync="dialogFormVisible" @close="resetForm">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="80px">
          <el-form-item label="资源预览">
            <el-image
              style="width: 100px; height: 100px"
              fit="contain"
              :src="curResourceInfo.url + curResourceInfo.path"
              :preview-src-list="[curResourceInfo.url + curResourceInfo.path]"
            />
          </el-form-item>
          <el-form-item label="资源名称" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="资源名称" />
          </el-form-item>
          <el-form-item label="资源描述" prop="description">
            <el-input v-model.trim="dialogFormData.description" placeholder="资源描述" />
          </el-form-item>
          <el-form-item label="资源类型">
            <span>{{ curResourceInfo.fileType | getResourceType }}</span>
          </el-form-item>
          <el-form-item label="上传时间">
            <span>{{ curResourceInfo.createdAt }}</span>
          </el-form-item>
          <el-form-item label="资源大小">
            <span>{{ bytesToSize(curResourceInfo.size) }}</span>
          </el-form-item>
          <el-form-item label="链接">
            <ResourceLink
              :name="curResourceInfo.title"
              :domain="curResourceInfo.url"
              :path="curResourceInfo.path"
            />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import { deleteResource, getResourceList, updateResource, uploadResource } from '@/api/content/resource'
import { resourceType } from '@/constant'
import { bytesToSize } from '@/utils'
import ResourceLink from '@/components/ResourceLink'

export default {
  name: 'Resource',
  filters: {
    getResourceType(type) {
      if (!type) return ''
      return resourceType.find(item => type === item.id).type
    }
  },
  components: { ResourceLink },
  data() {
    return {
      // 查询参数
      params: {
        type: '',
        pageNum: 1,
        pageSize: 12
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // dialog对话框
      submitLoading: false,
      dialogFormVisible: false,
      dialogFormData: {
        title: '',
        description: ''
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入资源名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 20 个字符', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入资源描述', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ]
      },
      curResourceInfo: {}
    }
  },
  computed: {
    resourceType() {
      return resourceType
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    bytesToSize,
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getResourceList(this.params)
        this.tableData = data.resources
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.resourceID
      this.dialogFormData.title = row.title
      this.dialogFormData.description = row.description
      this.curResourceInfo = row
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let msg = ''
          this.submitLoading = true
          try {
            const { message } = await updateResource(this.dialogFormData.ID, this.dialogFormData)
            msg = message
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          this.getTableData()
          this.$message({
            showClose: true,
            message: msg,
            type: 'success'
          })
        } else {
          return false
        }
      })
    },

    uploadResource(file) {
      const formData = new FormData()
      formData.append('file', file)
      uploadResource(formData).then(res => {
        this.params.pageNum = 1
        this.getTableData()
      })
      return false
    },

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        title: '',
        description: ''
      }
    },

    // 单个删除
    async singleDelete(resourceID) {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        let msg = ''
        try {
          const { message } = await deleteResource({ resourceID: resourceID })
          msg = message
        } finally {
          this.loading = false
        }

        this.getTableData()
        this.$message({
          showClose: true,
          message: msg,
          type: 'success'
        })
      })
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    },
    getResourceIcon(type) {
      return require(`@/assets/FileIcon/${type}.png`)
    }
  }
}
</script>

<style lang="scss" scoped>
  .resource-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    .resource-item {
      width: calc((100% - 50px)/6);
      @media screen and (min-width: 1920px) {
        width: calc((100% - 80px)/9);
      }
      @media screen and (min-width: 2560px) {
        width: calc((100% - 120px)/13);
      }
      position: relative;
      border-radius: 4px;
      border: 1px solid rgb(232, 228, 228);
      .resource-title {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 12px;
        font-weight: 500;
        padding: 4px 8px;
        text-align: center;
      }
      .hover-content {
        display: none;
        background-image: linear-gradient(to bottom, rgba(211, 214, 218, 1) 0%, rgba(211, 214, 218, 0.5) 20%,rgba(211, 214, 218, 0) 100%);
        position: absolute;
        height: 100%;
        width: 100%;
        top: 0;
        left: 0;
        padding: 12px;
        i {
          color: white;
          cursor: pointer;
        }
      }
      &:hover .hover-content {
        display: flex;
        justify-content: flex-end;
      }
    }
  }
</style>
