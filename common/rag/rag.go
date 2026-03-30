// Package rag 提供 rag 向量索引相关方法
package rag

import (
	"GoNexus/config"
	"context"
	"fmt"
	"os"

	redispkg "GoNexus/common/redis"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"

	embeddingark "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisindexer "github.com/cloudwego/eino-ext/components/indexer/redis"
)

// RAGIndexer RAG索引器
type RAGIndexer struct {
	embedding embedding.Embedder
	indexer   *redisindexer.Indexer
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
