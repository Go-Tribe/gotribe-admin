export const specTypeEnum = {
  text: 1,
  image: 2
}

export const specTypeMap = {
  [specTypeEnum.text]: '文字',
  [specTypeEnum.image]: '图片'
}

export const specTypeOptions = [
  {
    text: '文字',
    value: specTypeEnum.text
  },
  {
    text: '图片',
    value: specTypeEnum.image
  }
]

export const specItemStatusEnum = {
  enable: 1,
  disable: 2
}

export const specItemStatusMap = {
  [specItemStatusEnum.enable]: '启用',
  [specItemStatusEnum.disable]: '禁用'
}

export const specItemStatusOptions = [
  {
    text: '启用',
    value: specItemStatusEnum.enable
  },
  {
    text: '禁用',
    value: specItemStatusEnum.disable
  }
]

export const productStatusEnum = {
  enable: 1,
  disable: 2
}

export const productStatusMap = {
  [productStatusEnum.enable]: '上架',
  [productStatusEnum.disable]: '下架'
}

export const productStatusOptions = [
  {
    text: '上架',
    value: productStatusEnum.enable
  },
  {
    text: '下架',
    value: productStatusEnum.disable
  }
]
