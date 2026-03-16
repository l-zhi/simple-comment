import { FormEvent, useState } from "react";

interface Props {
  onSubmit: (title: string, content: string) => Promise<boolean>;
  submitting: boolean;
}

export function PostForm({ onSubmit, submitting }: Props) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const ok = await onSubmit(title.trim(), content.trim());
    if (ok) {
      setTitle("");
      setContent("");
    }
  };

  return (
    <form className="comment-form post-form" onSubmit={handleSubmit}>
      <div className="form-row">
        <label htmlFor="post-title">标题</label>
        <input
          id="post-title"
          type="text"
          maxLength={200}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="帖子标题"
          required
        />
      </div>
      <div className="form-row">
        <label htmlFor="post-content">正文</label>
        <textarea
          id="post-content"
          maxLength={50000}
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="写点什么…"
          required
          rows={4}
        />
      </div>
      <div className="form-actions">
        <button type="submit" disabled={submitting}>
          {submitting ? "发布中..." : "发布帖子"}
        </button>
      </div>
    </form>
  );
}
