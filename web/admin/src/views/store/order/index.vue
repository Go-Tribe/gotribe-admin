<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="订单号">
          <el-input v-model.trim="params.orderNumber" clearable placeholder="订单号" @clear="search" />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model.trim="params.userID" clearable placeholder="用户ID" @clear="search" />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="params.times"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
            value-format="yyyy-MM-dd HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select
            v-model="params.status"
            placeholder="订单状态"
            clearable
            @clear="search"
          >
            <el-option
              v-for="item in orderStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button icon="el-icon-download" type="primary" @click="exportData">导出</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table
        id="order_table"
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip prop="orderID" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="orderNumber" label="订单号" />
        <el-table-column label="商品信息">
          <template slot-scope="scope">
            <div class="flex items-center">
              <img :src="scope.row.productImage" style="height: 50px;width: 50px;">
              <div class="flex items-center ml-10">{{ scope.row.productName }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip label="用户信息">
          <template slot-scope="scope">
            {{ scope.row.user.nickname }} | {{ scope.row.user.userID }}
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="amountPay" label="实际支付" />
        <el-table-column show-overflow-tooltip label="支付方式">
          <template slot-scope="scope">
            <div>{{ payMethodMap[scope.row.orderType] }}</div>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="payTime" label="支付时间" />
        <el-table-column show-overflow-tooltip label="订单状态">
          <template slot-scope="scope">
            <div>{{ orderStatusMap[scope.row.status] }}</div>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" width="180">
          <template slot-scope="scope">
            <el-tooltip content="查看详情" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-view" circle type="primary" @click="showOrderDetail(scope.row.orderID)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.orderID)">
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
          <el-form-item label="实际支付金额" prop="amountPay">
            <el-input v-model.number="dialogFormData.amountPay" type="number" placeholder="实际支付金额" />
          </el-form-item>
          <el-form-item label="订单状态" prop="status">
            <el-select v-model="dialogFormData.status" placeholder="请选择" class="w-100">
              <el-option
                v-for="item in orderStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="备注" prop="remarkAdmin">
            <el-input v-model.trim="dialogFormData.remarkAdmin" type="textarea" placeholder="备注" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
      <OrderDetail ref="orderDetail" />
    </el-card>
  </div>
</template>

<script>
import { getOrderList, updateOrder, batchDeleteOrder } from '@/api/store/order'
import { payMethodMap, orderStatusMap, orderStatusOptions } from '@/constant/order'
import OrderDetail from './components/order-detail.vue'
import { exportData } from '@/utils/excel-export'

export default {
  name: 'Order',
  components: {
    OrderDetail
  },
  data() {
    return {
      payMethodMap,
      orderStatusMap,
      orderStatusOptions,
      // 查询参数
      params: {
        orderNumber: '',
        userID: '',
        times: [],
        startTime: '',
        endTime: '',
        status: '',
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
        amountPay: '',
        remarkAdmin: null,
        status: ''
      },
      dialogFormRules: {
        amountPay: [
          { required: true, type: 'number', message: '请输入实际支付金额', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择订单状态', trigger: 'blur' }
        ],
        remarkAdmin: [
          { required: true, message: '请输入备注', trigger: 'blur' }
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
      const header = ['ID', '订单号', '商品名称', '用户ID', '用户名', '实际支付', '支付方式', '支付时间', '订单状态']
      const body = this.tableData.map(item => {
        return [
          item.orderID,
          item.orderNumber,
          item.productName,
          item.user.userID,
          item.user.nickname,
          item.amountPay,
          payMethodMap[item.orderType],
          item.payTime,
          orderStatusMap[item.status]
        ]
      })
      const rows = [header, ...body]
      exportData(rows, '订单数据')
    },
    showOrderDetail(orderID) {
      this.$refs.orderDetail.showOrderDetail(orderID)
    },
    // 查询
    search() {
      this.params.pageNum = 1
      this.params.startTime = this.params.times?.[0] || ''
      this.params.endTime = this.params.times?.[1] || ''
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getOrderList(this.params)
        this.tableData = data.orders
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增订单'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = row

      this.dialogFormTitle = '修改订单'
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
            const params = this.dialogFormData

            const { message } = await updateOrder(params)
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
        amountPay: '',
        remarkAdmin: null,
        status: ''
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
          ids.push(x.orderID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteOrder({ orderIds: ids.join(',') })
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
        const { message } = await batchDeleteOrder({ orderIds: id })
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
