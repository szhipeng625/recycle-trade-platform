<template>
  <div class="login-page">
    <div class="login-box">
      <div class="login-header">
        <el-icon :size="40" color="#409eff"><Box /></el-icon>
        <h1>废旧物资回收交易平台</h1>
        <p>管理员 / 用户登录</p>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="0" @keyup.enter="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" size="large" :prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" :prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" style="width:100%" :loading="loading" @click="handleLogin">登 录</el-button>
        </el-form-item>
      </el-form>
      <div class="login-tips">
        <p>管理员账号: admin / admin123</p>
        <p>用户可自行注册账号</p>
      </div>
      <div class="register-link">
        没有账号？<el-link type="primary" @click="showRegister = true">立即注册</el-link>
      </div>
    </div>

    <el-dialog v-model="showRegister" title="用户注册" width="420px">
      <el-form :model="regForm" :rules="regRules" ref="regFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="regForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="regForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="regForm.name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="regForm.phone" placeholder="请输入手机号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRegister = false">取消</el-button>
        <el-button type="primary" :loading="regLoading" @click="handleRegister">注册</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Box } from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const formRef = ref()
const regFormRef = ref()
const loading = ref(false)
const regLoading = ref(false)
const showRegister = ref(false)

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const regForm = reactive({ username: '', password: '', name: '', phone: '' })
const regRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const handleLogin = () => {
  formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const res = await api.post('/auth/login', form)
      const data = res.data || {}
      if (!data.token) {
        ElMessage.error('登录响应异常：缺少 token')
        return
      }
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user || {}))
      ElMessage.success('登录成功')
      const role = data.user?.role || 'user'
      const target = role === 'admin' ? '/admin' : '/user'
      await router.push(target)
    } catch (e) {
      console.error('登录失败:', e)
      ElMessage.error(e.message || '登录失败，请重试')
    } finally {
      loading.value = false
    }
  })
}

const handleRegister = () => {
  regFormRef.value.validate(async (valid) => {
    if (!valid) return
    regLoading.value = true
    try {
      await api.post('/auth/register', regForm)
      ElMessage.success('注册成功，请登录')
      showRegister.value = false
      form.username = regForm.username
    } catch (e) {} finally {
      regLoading.value = false
    }
  })
}
</script>

<style scoped>
.login-page { height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.login-box { width: 400px; background: #fff; border-radius: 12px; padding: 40px; box-shadow: 0 20px 60px rgba(0,0,0,0.3); }
.login-header { text-align: center; margin-bottom: 30px; }
.login-header h1 { font-size: 22px; margin: 12px 0 6px; color: #333; }
.login-header p { color: #999; font-size: 14px; }
.login-tips { text-align: center; color: #999; font-size: 12px; line-height: 1.8; margin-top: 10px; }
.register-link { text-align: center; margin-top: 16px; font-size: 14px; color: #666; }
</style>
