-- 模型映射表：添加输入输出单价（每1M tokens）
ALTER TABLE `tb_model_manufacturer_mapping`
    ADD COLUMN `input_price` DECIMAL(10,8) NOT NULL DEFAULT 0 COMMENT '输入单价（每1M tokens）',
    ADD COLUMN `output_price` DECIMAL(10,8) NOT NULL DEFAULT 0 COMMENT '输出单价（每1M tokens）';

-- 对话记录表：添加输入输出费用
ALTER TABLE `tb_user_agent_dialogues`
    ADD COLUMN `input_cost` DECIMAL(12,8) NOT NULL DEFAULT 0 COMMENT '输入费用',
    ADD COLUMN `output_cost` DECIMAL(12,8) NOT NULL DEFAULT 0 COMMENT '输出费用';
