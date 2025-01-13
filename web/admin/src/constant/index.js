export const resourceType = [
  {
    type: '全部',
    id: 0
  },
  {
    type: '图片',
    id: 1
  },
  {
    type: '视频',
    id: 2
  },
  {
    type: '音频',
    id: 3
  },
  {
    type: '压缩包',
    id: 4
  },
  {
    type: '文档',
    id: 5
  },
  {
    type: '字体',
    id: 6
  },
  {
    type: '应用',
    id: 7
  }
]

export const sexMap = {
  'M': '男',
  'F': '女'
}

export const urlTypeEnum = {
  link: 1,
  article: 2,
  goods: 3
}

export const urlTypeMap = {
  [urlTypeEnum.link]: '链接',
  [urlTypeEnum.article]: '文章',
  [urlTypeEnum.goods]: '商品'
}

export const urlTypeOptions = [
  {
    id: urlTypeEnum.link,
    type: '链接'
  },
  {
    id: urlTypeEnum.article,
    type: '文章'
  },
  {
    id: urlTypeEnum.goods,
    type: '商品'
  }
]

export const publishStatusEnum = {
  unPublished: 1,
  published: 2
}

export const objectTypeEnum = {
  article: 1,
  goods: 2
}
