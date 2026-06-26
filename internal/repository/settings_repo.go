// Package repository 提供数据访问层
package repository

import (
	"database/sql"
	"fmt"
)

// SettingsRepo 设置数据访问
type SettingsRepo struct {
	db *sql.DB
}

// NewSettingsRepo 创建 SettingsRepo 实例
func NewSettingsRepo(db *sql.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

// GetAll 查询所有设置
// 返回 map[string]string，键为设置名，值为设置值
func (r *SettingsRepo) GetAll() (map[string]string, error) {
	rows, err := r.db.Query(`SELECT setting_key, setting_value FROM settings`)
	if err != nil {
		return nil, fmt.Errorf("查询设置列表失败: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("扫描设置行数据失败: %w", err)
		}
		settings[key] = value
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历设置行数据失败: %w", err)
	}
	return settings, nil
}

// Get 查询单个设置
func (r *SettingsRepo) Get(key string) (string, error) {
	var value string
	err := r.db.QueryRow(`SELECT setting_value FROM settings WHERE setting_key = ?`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("查询设置失败: %w", err)
	}
	return value, nil
}

// Set 更新或插入设置
// 使用 INSERT ... ON DUPLICATE KEY UPDATE 实现 upsert
func (r *SettingsRepo) Set(key, value string) error {
	_, err := r.db.Exec(
		`INSERT INTO settings (setting_key, setting_value)
		 VALUES (?, ?)
		 ON DUPLICATE KEY UPDATE setting_value = ?`,
		key, value, value)
	if err != nil {
		return fmt.Errorf("更新设置失败: %w", err)
	}
	return nil
}

// Delete 删除设置
func (r *SettingsRepo) Delete(key string) error {
	_, err := r.db.Exec(`DELETE FROM settings WHERE setting_key = ?`, key)
	if err != nil {
		return fmt.Errorf("删除设置失败: %w", err)
	}
	return nil
}
