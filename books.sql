CREATE TABLE books (
                                     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                                     created_at DATETIME(3) NULL DEFAULT NULL,
                                     updated_at DATETIME(3) NULL DEFAULT NULL,
                                     deleted_at DATETIME(3) NULL DEFAULT NULL,
                                     title VARCHAR(255) NOT NULL COMMENT '书名',
                                     count BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '数量',
                                     isbn VARCHAR(17) NOT NULL COMMENT '编码',
                                     author VARCHAR(100) NULL COMMENT '作者',
                                     content LONGTEXT NULL COMMENT '书本内容',
                                     summary TEXT NULL COMMENT '内容摘要',
                                     version INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',
                                     PRIMARY KEY (id),
                                     UNIQUE INDEX idx_isbn (isbn ASC),
                                     INDEX idx_books_deleted_at (deleted_at ASC)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci, COMMENT='图书表';