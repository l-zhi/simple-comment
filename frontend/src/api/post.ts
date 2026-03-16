import { Post } from "../types/post";

const BASE_URL =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) || "/api";

async function unwrap<T>(res: Response): Promise<T> {
  const body = await res.json();
  if (body && typeof body === "object" && "data" in body) {
    return body.data as T;
  }
  return body as T;
}

export async function getPosts(
  page = 1,
  pageSize = 10
): Promise<{ items: Post[]; total: number }> {
  const res = await fetch(
    `${BASE_URL}/posts?page=${page}&pageSize=${pageSize}`
  );
  if (!res.ok) {
    throw new Error("failed to fetch posts");
  }
  return unwrap(res);
}

export async function getPost(id: number): Promise<Post> {
  const res = await fetch(`${BASE_URL}/posts/${id}`);
  if (!res.ok) {
    throw new Error("failed to fetch post");
  }
  return unwrap(res);
}

export async function createPost(
  title: string,
  content: string
): Promise<Post> {
  const res = await fetch(`${BASE_URL}/posts`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, content }),
  });
  if (!res.ok) {
    throw new Error("failed to create post");
  }
  return unwrap(res);
}

export async function deletePost(id: number): Promise<void> {
  const res = await fetch(`${BASE_URL}/posts/${id}`, { method: "DELETE" });
  if (!res.ok) {
    throw new Error("failed to delete post");
  }
  const body = await res.json();
  if (body && typeof body === "object" && body.code !== 0) {
    throw new Error(body.msg || "failed to delete post");
  }
}
