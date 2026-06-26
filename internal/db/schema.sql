-- ============================================================
-- 藏叶 (ZangYe) 虚拟文件管理器 - 数据库建表脚本
-- 数据库：MySQL 8.0
-- ============================================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS zangye
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE zangye;

-- ============================================================
-- folders：文件夹表，parent_id 自引用实现树形结构
-- 类比：Java 中可以用 @ManyToOne 自关联，Go 中直接用 parent_id 字段
-- ============================================================
CREATE TABLE IF NOT EXISTS folders (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY COMMENT 'UUID 主键',
    name        VARCHAR(255)    NOT NULL             COMMENT '文件夹名称',
    parent_id   VARCHAR(36)     NULL                 COMMENT '父文件夹 ID，NULL 表示根目录',
    sort_order  INT             NOT NULL DEFAULT 0   COMMENT '排序序号',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件夹表';

-- ============================================================
-- files：文件表
-- paths TEXT -> JSON 数组，如 ["Z:\\backup\\a.jpg", "115:/xxx/a.jpg"]
-- file_type 自动从扩展名推断，如 "jpg", "mp4", "pdf"
-- deleted_at 非 NULL 表示已删除（软删除），移到回收站
-- ============================================================
CREATE TABLE IF NOT EXISTS files (
    id           VARCHAR(36)     NOT NULL PRIMARY KEY COMMENT 'UUID 主键',
    folder_id    VARCHAR(36)     NULL                 COMMENT '所属文件夹 ID',
    name         VARCHAR(255)    NOT NULL             COMMENT '文件名（含扩展名）',
    paths        TEXT            NULL                 COMMENT '文件路径 JSON 数组',
    file_type    VARCHAR(32)     NOT NULL             COMMENT '文件类型（扩展名）',
    file_size    BIGINT          NOT NULL DEFAULT 0   COMMENT '文件大小（字节）',
    description  TEXT            NULL                 COMMENT '文件描述',
    deleted_at   TIMESTAMP       NULL                 COMMENT '软删除时间，NULL=未删除',
    created_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_folder_id (folder_id),
    INDEX idx_file_type (file_type),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- ============================================================
-- tags：标签表，color 使用 HSL 随机生成
-- ============================================================
CREATE TABLE IF NOT EXISTS tags (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY COMMENT 'UUID 主键',
    name        VARCHAR(64)     NOT NULL             COMMENT '标签名称',
    color       VARCHAR(7)      NOT NULL DEFAULT '#4A90D9' COMMENT '标签颜色（十六进制）',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- ============================================================
-- file_tags：文件-标签多对多关联表
-- ============================================================
CREATE TABLE IF NOT EXISTS file_tags (
    file_id     VARCHAR(36)     NOT NULL COMMENT '文件 ID',
    tag_id      VARCHAR(36)     NOT NULL COMMENT '标签 ID',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (file_id, tag_id),
    INDEX idx_tag_id (tag_id),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件-标签关联表';

-- ============================================================
-- settings：键值对设置表，含打开方式配置
-- ============================================================
CREATE TABLE IF NOT EXISTS settings (
    setting_key     VARCHAR(128)    NOT NULL PRIMARY KEY COMMENT '设置键名',
    setting_value   TEXT            NULL                 COMMENT '设置值（JSON 格式）',
    updated_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设置表';

-- ============================================================
-- 初始化默认设置
-- ============================================================
INSERT IGNORE INTO settings (setting_key, setting_value) VALUES
    ('open_with_programs', '[]'),
    ('default_view', 'list'),
    ('items_per_page', '50');
