// Package service 提供业务逻辑层
package service

import (
	"time"
	"zangye/internal/model"
	"zangye/internal/repository"
	"zangye/internal/util"
)

// TagSvc 标签业务逻辑
type TagSvc struct {
	tagRepo     *repository.TagRepo
	fileTagRepo *repository.FileTagRepo
}

// NewTagSvc 创建 TagSvc 实例
func NewTagSvc(tagRepo *repository.TagRepo, fileTagRepo *repository.FileTagRepo) *TagSvc {
	return &TagSvc{
		tagRepo:     tagRepo,
		fileTagRepo: fileTagRepo,
	}
}

// GetAll 获取所有标签
func (s *TagSvc) GetAll() ([]*model.TagWithCount, error) {
	return s.tagRepo.FindAll()
}

// GetByID 根据 ID 获取标签
func (s *TagSvc) GetByID(id string) (*model.Tag, error) {
	return s.tagRepo.FindByID(id)
}

// Create 创建标签
// 如果标签名已存在，返回已存在的标签
func (s *TagSvc) Create(name string) (*model.Tag, error) {
	// 检查是否已存在同名标签
	existing, err := s.tagRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	now := time.Now()
	tag := &model.Tag{
		ID:        util.GenerateUUID(),
		Name:      name,
		Color:     util.GenerateRandomColor(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// Update 更新标签
func (s *TagSvc) Update(id, name, color string) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, nil
	}

	tag.Name = name
	if color != "" {
		tag.Color = color
	}
	tag.UpdatedAt = time.Now()

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// Delete 删除标签
func (s *TagSvc) Delete(id string) error {
	// 先删除关联
	if err := s.fileTagRepo.DeleteByTagID(id); err != nil {
		return err
	}
	return s.tagRepo.Delete(id)
}
