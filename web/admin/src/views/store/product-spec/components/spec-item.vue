<template>
  <el-card shadow="never">
    <el-form size="mini" :inline="true" class="demo-form-inline">
      <el-form-item>
        <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
      </el-form-item>
      <el-form-item>
        <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
      </el-form-item>
    </el-form>
    <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" height="500" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column show-overflow-tooltip prop="productSpecItemID" label="ID" width="150" />
      <el-table-column show-overflow-tooltip prop="title" label="规格值名称" />
      <el-table-column show-overflow-tooltip prop="sort" label="排序" />
      <el-table-column show-overflow-tooltip label="是否启用">
        <template slot-scope="scope">
          <el-tag
            size="small"
            :type="scope.row.enabled === specItemStatusEnum.enable ? 'success' : 'danger'"
          >{{ specItemStatusMap[scope.row.enabled] }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" align="center" width="180">
        <template slot-scope="scope">
          <el-tooltip content="启用" effect="dark" placement="top">
            <el-popconfirm title="确定启用吗？" @onConfirm="updateSpecItemStatus(scope.row, specItemStatusEnum.enable)">
              <el-button
                v-show="scope.row.enabled === specItemStatusEnum.disable"
                slot="reference"
                size="mini"
                icon="el-icon-turn-off"
                circle
                type="primary"
              />
            </el-popconfirm>
          </el-tooltip>
          <el-tooltip content="禁用" effect="dark" placement="top">
            <el-popconfirm title="确定禁用吗？" @onConfirm="updateSpecItemStatus(scope.row, specItemStatusEnum.disable)">
              <el-button
                v-show="scope.row.enabled === specItemStatusEnum.enable"
                slot="reference"
                size="mini"
                icon="el-icon-turn-off"
                circle
                type="danger"
              />
            </el-popconfirm>
          </el-tooltip>
          <el-tooltip content="编辑" effect="dark" placement="top">
            <el-button size="mini" class="ml-10" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
          </el-tooltip>
          <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
            <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.productSpecItemID)">
              <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
            </el-popconfirm>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" :modal="false" @close="resetForm">
      <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
        <el-form-item label="规格值名称" prop="title">
          <el-input v-model.trim="dialogFormData.title" placeholder="规格值名称" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input v-model.number="dialogFormData.sort" placeholder="排序" />
        </el-form-item>
        <el-form-item label="是否启用" prop="enabled">
          <el-radio-group v-model="dialogFormData.enabled">
            <el-radio
              v-for="item in specItemStatusOptions"
              :key="item.value"
              :label="item.value"
            >{{ item.text }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="mini" @click="cancelForm()">取 消</el-button>
        <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
      </div>
    </el-dialog>

  </el-card>
</template>

<script>
import { getSpecItemList, createSpecItem, updateSpecItem, batchDeleteSpecItem } from '@/api/store/product-spec'
import { specItemStatusEnum, specItemStatusOptions, specItemStatusMap } from '@/constant/store'

export default {
  name: 'ProductSpecItem',
  props: {
    specID: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      specItemStatusEnum,
      specItemStatusOptions,
      specItemStatusMap,
      // 表格数据
      tableData: [],
      loading: false,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        specID: this.specID,
        title: '',
        sort: '',
        enabled: specItemStatusEnum.enable
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
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    async updateSpecItemStatus(row, enabled) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await updateSpecItem(row.productSpecItemID, {
          ...row,
          enabled
        })
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
    // 获取表格数据
    async getTableData() {
      this.loading = true
      const params = {
        pageNum: 1,
        pageSize: 500,
        specID: this.specID
      }
      try {
        const { data } = await getSpecItemList(params)
        this.tableData = data.productSpecItems
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
        ID: row.productSpecItemID
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
              const { message } = await createSpecItem(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateSpecItem(this.dialogFormData.ID, this.dialogFormData)
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
        specID: this.specID,
        title: '',
        sort: '',
        enabled: specItemStatusEnum.enable
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
          ids.push(x.productSpecItemID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteSpecItem({ productSpecItemIds: ids.join(',') })
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
        const { message } = await batchDeleteSpecItem({ productSpecItemIds: id })
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
    }
  }
}
</script>
