import { Comment, RootWithPreview } from "../types/comment";

const REPLIES_PAGE_SIZE = 20;

interface Props {
  roots: RootWithPreview[];
  totalRoots: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  expandedReplies: Record<number, Comment[]>;
  loadingReplies: Record<number, boolean>;
  onExpandReplies: (parentId: number) => void;
  onLoadMoreReplies: (parentId: number) => void;
  onReply?: (comment: Comment) => void;
  onDelete?: (id: number) => void;
}

function CommentItem({
  c,
  isRoot,
  onReply,
  onDelete,
}: {
  c: Comment;
  isRoot?: boolean;
  onReply?: (comment: Comment) => void;
  onDelete?: (id: number) => void;
}) {
  const hasReplyTo = c.replyToUserName != null && c.replyToUserName !== "";

  return (
    <div className={`comment-item ${isRoot ? "comment-item--root" : "comment-item--reply"}`}>
      <div className="comment-header">
        <span className="nickname">{c.userName}</span>
        <span className="comment-header-right">
          <span className="created-at">
            {new Date(c.createdAt).toLocaleString("zh-CN")}
          </span>
          {onDelete && (
            <button
              type="button"
              className="comment-delete-btn"
              onClick={() => onDelete(c.id)}
            >
              删除
            </button>
          )}
        </span>
      </div>
      {hasReplyTo && (
        <div className="reply-to-block">
          <span className="reply-to-label">回复 {c.replyToUserName} 的评论：</span>
          {c.replyToContent && (
            <blockquote className="reply-to-content">{c.replyToContent}</blockquote>
          )}
        </div>
      )}
      <p className="content">{c.content}</p>
      {onReply && (
        <button
          type="button"
          className="comment-reply-btn"
          onClick={() => onReply(c)}
        >
          回复
        </button>
      )}
    </div>
  );
}

export function CommentList({
  roots,
  totalRoots,
  page,
  pageSize,
  onPageChange,
  expandedReplies,
  loadingReplies,
  onExpandReplies,
  onLoadMoreReplies,
  onReply,
  onDelete,
}: Props) {
  const totalPages = Math.max(1, Math.ceil(totalRoots / pageSize));
  const list = Array.isArray(roots) ? roots : [];

  if (list.length === 0 && page === 1) {
    return <p className="empty">还没有评论，快来抢沙发～</p>;
  }

  return (
    <div className="comment-list-wrap">
      <ul className="comment-list comment-list--roots">
        {list.map((root) => {
          const rootComment = root.comment;
          const replyCount = root.replyCount ?? 0;
          const expanded = expandedReplies[rootComment.id] ?? null;
          const loadedReplies = expanded !== null ? expanded : [];
          const hasMoreReplies = replyCount > loadedReplies.length;
          const loading = loadingReplies[rootComment.id] ?? false;
          const isExpanded = expanded !== null;

          return (
            <li key={rootComment.id} className="comment-root-li">
              <CommentItem
                c={rootComment}
                isRoot
                onReply={onReply}
                onDelete={onDelete}
              />

              <div className="comment-replies">
                {!isExpanded && replyCount > 0 && (
                  <button
                    type="button"
                    className="comment-more-link"
                    onClick={() => onExpandReplies(rootComment.id)}
                  >
                    查看全部 {replyCount} 条回复
                  </button>
                )}

                {isExpanded && (
                  <>
                    {loadedReplies.map((r) => (
                      <CommentItem
                        key={r.id}
                        c={r}
                        onReply={onReply}
                        onDelete={onDelete}
                      />
                    ))}
                    {loading && <p className="replies-loading">加载中...</p>}
                    {hasMoreReplies && !loading && (
                      <button
                        type="button"
                        className="comment-more-link"
                        onClick={() => onLoadMoreReplies(rootComment.id)}
                      >
                        加载更多（还有 {replyCount - loadedReplies.length} 条）
                      </button>
                    )}
                  </>
                )}
              </div>
            </li>
          );
        })}
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
            第 {page} / {totalPages} 页，共 {totalRoots} 条
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
