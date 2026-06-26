// Package db 提供数据库连接管理
package db

import (
	"database/sql"
	"fmt"
)

// schemaSQL 完整的建表 SQL
// 使用 IF NOT EXISTS 确保幂等执行
var schemaSQL = []string{
	`CREATE TABLE IF NOT EXISTS folders (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		parent_id VARCHAR(36) NULL,
		sort_order INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_parent_id (parent_id),
		INDEX idx_sort_order (sort_order)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS files (
		id VARCHAR(36) PRIMARY KEY,
		folder_id VARCHAR(36) NULL,
		name VARCHAR(255) NOT NULL,
		paths TEXT NULL,
		file_type VARCHAR(32) DEFAULT '',
		file_size BIGINT DEFAULT 0,
		description TEXT NULL,
		deleted_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_folder_id (folder_id),
		INDEX idx_file_type (file_type),
		INDEX idx_deleted_at (deleted_at),
		INDEX idx_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS tags (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(64) NOT NULL UNIQUE,
		color VARCHAR(7) DEFAULT '#FF6B6B',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS file_tags (
		file_id VARCHAR(36) NOT NULL,
		tag_id VARCHAR(36) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (file_id, tag_id),
		INDEX idx_tag_id (tag_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

	`CREATE TABLE IF NOT EXISTS settings (
		setting_key VARCHAR(128) PRIMARY KEY,
		setting_value TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
}

// Migrate 自动建表，保证数据库表存在
// 多次执行安全（幂等），不会重复创建已存在的表
func Migrate(db *sql.DB) error {
	for _, sql := range schemaSQL {
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("自动建表失败: %w\nSQL: %s", err, sql)
		}
	}
	return nil
}
