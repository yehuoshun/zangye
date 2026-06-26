// Package repository 提供数据访问层
package repository

import (
	"database/sql"
	"fmt"
	"zangye/internal/model"
)

// FileTagRepo 文件-标签关联数据访问
type FileTagRepo struct {
	db *sql.DB
}

// NewFileTagRepo 创建 FileTagRepo 实例
func NewFileTagRepo(db *sql.DB) *FileTagRepo {
	return &FileTagRepo{db: db}
}

// FindByFileID 查询文件的所有标签
func (r *FileTagRepo) FindByFileID(fileID string) ([]model.Tag, error) {
	rows, err := r.db.Query(
		`SELECT t.id, t.name, t.color, t.created_at, t.updated_at
		 FROM tags t
		 INNER JOIN file_tags ft ON t.id = ft.tag_id
		 WHERE ft.file_id = ?
		 ORDER BY t.name ASC`, fileID)
	if err != nil {
		return nil, fmt.Errorf("查询文件标签失败: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		t := model.Tag{}
		err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描标签行数据失败: %w", err)
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历标签行数据失败: %w", err)
	}
	return tags, nil
}

// AddTag 为文件添加标签
func (r *FileTagRepo) AddTag(fileID, tagID string) error {
	_, err := r.db.Exec(
		`INSERT IGNORE INTO file_tags (file_id, tag_id) VALUES (?, ?)`,
		fileID, tagID)
	if err != nil {
		return fmt.Errorf("添加文件标签失败: %w", err)
	}
	return nil
}

// RemoveTag 移除文件的某个标签
func (r *FileTagRepo) RemoveTag(fileID, tagID string) error {
	_, err := r.db.Exec(
		`DELETE FROM file_tags WHERE file_id = ? AND tag_id = ?`,
		fileID, tagID)
	if err != nil {
		return fmt.Errorf("移除文件标签失败: %w", err)
	}
	return nil
}

// SetTags 设置文件的标签列表（先清空再批量添加）
// 在事务中执行，确保原子性
func (r *FileTagRepo) SetTags(fileID string, tagIDs []string, tx *sql.Tx) error {
	// 使用事务确保原子操作
	var execer interface{ Exec(string, ...interface{}) (sql.Result, error) }
	if tx != nil {
		execer = tx
	} else {
		execer = r.db
	}

	// 先清空该文件的所有标签
	_, err := execer.Exec(`DELETE FROM file_tags WHERE file_id = ?`, fileID)
	if err != nil {
		return fmt.Errorf("清空文件标签失败: %w", err)
	}

	// 批量添加新标签
	for _, tagID := range tagIDs {
		_, err := execer.Exec(
			`INSERT INTO file_tags (file_id, tag_id) VALUES (?, ?)`,
			fileID, tagID)
		if err != nil {
			return fmt.Errorf("批量添加文件标签失败: %w", err)
		}
	}
	return nil
}

// DeleteByFileID 删除文件的所有标签关联
func (r *FileTagRepo) DeleteByFileID(fileID string) error {
	_, err := r.db.Exec(`DELETE FROM file_tags WHERE file_id = ?`, fileID)
	if err != nil {
		return fmt.Errorf("删除文件标签关联失败: %w", err)
	}
	return nil
}

// DeleteByTagID 删除标签的所有关联
func (r *FileTagRepo) DeleteByTagID(tagID string) error {
	_, err := r.db.Exec(`DELETE FROM file_tags WHERE tag_id = ?`, tagID)
	if err != nil {
		return fmt.Errorf("删除标签关联失败: %w", err)
	}
	return nil
}
