<template>
  <el-dialog title="订单详情" :visible.sync="dialogFormVisible" @close="resetForm">
    <div v-if="orderDetail.orderID" class="order-detail">
      <div class="order-detail-header">
        <div class="order-detail-header-info">
          <i class="el-icon-s-order order-detail-header-info-img" />
          <div class="order-detail-header-info-right">
            <div class="order-detail-header-info-right-title">{{ orderDetail.productName }}</div>
            <div class="order-detail-header-info-right-num">订单号：{{ orderDetail.orderNumber }}</div>
          </div>
        </div>
        <el-row class="mt-24">
          <el-col :span="6">
            <div class="order-detail-header-text">订单状态</div>
            <div class="order-detail-header-value">{{ orderStatusMap[orderDetail.status] }}</div>
          </el-col>
          <el-col :span="6">
            <div class="order-detail-header-text">实际支付</div>
            <div class="order-detail-header-value">¥ {{ orderDetail.amountPay }}</div>
          </el-col>
          <el-col :span="6">
            <div class="order-detail-header-text">支付方式</div>
            <div class="order-detail-header-value">{{ payMethodMap[orderDetail.payMethod] }}</div>
          </el-col>
          <el-col :span="6">
            <div class="order-detail-header-text">支付时间</div>
            <div class="order-detail-header-value">{{ orderDetail.payTime }}</div>
          </el-col>
        </el-row>
      </div>
      <el-tabs v-model="activeName">
        <el-tab-pane label="订单信息" name="first">
          <div class="order-detail-order-info">
            <div class="order-detail-order-info-title">用户信息</div>
            <el-row class="mt-16">
              <el-col :span="8">
                <span class="order-detail-order-info-text">用户名称：{{ orderDetail.user.nickname }}</span>
              </el-col>
              <el-col :span="8">
                <span class="order-detail-order-info-text">邮箱：{{ orderDetail.user.email }}</span>
              </el-col>
            </el-row>
          </div>
          <div class="order-detail-order-info">
            <div class="order-detail-order-info-title">收货信息</div>
            <el-row class="mt-16">
              <el-col :span="8">
                <span class="order-detail-order-info-text">收货人：{{ orderDetail.consigneeName }}</span>
              </el-col>
              <el-col :span="8">
                <span class="order-detail-order-info-text">收货电话：{{ orderDetail.consigneePhone }}</span>
              </el-col>
              <el-col :span="8">
                <span class="order-detail-order-info-text">收货地址：{{ orderAddress }}</span>
              </el-col>
            </el-row>
          </div>
          <div class="order-detail-order-info">
            <div class="order-detail-order-info-title">订单信息信息</div>
            <el-row class="mt-16">
              <el-col :span="8">
                <span class="order-detail-order-info-text">创建时间：{{ orderDetail.createdAt }}</span>
              </el-col>
              <el-col :span="8">
                <span class="order-detail-order-info-text">商品总数：{{ orderDetail.quantity }}</span>
              </el-col>
              <el-col :span="8">
                <span class="order-detail-order-info-text">商品总价：{{ orderDetail.amount }}</span>
              </el-col>
            </el-row>
            <el-row class="mt-16">
              <el-col :span="8">
                <span class="order-detail-order-info-text">实际支付：{{ orderDetail.amountPay }}</span>
              </el-col>
            </el-row>
          </div>
          <div class="order-detail-order-info">
            <div class="order-detail-order-info-title">买家留言</div>
            <div class="order-detail-order-info-text mt-16">{{ orderDetail.remark || '-' }}</div>
          </div>
          <div class="order-detail-order-info">
            <div class="order-detail-order-info-title">订单备注</div>
            <div class="order-detail-order-info-text mt-16">{{ orderDetail.remarkAdmin || '-' }}</div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="商品信息" name="second">
          <el-table
            :data="productList"
            border
            style="width: 100%"
          >
            <el-table-column label="商品信息">
              <template slot-scope="scope">
                <div class="order-detail-product-info">
                  <img :src="orderDetail.productImage">
                  <div class="order-detail-product-info-right">
                    <div class="order-detail-product-info-right-title">{{ scope.row.productName }}</div>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="amountPay" label="支付价格" width="180" />
            <el-table-column prop="quantity" width="180" label="购买数量" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="订单记录" name="third">
          <el-table
            :data="orderLogList"
            border
            style="width: 100%"
          >
            <el-table-column prop="orderID" label="订单ID" />
            <el-table-column prop="remark" label="订单操作" />
            <el-table-column prop="createdAt" label="操作时间" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>
</template>

<script>
import { getOrderDetail, getOrderLog } from '@/api/store/order'
import { payMethodMap, orderStatusMap } from '@/constant/order'

export default {
  name: 'OrderDetail',
  data() {
    return {
      payMethodMap,
      orderStatusMap,
      dialogFormVisible: false,
      orderDetail: {},
      activeName: 'first',
      orderLogList: []
    }
  },
  computed: {
    orderAddress() {
      return this.orderDetail.consigneeProvince +
        this.orderDetail.consigneeCity +
        this.orderDetail.consigneeDistrict +
        this.orderDetail.consigneeStreet
    },
    productList() {
      return this.orderDetail.orderID ? [this.orderDetail] : []
    }
  },
  methods: {
    showOrderDetail(orderID) {
      this.dialogFormVisible = true
      this.activeName = 'first'
      this.getOrderDetailData(orderID)
      this.getOrderLogData(orderID)
    },
    getOrderDetailData(orderID) {
      getOrderDetail(orderID).then(res => {
        this.orderDetail = res.data.order || {}
      })
    },
    getOrderLogData(orderID) {
      getOrderLog(orderID).then(res => {
        this.orderLogList = res.data.orderLogs || []
      })
    },
    resetForm() {
      this.orderDetail = {}
    }
  }
}
</script>

<style lang="scss" scoped>
.order-detail {
  &-header {
    padding-bottom: 30px;
    &-info {
      display: flex;
      &-img {
        font-size: 60px;
        color: rgb(24, 144, 255);
      }
      &-right {
        margin-left: 12px;
        padding-top: 4px;
        &-title {
          margin-bottom: 10px;
          font-weight: 500;
          font-size: 16px;
          line-height: 16px;
        }
        &-num {
          font-size: 13px;
          color: #606266;
        }
      }
    }
    &-text {
      margin-bottom: 12px;
      font-size: 13px;
      line-height: 13px;
      color: #666;
    }
    &-value {
      font-size: 14px;
      line-height: 14px;
      color: rgba(0,0,0,.85);
    }
  }
  &-order-info {
    padding: 25px 0;
    border-bottom: 1px dashed #eee;
    &-title {
      padding-left: 10px;
      border-left: 3px solid #1890ff;
      font-size: 15px;
      line-height: 15px;
      color: #303133;
    }
    &-text {
      font-size: 13px;
      color: #666;
    }
  }
  &-product-info {
    display: flex;
    img {
      height: 50px;
      width: 50px;
    }
    &-right {
      margin-left: 8px;
    }
  }
}
</style>
