-- 扩展 api_key 列长度以容纳 AES-256-CBC 加密后的密文
ALTER TABLE `tb_model_config` MODIFY COLUMN `api_key` VARCHAR(512) NOT NULL COMMENT 'API密钥';
