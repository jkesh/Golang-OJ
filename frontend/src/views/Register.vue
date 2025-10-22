<template>
  <div class="register-container">
    <div class="register-box">
      <div class="register-header">
        <h2>OJPlus 在线评测系统</h2>
        <p>创建新账户</p>
      </div>
      
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" class="register-form">
        <el-form-item prop="username">
          <el-input 
            v-model="registerForm.username" 
            prefix-icon="el-icon-user" 
            placeholder="用户名">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="email">
          <el-input 
            v-model="registerForm.email" 
            prefix-icon="el-icon-message" 
            placeholder="邮箱">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input 
            v-model="registerForm.password" 
            prefix-icon="el-icon-lock" 
            type="password" 
            placeholder="密码">
          </el-input>
        </el-form-item>
        
        <el-form-item prop="confirmPassword">
          <el-input 
            v-model="registerForm.confirmPassword" 
            prefix-icon="el-icon-lock" 
            type="password" 
            placeholder="确认密码"
            @keyup.enter="handleRegister">
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            :loading="loading" 
            class="register-button" 
            @click="handleRegister">
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>
        
        <div class="register-options">
          <span>已有账户? </span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

export default {
  name: 'Register',
  setup() {
    const store = useStore()
    const router = useRouter()
    const registerFormRef = ref(null)
    
    const loading = ref(false)
    
    const registerForm = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })
    
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else {
        if (registerForm.confirmPassword !== '') {
          // 这里不能直接调用validateField，因为可能还没有挂载
          if (registerFormRef.value) {
            registerFormRef.value.validateField('confirmPassword')
          }
        }
        callback()
      }
    }
    
    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    
    const registerRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' },
        { validator: validatePass, trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        { validator: validatePass2, trigger: 'blur' }
      ]
    }
    
    const handleRegister = () => {
      if (!registerFormRef.value) return
      
      registerFormRef.value.validate(valid => {
        if (valid) {
          loading.value = true
          
          const userData = {
            username: registerForm.username,
            email: registerForm.email,
            password: registerForm.password
          }
          
          store.dispatch('register', userData)
            .then(() => {
              ElMessage({
                message: '注册成功，请登录',
                type: 'success'
              })
              router.push('/login')
            })
            .catch(error => {
              console.error('注册错误:', error)
              ElMessage.error(error.response?.data?.message || '注册失败，请重试')
            })
            .finally(() => {
              loading.value = false
            })
        }
      })
    }
    
    return {
      registerForm,
      registerRules,
      registerFormRef,
      loading,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
}

.register-box {
  width: 400px;
  padding: 30px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 10px;
}

.register-header p {
  font-size: 14px;
  color: #909399;
}

.register-form {
  margin-top: 20px;
}

.register-button {
  width: 100%;
}

.register-options {
  display: flex;
  justify-content: center;
  margin-top: 15px;
  font-size: 14px;
  color: #606266;
}

.register-options a {
  color: #409eff;
  margin-left: 5px;
}
</style>