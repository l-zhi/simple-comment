import { FormEvent, useState } from "react";
import { Comment } from "../types/comment";
import { CreateCommentPayload } from "../api/comment";

interface Props {
  onSubmit: (payload: CreateCommentPayload) => Promise<boolean>;
  submitting: boolean;
  replyTo?: Comment | null;
  onCancelReply?: () => void;
}

export function CommentForm({
  onSubmit,
  submitting,
  replyTo,
  onCancelReply,
}: Props) {
  const [userName, setUserName] = useState("");
  const [content, setContent] = useState("");

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    // 直接回复哪条就传哪条（根或回复都可以）；所属楼由后端根据 parentId 自动计算 reply_root_id
    const parentId = replyTo ? replyTo.id : 0;
    const ok = await onSubmit({
      userName,
      content,
      parentId,
    });
    if (ok) {
      setUserName("");
      setContent("");
      onCancelReply?.();
    }
  };

  return (
    <form className="comment-form" onSubmit={handleSubmit}>
      {replyTo && (
        <div className="reply-hint">
          <span>回复 @{replyTo.userName}</span>
          {onCancelReply && (
            <button type="button" className="link" onClick={onCancelReply}>
              取消
            </button>
          )}
        </div>
      )}
      <div className="form-row">
        <label htmlFor="userName">昵称</label>
        <input
          id="userName"
          type="text"
          maxLength={50}
          value={userName}
          onChange={(e) => setUserName(e.target.value)}
          placeholder="请输入你的昵称"
          required
        />
      </div>
      <div className="form-row">
        <label htmlFor="content">评论内容</label>
        <textarea
          id="content"
          maxLength={2000}
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder={replyTo ? `回复 ${replyTo.userName}：` : "想说点什么？"}
          required
          rows={3}
        />
      </div>
      <div className="form-actions">
        <button type="submit" disabled={submitting}>
          {submitting ? "提交中..." : replyTo ? "发送回复" : "发表评论"}
        </button>
      </div>
    </form>
  );
}
