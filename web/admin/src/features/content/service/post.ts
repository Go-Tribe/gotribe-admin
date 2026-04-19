import { request } from '@/service'
import type { Post, PostListParams, PostListResponse, PostParams } from '../types/post'

/** 获取文章列表（分页、按 post_id / title 筛选） */
export const getPostList = async (params?: PostListParams): Promise<PostListResponse> => {
  return request.get<PostListResponse>('/api/post', { params })
}

/** 获取文章详情（编辑回显）；GET /api/post/:postID */
export const getPostDetail = async (postID: string): Promise<Post> => {
  const data = await request.get<{ post: Post }>(`/api/post/${postID}`)
  return data.post
}

/** 创建文章 */
export async function createPost(params: PostParams): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/post', params)
}

/** 更新文章（路径参数为 postID） */
export async function updatePost(
  postID: string,
  params: Partial<PostParams>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/post/${postID}`, params)
}

/** 删除文章（body: postIds 为逗号分隔字符串） */
export async function deletePost(postID: string): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/post', {
    data: { postIds: postID },
  })
}
