-- 模型策略配置表
CREATE TABLE IF NOT EXISTS `tb_agent_strategy_rules` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `rule_id` VARCHAR(50) NOT NULL,
    `api_key` VARCHAR(128) NOT NULL COMMENT 'ApiKey数据',
    `agent_model` VARCHAR(50) DEFAULT NULL COMMENT '模型类型',
    `agent_manufacturer` VARCHAR(50) DEFAULT NULL COMMENT '模型所属公司',
    `agent_generate_type` VARCHAR(50) DEFAULT NULL COMMENT '模型生成类型',
    `url` VARCHAR(128) DEFAULT NULL COMMENT '模型调用地址',
    `max_tokens` INT NOT NULL DEFAULT 0 COMMENT '控制生成回复长度',
    `stream` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用流式输出',
    `temperature` FLOAT NOT NULL DEFAULT 0 COMMENT '控制生成随机性的参数',
    `top_p` FLOAT NOT NULL DEFAULT 0 COMMENT '核采样阈值',
    `n` INT NOT NULL DEFAULT 1 COMMENT '候选结果数量',
    `presence_penalty` FLOAT NOT NULL DEFAULT 0 COMMENT '已出现词惩罚系数',
    `frequency_penalty` FLOAT NOT NULL DEFAULT 0 COMMENT '高频词惩罚系数',
    `route` VARCHAR(128) DEFAULT NULL COMMENT '模型调用地址路径',
    `is_enabled` INT NOT NULL DEFAULT 1 COMMENT '是否启用',
    `created_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `updated_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tb_agent_strategy_rules_rule_id` (`rule_id`),
    KEY `idx_tb_agent_strategy_rules_generate_type` (`agent_generate_type`),
    KEY `idx_tb_agent_strategy_rules_model_manufacturer` (`agent_model`, `agent_manufacturer`),
    KEY `idx_tb_agent_strategy_rules_manufacturer_enabled` (`agent_manufacturer`, `is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型策略配置表';

-- 模型类型与厂商映射表
CREATE TABLE IF NOT EXISTS `tb_model_manufacturer_mapping` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `model_type` VARCHAR(50) NOT NULL COMMENT '模型类型',
    `manufacturer` VARCHAR(50) NOT NULL COMMENT '模型所属公司',
    `description` TEXT DEFAULT NULL COMMENT '模型与厂家关系描述',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tb_model_manufacturer_mapping_model_type` (`model_type`),
    KEY `idx_tb_model_manufacturer_mapping_manufacturer` (`manufacturer`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型类型与厂商映射表';

-- 用户问答记录表
CREATE TABLE IF NOT EXISTS `tb_user_agent_dialogues` (
    `dialogue_id` VARCHAR(255) NOT NULL COMMENT '对话ID',
    `record_id` VARCHAR(255) NOT NULL COMMENT '问答记录ID',
    `user_id` VARCHAR(255) NOT NULL COMMENT '用户ID',
    `model` VARCHAR(255) NOT NULL COMMENT '模型标识',
    `user_content` TEXT DEFAULT NULL COMMENT '用户问题',
    `agent_content` TEXT DEFAULT NULL COMMENT '模型回答',
    `total_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '总消耗Token',
    `user_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '用户问题Token',
    `agent_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '模型回答Token',
    `created_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `updated_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_tb_user_agent_dialogues_dialogue_id` (`dialogue_id`),
    KEY `idx_tb_user_agent_dialogues_user_id` (`user_id`),
    KEY `idx_tb_user_agent_dialogues_dialogue_created_time` (`dialogue_id`, `created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户问答记录表';
