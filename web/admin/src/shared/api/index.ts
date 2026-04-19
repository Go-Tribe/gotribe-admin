/**
 * 共享 API 层
 *
 * 存放跨 Feature 使用的通用 API，打破 Feature 之间的直接耦合。
 * 其他 Feature 应从此处导入，而非直接引用具体 Feature 的内部模块。
 */

// 业务模块 API（project/user 被 content/operation/promotion/dashboard 等多个 Feature 使用）
export { getProjectList, getProject } from '@/features/business/service/project'
export { getUserList, getUserDetail } from '@/features/business/service/user'

// 内容模块 API（post 被 promotion 使用）
export { getPostList } from '@/features/content/service/post'

// 系统模块 API（config 被 auth/layout 使用）
export { getConfig } from '@/features/system/service/config'
export { getAdminInfo, changePassword, updateAdmin } from '@/features/system/service'

// 运营模块 API（scene 被 advertising 使用）
export { getSceneList } from '@/features/promotion/scene/service'
