-- 嵌套评论表结构：支持归属、层级、状态与点赞，便于高性能树形查询
DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,

  article_id INT UNSIGNED NOT NULL COMMENT '关联的文章/主题ID',
  user_id INT UNSIGNED NOT NULL COMMENT '评论者用户ID',
  user_name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '评论者昵称(冗余)',
  avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '评论者头像URL(冗余)',

  parent_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID，0为顶级；回复时=被回复的那条评论ID（可为根或回复）',
  reply_root_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属根评论ID；0表示本身是根评论',

  content TEXT NOT NULL COMMENT '评论正文',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '0-审核中 1-正常 2-隐藏',
  likes INT NOT NULL DEFAULT 0 COMMENT '点赞数',

  INDEX idx_article_created (article_id, created_at),
  INDEX idx_article_parent (article_id, parent_id),
  INDEX idx_parent_id (parent_id),
  INDEX idx_reply_root_id (reply_root_id),
  INDEX idx_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
