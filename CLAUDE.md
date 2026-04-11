# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 架构说明

### 后端 (Go — 模块名 `GoNexus`)

后端采用分层架构：`router → controller → service → dao`，共享基础设施位于 `common/`。

**请求流程：**
```
HTTP → router/ → middleware/jwt (鉴权) → controller/ → service/ → dao/ → MySQL/Redis
                                                      ↘ common/rabbitmq (异步写数据库)
```

**所有路由**前缀为 `/api/v1`。用户路由公开访问；AI 和文件路由需要 JWT（`Authorization: Bearer <token>`）。

```
/api/v1/user/...          # 注册、登录、验证码
/api/v1/ai/chat/...       # 会话、发送、流式、历史、TTS
/api/v1/file/...          # RAG 文档上传
```

**AI 子系统 (`common/aihelper/`)：**

核心抽象是 `AIHelper`，封装了模型和对话历史。全局单例 `AIHelperManager`（`manager.go`）维护二级映射 `map[username]map[sessionID]*AIHelper`，通过 RWMutex 实现用户会话隔离。全局单例 `AIModelFactory`（`factory.go`）将模型类型字符串映射到创建函数：

| 模型类型 | 实现 |
|---|---|
| `"1"` | DeepSeek |
| `"2"` | Qwen |
| `"3"` | Qwen + RAG（Redis 向量检索） |
| `"4"` | Qwen + MCP（通过 HTTP MCP 服务器调用天气工具） |

启动时（`main.go`），所有持久化消息从 MySQL 加载并重放到内存中的 `AIHelperManager`。

**异步持久化：** 聊天消息不会同步写入 MySQL。控制器将 `MessageMQParam` JSON 发布到 RabbitMQ；`common/rabbitmq/message.go` 消费队列并通过 `dao/message` 写入 MySQL。

**RAG：** 使用 `cloudwego/eino` 框架。文档通过 Qwen 嵌入模型向量化，并通过 RedisSearch 模块存储在 Redis 中（索引名模式：`rag_docs:<username>:idx`）。需要 `redislabs/redismod` 镜像——标准 Redis 不包含 RediSearch 模块。

**MCP：** `gonexus-mcp` 是独立二进制（入口：`common/mcp/main.go`），通过 HTTP SSE 在 `:8082/mcp` 暴露天气工具。主 `gonexus` 服务在选择模型类型 `"4"` 时作为客户端连接。

**TTS：** 百度云 TTS API（长文本异步）。流程：`CreateTTS` → 返回 `task_id` → 前端轮询 `QueryTTSTask` → 返回 `speech_url`。查询 API 需要字段 `task_ids`（数组）而非 `task_id`。

**配置 (`config/config.toml`)：** 单一 TOML 文件包含所有配置——MySQL、Redis、RabbitMQ、JWT、DeepSeek、Qwen、VoiceService（百度 TTS）。此文件已加入 `.gitignore`，部署时必须手动创建。

### 前端 (Vue 3 + Element Plus)

单页应用位于 `vue-frontend/src/`。主视图是 `views/AIChat.vue`，处理所有聊天功能：会话侧边栏（可折叠/可调整大小）、模型选择、流式 SSE 响应、TTS 播放和 RAG 文件上传。

**关键细节：**
- `src/utils/api.js`：Axios 实例，baseURL 为 `/api`，自动从 `localStorage` 附加 JWT，401 时重定向到 `/login`。
- 开发代理（`vue.config.js`）：`/api` → `http://localhost:8080/api/v1`
- 流式使用原生 `fetch`（非 Axios）从 `/api/v1/ai/chat/send-stream[-new-session]` 读取 SSE `data:` 行。
- 在 textarea 上显式处理 IME 合成（`compositionstart`/`compositionend`），防止移动端中文输入乱序。

### 基础设施

服务及其端口：

| 服务 | 端口 |
|---|---|
| gonexus（主 API） | 8080 |
| gonexus-mcp | 8082 |
| Redis (redismod) | 16379 |
| MySQL | 13306 |
| Nginx | 80 / 443 |

Nginx 挂载：静态 dist 从 `/opt/gonexus/dist`，配置从 `/data/nginx/conf/`。

两个编译后的二进制文件（`gonexus`、`gonexus-mcp`）通过 `.gitignore` 排除在 git 之外。
