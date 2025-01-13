<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="类型名称">
          <el-input v-model.trim="params.title" clearable placeholder="类型名称" @clear="search" />
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
        <el-table-column show-overflow-tooltip prop="productTypeID" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="title" label="类型名称" />
        <el-table-column show-overflow-tooltip prop="categoryID" label="分类编号" />
        <el-table-column show-overflow-tooltip prop="remark" label="备注" />
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.productTypeID)">
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
          <el-form-item label="类型名称" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="类型名称" />
          </el-form-item>
          <el-form-item label="分类" prop="categoryID">
            <treeselect
              v-model="dialogFormData.categoryID"
              :options="treeselectData"
              :normalizer="normalizer"
              style="width:100%;"
            />
          </el-form-item>
          <el-form-item label="规格" prop="specIds">
            <el-select v-model="dialogFormData.specIds" multiple placeholder="请选择" class="w-100">
              <el-option
                v-for="item in specList"
                :key="item.productSpecID"
                :label="item.title"
                :value="item.productSpecID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input v-model.trim="dialogFormData.remark" type="textarea" placeholder="备注" />
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
import { getTypeList, createType, updateType, batchDeleteType } from '@/api/store/product-type'
import { specTypeEnum, specTypeMap, specTypeOptions } from '@/constant/store'
import { getCategoryTree } from '@/api/store/product-category'
import { getSpecList } from '@/api/store/product-spec'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'

export default {
  name: 'ProductType',
  components: {
    Treeselect
  },
  data() {
    return {
      specTypeEnum,
      specTypeMap,
      specTypeOptions,
      // 查询参数
      params: {
        title: '',
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
        categoryID: null,
        specIds: [],
        remark: ''
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入类型名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        categoryID: [
          { required: true, message: '请选择分类', trigger: 'blur' }
        ],
        specIds: [
          {
            required: true,
            message: '请选择规格',
            trigger: 'blur',
            validator: (rule, value, callback) => {
              if (!value.length) {
                callback(new Error('请选择规格'))
              } else {
                callback()
              }
            }
          }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      treeselectData: [],
      specList: []
    }
  },
  created() {
    this.getTableData()
    this.getCategoryData()
    this.getSpecData()
  },
  methods: {
    async getSpecData() {
      const params = {
        pageNum: 1,
        pageSize: 500
      }
      const { data } = await getSpecList(params)
      this.specList = data.productSpecs
    },
    async getCategoryData() {
      const { data } = await getCategoryTree()
      this.treeselectData = data.productCategoryTree
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
        const { data } = await getTypeList(this.params)
        this.tableData = data.productTypes.map(item => {
          return {
            ...item,
            specIds: item.specIds.split(',')
          }
        })
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增类型'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = {
        ...row,
        ID: row.productTypeID
      }

      this.dialogFormTitle = '修改类型'
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
            const params = {
              ...this.dialogFormData,
              specIds: this.dialogFormData.specIds.join(','),
              categoryID: String(this.dialogFormData.categoryID)
            }
            if (this.dialogType === 'create') {
              const { message } = await createType(params)
              msg = message
            } else {
              const { message } = await updateType(this.dialogFormData.ID, params)
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
        categoryID: null,
        specIds: [],
        remark: ''
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
        const ids = []
        this.multipleSelection.forEach(x => {
          ids.push(x.productTypeID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteType({ productTypeIds: ids.join(',') })
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
    async singleDelete(id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteType({ productTypeIds: id })
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
    },
    // treeselect
    normalizer(node) {
      return {
        id: node.productCategoryID,
        label: node.title,
        children: node.children
      }
    }
  }
}
</script>
