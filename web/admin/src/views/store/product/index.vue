<template>
  <div>
    <Transition name="fade">
      <Create v-if="isCreate" :id="curId" @submit="submit" />
    </Transition>
    <el-card v-show="!isCreate" class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="商品标题">
          <el-input v-model.trim="params.title" clearable placeholder="商品标题" @clear="search" />
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
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column sortable prop="productID" align="center" width="80" label="ID" />
        <el-table-column show-overflow-tooltip prop="title" label="标题" />
        <el-table-column show-overflow-tooltip prop="description" label="描述" />
        <el-table-column label="状态" align="center" width="80">
          <template slot-scope="scope">
            <el-tag
              size="small"
              :type="scope.row.enable === productStatusEnum.enable ? 'success' : 'danger'"
            >{{ productStatusMap[scope.row.enable] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" width="160">
          <template slot-scope="scope">
            <!-- <el-tooltip content="上架" effect="dark" placement="top">
              <el-popconfirm title="确定上架吗？" @onConfirm="updateProductStatus(scope.row, productStatusEnum.enable)">
                <el-button
                  v-show="scope.row.enable === productStatusEnum.disable"
                  slot="reference"
                  size="mini"
                  icon="el-icon-turn-off"
                  circle
                  type="primary"
                />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip content="下架" effect="dark" placement="top">
              <el-popconfirm title="确定下架吗？" @onConfirm="updateProductStatus(scope.row, productStatusEnum.disable)">
                <el-button
                  v-show="scope.row.enable === productStatusEnum.enable"
                  slot="reference"
                  size="mini"
                  icon="el-icon-turn-off"
                  circle
                  type="danger"
                />
              </el-popconfirm>
            </el-tooltip> -->
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" class="ml-10" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.productID)">
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

    </el-card>
  </div>
</template>

<script>
import { getProductList, batchDeleteProduct, updateProduct } from '@/api/store/product'
import { getProjectList } from '@/api/business/project'
import { productStatusEnum, productStatusMap } from '@/constant/store'

import Create from './create.vue'

export default {
  name: 'Product',
  components: {
    Create
  },
  data() {
    return {
      productStatusEnum,
      productStatusMap,
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

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],

      isCreate: false,
      curId: '',
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
        const { data } = await getProductList(this.params)
        this.tableData = data.products
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    submit() {
      this.isCreate = false
      this.getTableData()
    },

    // 新增
    create() {
      this.curId = ''
      this.isCreate = true
    },

    // 修改
    update(row) {
      this.curId = row.productID
      this.isCreate = true
    },

    async updateProductStatus(row, enable) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await updateProduct({
          ...row,
          enable
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

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const productIds = []
        this.multipleSelection.forEach(x => {
          productIds.push(x.productID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteProduct({ productIds: productIds.join(',') })
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

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(productId) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteProduct({ productIds: productId })
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
