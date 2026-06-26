// Package service 提供业务逻辑层（Service Layer）
// 类比 Java 的 Service 层，处理业务逻辑，调用 repository 层
// Go 的 service 通常使用 struct 组合 repository，通过依赖注入组装
package service

import (
	"database/sql"
	"fmt"
	"time"
	"zangye/internal/model"
	"zangye/internal/repository"
	"zangye/internal/util"
)

// FolderSvc 文件夹业务逻辑
type FolderSvc struct {
	folderRepo *repository.FolderRepo
	fileRepo   *repository.FileRepo
	db         *sql.DB
}

// NewFolderSvc 创建 FolderSvc 实例
func NewFolderSvc(folderRepo *repository.FolderRepo, fileRepo *repository.FileRepo, db *sql.DB) *FolderSvc {
	return &FolderSvc{
		folderRepo: folderRepo,
		fileRepo:   fileRepo,
		db:         db,
	}
}

// GetTree 获取文件夹树形结构
// 一次性加载所有文件夹，在内存中组装成树
// 类比 Java 的递归组装树形菜单
func (s *FolderSvc) GetTree() ([]*model.FolderTreeItem, error) {
	folders, err := s.folderRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 构建 ID -> Folder 的映射
	// Go 的 map ≈ Java 的 HashMap
	folderMap := make(map[string]*model.FolderTreeItem)
	for _, f := range folders {
		// 查询该文件夹下的直接文件数量
		fileCount, _ := s.fileRepo.CountByFolderID(f.ID)
		folderMap[f.ID] = &model.FolderTreeItem{
			Folder:    *f,
			FileCount: fileCount,
		}
	}

	// 组装树形结构
	var roots []*model.FolderTreeItem
	for _, item := range folderMap {
		if item.ParentID == nil || *item.ParentID == "" {
			// 没有父节点 = 根节点
			roots = append(roots, item)
		} else {
			// 有父节点，挂到父节点的 Children 下
			if parent, ok := folderMap[*item.ParentID]; ok {
				parent.Children = append(parent.Children, &item.Folder)
			} else {
				// 父节点不存在（可能被删了），作为根节点
				roots = append(roots, item)
			}
		}
	}

	return roots, nil
}

// GetByID 根据 ID 获取文件夹
func (s *FolderSvc) GetByID(id string) (*model.Folder, error) {
	return s.folderRepo.FindByID(id)
}

// Create 创建文件夹
func (s *FolderSvc) Create(name string, parentID *string) (*model.Folder, error) {
	now := time.Now()
	folder := &model.Folder{
		ID:        util.GenerateUUID(),
		Name:      name,
		ParentID:  parentID,
		SortOrder: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.folderRepo.Create(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// Update 更新文件夹
func (s *FolderSvc) Update(id, name string, parentID *string) (*model.Folder, error) {
	folder, err := s.folderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, fmt.Errorf("文件夹不存在")
	}

	folder.Name = name
	folder.ParentID = parentID
	folder.UpdatedAt = time.Now()

	if err := s.folderRepo.Update(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// Delete 删除文件夹
// 注意：删除文件夹不会删除其下的文件，只是将文件的 folder_id 置空
func (s *FolderSvc) Delete(id string) error {
	return s.folderRepo.Delete(id)
}

// GetStats 获取文件夹递归统计信息
func (s *FolderSvc) GetStats(id string) (*model.FolderStats, error) {
	// 递归获取所有子文件夹 ID
	allFolders, err := s.folderRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 构建父子关系映射
	childrenMap := make(map[string][]string)
	for _, f := range allFolders {
		if f.ParentID != nil {
			childrenMap[*f.ParentID] = append(childrenMap[*f.ParentID], f.ID)
		}
	}

	// 递归收集所有子文件夹 ID
	var allIDs []string
	s.collectChildIDs(id, childrenMap, &allIDs)
	allIDs = append(allIDs, id) // 包含自身

	// 统计
	totalFolders := int64(len(allIDs) - 1) // 减去自身
	totalFiles, err := s.fileRepo.CountByFolderIDRecursive(allIDs)
	if err != nil {
		return nil, err
	}
	totalSize, err := s.fileRepo.SumSizeByFolderIDRecursive(allIDs)
	if err != nil {
		return nil, err
	}

	return &model.FolderStats{
		TotalFolders: totalFolders,
		TotalFiles:   totalFiles,
		TotalSize:    totalSize,
	}, nil
}

// collectChildIDs 递归收集所有子文件夹 ID
// Go 的切片是引用类型，通过指针传递修改
func (s *FolderSvc) collectChildIDs(parentID string, childrenMap map[string][]string, result *[]string) {
	if children, ok := childrenMap[parentID]; ok {
		for _, childID := range children {
			*result = append(*result, childID)
			s.collectChildIDs(childID, childrenMap, result)
		}
	}
}
