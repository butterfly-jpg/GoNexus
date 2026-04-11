<template>
  <div class="menu-container">
    <!-- 背景光球 -->
    <div class="bg-orb orb1"></div>
    <div class="bg-orb orb2"></div>
    <div class="bg-orb orb3"></div>

    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="brand">
        <span class="brand-icon">✦</span>
        <span class="brand-name">GoNexus</span>
      </div>
      <button class="logout-btn" @click="handleLogout">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
          <polyline points="16 17 21 12 16 7"/>
          <line x1="21" y1="12" x2="9" y2="12"/>
        </svg>
        退出登录
      </button>
    </header>

    <!-- 主内容 -->
    <main class="main">
      <div class="welcome-section">
        <div class="welcome-glow">✦</div>
        <h2 class="welcome-title">选择你的 AI 功能</h2>
        <p class="welcome-sub">基于前沿大模型，为你提供智能化的应用体验</p>
      </div>

      <div class="menu-grid">
        <!-- AI 聊天 -->
        <div class="menu-card" @click="$router.push('/ai-chat')">
          <div class="card-glow card-glow-purple"></div>
          <div class="card-inner">
            <div class="card-icon icon-purple">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </div>
            <div class="card-text">
              <h3>AI 聊天</h3>
              <p>多模型智能对话，支持 RAG 文件检索、MCP 工具与 TTS 文本转语音</p>
            </div>
            <div class="card-arrow">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="5" y1="12" x2="19" y2="12"/>
                <polyline points="12 5 19 12 12 19"/>
              </svg>
            </div>
          </div>
          <div class="card-tag">DeepSeek · Qwen · RAG · MCP</div>
        </div>

        <!-- 图像识别 -->
        <div class="menu-card" @click="$router.push('/image-recognition')">
          <div class="card-glow card-glow-cyan"></div>
          <div class="card-inner">
            <div class="card-icon icon-cyan">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <circle cx="8.5" cy="8.5" r="1.5"/>
                <polyline points="21 15 16 10 5 21"/>
              </svg>
            </div>
            <div class="card-text">
              <h3>图像识别</h3>
              <p>上传图片，AI 自动分析内容、识别物体与场景</p>
            </div>
            <div class="card-arrow">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="5" y1="12" x2="19" y2="12"/>
                <polyline points="12 5 19 12 12 19"/>
              </svg>
            </div>
          </div>
          <div class="card-tag">Vision · 多模态</div>
        </div>
      </div>

      <!-- 底部说明 -->
      <p class="footer-tip">点击卡片进入对应功能</p>
    </main>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'MenuView',
  setup() {
    const router = useRouter()

    const handleLogout = async () => {
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        localStorage.removeItem('token')
        ElMessage.success('退出登录成功')
        router.push('/login')
      } catch {
        // 用户取消
      }
    }

    return { handleLogout }
  }
}
</script>

<style scoped>
/* ===== 基础布局 ===== */
.menu-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #0d0d1a;
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: #e2e2f0;
}

/* ===== 背景光球 ===== */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.16;
  pointer-events: none;
}
.orb1 {
  width: 600px; height: 600px;
  background: radial-gradient(circle, #7c3aed, transparent);
  top: -180px; left: -120px;
  animation: orbFloat1 20s ease-in-out infinite;
}
.orb2 {
  width: 450px; height: 450px;
  background: radial-gradient(circle, #06b6d4, transparent);
  bottom: -80px; right: -60px;
  animation: orbFloat2 25s ease-in-out infinite;
}
.orb3 {
  width: 320px; height: 320px;
  background: radial-gradient(circle, #a855f7, transparent);
  top: 50%; left: 55%;
  animation: orbFloat3 17s ease-in-out infinite;
}
@keyframes orbFloat1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(80px, 50px) scale(1.08); }
  66% { transform: translate(-30px, 90px) scale(0.92); }
}
@keyframes orbFloat2 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-90px, -70px) scale(1.12); }
}
@keyframes orbFloat3 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(50px, -100px); }
}

/* ===== 顶部导航 ===== */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 40px;
  background: rgba(13, 13, 26, 0.8);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(124, 58, 237, 0.15);
  position: relative;
  z-index: 10;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.brand-icon {
  font-size: 22px;
  background: linear-gradient(135deg, #7c3aed, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: iconPulse 3s ease-in-out infinite;
}
@keyframes iconPulse {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.5); }
}
.brand-name {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #c4b5fd, #e0f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 18px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: #94a3b8;
  font-size: 13.5px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.25s ease;
}
.logout-btn:hover {
  background: rgba(239, 68, 68, 0.12);
  border-color: rgba(239, 68, 68, 0.3);
  color: #fca5a5;
  box-shadow: 0 0 16px rgba(239, 68, 68, 0.1);
}

/* ===== 主内容 ===== */
.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  position: relative;
  z-index: 5;
}

/* 欢迎区 */
.welcome-section {
  text-align: center;
  margin-bottom: 56px;
  animation: fadeUp 0.7s ease-out both;
}
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

.welcome-glow {
  font-size: 52px;
  background: linear-gradient(135deg, #7c3aed, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 20px;
  display: inline-block;
  animation: glowPulse 3s ease-in-out infinite;
}
@keyframes glowPulse {
  0%, 100% { filter: brightness(1) drop-shadow(0 0 10px rgba(124, 58, 237, 0.4)); }
  50% { filter: brightness(1.4) drop-shadow(0 0 28px rgba(6, 182, 212, 0.5)); }
}

.welcome-title {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 12px;
  background: linear-gradient(135deg, #f8fafc 0%, #c4b5fd 60%, #67e8f9 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.3px;
}

.welcome-sub {
  font-size: 15px;
  color: #4b5563;
  margin: 0;
}

/* ===== 卡片网格 ===== */
.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 28px;
  max-width: 800px;
  width: 100%;
  animation: fadeUp 0.7s ease-out 0.15s both;
}

/* 卡片 */
.menu-card {
  position: relative;
  background: rgba(20, 17, 38, 0.85);
  border: 1px solid rgba(124, 58, 237, 0.18);
  border-radius: 20px;
  padding: 28px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(12px);
}

/* 悬停时的光晕 */
.card-glow {
  position: absolute;
  width: 200px; height: 200px;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0;
  transition: opacity 0.4s ease;
  pointer-events: none;
}
.card-glow-purple { background: #7c3aed; top: -60px; right: -40px; }
.card-glow-cyan   { background: #06b6d4; top: -60px; right: -40px; }

.menu-card:hover .card-glow { opacity: 0.25; }
.menu-card:hover {
  transform: translateY(-6px) scale(1.01);
  border-color: rgba(124, 58, 237, 0.45);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(124, 58, 237, 0.2);
}

/* 卡片内部布局 */
.card-inner {
  display: flex;
  align-items: center;
  gap: 18px;
  position: relative;
  z-index: 1;
}

/* 图标 */
.card-icon {
  width: 64px; height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.35s ease;
}
.icon-purple {
  background: rgba(124, 58, 237, 0.18);
  border: 1px solid rgba(124, 58, 237, 0.3);
  color: #a78bfa;
}
.icon-cyan {
  background: rgba(6, 182, 212, 0.12);
  border: 1px solid rgba(6, 182, 212, 0.25);
  color: #67e8f9;
}
.menu-card:hover .icon-purple {
  background: rgba(124, 58, 237, 0.3);
  box-shadow: 0 0 20px rgba(124, 58, 237, 0.35);
  color: #c4b5fd;
}
.menu-card:hover .icon-cyan {
  background: rgba(6, 182, 212, 0.22);
  box-shadow: 0 0 20px rgba(6, 182, 212, 0.3);
  color: #a5f3fc;
}

/* 文字 */
.card-text {
  flex: 1;
  min-width: 0;
}
.card-text h3 {
  font-size: 18px;
  font-weight: 700;
  color: #e2e8f0;
  margin: 0 0 6px;
  transition: color 0.25s;
}
.menu-card:hover .card-text h3 { color: #f8fafc; }
.card-text p {
  font-size: 13px;
  color: #4b5563;
  margin: 0;
  line-height: 1.6;
  transition: color 0.25s;
}
.menu-card:hover .card-text p { color: #6b7280; }

/* 箭头 */
.card-arrow {
  color: #374151;
  flex-shrink: 0;
  transform: translateX(0);
  transition: all 0.3s ease;
}
.menu-card:hover .card-arrow {
  color: #a78bfa;
  transform: translateX(4px);
}

/* 标签 */
.card-tag {
  margin-top: 18px;
  font-size: 12px;
  color: #374151;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 6px 12px;
  display: inline-block;
  position: relative;
  z-index: 1;
  letter-spacing: 0.4px;
  transition: all 0.25s;
}
.menu-card:hover .card-tag {
  color: #6b7280;
  border-color: rgba(255, 255, 255, 0.1);
}

/* 底部提示 */
.footer-tip {
  margin-top: 44px;
  font-size: 13px;
  color: rgba(75, 85, 99, 0.7);
  animation: fadeUp 0.7s ease-out 0.3s both;
}
</style>
