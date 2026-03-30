// Package rag 提供 rag 向量索引相关方法
package rag

import (
	"GoNexus/config"
	"context"
	"fmt"
	"os"

	redispkg "GoNexus/common/redis"

	"github.com/cloudwego/eino/components/embedding"

	embeddingark "github.com/cloudwego/eino-ext/components/embedding/ark"
)

// RAGIndexer RAG索引器
type RAGIndexer struct {
	embedding embedding.Embedder
}

// DeleteIndex 删除指定文件的知识库索引
func DeleteIndex(ctx context.Context, filename string) error {
	if err := redispkg.DeleteRedisIndex(ctx, filename); err != nil {
		return fmt.Errorf("delete redis index failed. err: %v", err)
	}
	return nil
}

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

	// 3.2 创建索引器实例
}
