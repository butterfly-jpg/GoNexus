// Package rag 提供 rag 向量索引相关方法
package rag

import (
	"GoNexus/config"
	"context"
	"fmt"
	"os"

	redispkg "GoNexus/common/redis"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	embeddingark "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisindexer "github.com/cloudwego/eino-ext/components/indexer/redis"
	redisretriever "github.com/cloudwego/eino-ext/components/retriever/redis"
)

// RAGIndexer RAG索引器
type RAGIndexer struct {
	embedding embedding.Embedder
	indexer   *redisindexer.Indexer
}

// RAGQuery RAG查询器
type RAGQuery struct {
	embedding embedding.Embedder
	retriever retriever.Retriever
}

// DeleteIndex 删除指定文件的知识库索引
func DeleteIndex(ctx context.Context, filename string) error {
	if err := redispkg.DeleteRedisIndex(ctx, filename); err != nil {
		return fmt.Errorf("delete redis index failed. err: %v", err)
	}
	return nil
}

// NewRAGIndexer 获取RAG的向量生成器Embedding和索引器Indexer的实例
func NewRAGIndexer(filename, embeddingModel string) (*RAGIndexer, error) {
	ctx := context.Background()
	// 1. 创建Embedding组件,将文本文档转换为向量表示,使用ARK平台的模型生成向量
	// 1.1 设置向量生成器的配置
	embedConfig := &embeddingark.EmbeddingConfig{
		BaseURL: config.GetConfig().RagBaseUrl,
		APIKey:  os.Getenv("QWEN_API_KEY"), // 通义千问的API
		Model:   embeddingModel,
	}
	// 1.2 创建向量生成器实例
	embedder, err := embeddingark.NewEmbedder(ctx, embedConfig)
	if err != nil {
		return nil, fmt.Errorf("create embedder failed. err: %v", err)
	}
	// 2. 初始化 Redis 的向量索引结构,使用 Redisearch 模块
	if err = redispkg.InitRedisIndex(ctx, filename, config.GetConfig().RagDimension); err != nil {
		return nil, fmt.Errorf("init redisearch index failed. err: %v", err)
	}
	// 3. 创建Indexer组件,将文档及其向量表示存储到后端存储系统中,构建向量数据库以提供高效的检索能力
	// 3.1 设置索引器的配置
	indexConfig := &redisindexer.IndexerConfig{
		Client:    redispkg.Rdb,
		KeyPrefix: redispkg.GenerateIndexNamePrefix(filename), // 不同知识库对应不同的前缀,避免冲突
		BatchSize: 10,
		Embedding: embedder, // 将向量生成器交给索引器,在索引器写入文本时,可以自动完成向量计算
		// 自定义文档映射
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redisindexer.Hashes, error) {
			// 从文档的元数据中提取源信息（例如文件名、URL）
			source := ""
			if s, ok := doc.MetaData["source"].(string); ok {
				source = s
			}
			// 文档如何映射到Redis hashes
			return &redisindexer.Hashes{
				// redis key 由 "知识库名+文档块ID" 组成
				Key: fmt.Sprintf("%s%s", filename, doc.ID),
				// redis hash 中的字段
				Field2Value: map[string]redisindexer.FieldValue{
					"content": {
						Value:    doc.Content,
						EmbedKey: "vector", // 将 content 向量化然后保存到名为"vector"的字段中
					},
					// 一些辅助的元数据信息,不参与向量计算
					"metadata": {
						Value: source,
					},
				},
			}, nil
		},
	}
	// 3.2 创建索引器实例
	idx, err := redisindexer.NewIndexer(ctx, indexConfig)
	if err != nil {
		return nil, fmt.Errorf("create indexer failed. err: %v", err)
	}
	// 4. 返回封装好的RAGIndexer,包含向量生成器和索引器
	return &RAGIndexer{
		embedding: embedder,
		indexer:   idx,
	}, nil
}

// IndexFile 读取文件内容并创建向量索引
func (r *RAGIndexer) IndexFile(ctx context.Context, filePath string) error {
	// 1. 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file failed. err: %v", err)
	}
	// 2. 将文件内容转换成文档
	docs := []*schema.Document{
		{
			ID:      uuid.NewString(),
			Content: string(content),
			MetaData: map[string]any{
				"source": filePath,
			},
		},
	}
	// 3. 索引文档
	_, err = r.indexer.Store(ctx, docs)
	if err != nil {
		return fmt.Errorf("store documents failed. err: %v", err)
	}
	return nil
}

// NewRAGQuery 创建RAG查询器,用于向量检索和问答
func NewRAGQuery(ctx context.Context, username string) (*RAGQuery, error) {
	// 1. 创建embedding组件
	embedConfig := &embeddingark.EmbeddingConfig{
		BaseURL: config.GetConfig().RagBaseUrl,
		APIKey:  os.Getenv("QWEN_API_KEY"), // 通义千问的API
		Model:   config.GetConfig().RagEmbeddingModel,
	}
	embedder, err := embeddingark.NewEmbedder(ctx, embedConfig)
	if err != nil {
		return nil, fmt.Errorf("create embedder failed. err: %v", err)
	}
	// 2. 获取用户上传的文件（目前只支持用户上传一个文件）
	userDir := fmt.Sprintf("uploads/%s", username)
	files, err := os.ReadDir(userDir)
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no uploaded file found for user %s", username)
	}
	var filename string
	for _, f := range files {
		if !f.IsDir() {
			filename = f.Name()
			break
		}
	}
	if filename == "" {
		return nil, fmt.Errorf("no valid file found for user %s", username)
	}
	// 3. 创建retriever组件
	// 3.1 设置召回器的配置
	// 生成索引
	indexName := redispkg.GenerateIndexName(filename)
	retrieverConfig := &redisretriever.RetrieverConfig{
		Client:       redispkg.Rdb,
		Index:        indexName,
		Dialect:      2,
		ReturnFields: []string{"content", "metadata", "distance"},
		TopK:         5,
		VectorField:  "vector",
		Embedding:    embedder,
		DocumentConverter: func(ctx context.Context, doc redis.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       doc.ID,
				Content:  "",
				MetaData: map[string]any{},
			}
			for field, value := range doc.Fields {
				if field == "content" {
					resp.Content = value
				} else {
					resp.MetaData[field] = value
				}
			}
			return resp, nil
		},
	}
	// 3.2 创建召回器实例
	rtr, err := redisretriever.NewRetriever(ctx, retrieverConfig)
	if err != nil {
		return nil, fmt.Errorf("create retriever failed. err: %v", err)
	}
	// 4. 返回封装好的RAGQuery,包含向量生成器和召回器
	return &RAGQuery{
		embedding: embedder,
		retriever: rtr,
	}, nil
}

// RetrieveDocuments 检索相关文档
func (r *RAGQuery) RetrieveDocuments(ctx context.Context, query string) ([]*schema.Document, error) {
	docs, err := r.retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("retrieve documents failed. err: %v", err)
	}
	return docs, nil
}

// BuildRAGPrompt 构建包含检索文档的提示词
func BuildRAGPrompt(query string, docs []*schema.Document) string {
	if len(docs) == 0 {
		return query
	}
	contextText := ""
	for i, doc := range docs {
		contextText += fmt.Sprintf("[文档 %d]: %s\n\n", i+1, doc.Content)
	}
	prompt := fmt.Sprintf(`基于以下参考文档回答用户的问题。如果文档中没有相关信息，请说明无法找到相关信息。

参考文档：
%s

用户问题: %s

请提供准确、完整的回答: `, contextText, query)

	return prompt
}
