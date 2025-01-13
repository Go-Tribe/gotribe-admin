<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="用户ID">
          <el-input v-model.trim="params.userID" clearable placeholder="用户ID" @clear="search" />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model.trim="params.nickname" clearable placeholder="用户昵称" @clear="search" />
        </el-form-item>
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
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%">
        <el-table-column show-overflow-tooltip prop="id" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="userID" label="用户ID" />
        <el-table-column show-overflow-tooltip prop="nickname" label="用户昵称" />
        <el-table-column show-overflow-tooltip prop="point" label="积分变化" />
        <el-table-column show-overflow-tooltip prop="reason" label="来源" />
        <el-table-column show-overflow-tooltip prop="createdAt" label="创建时间" />
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
          <el-form-item label="项目" prop="projectID">
            <el-select
              v-model="dialogFormData.projectID"
              placeholder="请选择项目"
              class="w-100"
              @change="getUserData(null)"
            >
              <el-option
                v-for="item in projectList"
                :key="item.projectID"
                :label="item.title"
                :value="item.projectID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="用户" prop="userID">
            <el-select
              v-model="dialogFormData.userID"
              :disabled="!dialogFormData.projectID"
              placeholder="请选择用户"
              class="w-100"
            >
              <el-option
                v-for="item in userList"
                :key="item.userID"
                :label="item.nickname"
                :value="item.userID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="积分" prop="point">
            <el-input v-model.number="dialogFormData.point" controls-position="right" />
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
import { getPointList, createPoint } from '@/api/operation/point'
import { getProjectList } from '@/api/business/project'
import { getUserList } from '@/api/business/user'

export default {
  name: 'Point',
  data() {
    return {
      // 查询参数
      params: {
        userID: '',
        nickname: '',
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
        userID: '',
        projectID: '',
        point: ''
      },
      dialogFormRules: {
        userID: [
          { required: true, message: '请选择用户', trigger: 'blur' }
        ],
        projectID: [
          { required: true, message: '请选择项目', trigger: 'blur' }
        ],
        point: [
          { type: 'number', required: true, message: '请填写积分', trigger: 'blur' }
        ]
      },
      projectList: [],
      userList: []
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
    async getUserData(username) {
      const params = {
        pageNum: 1,
        pageSize: 50,
        projectID: this.dialogFormData.projectID,
        username
      }
      const { data } = await getUserList(params)
      this.userList = data.users
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
        const { data } = await getPointList(this.params)
        this.tableData = data.points
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增标签'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let msg = ''
          this.submitLoading = true
          try {
            const { message } = await createPoint(this.dialogFormData)
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
      this.dialogFormData = {
        userID: '',
        projectID: '',
        point: ''
      }
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
