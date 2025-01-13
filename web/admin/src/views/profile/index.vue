<template>
  <div class="app-container">
    <el-row :gutter="20">

      <el-col :span="6" :xs="24">
        <user-card :user="user" />
      </el-col>

      <el-col :span="18" :xs="24">
        <Account />
      </el-col>

    </el-row>
  </div>
</template>

<script>
import UserCard from './components/UserCard'
import Account from './components/Account'

export default {
  name: 'Profile',
  components: { UserCard, Account },
  data() {
    return {
      activeTab: 'activity'
    }
  }
}
</script>
