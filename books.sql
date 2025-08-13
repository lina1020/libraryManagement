CREATE TABLE books (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       created_at DATETIME(3),
                       updated_at DATETIME(3),
                       deleted_at DATETIME(3),
                       title VARCHAR(16) NOT NULL COMMENT '书名',
                       count BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '数量',
                       isbn VARCHAR(13) NOT NULL COMMENT '编码',
                       version INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',

                       UNIQUE KEY idx_isbn (isbn),
                       INDEX idx_title (title),
                       INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4, COMMENT='图书表';