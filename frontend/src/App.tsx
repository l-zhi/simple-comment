import { useCallback, useEffect, useState } from "react";
import { Comment, RootWithPreview } from "./types/comment";
import {
  getComments,
  getReplies,
  createComment,
  deleteComment,
  CreateCommentPayload,
} from "./api/comment";
import { CommentForm } from "./components/CommentForm";
import { CommentList } from "./components/CommentList";

const DEFAULT_ARTICLE_ID = 1;
const ROOTS_PAGE_SIZE = 10;
const REPLIES_LOAD_SIZE = 20;

function App() {
  const [roots, setRoots] = useState<RootWithPreview[]>([]);
  const [totalRoots, setTotalRoots] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [replyTo, setReplyTo] = useState<Comment | null>(null);

  const [expandedReplies, setExpandedReplies] = useState<
    Record<number, Comment[]>
  >({});
  const [loadingReplies, setLoadingReplies] = useState<Record<number, boolean>>(
    {}
  );

  const loadRoots = useCallback(async (p = page) => {
    try {
      setLoading(true);
      const data = await getComments(DEFAULT_ARTICLE_ID, p, ROOTS_PAGE_SIZE);
      setRoots(Array.isArray(data?.items) ? data.items : []);
      setTotalRoots(data?.total ?? 0);
      setPage(p);
    } catch (e) {
      setError("加载评论失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadRoots(1);
  }, []);

  const handlePageChange = useCallback(
    (newPage: number) => {
      loadRoots(newPage);
    },
    [loadRoots]
  );

  const handleExpandReplies = useCallback(async (parentId: number) => {
    setLoadingReplies((prev) => ({ ...prev, [parentId]: true }));
    try {
      const data = await getReplies(parentId, 0, REPLIES_LOAD_SIZE);
      setExpandedReplies((prev) => ({
        ...prev,
        [parentId]: data?.items ?? [],
      }));
    } finally {
      setLoadingReplies((prev) => ({ ...prev, [parentId]: false }));
    }
  }, []);

  const handleDelete = useCallback(
    async (id: number) => {
      try {
        setError(null);
        await deleteComment(id);
        await loadRoots(page);
        setExpandedReplies((prev) => {
          const next = { ...prev };
          Object.keys(next).forEach((key) => {
            const pid = Number(key);
            next[pid] = next[pid].filter((c) => c.id !== id);
            if (next[pid].length === 0) delete next[pid];
          });
          return next;
        });
      } catch (e) {
        setError("删除失败，请稍后重试");
      }
    },
    [loadRoots, page]
  );

  const handleLoadMoreReplies = useCallback(async (parentId: number) => {
    const current = expandedReplies[parentId] ?? [];
    setLoadingReplies((prev) => ({ ...prev, [parentId]: true }));
    try {
      const data = await getReplies(
        parentId,
        current.length,
        REPLIES_LOAD_SIZE
      );
      const next = data?.items ?? [];
      setExpandedReplies((prev) => ({
        ...prev,
        [parentId]: [...(prev[parentId] ?? []), ...next],
      }));
    } finally {
      setLoadingReplies((prev) => ({ ...prev, [parentId]: false }));
    }
  }, [expandedReplies]);

  const handleSubmit = async (payload: CreateCommentPayload) => {
    try {
      setSubmitting(true);
      setError(null);
      await createComment({
        ...payload,
        articleId: payload.articleId ?? DEFAULT_ARTICLE_ID,
      });
      await loadRoots(payload.parentId ? page : 1);
      setReplyTo(null);
      if (payload.parentId && expandedReplies[payload.parentId]) {
        setExpandedReplies((prev) => {
          const next = { ...prev };
          delete next[payload.parentId!];
          return next;
        });
      }
      return true;
    } catch (e) {
      setError("提交评论失败，请稍后重试");
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>简单评论系统 Demo</h1>
        <p className="subtitle">Go + React</p>
      </header>
      <main className="app-main">
        <section className="card">
          <h2>发表评论</h2>
          <CommentForm
            onSubmit={handleSubmit}
            submitting={submitting}
            replyTo={replyTo}
            onCancelReply={() => setReplyTo(null)}
          />
        </section>
        <section className="card">
          <h2>评论列表</h2>
          {loading ? (
            <p>加载中...</p>
          ) : (
            <CommentList
              roots={roots}
              totalRoots={totalRoots}
              page={page}
              pageSize={ROOTS_PAGE_SIZE}
              onPageChange={handlePageChange}
              expandedReplies={expandedReplies}
              loadingReplies={loadingReplies}
              onExpandReplies={handleExpandReplies}
              onLoadMoreReplies={handleLoadMoreReplies}
              onReply={setReplyTo}
              onDelete={handleDelete}
            />
          )}
        </section>
        {error && <div className="error-banner">{error}</div>}
      </main>
    </div>
  );
}

export default App;
