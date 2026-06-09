# Admin API 文档

基础地址：`http://{host}:9999/admin/v1`

## 策略规则管理 (Strategy)

### 列表查询

```
GET /admin/v1/strategies?page=1&page_size=20
```

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int32 | 否 | 页码，从 1 开始，默认 1 |
| page_size | int32 | 否 | 每页条数，默认 20，最大 100 |

**响应**

```json
{
  "strategies": [
    {
      "id": 1,
      "rule_id": "uuid-string",
      "api_key": "sk-****",
      "agent_model": "gpt-4",
      "agent_manufacturer": "openai",
      "agent_generate_type": "chat",
      "url": "https://api.openai.com/v1",
      "max_tokens": 2048,
      "stream": true,
      "temperature": 0.7,
      "top_p": 0.9,
      "n": 1,
      "presence_penalty": 0.0,
      "frequency_penalty": 0.0,
      "route": "primary",
      "is_enabled": 1
    }
  ],
  "total": 10
}
```

### 查询详情

```
GET /admin/v1/strategies/{rule_id}
```

**Path 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| rule_id | string | 是 | 策略规则 ID |

**响应**: Strategy 对象

### 创建

```
POST /admin/v1/strategies
```

**Body**

```json
{
  "api_key": "sk-xxx",
  "agent_model": "gpt-4",
  "agent_manufacturer": "openai",
  "agent_generate_type": "chat",
  "url": "https://api.openai.com/v1",
  "max_tokens": 2048,
  "stream": true,
  "temperature": 0.7,
  "top_p": 0.9,
  "n": 1,
  "presence_penalty": 0.0,
  "frequency_penalty": 0.0,
  "route": "primary"
}
```

必填字段: `api_key`, `agent_manufacturer`

**响应**: Strategy 对象

### 更新

```
PUT /admin/v1/strategies/{rule_id}
```

**Path 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| rule_id | string | 是 | 策略规则 ID |

**Body**

```json
{
  "api_key": "sk-xxx",
  "agent_model": "gpt-4",
  "agent_manufacturer": "openai",
  "agent_generate_type": "chat",
  "url": "https://api.openai.com/v1",
  "max_tokens": 2048,
  "stream": true,
  "temperature": 0.7,
  "top_p": 0.9,
  "n": 1,
  "presence_penalty": 0.0,
  "frequency_penalty": 0.0,
  "route": "primary",
  "is_enabled": 1
}
```

**响应**: Strategy 对象

### 删除

```
DELETE /admin/v1/strategies/{rule_id}
```

**响应**: 空

### Strategy 对象

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int32 | 自增主键 |
| rule_id | string | 策略规则唯一标识 |
| api_key | string | API 密钥（响应中已脱敏） |
| agent_model | string | 模型名称 |
| agent_manufacturer | string | 模型厂商 |
| agent_generate_type | string | 生成类型 |
| url | string | API 地址 |
| max_tokens | int32 | 最大 token 数 |
| stream | bool | 是否启用流式 |
| temperature | float | 温度参数 |
| top_p | float | Top-P 采样参数 |
| n | int32 | 生成候选数 |
| presence_penalty | float | 存在惩罚系数 |
| frequency_penalty | float | 频率惩罚系数 |
| route | string | 路由标识 |
| is_enabled | int32 | 是否启用：0-禁用，1-启用 |

---

## 模型映射管理 (Mapping)

### 列表查询

```
GET /admin/v1/mappings?manufacturer=openai
```

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| manufacturer | string | 否 | 厂商过滤条件，为空时返回全部 |

**响应**

```json
{
  "mappings": [
    {
      "id": 1,
      "model_type": "gpt-4",
      "manufacturer": "openai",
      "description": "OpenAI GPT-4 模型"
    }
  ]
}
```

### 查询详情

```
GET /admin/v1/mappings/{id}
```

**Path 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int32 | 是 | 映射 ID |

**响应**: ModelMapping 对象

### 创建

```
POST /admin/v1/mappings
```

**Body**

```json
{
  "model_type": "gpt-4",
  "manufacturer": "openai",
  "description": "OpenAI GPT-4 模型"
}
```

必填字段: `model_type`, `manufacturer`

**响应**: ModelMapping 对象

### 更新

```
PUT /admin/v1/mappings/{id}
```

**Body**

```json
{
  "model_type": "gpt-4",
  "manufacturer": "openai",
  "description": "OpenAI GPT-4 模型"
}
```

**响应**: ModelMapping 对象

### 删除

```
DELETE /admin/v1/mappings/{id}
```

**响应**: 空

### ModelMapping 对象

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int32 | 自增主键 |
| model_type | string | 模型类型标识 |
| manufacturer | string | 厂商名称 |
| description | string | 模型描述 |

---

## 错误码

| HTTP 状态码 | gRPC 码 | 说明 |
|------------|---------|------|
| 400 | InvalidArgument | 请求参数校验失败（缺少必填字段等） |
| 404 | NotFound | 资源不存在 |
| 500 | Internal | 服务端内部错误 |
