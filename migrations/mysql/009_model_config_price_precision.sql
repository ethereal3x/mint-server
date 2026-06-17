ALTER TABLE tb_model_config
    MODIFY COLUMN input_price DECIMAL(14,8) NOT NULL DEFAULT 0 COMMENT '输入单价（每1M tokens）',
    MODIFY COLUMN output_price DECIMAL(14,8) NOT NULL DEFAULT 0 COMMENT '输出单价（每1M tokens）';
