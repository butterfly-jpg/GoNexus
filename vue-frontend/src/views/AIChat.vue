<template>
  <div class="ai-chat-container">
    <!-- 浮动光球背景 -->
    <div class="bg-orb orb1"></div>
    <div class="bg-orb orb2"></div>
    <div class="bg-orb orb3"></div>

    <!-- 左侧会话列表 -->
    <div class="session-list">
      <div class="session-list-header">
        <div class="brand">
          <span class="brand-icon">✦</span>
          <span class="brand-name">GoNexus</span>
        </div>
        <button class="new-chat-btn" @click="createNewSession">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          新对话
        </button>
      </div>
      <div class="session-list-scroll">
        <ul class="session-list-ul">
          <li
            v-for="session in sessions"
            :key="session.id"
            :class="['session-item', { active: currentSessionId === session.id }]"
            @click="switchSession(session.id)"
          >
            <span class="session-dot"></span>
            <span class="session-title">{{ session.name || `会话 ${session.id}` }}</span>
          </li>
        </ul>
        <div v-if="sessions.length === 0" class="sidebar-empty">
          <p>暂无会话</p>
          <p>点击"新对话"开始</p>
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-section">
      <!-- 顶部工具栏 -->
      <div class="top-bar">
        <button class="icon-btn" @click="$router.push('/menu')" title="返回菜单">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
        </button>

        <div class="divider-v"></div>

        <div class="model-selector">
          <span class="label-tag">模型</span>
          <select v-model="selectedModel" class="model-select">
            <option value="1">DeepSeek</option>
            <option value="2">Qwen</option>
            <option value="3">Qwen RAG</option>
            <option value="4">Qwen MCP</option>
          </select>
        </div>

        <div class="stream-control">
          <span class="label-tag">流式响应</span>
          <label class="toggle-switch">
            <input type="checkbox" v-model="isStreaming" />
            <span class="toggle-track">
              <span class="toggle-thumb"></span>
            </span>
          </label>
        </div>

        <div class="toolbar-right">
          <button
            class="icon-btn"
            @click="syncHistory"
            :disabled="!currentSessionId || tempSession"
            title="同步历史数据"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M23 4v6h-6"/><path d="M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
          </button>
          <button
            class="icon-btn upload-icon-btn"
            @click="triggerFileUpload"
            :disabled="uploading"
            title="上传文档 (.md / .txt)"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
            <span v-if="uploading" class="btn-spinner"></span>
          </button>
        </div>

        <input
          ref="fileInput"
          type="file"
          accept=".md,.txt,text/markdown,text/plain"
          style="display: none"
          @change="handleFileUpload"
        />
      </div>

      <!-- 消息区域 -->
      <div class="chat-messages" ref="messagesRef">
        <div v-if="currentMessages.length === 0" class="empty-state">
          <div class="empty-glow">✦</div>
          <p class="empty-title">开始一段新的对话</p>
          <p class="empty-sub">选择模型，输入你的问题</p>
        </div>

        <div
          v-for="(message, index) in currentMessages"
          :key="index"
          :class="['msg-row', message.role === 'user' ? 'msg-row-user' : 'msg-row-ai']"
        >
          <div :class="['avatar', message.role === 'user' ? 'avatar-user' : 'avatar-ai']">
            {{ message.role === 'user' ? 'U' : 'AI' }}
          </div>
          <div :class="['bubble', message.role === 'user' ? 'bubble-user' : 'bubble-ai']">
            <div class="bubble-actions" v-if="message.role === 'assistant'">
              <button class="tts-btn" @click="playTTS(message.content)" title="语音播放">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z"/></svg>
              </button>
            </div>
            <div class="bubble-content" v-html="renderMarkdown(message.content)"></div>
            <div v-if="message.meta && message.meta.status === 'streaming'" class="typing-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区 -->
      <div class="chat-input-area">
        <div class="input-wrapper">
          <textarea
            v-model="inputMessage"
            placeholder="输入消息...  Enter 发送 / Shift+Enter 换行"
            @keydown.enter.exact.prevent="sendMessage"
            :disabled="loading"
            ref="messageInput"
            rows="1"
          ></textarea>
          <button
            type="button"
            :disabled="!inputMessage.trim() || loading"
            @click="sendMessage"
            class="send-btn"
          >
            <span v-if="!loading">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/></svg>
            </span>
            <span v-else class="loading-ring"></span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../utils/api'

export default {
  name: 'AIChat',
  setup() {
    const sessions = ref({})
    const currentSessionId = ref(null)
    const tempSession = ref(false)
    const currentMessages = ref([])
    const inputMessage = ref('')
    const loading = ref(false)
    const messagesRef = ref(null)
    const messageInput = ref(null)
    const selectedModel = ref('1')
    const isStreaming = ref(true)
    const uploading = ref(false)
    const fileInput = ref(null)

    const renderMarkdown = (text) => {
      if (!text && text !== '') return ''
      return String(text)
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/`(.*?)`/g, '<code>$1</code>')
        .replace(/\n/g, '<br>')
    }

    const playTTS = async (text) => {
      try {
        const createResponse = await api.post('/ai/chat/tts', { text })
        if (createResponse.data && createResponse.data.status_code === 1000 && createResponse.data.task_id) {
          const taskId = createResponse.data.task_id
          await new Promise(resolve => setTimeout(resolve, 5000))
          const maxAttempts = 30
          const pollInterval = 2000
          let attempts = 0
          const pollResult = async () => {
            const queryResponse = await api.get('/ai/chat/tts/query', { params: { task_id: taskId } })
            if (queryResponse.data && queryResponse.data.status_code === 1000) {
              const taskStatus = queryResponse.data.task_status
              if (taskStatus === 'Success' && queryResponse.data.task_result) {
                const audio = new Audio(queryResponse.data.task_result)
                audio.play()
                return true
              } else if (taskStatus === 'Running' || taskStatus === 'Created') {
                attempts++
                if (attempts < maxAttempts) {
                  await new Promise(resolve => setTimeout(resolve, pollInterval))
                  return await pollResult()
                } else {
                  ElMessage.error('语音合成超时')
                  return true
                }
              } else {
                ElMessage.error('语音合成失败')
                return true
              }
            }
            attempts++
            if (attempts < maxAttempts) {
              await new Promise(resolve => setTimeout(resolve, pollInterval))
              return await pollResult()
            } else {
              ElMessage.error('语音合成超时')
              return true
            }
          }
          await pollResult()
        } else {
          ElMessage.error('无法创建语音合成任务')
        }
      } catch (error) {
        console.error('TTS error:', error)
        ElMessage.error('请求语音接口失败')
      }
    }

    const loadSessions = async () => {
      try {
        const response = await api.get('/ai/chat/sessions')
        if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.sessions)) {
          const sessionMap = {}
          response.data.sessions.forEach(s => {
            const sid = String(s.sessionID)
            sessionMap[sid] = {
              id: sid,
              name: s.title || `会话 ${sid}`,
              messages: []
            }
          })
          sessions.value = sessionMap
        }
      } catch (error) {
        console.error('Load sessions error:', error)
      }
    }

    const createNewSession = () => {
      currentSessionId.value = 'temp'
      tempSession.value = true
      currentMessages.value = []
      nextTick(() => {
        if (messageInput.value) messageInput.value.focus()
      })
    }

    const switchSession = async (sessionId) => {
      if (!sessionId) return
      currentSessionId.value = String(sessionId)
      tempSession.value = false
      if (!sessions.value[sessionId].messages || sessions.value[sessionId].messages.length === 0) {
        try {
          const response = await api.post('/ai/chat/history', { sessionId: currentSessionId.value })
          if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.history)) {
            const messages = response.data.history.map(item => ({
              role: item.is_user ? 'user' : 'assistant',
              content: item.content
            }))
            sessions.value[sessionId].messages = messages
          }
        } catch (err) {
          console.error('Load history error:', err)
        }
      }
      currentMessages.value = [...(sessions.value[sessionId].messages || [])]
      await nextTick()
      scrollToBottom()
    }

    const syncHistory = async () => {
      if (!currentSessionId.value || tempSession.value) {
        ElMessage.warning('请选择已有会话进行同步')
        return
      }
      try {
        const response = await api.post('/ai/chat/history', { sessionId: currentSessionId.value })
        if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.history)) {
          const messages = response.data.history.map(item => ({
            role: item.is_user ? 'user' : 'assistant',
            content: item.content
          }))
          sessions.value[currentSessionId.value].messages = messages
          currentMessages.value = [...messages]
          await nextTick()
          scrollToBottom()
        } else {
          ElMessage.error('无法获取历史数据')
        }
      } catch (err) {
        console.error('Sync history error:', err)
        ElMessage.error('请求历史数据失败')
      }
    }

    const sendMessage = async () => {
      if (!inputMessage.value || !inputMessage.value.trim()) {
        ElMessage.warning('请输入消息内容')
        return
      }
      const userMessage = { role: 'user', content: inputMessage.value }
      const currentInput = inputMessage.value
      inputMessage.value = ''
      currentMessages.value.push(userMessage)
      await nextTick()
      scrollToBottom()
      try {
        loading.value = true
        if (isStreaming.value) {
          await handleStreaming(currentInput)
        } else {
          await handleNormal(currentInput)
        }
      } catch (err) {
        console.error('Send message error:', err)
        ElMessage.error('发送失败，请重试')
        if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value] && sessions.value[currentSessionId.value].messages) {
          const sessionArr = sessions.value[currentSessionId.value].messages
          if (sessionArr && sessionArr.length) sessionArr.pop()
        }
        currentMessages.value.pop()
      } finally {
        if (!isStreaming.value) {
          loading.value = false
        }
        await nextTick()
        scrollToBottom()
      }
    }

    async function handleStreaming(question) {
      const aiMessage = { role: 'assistant', content: '', meta: { status: 'streaming' } }
      const aiMessageIndex = currentMessages.value.length
      currentMessages.value.push(aiMessage)
      if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
        if (!sessions.value[currentSessionId.value].messages) sessions.value[currentSessionId.value].messages = []
        sessions.value[currentSessionId.value].messages.push({ role: 'assistant', content: '' })
      }
      const url = tempSession.value
        ? '/api/ai/chat/send-stream-new-session'
        : '/api/ai/chat/send-stream'
      const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      }
      const body = tempSession.value
        ? { question: question, modelType: selectedModel.value }
        : { question: question, modelType: selectedModel.value, sessionId: currentSessionId.value }
      try {
        const response = await fetch(url, { method: 'POST', headers, body: JSON.stringify(body) })
        if (!response.ok) {
          loading.value = false
          throw new Error('Network response was not ok')
        }
        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        let buffer = ''
        // eslint-disable-next-line no-constant-condition
        while (true) {
          const { done, value } = await reader.read()
          if (done) break
          const chunk = decoder.decode(value, { stream: true })
          buffer += chunk
          const lines = buffer.split('\n')
          buffer = lines.pop() || ''
          for (const line of lines) {
            const trimmedLine = line.trim()
            if (!trimmedLine) continue
            if (trimmedLine.startsWith('data:')) {
              const data = trimmedLine.slice(5).trim()
              if (data === '[DONE]') {
                loading.value = false
                currentMessages.value[aiMessageIndex].meta = { status: 'done' }
                currentMessages.value = [...currentMessages.value]
              } else if (data.startsWith('{')) {
                try {
                  const parsed = JSON.parse(data)
                  if (parsed.sessionId) {
                    const newSid = String(parsed.sessionId)
                    if (tempSession.value) {
                      sessions.value[newSid] = {
                        id: newSid,
                        name: '新会话',
                        messages: [...currentMessages.value]
                      }
                      currentSessionId.value = newSid
                      tempSession.value = false
                    }
                  }
                } catch (e) {
                  currentMessages.value[aiMessageIndex].content += data
                }
              } else {
                currentMessages.value[aiMessageIndex].content += data
              }
              currentMessages.value = [...currentMessages.value]
              await new Promise(resolve => {
                requestAnimationFrame(() => { scrollToBottom(); resolve() })
              })
            }
          }
        }
        loading.value = false
        currentMessages.value[aiMessageIndex].meta = { status: 'done' }
        currentMessages.value = [...currentMessages.value]
        if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
          const sessMsgs = sessions.value[currentSessionId.value].messages
          if (Array.isArray(sessMsgs) && sessMsgs.length) {
            const lastIndex = sessMsgs.length - 1
            if (sessMsgs[lastIndex] && sessMsgs[lastIndex].role === 'assistant') {
              sessMsgs[lastIndex].content = currentMessages.value[aiMessageIndex].content
            }
          }
        }
      } catch (err) {
        console.error('Stream error:', err)
        loading.value = false
        currentMessages.value[aiMessageIndex].meta = { status: 'error' }
        currentMessages.value = [...currentMessages.value]
        ElMessage.error('流式传输出错')
      }
    }

    async function handleNormal(question) {
      if (tempSession.value) {
        const response = await api.post('/ai/chat/send-new-session', {
          question: question,
          modelType: selectedModel.value
        })
        if (response.data && response.data.status_code === 1000) {
          const sessionId = String(response.data.sessionID)
          const aiMessage = { role: 'assistant', content: response.data.information || '' }
          sessions.value[sessionId] = {
            id: sessionId,
            name: '新会话',
            messages: [{ role: 'user', content: question }, aiMessage]
          }
          currentSessionId.value = sessionId
          tempSession.value = false
          currentMessages.value = [...sessions.value[sessionId].messages]
        } else {
          ElMessage.error(response.data?.status_msg || '发送失败')
          currentMessages.value.pop()
        }
      } else {
        const sessionMsgs = sessions.value[currentSessionId.value].messages
        sessionMsgs.push({ role: 'user', content: question })
        const response = await api.post('/ai/chat/send', {
          question: question,
          modelType: selectedModel.value,
          sessionId: currentSessionId.value
        })
        if (response.data && response.data.status_code === 1000) {
          const aiMessage = { role: 'assistant', content: response.data.information || '' }
          sessionMsgs.push(aiMessage)
          currentMessages.value = [...sessionMsgs]
        } else {
          ElMessage.error(response.data?.status_msg || '发送失败')
          sessionMsgs.pop()
          currentMessages.value.pop()
        }
      }
    }

    const scrollToBottom = () => {
      if (messagesRef.value) {
        try {
          messagesRef.value.scrollTop = messagesRef.value.scrollHeight
        } catch (e) { /* ignore */ }
      }
    }

    const triggerFileUpload = () => {
      if (fileInput.value) fileInput.value.click()
    }

    const handleFileUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) return
      const fileName = file.name.toLowerCase()
      if (!fileName.endsWith('.md') && !fileName.endsWith('.txt')) {
        ElMessage.error('只允许上传 .md 或 .txt 文件')
        if (fileInput.value) fileInput.value.value = ''
        return
      }
      try {
        uploading.value = true
        const formData = new FormData()
        formData.append('file', file)
        const response = await api.post('/file/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        if (response.data && response.data.status_code === 1000) {
          ElMessage.success('文件上传成功')
        } else {
          ElMessage.error(response.data?.status_msg || '上传失败')
        }
      } catch (error) {
        console.error('File upload error:', error)
        ElMessage.error('文件上传失败')
      } finally {
        uploading.value = false
        if (fileInput.value) fileInput.value.value = ''
      }
    }

    onMounted(() => {
      loadSessions()
    })

    return {
      sessions: computed(() => Object.values(sessions.value)),
      currentSessionId,
      tempSession,
      currentMessages,
      inputMessage,
      loading,
      messagesRef,
      messageInput,
      selectedModel,
      isStreaming,
      uploading,
      fileInput,
      renderMarkdown,
      playTTS,
      createNewSession,
      switchSession,
      syncHistory,
      sendMessage,
      triggerFileUpload,
      handleFileUpload
    }
  }
}
</script>

<style scoped>
/* ===== 基础布局 ===== */
.ai-chat-container {
  height: 100vh;
  display: flex;
  background: #0d0d1a;
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: #e2e2f0;
}

/* ===== 动态背景光球 ===== */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.18;
  pointer-events: none;
  z-index: 0;
}
.orb1 {
  width: 500px; height: 500px;
  background: radial-gradient(circle, #7c3aed, transparent);
  top: -150px; left: -100px;
  animation: orbFloat1 18s ease-in-out infinite;
}
.orb2 {
  width: 400px; height: 400px;
  background: radial-gradient(circle, #06b6d4, transparent);
  bottom: -100px; right: 200px;
  animation: orbFloat2 22s ease-in-out infinite;
}
.orb3 {
  width: 300px; height: 300px;
  background: radial-gradient(circle, #a855f7, transparent);
  top: 40%; right: -80px;
  animation: orbFloat3 15s ease-in-out infinite;
}
@keyframes orbFloat1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(60px, 40px) scale(1.1); }
  66% { transform: translate(-40px, 80px) scale(0.9); }
}
@keyframes orbFloat2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-80px, -60px) scale(1.15); }
}
@keyframes orbFloat3 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(40px, -80px); }
}

/* ===== 左侧会话列表 ===== */
.session-list {
  width: 260px;
  min-width: 260px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: rgba(13, 13, 26, 0.92);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(124, 58, 237, 0.2);
  position: relative;
  z-index: 10;
}

.session-list-header {
  padding: 24px 20px 20px;
  border-bottom: 1px solid rgba(124, 58, 237, 0.15);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.brand-icon {
  font-size: 20px;
  background: linear-gradient(135deg, #7c3aed, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: brandPulse 3s ease-in-out infinite;
}
@keyframes brandPulse {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.4); }
}
.brand-name {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #c4b5fd, #e0f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.new-chat-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 11px 0;
  background: linear-gradient(135deg, rgba(124, 58, 237, 0.8), rgba(139, 92, 246, 0.8));
  color: white;
  border: 1px solid rgba(167, 139, 250, 0.3);
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  position: relative;
  overflow: hidden;
}
.new-chat-btn::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1), transparent);
  opacity: 0;
  transition: opacity 0.25s;
}
.new-chat-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(124, 58, 237, 0.4);
}
.new-chat-btn:hover::after { opacity: 1; }

.session-list-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}
.session-list-scroll::-webkit-scrollbar { width: 4px; }
.session-list-scroll::-webkit-scrollbar-track { background: transparent; }
.session-list-scroll::-webkit-scrollbar-thumb {
  background: rgba(124, 58, 237, 0.3);
  border-radius: 4px;
}

.session-list-ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-left: 2px solid transparent;
  color: #94a3b8;
  font-size: 13.5px;
}
.session-item:hover {
  background: rgba(124, 58, 237, 0.1);
  color: #c4b5fd;
  border-left-color: rgba(124, 58, 237, 0.4);
}
.session-item.active {
  background: linear-gradient(90deg, rgba(124, 58, 237, 0.25), rgba(124, 58, 237, 0.05));
  color: #c4b5fd;
  border-left-color: #7c3aed;
  font-weight: 600;
}
.session-dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: currentColor;
  flex-shrink: 0;
  opacity: 0.6;
}
.session-item.active .session-dot {
  background: #a78bfa;
  opacity: 1;
  box-shadow: 0 0 6px rgba(167, 139, 250, 0.6);
}
.session-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.sidebar-empty {
  padding: 30px 20px;
  text-align: center;
  color: rgba(148, 163, 184, 0.4);
  font-size: 13px;
  line-height: 1.8;
}

/* ===== 右侧聊天区域 ===== */
.chat-section {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 5;
  overflow: hidden;
}

/* ===== 顶部工具栏 ===== */
.top-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(13, 13, 26, 0.85);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(124, 58, 237, 0.15);
  flex-shrink: 0;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px; height: 36px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
  position: relative;
}
.icon-btn:hover:not(:disabled) {
  background: rgba(124, 58, 237, 0.2);
  border-color: rgba(124, 58, 237, 0.4);
  color: #c4b5fd;
  box-shadow: 0 0 12px rgba(124, 58, 237, 0.2);
}
.icon-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}
.upload-icon-btn { width: auto; padding: 0 12px; gap: 6px; font-size: 13px; }

.divider-v {
  width: 1px;
  height: 24px;
  background: rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

.label-tag {
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}
.model-select {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #c4b5fd;
  padding: 7px 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;
}
.model-select:focus {
  border-color: rgba(124, 58, 237, 0.5);
  box-shadow: 0 0 10px rgba(124, 58, 237, 0.15);
}
.model-select option {
  background: #1a1a2e;
  color: #e2e2f0;
}

.stream-control {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 4px;
}

/* 自定义 Toggle 开关 */
.toggle-switch {
  position: relative;
  display: inline-block;
  cursor: pointer;
}
.toggle-switch input { display: none; }
.toggle-track {
  display: block;
  width: 42px; height: 23px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  position: relative;
  transition: all 0.3s ease;
}
.toggle-thumb {
  position: absolute;
  width: 17px; height: 17px;
  background: #94a3b8;
  border-radius: 50%;
  top: 2px; left: 2px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}
.toggle-switch input:checked + .toggle-track {
  background: rgba(124, 58, 237, 0.6);
  border-color: rgba(124, 58, 237, 0.5);
  box-shadow: 0 0 10px rgba(124, 58, 237, 0.3);
}
.toggle-switch input:checked + .toggle-track .toggle-thumb {
  transform: translateX(19px);
  background: #a78bfa;
  box-shadow: 0 0 8px rgba(167, 139, 250, 0.5);
}

.toolbar-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-spinner {
  width: 14px; height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* ===== 消息区域 ===== */
.chat-messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.chat-messages::-webkit-scrollbar { width: 6px; }
.chat-messages::-webkit-scrollbar-track { background: transparent; }
.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(124, 58, 237, 0.25);
  border-radius: 6px;
}
.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(124, 58, 237, 0.4);
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 60px 0;
  user-select: none;
}
.empty-glow {
  font-size: 48px;
  background: linear-gradient(135deg, #7c3aed, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: glowPulse 2.5s ease-in-out infinite;
}
@keyframes glowPulse {
  0%, 100% { filter: brightness(1) drop-shadow(0 0 8px rgba(124, 58, 237, 0.4)); }
  50% { filter: brightness(1.3) drop-shadow(0 0 20px rgba(6, 182, 212, 0.5)); }
}
.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: #c4b5fd;
  margin: 0;
}
.empty-sub {
  font-size: 14px;
  color: #4b5563;
  margin: 0;
}

/* 消息行 */
.msg-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  animation: msgAppear 0.25s ease-out;
}
@keyframes msgAppear {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.msg-row-user { flex-direction: row-reverse; }

/* 头像 */
.avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  letter-spacing: 0.3px;
}
.avatar-user {
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
  color: white;
  box-shadow: 0 0 12px rgba(124, 58, 237, 0.4);
}
.avatar-ai {
  background: linear-gradient(135deg, #0e7490, #0284c7);
  color: white;
  box-shadow: 0 0 12px rgba(6, 182, 212, 0.3);
}

/* 气泡 */
.bubble {
  max-width: 68%;
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 14.5px;
  line-height: 1.65;
  word-break: break-word;
  position: relative;
}
.bubble-user {
  background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
  color: white;
  border-radius: 16px 4px 16px 16px;
  box-shadow: 0 4px 20px rgba(124, 58, 237, 0.25);
}
.bubble-ai {
  background: rgba(30, 27, 50, 0.9);
  color: #e2e8f0;
  border: 1px solid rgba(124, 58, 237, 0.15);
  border-radius: 4px 16px 16px 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.bubble-actions {
  position: absolute;
  top: 8px; right: 10px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}
.bubble-ai:hover .bubble-actions { opacity: 1; }

.tts-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px; height: 26px;
  background: rgba(124, 58, 237, 0.3);
  border: 1px solid rgba(124, 58, 237, 0.3);
  border-radius: 6px;
  color: #a78bfa;
  cursor: pointer;
  transition: all 0.2s;
}
.tts-btn:hover {
  background: rgba(124, 58, 237, 0.5);
  color: #c4b5fd;
}

.bubble-content {
  white-space: pre-wrap;
  word-break: break-word;
}
.bubble-content code {
  background: rgba(0, 0, 0, 0.3);
  padding: 1px 5px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  color: #06b6d4;
}

/* 流式打字动画 */
.typing-dots {
  display: flex;
  gap: 5px;
  padding-top: 8px;
  align-items: center;
}
.typing-dots span {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: #7c3aed;
  animation: dotBounce 1.4s ease-in-out infinite;
}
.typing-dots span:nth-child(2) { animation-delay: 0.2s; background: #8b5cf6; }
.typing-dots span:nth-child(3) { animation-delay: 0.4s; background: #06b6d4; }
@keyframes dotBounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1.1); opacity: 1; }
}

/* ===== 输入区域 ===== */
.chat-input-area {
  padding: 20px 28px 24px;
  background: rgba(13, 13, 26, 0.85);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(124, 58, 237, 0.15);
  flex-shrink: 0;
}

.input-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1.5px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  padding: 12px 14px 12px 18px;
  transition: all 0.25s ease;
}
.input-wrapper:focus-within {
  border-color: rgba(124, 58, 237, 0.5);
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1), 0 8px 30px rgba(124, 58, 237, 0.08);
  background: rgba(255, 255, 255, 0.06);
}

.input-wrapper textarea {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e2e2f0;
  font-size: 14.5px;
  line-height: 1.6;
  resize: none;
  min-height: 22px;
  max-height: 160px;
  overflow-y: auto;
  font-family: inherit;
}
.input-wrapper textarea::placeholder { color: #374151; }
.input-wrapper textarea:disabled { opacity: 0.5; cursor: not-allowed; }

.send-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px; height: 40px;
  border-radius: 10px;
  border: none;
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
  box-shadow: 0 4px 15px rgba(124, 58, 237, 0.3);
}
.send-btn:hover:not(:disabled) {
  transform: scale(1.08);
  box-shadow: 0 6px 20px rgba(124, 58, 237, 0.5);
}
.send-btn:disabled {
  background: rgba(255, 255, 255, 0.06);
  box-shadow: none;
  cursor: not-allowed;
  color: #374151;
}

.loading-ring {
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.25);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
</style>
