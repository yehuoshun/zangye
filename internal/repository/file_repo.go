// Package repository 提供数据访问层
package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"zangye/internal/model"
)

// FileRepo 文件数据访问
type FileRepo struct {
	db *sql.DB
}

// NewFileRepo 创建 FileRepo 实例
func NewFileRepo(db *sql.DB) *FileRepo {
	return &FileRepo{db: db}
}

// baseColumns 基础查询列（复用避免重复编写）
// Go 的常量不能用切片，所以用函数返回
func (r *FileRepo) baseColumns() string {
	return `id, folder_id, name, paths, file_type, file_size, description, deleted_at, created_at, updated_at`
}

// Query 查询文件列表（支持分页、搜索、筛选）
// 动态构建 SQL 查询，类比 Java 的 MyBatis 动态 SQL
func (r *FileRepo) Query(q model.FileQuery) (*model.FileListResult, error) {
	// 构建 WHERE 条件
	var conditions []string
	var args []interface{} // Go 的 interface{} = any，类比 Java 的 Object

	// 回收站筛选
	if q.Trash {
		conditions = append(conditions, "deleted_at IS NOT NULL")
	} else {
		conditions = append(conditions, "deleted_at IS NULL")
	}

	// 按文件夹筛选
	if q.FolderID != nil && *q.FolderID != "" {
		conditions = append(conditions, "folder_id = ?")
		args = append(args, *q.FolderID)
	}

	// 按关键字搜索（匹配 name/paths/description）
	if q.Keyword != nil && *q.Keyword != "" {
		keyword := "%" + *q.Keyword + "%"
		conditions = append(conditions, "(name LIKE ? OR paths LIKE ? OR description LIKE ?)")
		args = append(args, keyword, keyword, keyword)
	}

	// 按文件类型筛选
	if q.FileType != nil && *q.FileType != "" {
		conditions = append(conditions, "file_type = ?")
		args = append(args, *q.FileType)
	}

	// 按标签筛选（通过 file_tags 关联查询）
	if q.TagID != nil && *q.TagID != "" {
		conditions = append(conditions, "id IN (SELECT file_id FROM file_tags WHERE tag_id = ?)")
		args = append(args, *q.TagID)
	}

	// 拼接 WHERE 子句
	whereClause := strings.Join(conditions, " AND ")

	// 排序
	orderBy := "created_at DESC" // 默认按创建时间降序
	switch q.OrderBy {
	case "name":
		orderBy = "name " + q.OrderDir
	case "file_type":
		orderBy = "file_type " + q.OrderDir
	case "file_size":
		orderBy = "file_size " + q.OrderDir
	case "created_at":
		orderBy = "created_at " + q.OrderDir
	}
	if q.OrderDir == "" {
		q.OrderDir = "DESC"
	}

	// 分页
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 50
	}
	offset := (q.Page - 1) * q.PageSize

	// 查询总数
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM files WHERE %s", whereClause)
	var total int64
	err := r.db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("统计文件数量失败: %w", err)
	}

	// 查询列表
	listSQL := fmt.Sprintf("SELECT %s FROM files WHERE %s ORDER BY %s LIMIT ? OFFSET ?",
		r.baseColumns(), whereClause, orderBy)
	queryArgs := append(args, q.PageSize, offset)
	rows, err := r.db.Query(listSQL, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("查询文件列表失败: %w", err)
	}
	defer rows.Close()

	files, err := scanFiles(rows)
	if err != nil {
		return nil, err
	}

	return &model.FileListResult{
		Files: files,
		Total: total,
		Page:  q.Page,
		Size:  q.PageSize,
	}, nil
}

// FindByID 根据 ID 查询文件
func (r *FileRepo) FindByID(id string) (*model.File, error) {
	row := r.db.QueryRow(
		`SELECT `+r.baseColumns()+` FROM files WHERE id = ?`, id)
	return scanFile(row)
}

// Create 创建文件
func (r *FileRepo) Create(file *model.File) error {
	_, err := r.db.Exec(
		`INSERT INTO files (id, folder_id, name, paths, file_type, file_size, description, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		file.ID, file.FolderID, file.Name, file.Paths, file.FileType,
		file.FileSize, file.Description, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	return nil
}

// Update 更新文件
func (r *FileRepo) Update(file *model.File) error {
	_, err := r.db.Exec(
		`UPDATE files SET folder_id = ?, name = ?, paths = ?, file_type = ?,
		 file_size = ?, description = ?, updated_at = ? WHERE id = ?`,
		file.FolderID, file.Name, file.Paths, file.FileType,
		file.FileSize, file.Description, file.UpdatedAt, file.ID)
	if err != nil {
		return fmt.Errorf("更新文件失败: %w", err)
	}
	return nil
}

// SoftDelete 软删除文件（移到回收站）
func (r *FileRepo) SoftDelete(id string, deletedAt interface{}) error {
	_, err := r.db.Exec(`UPDATE files SET deleted_at = ? WHERE id = ?`, deletedAt, id)
	if err != nil {
		return fmt.Errorf("软删除文件失败: %w", err)
	}
	return nil
}

// Restore 恢复软删除的文件
func (r *FileRepo) Restore(id string) error {
	_, err := r.db.Exec(`UPDATE files SET deleted_at = NULL WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("恢复文件失败: %w", err)
	}
	return nil
}

// HardDelete 彻底删除文件
func (r *FileRepo) HardDelete(id string) error {
	_, err := r.db.Exec(`DELETE FROM files WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("彻底删除文件失败: %w", err)
	}
	return nil
}

// CountByFolderID 统计指定文件夹下的直接文件数量
func (r *FileRepo) CountByFolderID(folderID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM files WHERE folder_id = ? AND deleted_at IS NULL`, folderID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计文件数量失败: %w", err)
	}
	return count, nil
}

// CountAll 统计所有文件数量（不含回收站）
func (r *FileRepo) CountAll() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM files WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计文件总数失败: %w", err)
	}
	return count, nil
}

// CountByType 按文件类型统计数量
func (r *FileRepo) CountByType(fileType string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM files WHERE file_type = ? AND deleted_at IS NULL`, fileType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("按类型统计文件数量失败: %w", err)
	}
	return count, nil
}

// CountByCategory 按文件大类统计数量
// categoryPrefixes 是扩展名前缀列表，如 []string{"jpg","jpeg","png"} 表示图片
func (r *FileRepo) CountByCategory(categoryPrefixes []string) (int64, error) {
	if len(categoryPrefixes) == 0 {
		return 0, nil
	}
	// 构建 IN 子句的占位符
	placeholders := make([]string, len(categoryPrefixes))
	args := make([]interface{}, len(categoryPrefixes))
	for i, p := range categoryPrefixes {
		placeholders[i] = "?"
		args[i] = p
	}
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM files WHERE file_type IN (%s) AND deleted_at IS NULL`,
		strings.Join(placeholders, ","))
	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("按大类统计文件数量失败: %w", err)
	}
	return count, nil
}

// SumAllSize 计算所有文件总大小
func (r *FileRepo) SumAllSize() (int64, error) {
	var total sql.NullInt64
	err := r.db.QueryRow(`SELECT SUM(file_size) FROM files WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("统计文件总大小失败: %w", err)
	}
	if total.Valid {
		return total.Int64, nil
	}
	return 0, nil
}

// CountByFolderIDRecursive 递归统计文件夹及其子文件夹下的文件数量
func (r *FileRepo) CountByFolderIDRecursive(folderIDs []string) (int64, error) {
	if len(folderIDs) == 0 {
		return 0, nil
	}
	placeholders := make([]string, len(folderIDs))
	args := make([]interface{}, len(folderIDs))
	for i, id := range folderIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM files WHERE folder_id IN (%s) AND deleted_at IS NULL`,
		strings.Join(placeholders, ","))
	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("递归统计文件数量失败: %w", err)
	}
	return count, nil
}

// SumSizeByFolderIDRecursive 递归统计文件夹及其子文件夹下的文件总大小
func (r *FileRepo) SumSizeByFolderIDRecursive(folderIDs []string) (int64, error) {
	if len(folderIDs) == 0 {
		return 0, nil
	}
	placeholders := make([]string, len(folderIDs))
	args := make([]interface{}, len(folderIDs))
	for i, id := range folderIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(
		`SELECT COALESCE(SUM(file_size), 0) FROM files WHERE folder_id IN (%s) AND deleted_at IS NULL`,
		strings.Join(placeholders, ","))
	var total sql.NullInt64
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("递归统计文件大小失败: %w", err)
	}
	if total.Valid {
		return total.Int64, nil
	}
	return 0, nil
}

// scanFiles 扫描多行数据到 File 切片
func scanFiles(rows *sql.Rows) ([]model.File, error) {
	var files []model.File
	for rows.Next() {
		f, err := scanFileRow(rows)
		if err != nil {
			return nil, err
		}
		files = append(files, *f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历文件行数据失败: %w", err)
	}
	return files, nil
}

// scanFile 扫描单行数据到 File
func scanFile(row *sql.Row) (*model.File, error) {
	f := &model.File{}
	var folderID, paths, description sql.NullString
	var deletedAt sql.NullTime
	err := row.Scan(&f.ID, &folderID, &f.Name, &paths, &f.FileType,
		&f.FileSize, &description, &deletedAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("扫描文件行数据失败: %w", err)
	}
	if folderID.Valid {
		f.FolderID = &folderID.String
	}
	f.Paths = paths
	f.Description = description
	if deletedAt.Valid {
		f.DeletedAt = &deletedAt.Time
	}
	return f, nil
}

// scanFileRow 从 *sql.Rows 扫描单行数据到 File
func scanFileRow(rows *sql.Rows) (*model.File, error) {
	f := &model.File{}
	var folderID, paths, description sql.NullString
	var deletedAt sql.NullTime
	err := rows.Scan(&f.ID, &folderID, &f.Name, &paths, &f.FileType,
		&f.FileSize, &description, &deletedAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("扫描文件行数据失败: %w", err)
	}
	if folderID.Valid {
		f.FolderID = &folderID.String
	}
	f.Paths = paths
	f.Description = description
	if deletedAt.Valid {
		f.DeletedAt = &deletedAt.Time
	}
	return f, nil
}
