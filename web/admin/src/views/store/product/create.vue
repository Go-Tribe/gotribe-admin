<template>
  <el-card class="m-10 product-create" shadow="always">
    <div class="create-header">
      <div class="title">{{ title }}</div>
      <el-button type="primary" @click="$emit('submit')">返回</el-button>
    </div>
    <el-card>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" disabled name="1" />
        <el-tab-pane label="商品规格" disabled name="2" />
        <el-tab-pane label="商品详情" disabled name="3" />
      </el-tabs>
      <template v-if="activeTab === '1'">
        <el-form ref="basicForm" :rules="basicFormRules" :model="basicForm" label-width="80px">
          <el-form-item label="标题" prop="title">
            <el-input v-model="basicForm.title" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="basicForm.description" type="textarea" />
          </el-form-item>
          <el-form-item label="商品图" prop="images">
            <ResourceSelectV2 v-model="basicForm.images" multi />
          </el-form-item>
          <el-form-item label="分类" prop="categoryID">
            <treeselect
              v-model="basicForm.categoryID"
              :options="optionsMap.treeselectData"
              :normalizer="normalizer"
              class="w-100"
              @input="treeselectInput"
            />
          </el-form-item>
          <el-form-item label="项目" prop="projectID">
            <el-select v-model="basicForm.projectID" class="w-100" placeholder="请选择项目">
              <el-option
                v-for="item in optionsMap.projectList"
                :key="item.projectID"
                :label="item.title"
                :value="item.projectID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="限购" prop="buyLimit">
            <el-input-number
              v-model="basicForm.buyLimit"
              controls-position="right"
              :min="1"
            />
          </el-form-item>
          <el-form-item label="发布状态" prop="enable">
            <el-radio-group v-model="basicForm.enable">
              <el-radio
                v-for="item in productStatusOptions"
                :key="item.value"
                :label="item.value"
              >{{ item.text }}</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </template>
      <template v-if="activeTab === '2'">
        <el-form ref="basicForm" :rules="basicFormRules" :model="basicForm.sku" label-width="80px">
          <el-form-item label="规格名称" prop="skuTitle">
            <el-input v-model="basicForm.sku.skuTitle" />
          </el-form-item>
          <el-form-item label="商品价格" prop="unitPrice">
            <el-input v-model.number="basicForm.sku.unitPrice" type="number" />
          </el-form-item>
          <el-form-item label="成本价格" prop="costPrice">
            <el-input v-model.number="basicForm.sku.costPrice" type="number" />
          </el-form-item>
          <el-form-item label="市场价格" prop="marketPrice">
            <el-input v-model.number="basicForm.sku.marketPrice" type="number" />
          </el-form-item>
          <el-form-item label="库存" prop="quantity">
            <el-input v-model.number="basicForm.sku.quantity" />
          </el-form-item>
          <el-form-item label="积分价格" prop="unitPoint">
            <el-input v-model.number="basicForm.sku.unitPoint" type="number" />
          </el-form-item>
        </el-form>
      </template>
      <MdEditor
        v-if="(!id || basicForm.content) && activeTab === '3'"
        ref="mdEditor"
        class="article-editor"
        :md-content="basicForm.content"
      />
      <div class="operate-btn">
        <el-button
          v-show="activeTab !== '1'"
          type="primary"
          @click="handlePrevClick"
        >上一步</el-button>
        <el-button
          v-show="activeTab !== '3'"
          type="primary"
          @click="handleNextClick"
        >下一步</el-button>
        <el-button
          v-show="activeTab === '3'"
          type="primary"
          @click="submit"
        >提交</el-button>
      </div>
    </el-card>
  </el-card>
</template>

<script>
import MdEditor from '@/components/MdEditor'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { createProduct, updateProduct, getProductDetail } from '@/api/store/product'
import ResourceSelectV2 from '@/components/ResourceSelectV2'
import { getCategoryTree } from '@/api/store/product-category'
import { getProjectList } from '@/api/business/project'
import { productStatusOptions, productStatusEnum } from '@/constant/store'
export default {
  name: 'CreateArticle',
  components: {
    MdEditor,
    Treeselect,
    ResourceSelectV2
  },
  props: {
    id: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      productStatusOptions,
      activeTab: '1',
      title: '新建商品',
      basicForm: {
        title: '',
        content: '',
        htmlContent: '',
        description: '',
        images: [],
        categoryID: null,
        projectID: '',
        buyLimit: 1,
        enable: productStatusEnum.enable,
        sku: {
          skuTitle: '',
          costPrice: 0,
          marketPrice: 0,
          unitPrice: 0,
          unitPoint: 0,
          quantity: 0
        }
      },
      basicFormRules: {
        title: [
          { required: true, message: '请输入标题', trigger: 'blur' }
        ],
        categoryID: [
          { required: true, message: '请选择分类', trigger: 'blur' }
        ],
        projectID: [
          { required: true, message: '请选择项目', trigger: 'blur' }
        ],
        buyLimit: [
          { required: true, message: '请填写限购', trigger: 'blur' }
        ],
        images: [
          {
            required: true,
            validator: (rule, value, callback) => {
              if (!value.length) {
                callback(new Error('请上传商品图'))
              } else {
                callback()
              }
            }
          }
        ],
        skuTitle: [
          { required: true, message: '请填写规格名称', trigger: 'blur' }
        ],
        costPrice: [
          { required: true, type: 'number', message: '请填写成本价格', trigger: 'blur' }
        ],
        unitPrice: [
          { required: true, type: 'number', message: '请填写商品价格', trigger: 'blur' }
        ],
        marketPrice: [
          { required: true, type: 'number', message: '请填写市场价格', trigger: 'blur' }
        ],
        quantity: [
          { required: true, type: 'number', message: '请填写库存', trigger: 'blur' }
        ],
        unitPoint: [
          { required: true, type: 'number', message: '请填写积分价格', trigger: 'blur' }
        ]
      },
      optionsMap: {
        treeselectData: [],
        projectList: []
      },
      specDetail: []
    }
  },
  created() {
    this.getProductDetail()
    this.getCategoryData()
    this.getProjectData()
  },
  methods: {
    async getCategoryData() {
      this.loading = true
      const { data } = await getCategoryTree()
      this.optionsMap.treeselectData = data.productCategoryTree
    },
    async getProjectData() {
      const params = {
        title: '',
        description: '',
        pageNum: 1,
        pageSize: 50
      }
      const { data } = await getProjectList(params)
      this.optionsMap.projectList = data.projects
      !this.id && (this.basicForm.projectID = data.projects[0]?.projectID)
    },
    getProductDetail() {
      if (this.id) {
        this.title = '编辑商品'
        getProductDetail(this.id).then(res => {
          if (res.code === 200) {
            this.basicForm = {
              ...res.data.product,
              sku: res.data.product?.sku?.[0] || {}
            }
          }
        })
      }
    },
    submit() {
      this.basicForm.content = this.$refs.mdEditor.getMarkdown()
      this.basicForm.htmlContent = this.$refs.mdEditor.getHtml()
      if (!this.basicForm.content) {
        this.$message({
          message: '请填写商品详情',
          type: 'warning'
        })
        return
      }
      const productMethod = this.id ? updateProduct : createProduct
      productMethod({
        ...this.basicForm,
        sku: [this.basicForm.sku]
      }).then(res => {
        this.$message({
          message: `${this.id ? '编辑' : '新建'}成功`,
          type: 'success'
        })
        setTimeout(() => {
          this.$emit('submit')
        }, 1000)
      })
    },
    // treeselect
    normalizer(node) {
      return {
        id: node.productCategoryID,
        label: node.title,
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
    },
    handleNextClick() {
      this.$refs['basicForm'].validate(valid => {
        if (valid) {
          this.activeTab = String(Number(this.activeTab) + 1)
        }
      })
    },
    handlePrevClick() {
      this.activeTab = String(Number(this.activeTab) - 1)
    }
  }
}
</script>

<style lang="scss" scoped>
.create-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  .title {
    font-weight: 500;
    font-size: 18px;
  }
}
.step {
  width: 60%;
  margin: 0 auto 16px;
}
</style>
<style lang="scss">
.article-editor {
  height: calc(100vh - 400px) !important;
  margin-bottom: 50px;
}
.product-create .el-tabs__item.is-disabled {
  color: #303133;
}
</style>
