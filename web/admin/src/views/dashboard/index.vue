<template>
  <div>
    <el-card class="dashboard-container" shadow="always">
      <div class="user-info">
        <img :src="user.avatar" class="user-info-avatar">
        <div class="user-info-right">
          <div class="user-info-right-title">欢迎，{{ user.nickname }}，祝你开心每一天！</div>
        </div>
      </div>
    </el-card>
    <el-card class="operate-container" shadow="always">
      <div slot="header" class="clearfix">
        <i class="el-icon-s-operation" style="margin-right: 4px;" />
        <span>快捷操作</span>
      </div>
      <div class="operate-list">
        <router-link
          v-for="item in operateList"
          :key="item.name"
          :to="item.url"
          class="operate-item"
        >
          <i :class="item.icon" :style="{ color: item.color }" />
          <span>{{ item.name }}</span>
        </router-link>
      </div>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'Dashboard',
  data() {
    return {
      user: {},
      operateList: [
        {
          icon: 'el-icon-user-solid',
          name: '用户管理',
          url: '/business/user',
          color: 'rgb(105, 192, 255)'
        },
        {
          icon: 'el-icon-tickets',
          name: '文章管理',
          url: '/content/article',
          color: 'rgb(255, 214, 102)'
        },
        {
          icon: 'el-icon-s-custom',
          name: '角色管理',
          url: '/system/role',
          color: 'rgb(149, 222, 100)'
        },
        {
          icon: 'el-icon-menu',
          name: '菜单管理',
          url: '/system/menu',
          color: 'rgb(255, 156, 110)'
        }
      ]
    }
  },
  computed: {
    ...mapGetters([
      'nickname',
      'avatar',
      'roles'
    ])
  },
  created() {
    this.getUser()
  },
  methods: {
    getUser() {
      this.user = {
        nickname: this.nickname,
        role: this.roles.join(' | '),
        avatar: this.avatar
      }
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.dashboard-container {
  margin: 10px;
  .user-info {
    display: flex;
    gap: 20px;
    &-avatar {
      height: 80px;
      width: 80px;
      border-radius: 50%;
    }
    &-right {
      &-title {
        font-size: 20px;
        font-weight: 500;
      }
    }
  }
}
.operate-container {
  margin: 10px;
  .operate-list {
    display: flex;
    padding: 20px;
    gap: 60px;
    .operate-item {
      display: flex;
      gap: 24px;
      flex-direction: column;
      align-items: center;
      color: #515a6e;
      cursor: pointer;
      span {
        font-size: 14px;
      }
      &:hover span {
        color: #57a3f3;
      }
    }
  }
}
</style>

