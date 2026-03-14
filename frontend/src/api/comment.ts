import { Comment, RootWithPreview } from "../types/comment";

const BASE_URL =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) || "/api";

async function unwrap<T>(res: Response): Promise<T> {
  const body = await res.json();
  if (body && typeof body === "object" && "data" in body) {
    return body.data as T;
  }
  return body as T;
}

/** 一级评论分页，每条约带 replyCount 与最多 2 条预览回复 */
export async function getComments(
  articleId = 1,
  page = 1,
  pageSize = 10
): Promise<{ items: RootWithPreview[]; total: number }> {
  const res = await fetch(
    `${BASE_URL}/comments?articleId=${articleId}&page=${page}&pageSize=${pageSize}`
  );
  if (!res.ok) {
    throw new Error("failed to fetch comments");
  }
  return unwrap(res);
}

/** 某条一级评论下的二级回复分页（查看更多 / 加载更多） */
export async function getReplies(
  parentId: number,
  offset = 0,
  limit = 20
): Promise<{ items: Comment[]; total: number }> {
  const res = await fetch(
    `${BASE_URL}/comments/replies?parentId=${parentId}&offset=${offset}&limit=${limit}`
  );
  if (!res.ok) {
    throw new Error("failed to fetch replies");
  }
  return unwrap(res);
}

export type CreateCommentPayload = {
  articleId?: number;
  userId?: number;
  userName: string;
  avatar?: string;
  parentId?: number; // 0=根评论；回复时=被回复的那条评论 id
  content: string;
};

export async function deleteComment(id: number): Promise<void> {
  const res = await fetch(`${BASE_URL}/comments/${id}`, { method: "DELETE" });
  if (!res.ok) {
    throw new Error("failed to delete comment");
  }
  const body = await res.json();
  if (body && typeof body === "object" && body.code !== 0) {
    throw new Error(body.msg || "failed to delete comment");
  }
}

export async function createComment(
  payload: CreateCommentPayload
): Promise<Comment> {
  const res = await fetch(`${BASE_URL}/comments`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      articleId: payload.articleId ?? 1,
      userId: payload.userId ?? 0,
      userName: payload.userName,
      avatar: payload.avatar ?? "",
      parentId: payload.parentId ?? 0,
      content: payload.content,
    }),
  });
  if (!res.ok) {
    throw new Error("failed to create comment");
  }
  return unwrap(res);
}
