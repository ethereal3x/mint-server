-- 文件上传记录表
CREATE TABLE IF NOT EXISTS `tb_file_uploads` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `object_name` VARCHAR(255) NOT NULL COMMENT '存储对象名',
    `original_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
    `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    `content_type` VARCHAR(100) NOT NULL COMMENT 'MIME 类型',
    `url` VARCHAR(512) NOT NULL COMMENT '访问 URL',
    `upload_id` VARCHAR(255) DEFAULT NULL COMMENT '分片上传任务 ID',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/uploading/completed/canceled/failed',
    `uploaded_size` BIGINT NOT NULL DEFAULT 0 COMMENT '已上传字节数',
    `user_id` VARCHAR(255) NOT NULL COMMENT '用户 ID',
    `created_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `updated_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_object_name` (`object_name`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件上传记录表';
