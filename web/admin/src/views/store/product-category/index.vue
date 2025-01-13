<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" row-key="productCategoryID" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip prop="productCategoryID" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="title" label="分类名称" width="150" />
        <el-table-column show-overflow-tooltip prop="description" label="分类描述" />
        <el-table-column show-overflow-tooltip prop="icon" label="图标" />
        <el-table-column show-overflow-tooltip prop="path" label="路由地址" />
        <el-table-column show-overflow-tooltip prop="sort" label="排序" align="center" width="80" />
        <el-table-column show-overflow-tooltip prop="hidden" label="隐藏" align="center" width="80">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.hidden === 1 ? 'danger':'success'">{{ scope.row.hidden === 1 ? '是':'否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.productCategoryID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="580px" @close="resetForm">
        <el-form ref="dialogForm" :inline="true" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="80px">
          <el-form-item label="分类名称" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="分类名称" style="width: 440px" />
          </el-form-item>
          <el-form-item label="分类描述" prop="description">
            <el-input v-model.trim="dialogFormData.description" placeholder="分类描述" style="width: 220px" />
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
          <el-form-item label="排序" prop="sort">
            <el-input-number v-model.number="dialogFormData.sort" controls-position="right" :min="1" :max="999" />
          </el-form-item>
          <el-form-item label="图标" prop="icon">
            <ResourceSelect v-model.trim="dialogFormData.icon" placeholder="请填入图标链接" />
          </el-form-item>
          <!-- <el-form-item label="图标" prop="icon">
            <el-popover
              placement="bottom-start"
              width="450"
              trigger="click"
              @show="$refs['iconSelect'].reset()"
            >
              <IconSelect ref="iconSelect" @selected="selected" />
              <el-input slot="reference" v-model="dialogFormData.icon" style="width: 440px;" placeholder="点击选择图标" readonly>
                <svg-icon v-if="dialogFormData.icon" slot="prefix" :icon-class="dialogFormData.icon" class="el-input__icon" style="height: 32px;width: 16px;" />
                <i v-else slot="prefix" class="el-icon-search el-input__icon" />
              </el-input>
            </el-popover>
          </el-form-item> -->
          <el-form-item label="链接" prop="path">
            <el-input v-model.trim="dialogFormData.path" placeholder="链接" style="width: 440px" />
          </el-form-item>
          <el-form-item label="隐藏" prop="hidden">
            <el-radio-group v-model="dialogFormData.hidden">
              <el-radio-button label="是" />
              <el-radio-button label="否" />
            </el-radio-group>
          </el-form-item>
          <el-form-item label="上级分类" prop="parentId">
            <treeselect
              v-model="dialogFormData.parentId"
              :options="treeselectData"
              :normalizer="normalizer"
              style="width:440px"
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
import ResourceSelect from '@/components/ResourceSelect'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { getCategoryTree, createCategory, updateCategory, batchDeleteCategory } from '@/api/store/product-category'
import { getProjectList } from '@/api/business/project'

export default {
  name: 'ProductCategory',
  components: {
    ResourceSelect,
    Treeselect
  },
  data() {
    return {
      // 表格数据
      tableData: [],
      loading: false,

      // 上级目录数据
      treeselectData: [],

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        title: '',
        description: '',
        icon: '',
        path: '',
        sort: 999,
        hidden: '否',
        parentId: 0,
        productCategoryID: '',
        projectID: ''
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入分类名称', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入分类描述', trigger: 'blur' }
        ],
        parentId: [
          { required: true, message: '请选择上级分类', trigger: 'change' }
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
    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getCategoryTree()
        this.tableData = data.productCategoryTree
        this.treeselectData = [{ id: 0, title: '顶级分类', children: data.productCategoryTree }]
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增分类'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = {
        ...row,
        hidden: row.hidden === 1 ? '是' : '否',
        parentId: row.parentID
      }

      this.dialogFormTitle = '修改分类'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true

          if (this.dialogFormData.id === this.dialogFormData.parentId) {
            return this.$message({
              showClose: true,
              message: '不能选择自身作为自己的上级目录',
              type: 'error'
            })
          }

          if (this.dialogFormData.component === '') {
            this.dialogFormData.component = 'Layout'
          }

          this.dialogFormData.hidden = this.dialogFormData.hidden === '是' ? 1 : 2

          const dialogFormDataCopy = this.dialogFormData
          let msg = ''
          try {
            if (this.dialogType === 'create') {
              const { message } = await createCategory(dialogFormDataCopy)
              msg = message
            } else {
              const { message } = await updateCategory(dialogFormDataCopy.productCategoryID, dialogFormDataCopy)
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
        icon: '',
        path: '',
        sort: 999,
        hidden: '否',
        parentId: 0,
        productCategoryID: '',
        projectID: ''
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
        const categoryIds = []
        this.multipleSelection.forEach(x => {
          categoryIds.push(x.productCategoryID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteCategory({ productCategoryIds: categoryIds.join(',') })
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
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteCategory({ productCategoryIds: String(Id) })
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

    // 选中图标
    selected(name) {
      this.dialogFormData.icon = name
    },

    // treeselect
    normalizer(node) {
      return {
        id: node.id,
        label: node.title,
        children: node.children
      }
    }

  }
}
</script>
