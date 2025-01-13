<template>
  <el-tabs type="border-card" class="m-10">
    <el-tab-pane label="推广场景">
      <Scene />
    </el-tab-pane>
    <el-tab-pane label="推广内容">
      <Advertise />
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import Scene from './components/scene.vue'
import Advertise from './components/advertise.vue'

export default {
  name: 'Promotion',
  components: {
    Scene,
    Advertise
  }
}
</script>
