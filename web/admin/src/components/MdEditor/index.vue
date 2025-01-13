<template>
  <div id="markdown-container" />
</template>

<script>
import 'cherry-markdown/dist/cherry-markdown.css'
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import { uploadResource } from '@/api/content/resource'
export default {
  name: 'MdEditor',
  props: {
    mdContent: String
  },
  data() {
    return {
      mdEditor: null
    }
  },
  mounted() {
    this.initEditor()
  },
  methods: {
    initEditor() {
      this.mdEditor = new Cherry({
        id: 'markdown-container',
        value: this.mdContent,
        fileUpload: this.fileUpload
      })
    },
    getMarkdown() {
      return this.mdEditor.getMarkdown()
    },
    getHtml() {
      return this.mdEditor.getHtml()
    },
    fileUpload(file, callback) {
      const formData = new FormData()
      formData.append('file', file)
      uploadResource(formData).then(res => {
        callback(res.data.upload.domain + res.data.upload.key)
      })
    }
  }
}
</script>

<style>

</style>
