ALTER TABLE tb_model_config ADD COLUMN model_capabilities VARCHAR(512) NOT NULL DEFAULT '' AFTER description;

UPDATE tb_model_config
SET model_capabilities = CASE
    WHEN supports_multimodal = 1 THEN 'MODEL_CAPABILITY_TEXT_CHAT,MODEL_CAPABILITY_IMAGE_UNDERSTANDING'
    ELSE 'MODEL_CAPABILITY_TEXT_CHAT'
END
WHERE model_capabilities = '';
