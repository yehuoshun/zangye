// Package service 提供业务逻辑层
package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"zangye/internal/model"
	"zangye/internal/repository"
	"zangye/internal/util"
)

// FileSvc 文件业务逻辑
type FileSvc struct {
	fileRepo    *repository.FileRepo
	folderRepo  *repository.FolderRepo
	tagRepo     *repository.TagRepo
	fileTagRepo *repository.FileTagRepo
	db          *sql.DB
}

// NewFileSvc 创建 FileSvc 实例
func NewFileSvc(
	fileRepo *repository.FileRepo,
	folderRepo *repository.FolderRepo,
	tagRepo *repository.TagRepo,
	fileTagRepo *repository.FileTagRepo,
	db *sql.DB,
) *FileSvc {
	return &FileSvc{
		fileRepo:    fileRepo,
		folderRepo:  folderRepo,
		tagRepo:     tagRepo,
		fileTagRepo: fileTagRepo,
		db:          db,
	}
}

// Query 查询文件列表
func (s *FileSvc) Query(q model.FileQuery) (*model.FileListResult, error) {
	result, err := s.fileRepo.Query(q)
	if err != nil {
		return nil, err
	}

	// 为每个文件加载标签
	for i := range result.Files {
		tags, err := s.fileTagRepo.FindByFileID(result.Files[i].ID)
		if err == nil {
			result.Files[i].Tags = tags
		}
	}

	return result, nil
}

// GetByID 根据 ID 获取文件详情（含标签）
func (s *FileSvc) GetByID(id string) (*model.File, error) {
	file, err := s.fileRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, nil
	}

	// 加载标签
	tags, err := s.fileTagRepo.FindByFileID(id)
	if err == nil {
		file.Tags = tags
	}

	return file, nil
}

// Create 创建文件
func (s *FileSvc) Create(folderID *string, name, paths, description string, fileSize int64) (*model.File, error) {
	now := time.Now()
	fileType := util.GetFileType(name)

	file := &model.File{
		ID:        util.GenerateUUID(),
		FolderID:  folderID,
		Name:      name,
		FileType:  fileType,
		FileSize:  fileSize,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 处理 paths（JSON 数组字符串）
	if paths != "" {
		file.Paths = sql.NullString{String: paths, Valid: true}
	}

	// 处理 description
	if description != "" {
		file.Description = sql.NullString{String: description, Valid: true}
	}

	if err := s.fileRepo.Create(file); err != nil {
		return nil, err
	}
	return file, nil
}

// Update 更新文件
func (s *FileSvc) Update(id string, folderID *string, name, paths, description string, fileSize int64) (*model.File, error) {
	file, err := s.fileRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("文件不存在")
	}

	file.FolderID = folderID
	file.Name = name
	file.FileType = util.GetFileType(name)
	file.FileSize = fileSize
	file.UpdatedAt = time.Now()

	if paths != "" {
		file.Paths = sql.NullString{String: paths, Valid: true}
	}
	if description != "" {
		file.Description = sql.NullString{String: description, Valid: true}
	}

	if err := s.fileRepo.Update(file); err != nil {
		return nil, err
	}

	// 重新加载标签
	tags, _ := s.fileTagRepo.FindByFileID(id)
	file.Tags = tags

	return file, nil
}

// SoftDelete 软删除文件（移到回收站）
func (s *FileSvc) SoftDelete(id string) error {
	now := time.Now()
	return s.fileRepo.SoftDelete(id, now)
}

// Restore 恢复文件
func (s *FileSvc) Restore(id string) error {
	return s.fileRepo.Restore(id)
}

// HardDelete 彻底删除文件
func (s *FileSvc) HardDelete(id string) error {
	// 先删除标签关联
	if err := s.fileTagRepo.DeleteByFileID(id); err != nil {
		return err
	}
	return s.fileRepo.HardDelete(id)
}

// GetTrashFiles 获取回收站文件列表
func (s *FileSvc) GetTrashFiles() ([]model.File, error) {
	q := model.FileQuery{
		Trash:    true,
		OrderBy:  "created_at",
		OrderDir: "DESC",
		Page:     1,
		PageSize: 1000,
	}
	result, err := s.fileRepo.Query(q)
	if err != nil {
		return nil, err
	}
	return result.Files, nil
}

// GetTags 获取文件的标签列表
func (s *FileSvc) GetTags(fileID string) ([]model.Tag, error) {
	return s.fileTagRepo.FindByFileID(fileID)
}

// SetTags 设置文件的标签列表
func (s *FileSvc) SetTags(fileID string, tagIDs []string) error {
	// 使用事务
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	// defer 处理事务回滚
	defer tx.Rollback()

	if err := s.fileTagRepo.SetTags(fileID, tagIDs, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// AddTag 为文件添加单个标签
func (s *FileSvc) AddTag(fileID, tagID string) error {
	return s.fileTagRepo.AddTag(fileID, tagID)
}

// RemoveTag 移除文件的单个标签
func (s *FileSvc) RemoveTag(fileID, tagID string) error {
	return s.fileTagRepo.RemoveTag(fileID, tagID)
}

// GetContentPath 获取文件的实际路径（遍历 paths JSON 数组，返回第一个存在的路径）
// 用于文件预览
func (s *FileSvc) GetContentPath(fileID string) (string, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return "", err
	}
	if file == nil {
		return "", fmt.Errorf("文件不存在")
	}

	// 解析 paths JSON 数组
	if !file.Paths.Valid || file.Paths.String == "" {
		return "", fmt.Errorf("文件路径为空")
	}

	var paths []string
	if err := json.Unmarshal([]byte(file.Paths.String), &paths); err != nil {
		return "", fmt.Errorf("解析文件路径失败: %w", err)
	}

	// 遍历路径，返回第一个存在的
	for _, p := range paths {
		// 去除首尾空格和引号
		p = strings.Trim(p, "\" ")
		if p != "" && util.FileExists(p) {
			return p, nil
		}
	}

	return "", fmt.Errorf("所有文件路径均不存在")
}
