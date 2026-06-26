// Package repository 提供数据访问层（Data Access Layer）
// 类比 Java 的 DAO 层或 MyBatis Mapper
// 每个 repository 只做 SQL 操作，不处理业务逻辑
package repository

import (
	"database/sql"
	"fmt"
	"zangye/internal/model"
)

// FolderRepo 文件夹数据访问
// Go 的 struct 组合 ≈ Java 的继承，但更灵活
type FolderRepo struct {
	db *sql.DB // 数据库连接，依赖注入
}

// NewFolderRepo 创建 FolderRepo 实例
// 类比 Java 的构造函数注入
func NewFolderRepo(db *sql.DB) *FolderRepo {
	return &FolderRepo{db: db}
}

// FindAll 查询所有文件夹（用于构建树形结构）
// 一次性加载所有文件夹，在内存中组装树
func (r *FolderRepo) FindAll() ([]*model.Folder, error) {
	// 参数化查询，防止 SQL 注入
	// Go 的 ? 占位符 ≈ Java PreparedStatement 的 ?
	rows, err := r.db.Query(
		`SELECT id, name, parent_id, sort_order, created_at, updated_at
		 FROM folders ORDER BY sort_order ASC, name ASC`)
	if err != nil {
		return nil, fmt.Errorf("查询文件夹列表失败: %w", err)
	}
	// defer ≈ Java 的 finally，用于关闭资源
	// 记住：永远不要忽略 error
	defer rows.Close()

	return scanFolders(rows)
}

// FindByID 根据 ID 查询文件夹
func (r *FolderRepo) FindByID(id string) (*model.Folder, error) {
	row := r.db.QueryRow(
		`SELECT id, name, parent_id, sort_order, created_at, updated_at
		 FROM folders WHERE id = ?`, id)

	folder, err := scanFolder(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没找到返回 nil, nil，不返回 error
		}
		return nil, fmt.Errorf("查询文件夹失败: %w", err)
	}
	return folder, nil
}

// Create 创建文件夹
func (r *FolderRepo) Create(folder *model.Folder) error {
	_, err := r.db.Exec(
		`INSERT INTO folders (id, name, parent_id, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		folder.ID, folder.Name, folder.ParentID, folder.SortOrder,
		folder.CreatedAt, folder.UpdatedAt)
	if err != nil {
		return fmt.Errorf("创建文件夹失败: %w", err)
	}
	return nil
}

// Update 更新文件夹
func (r *FolderRepo) Update(folder *model.Folder) error {
	_, err := r.db.Exec(
		`UPDATE folders SET name = ?, parent_id = ?, sort_order = ?, updated_at = ?
		 WHERE id = ?`,
		folder.Name, folder.ParentID, folder.SortOrder, folder.UpdatedAt, folder.ID)
	if err != nil {
		return fmt.Errorf("更新文件夹失败: %w", err)
	}
	return nil
}

// Delete 删除文件夹
func (r *FolderRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM folders WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除文件夹失败: %w", err)
	}
	return nil
}

// CountByParentID 统计指定父文件夹下的子文件夹数量
func (r *FolderRepo) CountByParentID(parentID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM folders WHERE parent_id = ?`, parentID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计子文件夹数量失败: %w", err)
	}
	return count, nil
}

// scanFolders 扫描多行数据到 Folder 切片
// 切片 slice[] ≈ ArrayList，但底层是指向数组的指针
func scanFolders(rows *sql.Rows) ([]*model.Folder, error) {
	var folders []*model.Folder
	for rows.Next() {
		f := &model.Folder{}
		// 注意：ParentID 是 *string，需要特殊处理
		// Go 的 sql.Scan 不能直接扫描到 *string，需要用中间变量
		var parentID sql.NullString
		err := rows.Scan(&f.ID, &f.Name, &parentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描文件夹行数据失败: %w", err)
		}
		if parentID.Valid {
			f.ParentID = &parentID.String
		}
		folders = append(folders, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历文件夹行数据失败: %w", err)
	}
	return folders, nil
}

// scanFolder 扫描单行数据到 Folder
func scanFolder(row *sql.Row) (*model.Folder, error) {
	f := &model.Folder{}
	var parentID sql.NullString
	err := row.Scan(&f.ID, &f.Name, &parentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		f.ParentID = &parentID.String
	}
	return f, nil
}
