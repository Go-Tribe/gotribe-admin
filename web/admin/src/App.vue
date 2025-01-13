<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import store from '@/store'
export default {
  name: 'App',
  mounted() {
    store.dispatch('app/getSystemConfig')
  }
}
</script>
