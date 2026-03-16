import { useCallback, useEffect, useState } from "react";
import { Comment, RootWithPreview } from "./types/comment";
import { Post } from "./types/post";
import {
  getComments,
  getReplies,
  createComment,
  deleteComment,
  CreateCommentPayload,
} from "./api/comment";
import {
  getPosts,
  createPost,
  deletePost,
} from "./api/post";
import { CommentForm } from "./components/CommentForm";
import { CommentList } from "./components/CommentList";
import { PostForm } from "./components/PostForm";
import { PostList } from "./components/PostList";

const ROOTS_PAGE_SIZE = 10;
const REPLIES_LOAD_SIZE = 20;
const POSTS_PAGE_SIZE = 10;

function App() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [totalPosts, setTotalPosts] = useState(0);
  const [postPage, setPostPage] = useState(1);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [submittingPost, setSubmittingPost] = useState(false);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);

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

  const articleId = selectedPost?.id ?? 0;

  const loadPosts = useCallback(async (p = 1) => {
    try {
      setLoadingPosts(true);
      const data = await getPosts(p, POSTS_PAGE_SIZE);
      const items = Array.isArray(data?.items) ? data.items : [];
      setPosts(items);
      setTotalPosts(data?.total ?? 0);
      setPostPage(p);
      setSelectedPost((prev) => {
        if (prev && items.some((x) => x.id === prev.id)) return prev;
        return items[0] ?? null;
      });
    } catch (e) {
      setError("加载帖子列表失败");
    } finally {
      setLoadingPosts(false);
    }
  }, []);

  useEffect(() => {
    loadPosts(1);
  }, [loadPosts]);

  const loadRoots = useCallback(
    async (p = page) => {
      if (articleId <= 0) {
        setRoots([]);
        setTotalRoots(0);
        return;
      }
      try {
        setLoading(true);
        const data = await getComments(articleId, p, ROOTS_PAGE_SIZE);
        setRoots(Array.isArray(data?.items) ? data.items : []);
        setTotalRoots(data?.total ?? 0);
        setPage(p);
      } catch (e) {
        setError("加载评论失败，请稍后重试");
      } finally {
        setLoading(false);
      }
    },
    [articleId, page]
  );

  useEffect(() => {
    if (articleId > 0) {
      loadRoots(1);
    } else {
      setRoots([]);
      setTotalRoots(0);
      setExpandedReplies({});
    }
  }, [articleId]);

  const handlePostPageChange = useCallback(
    (newPage: number) => {
      loadPosts(newPage);
    },
    [loadPosts]
  );

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

  const handleDeleteComment = useCallback(
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

  const handleDeletePost = useCallback(
    async (id: number) => {
      try {
        setError(null);
        await deletePost(id);
        if (selectedPost?.id === id) {
          setSelectedPost(null);
          setRoots([]);
          setTotalRoots(0);
          setExpandedReplies({});
        }
        await loadPosts(1);
      } catch (e) {
        setError("删除帖子失败");
      }
    },
    [selectedPost, loadPosts]
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

  const handleCreatePost = async (title: string, content: string) => {
    try {
      setSubmittingPost(true);
      setError(null);
      const post = await createPost(title, content);
      await loadPosts(1);
      setSelectedPost(post);
      return true;
    } catch (e) {
      setError("发布帖子失败");
      return false;
    } finally {
      setSubmittingPost(false);
    }
  };

  const handleSubmitComment = async (payload: CreateCommentPayload) => {
    try {
      setSubmitting(true);
      setError(null);
      await createComment({
        ...payload,
        articleId: payload.articleId ?? articleId,
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
        <p className="subtitle">帖子 · 评论 · Go + React</p>
      </header>
      <main className="app-main">
        <section className="card card-posts">
          <h2>帖子</h2>
          <PostForm onSubmit={handleCreatePost} submitting={submittingPost} />
          <PostList
            posts={posts}
            total={totalPosts}
            page={postPage}
            pageSize={POSTS_PAGE_SIZE}
            loading={loadingPosts}
            selectedId={selectedPost?.id ?? null}
            onSelect={setSelectedPost}
            onDelete={handleDeletePost}
            onPageChange={handlePostPageChange}
          />
        </section>

        <section className="card">
          {selectedPost ? (
            <>
              <h2>《{selectedPost.title}》— 评论</h2>
              <div className="post-detail">
                <p className="post-detail-content">{selectedPost.content}</p>
              </div>
              <CommentForm
                onSubmit={handleSubmitComment}
                submitting={submitting}
                replyTo={replyTo}
                onCancelReply={() => setReplyTo(null)}
              />
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
                  onDelete={handleDeleteComment}
                />
              )}
            </>
          ) : (
            <p className="empty">请选择或发布一篇帖子以查看评论</p>
          )}
        </section>
        {error && <div className="error-banner">{error}</div>}
      </main>
    </div>
  );
}

export default App;
