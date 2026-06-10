-- 模型配置表（合并策略规则 + 模型映射）
CREATE TABLE IF NOT EXISTS `tb_model_config` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `model_type` VARCHAR(50) NOT NULL COMMENT '模型标识',
    `manufacturer` VARCHAR(50) NOT NULL COMMENT '厂商',
    `description` TEXT DEFAULT NULL COMMENT '模型描述',
    `input_price` DECIMAL(10,8) NOT NULL DEFAULT 0 COMMENT '输入单价（每1M tokens）',
    `output_price` DECIMAL(10,8) NOT NULL DEFAULT 0 COMMENT '输出单价（每1M tokens）',
    `api_key` VARCHAR(512) NOT NULL COMMENT 'API密钥',
    `url` VARCHAR(128) DEFAULT NULL COMMENT 'API地址',
    `max_tokens` INT NOT NULL DEFAULT 0 COMMENT '最大token数',
    `stream` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用流式',
    `temperature` FLOAT NOT NULL DEFAULT 0 COMMENT '温度参数',
    `top_p` FLOAT NOT NULL DEFAULT 0 COMMENT '核采样阈值',
    `n` INT NOT NULL DEFAULT 1 COMMENT '候选结果数量',
    `presence_penalty` FLOAT NOT NULL DEFAULT 0 COMMENT '已出现词惩罚系数',
    `frequency_penalty` FLOAT NOT NULL DEFAULT 0 COMMENT '高频词惩罚系数',
    `agent_generate_type` VARCHAR(50) DEFAULT NULL COMMENT '生成类型',
    `route` VARCHAR(128) DEFAULT NULL COMMENT '路由标识',
    `is_enabled` INT NOT NULL DEFAULT 1 COMMENT '是否启用',
    `created_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `updated_time` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_model_type` (`model_type`),
    KEY `idx_manufacturer` (`manufacturer`),
    KEY `idx_manufacturer_enabled` (`manufacturer`, `is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型配置表';

-- 从旧表迁移数据（model_type 每个厂商取第一条启用的策略配置）
INSERT INTO `tb_model_config` (
    `model_type`, `manufacturer`, `description`,
    `input_price`, `output_price`,
    `api_key`, `url`, `max_tokens`, `stream`,
    `temperature`, `top_p`, `n`, `presence_penalty`, `frequency_penalty`,
    `agent_generate_type`, `route`, `is_enabled`
)
SELECT
    m.model_type, m.manufacturer, m.description,
    m.input_price, m.output_price,
    COALESCE(s.api_key, ''), s.url, COALESCE(s.max_tokens, 0), COALESCE(s.stream, 1),
    COALESCE(s.temperature, 0), COALESCE(s.top_p, 0), COALESCE(s.n, 1),
    COALESCE(s.presence_penalty, 0), COALESCE(s.frequency_penalty, 0),
    s.agent_generate_type, s.route, COALESCE(s.is_enabled, 1)
FROM `tb_model_manufacturer_mapping` m
LEFT JOIN `tb_agent_strategy_rules` s
    ON m.manufacturer = s.agent_manufacturer
    AND s.id = (
        SELECT MIN(s2.id) FROM `tb_agent_strategy_rules` s2
        WHERE s2.agent_manufacturer = m.manufacturer AND s2.is_enabled = 1
    );

-- 删除旧表
DROP TABLE IF EXISTS `tb_agent_strategy_rules`;
DROP TABLE IF EXISTS `tb_model_manufacturer_mapping`;
