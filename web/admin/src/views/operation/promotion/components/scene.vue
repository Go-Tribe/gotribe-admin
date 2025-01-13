<template>
  <div>
    <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
      <el-form-item label="项目">
        <el-select v-model="params.projectID" clearable placeholder="项目" @clear="search">
          <el-option
            v-for="item in projectList"
            :key="item.projectID"
            :label="item.title"
            :value="item.projectID"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
      </el-form-item>
      <el-form-item>
        <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column show-overflow-tooltip prop="adSceneID" label="ID" width="150" />
      <el-table-column show-overflow-tooltip sortable prop="title" label="推广场景名称" />
      <el-table-column show-overflow-tooltip sortable prop="description" label="推广场景描述" />
      <el-table-column show-overflow-tooltip sortable prop="projectTitle" label="所属项目" />
      <el-table-column show-overflow-tooltip sortable prop="createdAt" label="创建时间" />
      <el-table-column fixed="right" label="操作" align="center" width="120">
        <template slot-scope="scope">
          <el-tooltip content="编辑" effect="dark" placement="top">
            <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
          </el-tooltip>
          <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
            <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.adSceneID)">
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
        <el-form-item label="推广场景名称" prop="title">
          <el-input v-model.trim="dialogFormData.title" placeholder="推广场景名称" />
        </el-form-item>
        <el-form-item label="推广场景描述" prop="description">
          <el-input v-model.trim="dialogFormData.description" placeholder="推广场景描述" />
        </el-form-item>
        <el-form-item label="项目" prop="projectID">
          <el-select v-model="dialogFormData.projectID" placeholder="请选择项目" class="w-100">
            <el-option
              v-for="item in projectList"
              :key="item.projectID"
              :label="item.title"
              :value="item.projectID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="mini" @click="cancelForm()">取 消</el-button>
        <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { getSceneList, createScene, updateScene, batchDeleteScene } from '@/api/operation/promotion'
import { getProjectList } from '@/api/business/project'

export default {
  name: 'Scene',
  data() {
    return {
      // 查询参数
      params: {
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
        projectID: ''
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入推广场景名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        description: [
          { max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        projectID: [
          { required: true, message: '请选择项目', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      projectList: []
    }
  },
  created() {
    this.getTableData()
    this.getProjectData()
  },
  methods: {
    async getProjectData() {
      const params = {
        pageNum: 1,
        pageSize: 50
      }
      const { data } = await getProjectList(params)
      this.projectList = data.projects
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
        const { data } = await getSceneList(this.params)
        this.tableData = data.adScenes || []
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增推广场景'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = row

      this.dialogFormTitle = '修改推广场景'
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
            const apiMethod = this.dialogType === 'create' ? createScene : updateScene
            const { message } = await apiMethod(this.dialogFormData)
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

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const adSceneIds = []
        this.multipleSelection.forEach(x => {
          adSceneIds.push(x.adSceneID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteScene({ adScenesIds: adSceneIds.join(',') })
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
    async singleDelete(adSceneId) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteScene({ adScenesIds: adSceneId })
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
