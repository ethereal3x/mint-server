-- 模型配置绑定用户，添加 user_id 字段
ALTER TABLE tb_model_config ADD COLUMN user_id VARCHAR(255) NOT NULL DEFAULT '' AFTER id;

-- 替换 model_type 单列唯一索引为 (model_type, user_id) 复合唯一索引
ALTER TABLE tb_model_config DROP INDEX model_type;
ALTER TABLE tb_model_config ADD UNIQUE INDEX idx_model_type_user (model_type, user_id);

-- user_id 查询索引
ALTER TABLE tb_model_config ADD INDEX idx_user_id (user_id);
