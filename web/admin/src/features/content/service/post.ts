import { request } from '@/service'
import type { Post, PostListParams, PostListResponse, PostParams } from '../types/post'

/** 获取文章列表（分页、按 id / title 筛选） */
export const getPostList = async (params?: PostListParams): Promise<PostListResponse> => {
  const data = await request.get<PostListResponse>('/api/post', { params })
  // 兼容后端字段名：优先 id，其次 ID，最后尝试将旧版 postID 转为数字
  const rawPosts = (data as any).posts || []
  const posts: Post[] = rawPosts.map((post: any) => ({
    ...post,
    id: post.id ?? post.ID ?? (post.postID ? Number(post.postID) : 0),
    slug: post.slug ?? '',
  }))
  return { posts, total: data.total ?? 0 }
}

/** 获取文章详情（编辑回显）；GET /api/post/:id */
export const getPostDetail = async (id: number): Promise<Post> => {
  const data = await request.get<{ post: Post }>(`/api/post/${id}`)
  const rawPost = (data as any).post
  const post: Post = {
    ...rawPost,
    id: rawPost.id ?? rawPost.ID ?? (rawPost.postID ? Number(rawPost.postID) : 0),
    slug: rawPost.slug ?? '',
  }
  return post
}

/** 创建文章 */
export async function createPost(params: PostParams): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/post', params)
}

/** 更新文章（路径参数为 id） */
export async function updatePost(
  id: number,
  params: Partial<PostParams>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/post/${id}`, params)
}

/** 删除文章（body: postIds 为逗号分隔字符串） */
export async function deletePost(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/post', {
    data: { postIds: String(id) },
  })
}
