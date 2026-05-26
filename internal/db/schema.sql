-- 藏叶 数据库初始化

CREATE TABLE IF NOT EXISTS settings (
    `key`   VARCHAR(128) PRIMARY KEY,
    `value` TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO settings (`key`, `value`) VALUES
    ('version', '1'),
    ('theme', 'dark'),
    ('layout', 'sidebar');