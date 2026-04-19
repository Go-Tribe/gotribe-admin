/**
 * 共享类型层
 *
 * 存放跨 Feature 使用的通用类型定义，打破 Feature 之间的直接耦合。
 */

// 业务模块类型
export type { User } from '@/features/business/types/user'
export type { Project } from '@/features/business/types/project'

// 内容模块类型
export type { Post, PostParams, PostListResponse } from '@/features/content/types/post'

// 系统模块类型
export type { Admin } from '@/features/system/types/admin'

// 运营模块类型
export type { Scene } from '@/features/promotion/scene/types'
