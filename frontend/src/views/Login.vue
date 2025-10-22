<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>OJPlus 在线评测系统</h2>
        <p>登录您的账户</p>
      </div>
      
      <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form">
        <el-form-item prop="username">
          <el-input 
            v-model="loginForm.username" 
            prefix-icon="el-icon-user" 
            placeholder="用户名">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input 
            v-model="loginForm.password" 
            prefix-icon="el-icon-lock" 
            type="password" 
            placeholder="密码"
            @keyup.enter="handleLogin">
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            :loading="loading" 
            class="login-button" 
            @click="handleLogin">
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
        
        <div class="login-options">
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          <el-link type="primary">忘记密码?</el-link>
        </div>
        
        <div class="register-link">
          <span>还没有账户? </span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'Login',
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    const loginFormRef = ref(null)
    
    const loading = ref(false)
    const rememberMe = ref(false)
    
    const loginForm = reactive({
      username: '',
      password: ''
    })
    
    const loginRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
      ]
    }
    
    const handleLogin = () => {
      if (!loginFormRef.value) return
      
      loginFormRef.value.validate(valid => {
        if (valid) {
          loading.value = true
          
          store.dispatch('login', loginForm)
            .then(() => {
              ElMessage({
                message: '登录成功',
                type: 'success'
              })
              
              const redirectPath = route.query.redirect || '/dashboard'
              router.push(redirectPath)
            })
            .catch(error => {
              console.error('登录错误:', error)
              ElMessage.error(error?.response?.data?.message || '登录失败，请重试')
            })
            .finally(() => {
              loading.value = false
            })
        } else {
          return false
        }
      })
    }
    
    return {
      loginFormRef,
      loginForm,
      loginRules,
      loading,
      rememberMe,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.login-header p {
  margin: 10px 0 0;
  font-size: 14px;
  color: #909399;
}

.login-form {
  margin-top: 30px;
}

.login-button {
  width: 100%;
  margin-top: 10px;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
}

.register-link {
  text-align: center;
  margin-top: 15px;
  font-size: 14px;
  color: #606266;
}

.register-link a {
  color: #409eff;
  margin-left: 5px;
}
</style>