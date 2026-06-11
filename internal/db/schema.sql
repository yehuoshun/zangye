-- ============================================================
-- 藏叶 数据库初始化脚本
-- 数据库名: zang_ye
-- 字符集: utf8mb4（支持 emoji 等 4 字节字符）
-- 引擎: InnoDB（支持事务和外键）
-- ============================================================
-- 说明：所有 CREATE TABLE 使用 IF NOT EXISTS，确保幂等执行。
--       程序启动时自动执行此脚本，无需手动运行。
-- ============================================================

-- -----------------------------------------------------------
-- folders: 文件夹表（树形结构）
-- 一个文件夹可以包含子文件夹（通过 parent_id 自引用），
-- 形成一个树形目录结构来组织文件。
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS folders (
    id         VARCHAR(36) PRIMARY KEY,                       -- 文件夹唯一标识（UUID）
    name       VARCHAR(255) NOT NULL,                         -- 文件夹名称
    icon       VARCHAR(10) DEFAULT '📁',                      -- 文件夹图标（emoji）
    parent_id  VARCHAR(36) DEFAULT NULL,                      -- 父文件夹 ID（NULL 表示根文件夹）
    sort_order INT DEFAULT 0,                                 -- 排序序号（数字越小越靠前）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,           -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- 更新时间（自动更新）
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- tags: 标签表
-- 标签可以附加到任意文件上，用于分类和检索。
-- 标签名全局唯一。
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS tags (
    id         VARCHAR(36) PRIMARY KEY,                       -- 标签唯一标识（UUID）
    name       VARCHAR(128) NOT NULL UNIQUE,                  -- 标签名称（全局唯一）
    color      VARCHAR(32) DEFAULT 'gray',                    -- 标签颜色（用于 UI 展示）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP            -- 创建时间
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- files: 文件表
-- 记录文件的基本元信息，每个文件必须属于一个集合。
-- 集合删除时级联删除其下所有文件（ON DELETE CASCADE）。
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS files (
    id              VARCHAR(36) PRIMARY KEY,                   -- 文件唯一标识（UUID）
    folder_id       VARCHAR(36) NOT NULL,                      -- 所属文件夹 ID
    path            VARCHAR(1024) NOT NULL,                    -- 文件路径
    display_name    VARCHAR(512) DEFAULT NULL,                 -- 显示名称（可不同于文件名）
    file_size       BIGINT DEFAULT 0,                          -- 文件大小（字节）
    mime_type       VARCHAR(255) DEFAULT NULL,                 -- MIME 类型
    sort_order      INT DEFAULT 0,                             -- 排序序号
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,       -- 创建时间
    FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE  -- 文件夹删除时级联删除文件
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- file_tags: 文件-标签关联表（多对多关系）
-- 联合主键确保同一文件不会被重复添加同一标签。
-- 文件或标签删除时，关联记录自动级联删除。
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS file_tags (
    file_id VARCHAR(36) NOT NULL,                              -- 文件 ID
    tag_id  VARCHAR(36) NOT NULL,                              -- 标签 ID
    PRIMARY KEY (file_id, tag_id),                             -- 联合主键，防止重复关联
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,  -- 文件删除时级联删除关联
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE     -- 标签删除时级联删除关联
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- settings: 系统设置表（键值对）
-- 存储应用级别的配置项，如主题、版本号等。
-- -----------------------------------------------------------
CREATE TABLE IF NOT EXISTS settings (
    `key`   VARCHAR(128) PRIMARY KEY,                          -- 设置键名
    `value` TEXT NOT NULL                                      -- 设置值
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------
-- 默认设置数据
-- 使用 INSERT IGNORE 确保只插入一次，不会覆盖已有数据。
-- -----------------------------------------------------------
INSERT IGNORE INTO settings (`key`, `value`) VALUES
    ('version', '1'),           -- 数据库版本号
    ('theme', 'dark'),          -- 默认主题（dark/light）
    ('layout', 'sidebar');      -- 默认布局模式