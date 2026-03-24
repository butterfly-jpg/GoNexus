package aihelper

import (
	"context"
	"sync"
)

// AIHelperManager AIHelper管理器数据结构，管理用户-会话-AIHelper的映射关系
type AIHelperManager struct {
	helpers map[string]map[string]*AIHelper // map[用户账号]map[会话ID]*AIHelper
	mu      sync.RWMutex
}

// NewAIHelperManager 创建新的AIHelper管理实例
func NewAIHelperManager() *AIHelperManager {
	return &AIHelperManager{
		helpers: make(map[string]map[string]*AIHelper),
		mu:      sync.RWMutex{},
	}
}

var (
	// globalManager 全局AIHelper单例实例
	globalManager *AIHelperManager
	once          sync.Once
)

// GetGlobalManager 获取全局AIHelper单例实例
func GetGlobalManager() *AIHelperManager {
	once.Do(func() {
		globalManager = NewAIHelperManager()
	})
	return globalManager
}

// GetOrCreateAIHelper 获取或创建AIHelper实例方法
func (m *AIHelperManager) GetOrCreateAIHelper(username, sessionID, modelType string, config map[string]interface{}) (*AIHelper, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 1. 获取用户的会话映射
	sessionHelper, exists := m.helpers[username]
	if !exists {
		sessionHelper = make(map[string]*AIHelper)
		m.helpers[username] = sessionHelper
	}
	// 2. 检查用户会话是否存在
	helper, exists := sessionHelper[sessionID]
	if exists {
		return helper, nil
	}
	// 3. 基于工厂模式创建新AIHelper实例
	factory := GetGlobalFactory()
	helper, err := factory.CreateAIHelper(context.Background(), modelType, sessionID, config)
	if err != nil {
		return nil, err
	}
	sessionHelper[sessionID] = helper
	return helper, nil
}
