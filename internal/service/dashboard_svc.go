// Package service 提供业务逻辑层
package service

import (
	"zangye/internal/repository"
	"zangye/internal/util"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	FolderCount int64 `json:"folder_count"` // 文件夹总数
	FileCount   int64 `json:"file_count"`   // 文件总数
	ImageCount  int64 `json:"image_count"`  // 图片文件数
	VideoCount  int64 `json:"video_count"`  // 视频文件数
	AudioCount  int64 `json:"audio_count"`  // 音频文件数
	OtherCount  int64 `json:"other_count"`  // 其他文件数
	TotalSize   int64 `json:"total_size"`   // 总大小（字节）
	SizeText    string `json:"size_text"`   // 格式化后的总大小
}

// DashboardSvc 仪表盘业务逻辑
type DashboardSvc struct {
	folderRepo *repository.FolderRepo
	fileRepo   *repository.FileRepo
}

// NewDashboardSvc 创建 DashboardSvc 实例
func NewDashboardSvc(folderRepo *repository.FolderRepo, fileRepo *repository.FileRepo) *DashboardSvc {
	return &DashboardSvc{
		folderRepo: folderRepo,
		fileRepo:   fileRepo,
	}
}

// GetStats 获取仪表盘统计数据
func (s *DashboardSvc) GetStats() (*DashboardStats, error) {
	// 查询文件夹总数
	folderCount, err := s.folderRepo.CountByParentID("")
	if err != nil {
		return nil, err
	}

	// 查询文件总数
	fileCount, err := s.fileRepo.CountAll()
	if err != nil {
		return nil, err
	}

	// 按类型统计
	imageCount, _ := s.fileRepo.CountByCategory([]string{"jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "ico"})
	videoCount, _ := s.fileRepo.CountByCategory([]string{"mp4", "avi", "mov", "wmv", "flv", "mkv", "webm"})
	audioCount, _ := s.fileRepo.CountByCategory([]string{"mp3", "wav", "ogg", "flac", "aac", "wma"})
	otherCount := fileCount - imageCount - videoCount - audioCount
	if otherCount < 0 {
		otherCount = 0
	}

	// 总大小
	totalSize, err := s.fileRepo.SumAllSize()
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		FolderCount: folderCount,
		FileCount:   fileCount,
		ImageCount:  imageCount,
		VideoCount:  videoCount,
		AudioCount:  audioCount,
		OtherCount:  otherCount,
		TotalSize:   totalSize,
		SizeText:    util.FormatFileSize(totalSize),
	}, nil
}
