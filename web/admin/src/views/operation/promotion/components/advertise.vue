<template>
  <div>
    <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
      <el-form-item label="推广场景">
        <el-select
          v-model="params.sceneID"
          placeholder="请选择推广场景"
          clearable
          @clear="search"
        >
          <el-option
            v-for="item in sceneList"
            :key="item.adSceneID"
            :label="item.title"
            :value="item.adSceneID"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="发布状态">
        <el-select
          v-model="params.status"
          placeholder="请选择发布状态"
          clearable
          @clear="search"
        >
          <el-option label="已发布" :value="publishStatusEnum.published" />
          <el-option label="未发布" :value="publishStatusEnum.unPublished" />
        </el-select>
      </el-form-item>
      <el-form-item label="推广内容名称">
        <el-input v-model.trim="params.title" clearable placeholder="推广内容名称" @clear="search" />
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
      <el-table-column show-overflow-tooltip prop="adID" label="ID" width="150" />
      <el-table-column show-overflow-tooltip prop="title" label="推广内容名称" />
      <el-table-column show-overflow-tooltip prop="description" label="推广内容描述" />
      <el-table-column label="图片">
        <template slot-scope="scope">
          <img class="cover-img" :src="scope.row.image">
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip prop="SceneTitle" label="推广场景" />
      <el-table-column show-overflow-tooltip prop="url" label="外链" />
      <el-table-column show-overflow-tooltip sortable prop="sort" label="排序" />
      <el-table-column label="发布状态">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.status === 2" type="success">已发布</el-tag>
          <el-tag v-if="scope.row.status === 1" type="info">未发布</el-tag>
        </template>
      </el-table-column>
      <el-table-column show-overflow-tooltip sortable prop="createdAt" label="创建时间" />
      <el-table-column fixed="right" label="操作" align="center" width="180">
        <template slot-scope="scope">
          <el-tooltip content="编辑" effect="dark" placement="top">
            <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
          </el-tooltip>
          <el-tooltip content="发布" effect="dark" placement="top">
            <el-popconfirm title="确定发布吗？" @onConfirm="updateAdStatus(scope.row, 2)">
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
            <el-popconfirm title="确定下线吗？" @onConfirm="updateAdStatus(scope.row, 1)">
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
          <el-tooltip class="ml-10" content="删除" effect="dark" placement="top">
            <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.adID)">
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
        <el-form-item label="推广内容名称" prop="title">
          <el-input v-model.trim="dialogFormData.title" placeholder="推广内容名称" />
        </el-form-item>
        <el-form-item label="推广内容描述" prop="description">
          <el-input v-model.trim="dialogFormData.description" placeholder="推广内容描述" />
        </el-form-item>
        <el-form-item label="推广场景" prop="sceneID">
          <el-select
            v-model="dialogFormData.sceneID"
            placeholder="请选择推广场景"
            class="w-100"
          >
            <el-option
              v-for="item in sceneList"
              :key="item.adSceneID"
              :label="item.title"
              :value="item.adSceneID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="urlType">
          <el-radio-group v-model="dialogFormData.urlType">
            <el-radio
              v-for="item in urlTypeOptions"
              :key="item.id"
              :label="item.id"
            >{{ item.type }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="urlTypeMap[dialogFormData.urlType]" prop="url">
          <el-input v-model.trim="dialogFormData.url" :placeholder="urlTypeMap[dialogFormData.urlType]" />
        </el-form-item>
        <el-form-item label="图片" prop="image">
          <ResourceSelect v-model="dialogFormData.image" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="dialogFormData.sort" controls-position="right" :min="1" />
        </el-form-item>
        <el-form-item label="是否发布" prop="status">
          <el-radio-group v-model="dialogFormData.status">
            <el-radio :label="publishStatusEnum.published">是</el-radio>
            <el-radio :label="publishStatusEnum.unPublished">否</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="mini" @click="cancelForm()">取 消</el-button>
        <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm(1)">保 存</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { getAdList, createAd, updateAd, batchDeleteAd, getSceneList } from '@/api/operation/promotion'
import { urlTypeOptions, urlTypeMap, publishStatusEnum, urlTypeEnum } from '@/constant'
import ResourceSelect from '@/components/ResourceSelect'
import { validateURL } from '@/utils/formValidate'

export default {
  name: 'Advertise',
  components: {
    ResourceSelect
  },
  data() {
    return {
      urlTypeOptions,
      urlTypeMap,
      publishStatusEnum,
      // 查询参数
      params: {
        pageNum: 1,
        pageSize: 10,
        sceneID: '',
        status: '',
        title: ''
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
        sceneID: '',
        url: '',
        urlType: urlTypeEnum.link,
        image: '',
        sort: 1,
        status: publishStatusEnum.published
      },
      dialogFormRules: {
        title: [
          { required: true, message: '请输入推广内容名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入推广内容描述', trigger: 'blur' },
          { max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        sceneID: [
          { required: true, message: '请选择推广场景', trigger: 'blur' }
        ],
        url: [
          { required: true, message: '请填写链接', trigger: 'blur' },
          { validator: validateURL, message: '请填写正确的链接地址', trigger: 'blur' }
        ],
        urlType: [
          { required: true, message: '请选择类型', trigger: 'blur' }
        ],
        image: [
          { required: true, message: '请填写图片链接', trigger: 'blur' }
        ],
        sort: [
          { required: true, message: '请填写排序', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择发布状态', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      sceneList: []
    }
  },
  created() {
    this.getTableData()
    this.getSceneData()
  },
  methods: {
    async getSceneData() {
      const params = {
        pageNum: 1,
        pageSize: 50
      }
      const { data } = await getSceneList(params)
      this.sceneList = data.adScenes
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
        const { data } = await getAdList(this.params)
        this.tableData = data.ads || []
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增推广内容'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData = row

      this.dialogFormTitle = '修改推广内容'
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
            const apiMethod = this.dialogType === 'create' ? createAd : updateAd
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

    // 发布
    async updateAdStatus(row, status) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await updateAd({
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
        sceneID: '',
        url: '',
        urlType: urlTypeEnum.link,
        image: '',
        sort: 1,
        status: publishStatusEnum.published
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
        const adIds = []
        this.multipleSelection.forEach(x => {
          adIds.push(x.adID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteAd({ adsIds: adIds.join(',') })
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
    async singleDelete(adId) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteAd({ adsIds: adId })
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
<style lang="scss" scoped>
.cover-img {
  height: 100px;
  width: 100px;
}
</style>
