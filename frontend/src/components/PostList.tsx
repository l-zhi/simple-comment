import { Post } from "../types/post";

interface Props {
  posts: Post[];
  total: number;
  page: number;
  pageSize: number;
  loading: boolean;
  selectedId: number | null;
  onSelect: (post: Post) => void;
  onDelete: (id: number) => void;
  onPageChange: (page: number) => void;
}

export function PostList({
  posts,
  total,
  page,
  pageSize,
  loading,
  selectedId,
  onSelect,
  onDelete,
  onPageChange,
}: Props) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const list = Array.isArray(posts) ? posts : [];

  if (loading) {
    return <p className="empty">加载中...</p>;
  }
  if (list.length === 0) {
    return <p className="empty">暂无帖子，请先发布一篇～</p>;
  }

  return (
    <div className="post-list-wrap">
      <ul className="post-list">
        {list.map((post) => (
          <li
            key={post.id}
            className={`post-item ${selectedId === post.id ? "post-item--selected" : ""}`}
          >
            <button
              type="button"
              className="post-item-main"
              onClick={() => onSelect(post)}
            >
              <span className="post-item-title">{post.title || "（无标题）"}</span>
              <span className="post-item-meta">
                {new Date(post.createdAt).toLocaleString("zh-CN")}
              </span>
            </button>
            <button
              type="button"
              className="post-item-delete comment-delete-btn"
              onClick={(e) => {
                e.stopPropagation();
                onDelete(post.id);
              }}
            >
              删除
            </button>
          </li>
        ))}
      </ul>
      {totalPages > 1 && (
        <div className="comment-pagination">
          <button
            type="button"
            className="pagination-btn"
            disabled={page <= 1}
            onClick={() => onPageChange(page - 1)}
          >
            上一页
          </button>
          <span className="pagination-info">
            第 {page} / {totalPages} 页，共 {total} 篇
          </span>
          <button
            type="button"
            className="pagination-btn"
            disabled={page >= totalPages}
            onClick={() => onPageChange(page + 1)}
          >
            下一页
          </button>
        </div>
      )}
    </div>
  );
}
