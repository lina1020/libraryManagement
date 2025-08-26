CREATE TABLE IF NOT EXISTS users (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(32) NOT NULL UNIQUE COMMENT '用户名',
                       password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
                       role ENUM('user', 'admin') NOT NULL DEFAULT 'user' COMMENT '用户角色',
                       created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
                       updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
                       deleted_at DATETIME(3) NULL DEFAULT NULL,

                       INDEX idx_username (username),
                       INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;