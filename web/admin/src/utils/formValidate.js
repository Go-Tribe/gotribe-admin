export const validateURL = (rule, value, callback) => {
  const urlRegex = /^(https?:\/\/)?([\w-]+\.)+[\w-]+(:\d+)?(\/[\w- ./?%&=]*)?$/

  // 如果输入为空，可以视为必填校验，否则检查格式
  if (!value) {
    callback(new Error('请输入链接'))
  } else if (!urlRegex.test(value)) {
    callback(new Error('请输入有效的链接'))
  } else {
    callback() // 校验通过
  }
}
