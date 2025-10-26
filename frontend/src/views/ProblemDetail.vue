<template>
  <div class="problem-detail-container">
    <el-container>
      <Sidebar @menu-selected="onMenuSelect" />
      
      <el-container>
        <el-header class="el-header">
          <div class="header-container">
            <div class="header-title">OJPlus 在线评测系统</div>
            <div class="header-right">
              <el-dropdown trigger="click" @command="handleCommand">
                <span class="el-dropdown-link">
                  {{ username }}<i class="el-icon-arrow-down el-icon--right"></i>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                    <el-dropdown-item command="settings">设置</el-dropdown-item>
                    <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>
        
        <el-main class="el-main">
          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="10" animated />
          </div>
          
          <div v-else>
            <el-card class="problem-card">
              <template #header>
                <div class="problem-header">
                  <h2>{{ problem.title }}</h2>
                  <el-tag :type="getDifficultyType(problem.difficulty)">
                    {{ problem.difficulty }}
                  </el-tag>
                </div>
              </template>
              
              <div class="problem-content">
                <div class="section">
                  <h3>题目描述</h3>
                  <p>{{ problem.description }}</p>
                </div>
                
                <el-divider></el-divider>
                
                <div class="section">
                  <h3>限制条件</h3>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-card class="limit-card">
                        <div class="limit-item">
                          <span class="limit-label">时间限制:</span>
                          <span class="limit-value">{{ problem.time_limit }} ms</span>
                        </div>
                      </el-card>
                    </el-col>
                    <el-col :span="12">
                      <el-card class="limit-card">
                        <div class="limit-item">
                          <span class="limit-label">内存限制:</span>
                          <span class="limit-value">{{ problem.memory_limit }} MB</span>
                        </div>
                      </el-card>
                    </el-col>
                  </el-row>
                </div>
                
                <el-divider></el-divider>

              </div>
            </el-card>
            
            <el-card class="code-card">
              <template #header>
                <div class="code-header">
                  <span>代码编辑器</span>
                </div>
              </template>
              
              <div class="code-editor-container">
                <el-select v-model="selectedLanguage" placeholder="选择编程语言" style="width: 150px; margin-bottom: 15px;">
                  <el-option
                    v-for="language in languages"
                    :key="language.value"
                    :label="language.label"
                    :value="language.value">
                  </el-option>
                </el-select>
                
                <el-input
                  type="textarea"
                  :rows="15"
                  placeholder="请输入代码"
                  v-model="code"
                  class="code-editor">
                </el-input>
                
                <div class="submit-button-container">
                  <el-button type="primary" @click="submitCode">提交代码</el-button>
                </div>
              </div>
            </el-card>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Sidebar from "../../components/Sidebar.vue"
import problemsApi from "../api/problems";

export default {
  name: 'ProblemDetail',
  components: { Sidebar },
  setup() {
    const store = useStore()
    const route = useRoute()
    const router = useRouter()
    
    const problem = ref({
      id: '',
      title: '',
      description: '',
      difficulty: 'medium',
      time_limit: 1000,
      memory_limit: 256,
      tags: ''
    })
    
    const loading = ref(true)
    const code = ref('')
    const selectedLanguage = ref('cpp')
    
    const languages = [
      { value: 'cpp', label: 'C++' },
      { value: 'c', label: 'C' },
      { value: 'java', label: 'Java' },
      { value: 'python', label: 'Python' },
      { value: 'go', label: 'Go' }
    ]
    
    const username = computed(() => {
      const user = store.getters.currentUser
      return user ? user.username : '用户'
    })
    
    // 获取题目详情
    const fetchProblem = async () => {
      loading.value = true
      try {
        const response = await problemsApi.getProblem(route.params.id)
        problem.value = response.problem
        setTimeout(() => {
          problem.value = {
            id: route.params.id,
            title: problem.value.title,
            description: problem.value.description,
            difficulty: problem.value.difficulty,
            time_limit: problem.value.time_limit,
            memory_limit: problem.value.memory_limit,
          }
          loading.value = false
        }, 500)
      } catch (error) {
        console.error('获取题目详情失败:', error)
        ElMessage.error('获取题目详情失败')
        loading.value = false
      }
    }
    
    // 获取难度类型
    const getDifficultyType = (difficulty) => {
      const typeMap = {
        'Easy': 'success',
        'Medium': 'warning',
        'Hard': 'danger'
      }
      return typeMap[difficulty] || 'info'
    }
    
    // 提交代码
    const submitCode = () => {
      if (!code.value.trim()) {
        ElMessage.warning('请输入代码')
        return
      }
      
      ElMessage.success('代码提交成功，正在评测中...')
      // 实际开发中应该调用提交代码的API接口
    }
    
    const handleCommand = (command) => {
      if (command === 'logout') {
        ElMessageBox.confirm(
          '确定要退出登录吗?',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).then(() => {
          store.dispatch('logout')
          router.push('/login')
          ElMessage({
            type: 'success',
            message: '已成功退出登录'
          })
        }).catch(() => {})
      } else if (command === 'profile') {
        ElMessage('个人信息功能开发中')
      } else if (command === 'settings') {
        ElMessage('设置功能开发中')
      }
    }
    
    const onMenuSelect = (index) => {
      console.log(`Selected menu item with index: ${index}`)
    }
    
    onMounted(() => {
      fetchProblem()
    })
    
    return {
      problem,
      loading,
      code,
      selectedLanguage,
      languages,
      username,
      getDifficultyType,
      submitCode,
      handleCommand,
      onMenuSelect
    }
  }
}
</script>

<style scoped>
.problem-detail-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1px;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.1);
}

.header-title {
  font-size: 18px;
  font-weight: bold;
}

.header-right {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
  display: flex;
  align-items: center;
}

.el-main {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.loading-container {
  padding: 40px 0;
  text-align: center;
}

.problem-card {
  margin-bottom: 20px;
}

.problem-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.problem-content .section {
  margin-bottom: 20px;
}

.limit-card {
  text-align: center;
}

.limit-item {
  display: flex;
  justify-content: space-between;
}

.limit-label {
  font-weight: bold;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-item {
  margin-right: 10px;
}

.code-card {
  margin-bottom: 20px;
}

.code-editor {
  font-family: 'Courier New', monospace;
}

.submit-button-container {
  margin-top: 15px;
  text-align: center;
}
</style>