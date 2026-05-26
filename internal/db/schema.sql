-- 藏叶 数据库初始化

CREATE TABLE IF NOT EXISTS collections (
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    icon       VARCHAR(10) DEFAULT '📁',
    parent_id  VARCHAR(36) DEFAULT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS tags (
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(128) NOT NULL UNIQUE,
    color      VARCHAR(32) DEFAULT 'gray',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS files (
    id              VARCHAR(36) PRIMARY KEY,
    collection_id   VARCHAR(36) NOT NULL,
    path            VARCHAR(1024) NOT NULL,
    display_name    VARCHAR(512) DEFAULT NULL,
    file_size       BIGINT DEFAULT 0,
    mime_type       VARCHAR(255) DEFAULT NULL,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS file_tags (
    file_id VARCHAR(36) NOT NULL,
    tag_id  VARCHAR(36) NOT NULL,
    PRIMARY KEY (file_id, tag_id),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS settings (
    `key`   VARCHAR(128) PRIMARY KEY,
    `value` TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT IGNORE INTO settings (`key`, `value`) VALUES
    ('version', '1'),
    ('theme', 'dark'),
    ('layout', 'sidebar');