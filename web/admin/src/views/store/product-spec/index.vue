<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="规格名称">
          <el-input v-model.trim="params.title" clearable placeholder="规格名称" @clear="search" />
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
        <el-table-column show-overflow-tooltip prop="productSpecID" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="title" label="规格名称" />
        <el-table-column show-overflow-tooltip label="显示类型">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.format === specTypeEnum.text ? 'primary':'success'">{{ specTypeMap[scope.row.format] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="sort" label="排序" />
        <el-table-column show-overflow-tooltip prop="remark" label="备注" />
        <el-table-column fixed="right" label="操作" align="center" width="180">
          <template slot-scope="scope">
            <el-tooltip content="规格选项" effect="dark" placement="top">
              <el-button
                class="ml-10"
                size="mini"
                icon="el-icon-files"
                circle
                type="primary"
                @click="editSpecItem(scope.row.productSpecID)"
              />
            </el-tooltip>
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.productSpecID)">
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
          <el-form-item label="规格名称" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="规格名称" />
          </el-form-item>
          <el-form-item label="显示类型" prop="specType">
            <el-radio-group v-model="dialogFormData.format" :disabled="dialogType === 'update'">
              <el-radio
                v-for="item in specTypeOptions"
                :key="item.value"
                :label="item.value"
              >{{ item.text }}</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="dialogFormData.format === specTypeEnum.image" label="图片" prop="image">
            <ResourceSelectV2 v-model="dialogFormData.image" :modal="false" />
          </el-form-item>
          <el-form-item label="排序" prop="sort">
            <el-input v-model.number="dialogFormData.sort" type="number" placeholder="排序" />
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

      <el-dialog title="规格选项" :visible.sync="dialogSpecItemVisible" width="80%">
        <SpecItem v-if="dialogSpecItemVisible" :spec-i-d="curSpecID" />
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import { getSpecList, createSpec, updateSpec, batchDeleteSpec } from '@/api/store/product-spec'
import { specTypeEnum, specTypeMap, specTypeOptions } from '@/constant/store'
import SpecItem from './components/spec-item.vue'
import ResourceSelectV2 from '@/components/ResourceSelectV2'

export default {
  name: 'ProductSpec',
  components: {
    SpecItem,
    ResourceSelectV2
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
        format: specTypeEnum.text,
        sort: '',
        remark: '',
        image: ''
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入规格名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        format: [
          { required: true, message: '请选择显示类型', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      curSpecID: '',
      dialogSpecItemVisible: false
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    editSpecItem(specID) {
      this.curSpecID = specID
      this.dialogSpecItemVisible = true
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
        const { data } = await getSpecList(this.params)
        this.tableData = data.productSpecs
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增规格'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = {
        ...row,
        ID: row.productSpecID
      }

      this.dialogFormTitle = '修改规格'
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
              const { message } = await createSpec(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateSpec(this.dialogFormData.ID, this.dialogFormData)
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
        format: specTypeEnum.text,
        sort: '',
        remark: '',
        image: ''
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
          ids.push(x.productSpecID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteSpec({ productSpecIds: ids.join(',') })
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
        const { message } = await batchDeleteSpec({ productSpecIds: id })
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
