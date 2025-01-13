<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="ID">
          <el-input v-model.trim="params.userID" clearable placeholder="ID" @clear="search" />
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="params.projectID" clearable placeholder="项目" @clear="search">
            <el-option
              v-for="item in optionsMap.projectList"
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

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="createdAt" label="头像">
          <template slot-scope="scope">
            <img class="avatar-img" :src="scope.row.avatarURL">
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="userID" label="ID" />
        <el-table-column show-overflow-tooltip sortable prop="createdAt" label="创建时间" />
        <el-table-column show-overflow-tooltip prop="username" label="用户名" />
        <el-table-column show-overflow-tooltip prop="nickname" label="用户昵称" />
        <el-table-column show-overflow-tooltip prop="point" label="积分" />
        <el-table-column show-overflow-tooltip prop="birthday" label="生日" />
        <el-table-column show-overflow-tooltip prop="sex" label="性别">
          <template slot-scope="scope">
            {{ sexMap[scope.row.sex] }}
          </template>
        </el-table-column>
        <!-- <el-table-column show-overflow-tooltip sortable prop="password" label="密码" />
        <el-table-column show-overflow-tooltip sortable prop="email" label="邮箱" />
        <el-table-column show-overflow-tooltip sortable prop="phone" label="手机号" /> -->
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
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
          <el-form-item v-if="!dialogFormData.ID" label="项目" prop="projectID">
            <el-select v-model="dialogFormData.projectID" placeholder="请选择项目">
              <el-option
                v-for="item in optionsMap.projectList"
                :key="item.projectID"
                :label="item.title"
                :value="item.projectID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="头像">
            <img :src="dialogFormData.avatarURL" class="avatar-img">
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input v-model.trim="dialogFormData.username" :disabled="!!dialogFormData.ID" placeholder="用户名" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model.trim="dialogFormData.nickname" placeholder="昵称" />
          </el-form-item>
          <el-form-item label="密码" :prop="dialogFormData.ID ? 'editPassword' : 'password'">
            <el-input v-model.trim="dialogFormData.password" placeholder="密码" />
          </el-form-item>
          <el-form-item v-if="!dialogFormData.ID" label="邮箱" prop="email">
            <el-input v-model.trim="dialogFormData.email" placeholder="邮箱" />
          </el-form-item>
          <el-form-item v-if="!dialogFormData.ID" label="手机号" prop="phone">
            <el-input v-model.trim="dialogFormData.phone" placeholder="手机号" />
          </el-form-item>
          <el-form-item v-if="dialogFormData.ID" label="性别">
            <el-input :value="sexMap[dialogFormData.sex]" disabled />
          </el-form-item>
          <el-form-item v-if="dialogFormData.ID" label="生日">
            <el-input :value="dialogFormData.birthday" disabled />
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
import { getUserList, createUser, updateUser } from '@/api/business/user'
import { getProjectList } from '@/api/business/project'
import { sexMap } from '@/constant'

export default {
  name: 'User',
  data() {
    return {
      sexMap,
      // 查询参数
      params: {
        username: '',
        userID: '',
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
        username: '',
        nickname: '',
        password: '',
        email: '',
        phone: '',
        projectID: ''
      },
      dialogFormRules: {
        username: [
          {
            validator: (rule, value, callback) => {
              // 使用正则表达式检查是否只包含数字和字符
              const regex = /^[0-9a-zA-Z]+$/
              if (!value) {
                callback(new Error('请输入用户名'))
              } else if (!regex.test(value)) {
                callback(new Error('用户名只能是数字和字符'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
            required: true
          }
        ],
        projectID: [
          { required: true, message: '请选择项目', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '请输入用户昵称', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
        ],
        password: [
          { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              // 使用正则表达式检查是否只包含数字和字符
              const regex = /^[0-9a-zA-Z]+$/
              if (!value) {
                callback(new Error('请输入密码'))
              } else if (!regex.test(value)) {
                callback(new Error('密码只能是数字和字符'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
            required: true
          }
        ],
        editPassword: [
          {
            validator: (rule, value, callback) => {
              // 使用正则表达式检查是否只包含数字和字符
              const regex = /^[0-9a-zA-Z]+$/
              if (!this.dialogFormData.password) {
                callback()
                return
              }
              if (this.dialogFormData.password.length < 6 || this.dialogFormData.password.length > 20) {
                callback(new Error('长度在 6 到 20 个字符'))
              } else if (!regex.test(this.dialogFormData.password)) {
                callback(new Error('密码只能是数字和字符'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        email: [
          {
            validator: (rule, value, callback) => {
              const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/
              if (!value) {
                callback(new Error('请输入邮箱'))
              } else if (!emailRegex.test(value)) {
                callback(new Error('请输入有效的邮箱地址'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
            required: true
          }
        ],
        phone: [
          {
            validator: (rule, value, callback) => {
              const phoneRegex = /^[1][3-9][0-9]{9}$/
              if (!value) {
                callback(new Error('请输入手机号'))
              } else if (!phoneRegex.test(value)) {
                callback(new Error('请输入有效的手机号'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
            required: true
          }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      optionsMap: {
        projectList: []
      }
    }
  },
  created() {
    this.getTableData()
    this.getProjectData()
  },
  methods: {
    async getProjectData() {
      const params = {
        title: '',
        description: '',
        pageNum: 1,
        pageSize: 50
      }
      const { data } = await getProjectList(params)
      this.optionsMap.projectList = data.projects
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
        const { data } = await getUserList(this.params)
        this.tableData = data.users
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormData.ID = ''
      this.dialogFormTitle = '新增用户'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = {
        ...row,
        ID: row.userID
      }

      this.dialogFormTitle = '修改用户'
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
              const { message } = await createUser(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateUser(this.dialogFormData.ID, this.dialogFormData)
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
      this.dialogFormData = {
        username: '',
        nickname: '',
        password: '',
        email: '',
        phone: '',
        projectID: ''
      }
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
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

<style scoped>
  .avatar-img {
    height: 40px;
    width: 40px;
    border-radius: 50%;
  }
</style>
