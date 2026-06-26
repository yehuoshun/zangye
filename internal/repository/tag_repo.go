// Package repository 提供数据访问层
package repository

import (
	"database/sql"
	"fmt"
	"zangye/internal/model"
)

// TagRepo 标签数据访问
type TagRepo struct {
	db *sql.DB
}

// NewTagRepo 创建 TagRepo 实例
func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

// FindAll 查询所有标签（含文件数量统计）
func (r *TagRepo) FindAll() ([]*model.TagWithCount, error) {
	rows, err := r.db.Query(
		`SELECT t.id, t.name, t.color, t.created_at, t.updated_at,
		        COUNT(ft.file_id) AS file_count
		 FROM tags t
		 LEFT JOIN file_tags ft ON t.id = ft.tag_id
		 GROUP BY t.id, t.name, t.color, t.created_at, t.updated_at
		 ORDER BY t.name ASC`)
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}
	defer rows.Close()

	var tags []*model.TagWithCount
	for rows.Next() {
		t := &model.TagWithCount{}
		err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.UpdatedAt, &t.FileCount)
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

// FindByID 根据 ID 查询标签
func (r *TagRepo) FindByID(id string) (*model.Tag, error) {
	row := r.db.QueryRow(
		`SELECT id, name, color, created_at, updated_at FROM tags WHERE id = ?`, id)
	t := &model.Tag{}
	err := row.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("查询标签失败: %w", err)
	}
	return t, nil
}

// FindByName 根据名称查询标签
func (r *TagRepo) FindByName(name string) (*model.Tag, error) {
	row := r.db.QueryRow(
		`SELECT id, name, color, created_at, updated_at FROM tags WHERE name = ?`, name)
	t := &model.Tag{}
	err := row.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("按名称查询标签失败: %w", err)
	}
	return t, nil
}

// Create 创建标签
func (r *TagRepo) Create(tag *model.Tag) error {
	_, err := r.db.Exec(
		`INSERT INTO tags (id, name, color, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		tag.ID, tag.Name, tag.Color, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		return fmt.Errorf("创建标签失败: %w", err)
	}
	return nil
}

// Update 更新标签
func (r *TagRepo) Update(tag *model.Tag) error {
	_, err := r.db.Exec(
		`UPDATE tags SET name = ?, color = ?, updated_at = ? WHERE id = ?`,
		tag.Name, tag.Color, tag.UpdatedAt, tag.ID)
	if err != nil {
		return fmt.Errorf("更新标签失败: %w", err)
	}
	return nil
}

// Delete 删除标签
func (r *TagRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM tags WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除标签失败: %w", err)
	}
	return nil
}
