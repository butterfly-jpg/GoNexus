package aihelper

import (
	"context"
	"fmt"
	"sync"
)

// ModelCreator 定义模型创建函数
type ModelCreator func(ctx context.Context, config map[string]interface{}) (AIModel, error)

// AIModelFactory AI模型工厂
type AIModelFactory struct {
	creators map[string]ModelCreator
	mu       sync.RWMutex
}

var (
	globalFactory *AIModelFactory
	factoryOnce   sync.Once
)

// GetGlobalFactory 获取全局AIHelper实例
func GetGlobalFactory() *AIModelFactory {
	factoryOnce.Do(func() {
		globalFactory = &AIModelFactory{
			creators: make(map[string]ModelCreator),
			mu:       sync.RWMutex{},
		}
		// 注册AI模型提升扩展能力
		globalFactory.registerCreators()
	})
	return globalFactory
}

// registerCreators 注册模型方法
func (f *AIModelFactory) registerCreators() {
	// 注册DeepSeek大模型
	f.creators["1"] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewDeepSeekModel(ctx)
	}
	// 注册Qwen大模型
	f.creators["2"] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		return NewQwenModel(ctx)
	}
	// 注册QwenRag模型
	f.creators["3"] = func(ctx context.Context, config map[string]interface{}) (AIModel, error) {
		username, ok := config["username"].(string)
		if !ok {
			return nil, fmt.Errorf("qwen RAG model requires username")
		}
		return NewQwenRAGModel(ctx, username)
	}
}

// CreateAIHelper 创建AIHelper方法
func (f *AIModelFactory) CreateAIHelper(ctx context.Context, modelType, sessionID string, config map[string]interface{}) (*AIHelper, error) {
	model, err := f.CreateAIModel(ctx, modelType, config)
	if err != nil {
		return nil, err
	}
	return NewAIHelper(model, sessionID), nil
}

// CreateAIModel 根据类型创建AI模型
func (f *AIModelFactory) CreateAIModel(ctx context.Context, modelType string, config map[string]interface{}) (AIModel, error) {
	f.mu.RLock()
	modelCreator, exists := f.creators[modelType]
	f.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("unsupported model type %s", modelType)
	}
	return modelCreator(ctx, config)
}

// RegisterModel 可扩展注册
func (f *AIModelFactory) RegisterModel(modelType string, creator ModelCreator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.creators[modelType] = creator
}
