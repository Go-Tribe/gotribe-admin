<template>
  <section class="app-main">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key" />
      </keep-alive>
    </transition>
    <div class="copyright">
      © {{ currentYear }}
      由
      <a
        href="https://www.gotribe.cn"
        target="_blank"
      >GoTribe</a>
      &
      <a
        href="https://www.mactribe.cn"
        target="_blank"
      >微椒网络</a>
      强力驱动
    </div>
  </section>
</template>

<script>
export default {
  name: 'AppMain',
  data() {
    return {
      currentYear: ''
    }
  },
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    }
  },
  mounted() {
    this.getCurYear()
  },
  methods: {
    getCurYear() {
      const currentDate = new Date()
      this.currentYear = currentDate.getFullYear()
    }
  }
}
</script>

<style lang="scss" scoped>
.app-main {
  /* 50= navbar  50  */
  min-height: calc(100vh - 50px);
  width: 100%;
  position: relative;
  overflow: hidden;
}

.fixed-header+.app-main {
  padding-top: 50px;
}

.hasTagsView {
  .app-main {
    /* 84 = navbar + tags-view = 50 + 34 */
    min-height: calc(100vh - 84px);
  }

  .fixed-header+.app-main {
    padding-top: 84px;
  }
}
::v-deep .copyright {
  font-size: 12px;
  padding-bottom: 26px;
  text-align: center;
  margin-top: 12px;
  color: #606266;
  a {
    text-decoration: underline;
    font-size: 12px;
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
