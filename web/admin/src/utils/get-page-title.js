import store from '@/store'

export default function getPageTitle(pageTitle) {
  const title = store.getters.systemConfig.title || 'Vue Element Admin'
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
