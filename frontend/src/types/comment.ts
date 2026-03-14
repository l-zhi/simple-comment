export interface Comment {
  id: number;
  createdAt: string;
  updatedAt: string;
  articleId: number;
  userId: number;
  userName: string;
  avatar: string;
  parentId: number; // 0=根评论；回复时=被回复的那条评论 id（可为根或回复）
  replyRootId: number; // 0=根评论；回复时=所属根评论 id
  replyToUserName?: string;
  replyToContent?: string;
  content: string;
  status: number;
  likes: number;
}

/** 一级评论 + 回复总数 + 最多 2 条预览回复 */
export interface RootWithPreview {
  comment: Comment;
  replyCount: number;
  replies: Comment[];
}

export interface CommentTree extends Comment {
  replies?: CommentTree[];
}
