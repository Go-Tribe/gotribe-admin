import * as xlsx from 'xlsx'

export const exportData = (data, name) => {
  const worksheet = xlsx.utils.aoa_to_sheet(data)
  const workbook = xlsx.utils.book_new()
  xlsx.utils.book_append_sheet(workbook, worksheet, '表格数据')
  xlsx.writeFile(workbook, `${name || 'list'}.xlsx`)
}
