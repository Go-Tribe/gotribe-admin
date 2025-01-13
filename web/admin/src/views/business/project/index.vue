<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="项目名称">
          <el-input v-model.trim="params.title" clearable placeholder="项目名称" @clear="search" />
        </el-form-item>
        <el-form-item label="ID">
          <el-input v-model.trim="params.projectID" clearable placeholder="ID" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-download" type="primary" @click="exportData">导出</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            :disabled="multipleSelection.length === 0"
            :loading="loading"
            icon="el-icon-delete"
            type="danger"
            @click="batchDelete"
          >批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="projectID" label="ID" />
        <el-table-column show-overflow-tooltip sortable prop="title" label="项目名称" />
        <el-table-column show-overflow-tooltip sortable prop="description" label="项目描述" />
        <el-table-column show-overflow-tooltip sortable prop="createdAt" label="创建时间" />
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.projectID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" @close="resetForm">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="项目别名" prop="name">
            <el-input v-model.trim="dialogFormData.name" :disabled="!!dialogFormData.ID" placeholder="项目别名" />
          </el-form-item>
          <el-form-item label="Meta标题" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="Meta标题" />
          </el-form-item>
          <el-form-item label="Meta描述" prop="description">
            <el-input v-model.trim="dialogFormData.description" type="textarea" placeholder="Meta描述" />
          </el-form-item>
          <el-form-item label="Meta关键词" prop="keywords">
            <el-input v-model.trim="dialogFormData.keywords" placeholder="Meta关键词" />
          </el-form-item>
          <el-form-item label="网站图标" prop="favicon">
            <ResourceSelect v-model="dialogFormData.favicon" :modal="false" />
          </el-form-item>
          <el-form-item label="Nav图标" prop="navImage">
            <ResourceSelect v-model="dialogFormData.navImage" :modal="false" />
          </el-form-item>
          <el-form-item label="项目域名" prop="domain">
            <el-input v-model.trim="dialogFormData.domain" placeholder="项目域名" />
          </el-form-item>
          <el-form-item label="内容链接" prop="postUrl">
            <el-input v-model.trim="dialogFormData.postUrl" placeholder="内容链接" />
          </el-form-item>
          <el-form-item label="icp备案号" prop="icp">
            <el-input v-model.trim="dialogFormData.icp" placeholder="icp备案号" />
          </el-form-item>
          <el-form-item label="公安备案号" prop="publicSecurity">
            <el-input v-model.trim="dialogFormData.publicSecurity" placeholder="公安备案号" />
          </el-form-item>
          <el-form-item label="项目归属者" prop="author">
            <el-input v-model.trim="dialogFormData.author" placeholder="项目归属者" />
          </el-form-item>
          <el-form-item label="第三方js" prop="baiduAnalytics">
            <el-input v-model.trim="dialogFormData.baiduAnalytics" type="textarea" placeholder="第三方js" />
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
import { getProjectList, createProject, updateProject, batchDeleteProject } from '@/api/business/project'
import { exportData } from '@/utils/excel-export'
import ResourceSelect from '@/components/ResourceSelect'

export default {
  name: 'Project',
  components: { ResourceSelect },
  data() {
    return {
      // 查询参数
      params: {
        title: '',
        projectID: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        title: '',
        description: '',
        name: '',
        keywords: '',
        domain: '',
        postUrl: '',
        icp: '',
        author: '',
        navImage: '',
        favicon: '',
        baiduAnalytics: '',
        publicSecurity: ''
      },
      dialogFormRules: {
        name: [
          { required: true, message: '请输入项目别名', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    exportData() {
      exportData(this.tableData, '项目')
    },
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getProjectList(this.params)
        this.tableData = data.projects
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增接口'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.projectID
      this.dialogFormData.title = row.title
      this.dialogFormData.description = row.description
      this.dialogFormData.name = row.name
      this.dialogFormData.keywords = row.keywords
      this.dialogFormData.domain = row.domain
      this.dialogFormData.postUrl = row.postUrl
      this.dialogFormData.icp = row.icp
      this.dialogFormData.author = row.author
      this.dialogFormData.publicSecurity = row.publicSecurity
      this.dialogFormData.navImage = row.navImage
      this.dialogFormData.baiduAnalytics = row.baiduAnalytics
      this.dialogFormData.favicon = row.favicon

      this.dialogFormTitle = '修改项目'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let msg = ''
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              const { message } = await createProject(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateProject(this.dialogFormData.ID, this.dialogFormData)
              msg = message
            }
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

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        title: '',
        description: '',
        name: '',
        keywords: '',
        domain: '',
        postUrl: '',
        icp: '',
        author: '',
        navImage: '',
        favicon: '',
        baiduAnalytics: '',
        publicSecurity: ''
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const projectIds = []
        console.log(22, this.multipleSelection)
        this.multipleSelection.forEach(x => {
          projectIds.push(x.projectId)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteProject({ projectIds: projectIds.join(',') })
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
      }).catch(() => {
        this.$message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(projectId) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteProject({ projectIds: projectId })
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
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    }
  }
}
</script>
