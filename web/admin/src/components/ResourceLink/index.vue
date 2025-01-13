/* eslint-disable */
<template>
  <div class="resource-link-list">
    <div class="link-item">
      <div class="left-content">
        <div class="item-title">URL</div>
        <div class="item-desc">{{ path }}</div>
      </div>
      <el-button class="copt-btn" @click="copyText('url')">复制</el-button>
    </div>
    <div class="link-item">
      <div class="left-content">
        <div class="item-title">HTML</div>
        <div class="item-desc">{{ `\<\img src="${path}" alt="${name}" />` }}</div>
      </div>
      <el-button class="copt-btn" @click="copyText('html')">复制</el-button>
    </div>
    <div class="link-item">
      <div class="left-content">
        <div class="item-title">Markdown</div>
        <div class="item-desc">{{ `![${name}](${path})` }}</div>
      </div>
      <el-button class="copt-btn" @click="copyText('markdown')">复制</el-button>
    </div>
  </div>
</template>

<script>
import * as clipboard from 'clipboard-polyfill'
export default {
  name: 'ResourceLink',
  props: {
    domain: String,
    path: String,
    name: String
  },
  methods: {
    copyText(type) {
      const imgUrl = this.domain + this.path
      let textToCopy = ''
      switch (type) {
        case 'url':
          textToCopy = imgUrl
          break
        case 'html':
          textToCopy = `<img src="${imgUrl}" alt="${this.name}" />`
          break
        case 'markdown':
          textToCopy = `![${this.name}](${imgUrl})`
          break
      }
      clipboard.writeText(textToCopy)
        .then(() => {
          this.$message.success('复制成功')
        })
    }
  }
}
</script>

<style lang="scss">
.resource-link-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 500px;
  .link-item {
    display: flex;
    justify-content: space-between;
    border: 1px solid rgb(221, 218, 218);
    padding: 8px 12px;
    align-items: center;
    border-radius: 4px;
    line-height: 24px;
    .item-title {
      font-weight: 500;
    }
    .item-desc {
      font-size: 13px;
      color: gray;
    }
    .el-button {
      height: 32px;
      margin-left: 12px;
    }
  }
}

</style>
