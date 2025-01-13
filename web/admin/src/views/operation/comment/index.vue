<template>
  <div>
    <el-card class="m-10" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="评论人">
          <el-input v-model.trim="params.nickname" clearable placeholder="评论人" @clear="search" />
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
        <el-form-item label="审核状态">
          <el-select
            v-model="params.status"
            placeholder="请选择审核状态"
            clearable
            @clear="search"
          >
            <el-option label="已通过" :value="publishStatusEnum.published" />
            <el-option label="待审核" :value="publishStatusEnum.unPublished" />
          </el-select>
        </el-form-item>
        <el-form-item label="评论对象类型">
          <el-select
            v-model="params.objectType"
            placeholder="请选择审核状态"
            clearable
            @clear="search"
          >
            <el-option label="文章" :value="objectTypeEnum.article" />
            <el-option label="商品" :value="objectTypeEnum.goods" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%">
        <el-table-column show-overflow-tooltip prop="commentID" label="ID" width="150" />
        <el-table-column show-overflow-tooltip prop="comment" label="评论内容" />
        <el-table-column show-overflow-tooltip prop="nickname" label="评论人" />
        <el-table-column show-overflow-tooltip prop="objectID" label="子评论ID" />
        <el-table-column label="发布状态">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 2" type="success">审核通过</el-tag>
            <el-tag v-if="scope.row.status === 1" type="info">待审核</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="createdAt" label="创建时间" />
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="审核" effect="dark" placement="top">
              <el-popconfirm title="确定审核通过吗？" @onConfirm="updateCommentStatus(scope.row, publishStatusEnum.published)">
                <el-button
                  v-show="scope.row.status === 1"
                  slot="reference"
                  class="ml-10"
                  size="mini"
                  icon="el-icon-turn-off"
                  circle
                  type="primary"
                />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip content="下线" effect="dark" placement="top">
              <el-popconfirm title="确定下线吗？" @onConfirm="updateCommentStatus(scope.row, publishStatusEnum.unPublished)">
                <el-button
                  v-show="scope.row.status === 2"
                  slot="reference"
                  class="ml-10"
                  size="mini"
                  icon="el-icon-turn-off"
                  circle
                  type="danger"
                />
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
import { getCommentList, updateComment } from '@/api/operation/comment'
import { getProjectList } from '@/api/business/project'
import { publishStatusEnum, objectTypeEnum } from '@/constant'

export default {
  name: 'Comment',
  data() {
    return {
      publishStatusEnum,
      objectTypeEnum,
      // 查询参数
      params: {
        nickname: '',
        projectID: '',
        status: '',
        objectType: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // 删除按钮弹出框
      popoverVisible: false,
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
        const { data } = await getCommentList(this.params)
        this.tableData = data.comments
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 审核
    async updateCommentStatus(row, status) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await updateComment({
          ...row,
          status
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
