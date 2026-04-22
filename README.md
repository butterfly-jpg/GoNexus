# GoNexus 项目开发待办清单 (Todo List)

## 📋 项目概况
- **项目名称**: GoNexus
- **核心功能**: 用户认证 (JWT)、AI 聊天 (流式/多会话/RAG/MCP)、图像识别 (ONNX)
- **技术栈**: Go, MySQL, Redis, RabbitMQ, ONNX Runtime, SSE, JWT

---

## 📺 项目展示
<video src="doc/录屏2026-04-23 00.03.58.mov" controls="controls" width="100%">
  您的浏览器不支持视频播放。
</video>

## 🛠️ 模块一：基础架构与公共组件
- [✅] **项目初始化**
    - [✅] 初始化 Go Module
    - [✅] 设计项目目录结构 (common,config,dao,model,service,util,middleware,controller,router,doc)
    - [✅] 配置管理加载 (.toml)
- [✅] **数据库设计 (MySQL)**
    - [✅] 设计 `User` 表，存储用户注册登录信息
    - [✅] 设计 `Message` 表，存储用户与AI会话消息内容
    - [✅] 设计 `Session` 表，存储用户创建的会话
- [✅] **MySQL 客户端封装**
    - [✅] 连接池配置
- [✅] **Redis 客户端封装**
    - [✅] 连接池配置
    - [✅] 封装验证码存储方法
- [✅] **RabbitMQ 客户端封装**
    - [✅] 定义消息队列结构 (Exchange, Queue, RoutingKey)
    - [✅] 实现生产者 (发送消息)
    - [✅] 实现消费者 (异步写入数据库)
- [✅] **通用工具库**
    - [✅] 11位随机账号或随机验证码生成器
    - [✅] 密码加密/token生成
    - [✅] 全局错误码定义

---

## 🔐 模块二：用户模块 (Auth Service)
- [✅] **注册流程**
    - [✅] 实现发送邮箱验证码接口 (生成验证码 -> 存Redis -> 发邮件)
    - [✅] 实现注册接口
        - [✅] 检验用户是否已存在（通过邮箱）
        - [✅] 校验验证码 (读取Redis并删除)
        - [✅] 生成11位随机账号
        - [✅] 哈希密码并存入数据库
        - [✅] 发送用户名到邮箱
        - [✅] 生成JWT Token并返回
- [✅] **登录流程**
    - [✅] 实现登录接口
        - [✅] 验证账号存在性
        - [✅] 比对密码哈希
        - [✅] 生成 JWT Token
- [✅] **认证中间件**
    - [✅] 获取请求头中 Authorization 字段值或 URL 中的token参数
    - [✅] 验证token是否有效
    - [✅] 将用户名存储到gin上下文中传递给后续业务逻辑
- [✅] **测试**
    - [✅] 测试发送验证码功能
    - [✅] 测试用户注册功能
    - [✅] 测试用户登录功能

---

## 💬 模块三：AI 聊天对话系统 (Chat Service)
- [✅] **会话管理 (Controller层接口)**
    - [✅] 实现通过发送问题获取AI同步回复来创建新会话接口 CreateSessionAndSendMessage
    - [✅] 实现通过发送问题获取AI流式回复来创建新会话接口 CreateStreamSessionAndSendMessage
    - [✅] 实现基于当前会话窗口与AI同步聊天接口 ChatSend
    - [✅] 实现基于当前会话窗口与AI流式聊天接口 ChatStreamSend
    - [✅] 实现获取用户会话列表接口 GetUserSessionByUsername
    - [✅] 实现获取特定会话历史记录接口 ChatHistory
- [✅] **AI架构层**
    - [✅] **AI模型工厂**
        - [✅] 设计存储 模型-创建函数 的数据结构AIModelFactory，根据模型类型获取对应的模型创建方法
        - [✅] 实现CreateAIModel方法
        - [✅] 实现CreateAIHelper方法，根据类型创建AIHelper实例
        - [✅] 实现RegisterModel方法，允许运行时注册新模型类型，动态扩展支持的AI服务
    - [✅] **AIHelper管理器** 
        - [✅] 设计存储 用户-会话-AIHelper 映射关系的数据结构AIHelperManager，以支持多用户会话隔离
        - [✅] 实现GetOrCreateAIHelper方法，获取/创建AIHelper实例
        ~~- [ ] 实现RemoveAIHelper方法，移除指定AIHelper实例~~
        - [✅] 实现GetAIHelper方法，获取现存指定AIHelper实例
        - [✅] 实现GetUserSession方法，获取用户所有的会话ID列表
        - [✅] 实现GetGlobalManager方法，基于单例模式返回AIHelper单例实例，提供AIHelper的全局统一管理入口
- [✅] **同步对话接口**
    - [✅] 接收用户消息
    - [✅] 调用 AI 模型获取完整回复
    - [✅] 构造异步消息发送至 RabbitMQ (存储历史)
    - [✅] 返回完整 JSON 响应
- [✅] **流式对话接口 (SSE)**
    - [✅] SSE配置（Content-Type、Cache-Control、Connection、Access-Control-Allow-Origin、X-Accel-Buffering）
    - [✅] 建立SSE连接
    - [✅] 实时接收 AI 模型流式片段并推送给前端
    - [✅] 发送结束标记 "data: [DONE]\n\n"
- [✅] **异步存储消费者**
    - [✅] 监听消息队列
    - [✅] 解析消息并批量/单条写入（只实现了单条写入）
    ~~- [ ] 错误重试机制~~
- [✅] **测试**
    - [✅] 测试多会话隔离性
    - [✅] 测试流式输出在浏览器的表现

---

## 📄 模块四：增强检索生成（RAG）
- [✅] **文件上传服务**
    - [✅] 实现文件上传接口 UploadingRagFile
    - [✅] 基于 eino 框架和 Redisearcb 模块实现 Embedding 实例初始化
    - [✅] 基于 eino 框架和 Redisearcb 模块实现 Indexer 实例初始化
    - [✅] 基于 eino 框架和 Redisearch 模块实现 Retriever 实例初始化
- [✅] **RAG服务**
    - [✅] 实现文件解析（未实现文本切块,目前只简单处理成一个文档）
    - [✅] 实现向量化处理
    - [✅] 实现向量入库功能（采用redis的redisearch模块实现存储向量）
    - [✅] 实现向量索引和检索筛选功能
    - [✅] 实现RAG模式下的同步对话方法 GenerateResponse
    - [✅] 实现RAG模式下的流式对话方法 StreamResponse
- [✅] **测试**
    - [✅] 测试流式输出在浏览器的表现

## 🔧 模块五：模型上下文协议（MCP）
- [✅] **MCP服务端**
    - [✅] 实现获取MCP服务端实例方法 NewMCPServer 
- [✅] **MCP客户端**
    - [✅] 实现获取MCP客户端实例方法 NewMCPClient
    - [✅] 创建MCP模型实例方法 NewMCPModel
    - [✅] 实现基于MCP的生成响应方法 GenerateResponse
    - [✅] 实现基于MCP的流式响应方法 StreamResponse
- [✅] **测试**
    - [✅] 测试流式输出在浏览器的表现

## 🎙️ 模块六：文本转语音服务（TTS）
- [✅] **TTS服务**
    - [✅] 实现初始化TTS服务实例方法 NewTTSService
    - [✅] 实现文本转语音的业务方法 CreateTTS
    - [✅] 实现获取百度云TTS服务token方法 GetAccessToken
    - [✅] 实现查询TTS任务方法 QueryTTSFull
- [✅] **前端定时轮询TTS后端接口** 
    - [✅] 实现前端定时轮询功能

## 🖼️ 模块七：图像识别系统 (Vision Service)
- [ ] **模型资源准备**
    - [ ] 下载 MobileNetV2 ONNX 模型文件
    - [ ] 准备 ImageNet 标签映射文件 (1000类)
- [✅] **图像处理流水线**
    - [✅] 图片上传接口 (限制格式: JPEG/PNG/GIF, 大小限制)
    - [✅] 图片解码 (Go image package)
    - [✅] 图像预处理：
        - [✅] Resize 至 224x224
        - [✅] 归一化 (Normalize)
        - [✅] 维度变换 (HWC -> CHW)
- [✅] **ONNX 推理引擎**
    - [✅] 集成 `onnxruntime-go`
    - [✅] 加载模型会话 (Session)
    - [✅] 执行推理 (Inference)
    - [✅] 解析输出张量，获取概率最高的类别索引
- [✅] **结果返回**
    - [✅] 根据索引查找标签名称
- [✅] **资源管理机制**
    - [✅] 提供完善的资源清理机制，确保长时间运行下内存安全
    - [✅] 预分配输入输出张量减少内存分配开销
    - [✅] 全局单例管理 ONNX 环境提升初始化效率

---

## 🚀 模块八：部署与优化
- [✅] **容器化**
    - [✅] 编写 `Dockerfile` (多阶段构建减小体积)
    - [✅] 编写 `docker-compose.yml` (编排 Go, Redis, RabbitMQ, MySQL)
- [✅] **文档**
    - [✅] 编写部署指南
- [✅] **功能优化**
    - [✅] 配置模块优化点
         - [✅] 将所有的配置参数(包括数据库和大模型相关配置)都放在config.toml文件里
    - [✅] 前端界面模块优化点 
         - [✅] 优化手机端界面展示效果
         - [✅] 前端发送完的消息内容顺序错乱(确定是IME输入法合成事件竞态问题)
         - [✅] 每次查询历史会话前端都会自动发送两条空白消息bug
         - [✅] 用户创建新会话后,使用用户提的第一个问题为会话title
         - [✅] 用户历史会话按照时间顺序降序排列
         - [✅] 支持用户删除此对话功能(目前实现的删除策略是会话采用软删除,消息采用硬删除,目的是学习gorm框架的软硬删除机制)
         - [✅] 用户点击新建对话后,用户输入问题后会立即刷新界面在左侧一栏显示该会话
         - [✅] 用户点击删除此会话的时候会有弹窗警告：删除后，这条对话记录将无法找回，其中包含的文件也将一并被删除。确定删除此对话？
         - [✅] 支持用户自主切换主题亮度(明亮和漆黑)
         ~~- [ ] TTS 功能做的更显眼一点~~
    - [✅] 注册登录模块优化点
        - [✅] 用户既能支持邮箱登录也能支持账号登录
    - [ ] RAG 模块优化点
        - [ ] 支持上传处理长文档(目前上传的文件内容被简单视作为一个文档,引入文本切分策略,如基于语义的切分或重叠切分,将切分得到的每个chunk单独向量话并存储)
        - [ ] 重排序功能,引入更复杂的模型和算法对初始检索结果进行二次排序,来提高检索准确度
        - [ ] 支持向量化失败数据回滚(基于Lua脚本实现事务回滚,在向量化失败后原子的删除已写入redis的数据,避免数据不一致)
        - [ ] 上传文件和向量化过程由同步改为异步处理,支持高并发和大文件
        - [ ] 引入上传文件大小限制在10MB以内,避免用户上传过大文件导致系统资源耗尽
---