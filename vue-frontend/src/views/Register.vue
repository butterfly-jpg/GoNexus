<template>
  <div class="auth-container">
    <!-- 背景光球 -->
    <div class="bg-orb orb1"></div>
    <div class="bg-orb orb2"></div>
    <div class="bg-orb orb3"></div>

    <!-- 卡片 -->
    <div class="auth-card">
      <!-- 品牌头部 -->
      <div class="card-header">
        <div class="logo">
          <span class="logo-icon">✦</span>
          <span class="logo-text">GoNexus</span>
        </div>
        <h2 class="card-title">创建账号</h2>
        <p class="card-sub">填写信息以完成注册</p>
      </div>

      <!-- 表单 -->
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-position="top"
        class="auth-form"
      >
        <el-form-item label="邮箱" prop="email">
          <div class="input-wrap">
            <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/>
            </svg>
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱"
            />
          </div>
        </el-form-item>

        <el-form-item label="验证码" prop="captcha">
          <div class="captcha-row">
            <div class="input-wrap captcha-input">
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
              <el-input
                v-model="registerForm.captcha"
                placeholder="请输入验证码"
              />
            </div>
            <button
              type="button"
              class="code-btn"
              :class="{ disabled: countdown > 0 || codeLoading }"
              :disabled="countdown > 0 || codeLoading"
              @click="sendCode"
            >
              <span v-if="codeLoading" class="loading-ring-sm"></span>
              <span v-else-if="countdown > 0">{{ countdown }}s</span>
              <span v-else>发送验证码</span>
            </button>
          </div>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <div class="input-wrap">
            <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
            <el-input
              v-model="registerForm.password"
              placeholder="请输入密码（不少于6位）"
              type="password"
              show-password
            />
          </div>
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <div class="input-wrap">
            <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
            </svg>
            <el-input
              v-model="registerForm.confirmPassword"
              placeholder="请再次输入密码"
              type="password"
              show-password
            />
          </div>
        </el-form-item>

        <el-form-item>
          <button
            type="button"
            class="submit-btn"
            :class="{ loading: loading }"
            :disabled="loading"
            @click="handleRegister"
          >
            <span v-if="!loading">注册</span>
            <span v-else class="loading-ring"></span>
          </button>
        </el-form-item>
      </el-form>

      <div class="card-footer">
        <span>已有账号？</span>
        <button class="link-btn" @click="$router.push('/login')">立即登录</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../utils/api'

export default {
  name: 'RegisterView',
  setup() {
    const router = useRouter()
    const registerFormRef = ref()
    const loading = ref(false)
    const codeLoading = ref(false)
    const countdown = ref(0)

    const registerForm = reactive({
      email: '',
      captcha: '',
      password: '',
      confirmPassword: ''
    })

    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }

    const registerRules = {
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        { validator: validateConfirmPassword, trigger: 'blur' }
      ]
    }

    const sendCode = async () => {
      if (!registerForm.email) {
        ElMessage.warning('请先输入邮箱')
        return
      }
      try {
        codeLoading.value = true
        const response = await api.post('/user/captcha', { email: registerForm.email })
        if (response.data.status_code === 1000) {
          ElMessage.success('验证码发送成功')
          countdown.value = 60
          const timer = setInterval(() => {
            countdown.value--
            if (countdown.value <= 0) clearInterval(timer)
          }, 1000)
        } else {
          ElMessage.error(response.data.status_msg || '验证码发送失败')
        }
      } catch (error) {
        console.error('Send code error:', error)
        ElMessage.error('验证码发送失败，请重试')
      } finally {
        codeLoading.value = false
      }
    }

    const handleRegister = async () => {
      try {
        await registerFormRef.value.validate()
        loading.value = true
        const response = await api.post('/user/register', {
          email: registerForm.email,
          captcha: registerForm.captcha,
          password: registerForm.password
        })
        if (response.data.status_code === 1000) {
          ElMessage.success('注册成功，请登录')
          router.push('/login')
        } else {
          ElMessage.error(response.data.status_msg || '注册失败')
        }
      } catch (error) {
        console.error('Register error:', error)
        ElMessage.error('注册失败，请重试')
      } finally {
        loading.value = false
      }
    }

    return {
      registerFormRef, loading, codeLoading, countdown,
      registerForm, registerRules, sendCode, handleRegister
    }
  }
}
</script>

<style scoped>
/* ===== 容器 ===== */
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0d0d1a;
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* ===== 背景光球 ===== */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.18;
  pointer-events: none;
}
.orb1 {
  width: 500px; height: 500px;
  background: radial-gradient(circle, #7c3aed, transparent);
  top: -150px; right: -80px;
  animation: orbFloat1 20s ease-in-out infinite;
}
.orb2 {
  width: 400px; height: 400px;
  background: radial-gradient(circle, #06b6d4, transparent);
  bottom: -80px; left: -60px;
  animation: orbFloat2 24s ease-in-out infinite;
}
.orb3 {
  width: 260px; height: 260px;
  background: radial-gradient(circle, #a855f7, transparent);
  top: 30%; left: 65%;
  animation: orbFloat3 16s ease-in-out infinite;
}
@keyframes orbFloat1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-60px, 70px) scale(1.08); }
}
@keyframes orbFloat2 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(80px, -60px) scale(1.1); }
}
@keyframes orbFloat3 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-40px, -70px); }
}

/* ===== 卡片 ===== */
.auth-card {
  width: 440px;
  background: rgba(18, 15, 35, 0.9);
  border: 1px solid rgba(124, 58, 237, 0.2);
  border-radius: 24px;
  padding: 36px 36px 32px;
  backdrop-filter: blur(20px);
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(124, 58, 237, 0.1);
  position: relative;
  z-index: 10;
  animation: cardIn 0.6s cubic-bezier(0.4, 0, 0.2, 1) both;
}
@keyframes cardIn {
  from { opacity: 0; transform: translateY(24px) scale(0.97); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}
.auth-card::before {
  content: '';
  position: absolute;
  top: 0; left: 20%; right: 20%;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(124, 58, 237, 0.6), rgba(6, 182, 212, 0.4), transparent);
  border-radius: 1px;
}

/* ===== 头部 ===== */
.card-header {
  text-align: center;
  margin-bottom: 28px;
}
.logo {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 18px;
}
.logo-icon {
  font-size: 20px;
  background: linear-gradient(135deg, #7c3aed, #06b6d4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: iconPulse 3s ease-in-out infinite;
}
@keyframes iconPulse {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.6); }
}
.logo-text {
  font-size: 17px;
  font-weight: 700;
  background: linear-gradient(135deg, #c4b5fd, #e0f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.card-title {
  font-size: 26px;
  font-weight: 700;
  margin: 0 0 8px;
  color: #f1f5f9;
}
.card-sub {
  font-size: 14px;
  color: #4b5563;
  margin: 0;
}

/* ===== 表单 ===== */
.auth-form :deep(.el-form-item__label) {
  color: #94a3b8 !important;
  font-size: 13px;
  font-weight: 500;
  padding-bottom: 6px;
  line-height: 1;
}
.auth-form :deep(.el-form-item) {
  margin-bottom: 18px;
}
.auth-form :deep(.el-form-item__error) {
  color: #f87171;
  font-size: 12px;
  padding-top: 4px;
}

/* 输入框容器 */
.input-wrap {
  position: relative;
  width: 100%;
}
.input-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #374151;
  z-index: 2;
  pointer-events: none;
}
.input-wrap :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px !important;
  box-shadow: none !important;
  padding-left: 36px !important;
  transition: all 0.25s ease;
}
.input-wrap :deep(.el-input__wrapper:hover) {
  border-color: rgba(124, 58, 237, 0.4) !important;
}
.input-wrap :deep(.el-input__wrapper.is-focus) {
  border-color: rgba(124, 58, 237, 0.6) !important;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}
.input-wrap :deep(.el-input__inner) {
  color: #e2e8f0 !important;
  font-size: 14px !important;
  background: transparent !important;
  height: 44px !important;
}
.input-wrap :deep(.el-input__inner::placeholder) {
  color: #374151 !important;
}
.input-wrap :deep(.el-input__suffix) {
  color: #4b5563 !important;
}

/* 验证码行 */
.captcha-row {
  display: flex;
  gap: 10px;
  align-items: center;
}
.captcha-input {
  flex: 1;
}

.code-btn {
  height: 44px;
  padding: 0 14px;
  white-space: nowrap;
  background: rgba(124, 58, 237, 0.15);
  border: 1px solid rgba(124, 58, 237, 0.3);
  border-radius: 10px;
  color: #a78bfa;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 96px;
  flex-shrink: 0;
}
.code-btn:hover:not(.disabled) {
  background: rgba(124, 58, 237, 0.28);
  border-color: rgba(124, 58, 237, 0.5);
  color: #c4b5fd;
  box-shadow: 0 0 12px rgba(124, 58, 237, 0.2);
}
.code-btn.disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

/* 提交按钮 */
.submit-btn {
  width: 100%;
  height: 48px;
  background: linear-gradient(135deg, #7c3aed, #6d28d9);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 24px rgba(124, 58, 237, 0.35);
  margin-top: 4px;
  position: relative;
  overflow: hidden;
}
.submit-btn::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1), transparent);
  opacity: 0;
  transition: opacity 0.25s;
}
.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 32px rgba(124, 58, 237, 0.5);
}
.submit-btn:hover:not(:disabled)::after { opacity: 1; }
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-ring {
  width: 20px; height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
.loading-ring-sm {
  width: 14px; height: 14px;
  border: 2px solid rgba(167, 139, 250, 0.3);
  border-top-color: #a78bfa;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* ===== 底部链接 ===== */
.card-footer {
  text-align: center;
  font-size: 13.5px;
  color: #4b5563;
  margin-top: 16px;
}
.link-btn {
  background: none;
  border: none;
  color: #a78bfa;
  font-size: 13.5px;
  font-weight: 600;
  cursor: pointer;
  padding: 0 4px;
  transition: color 0.2s;
}
.link-btn:hover { color: #c4b5fd; }
</style>