ALTER TABLE tb_model_config ADD COLUMN supports_multimodal TINYINT(1) NOT NULL DEFAULT 0 AFTER is_enabled;
