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
                    {{ getDifficultyText(problem.difficulty) }}
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

                <div class="section" v-if="problem.tags">
                  <h3>标签</h3>
                  <div class="tags">
                    <el-tag 
                      v-for="(tag, index) in problem.tags.split(',')" 
                      :key="index" 
                      class="tag-item"
                      size="small"
                    >
                      {{ tag.trim() }}
                    </el-tag>
                  </div>
                </div>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Sidebar from "../../components/Sidebar.vue"
import problemsApi from "../../src/api/problems.js"

export default {
  name: 'ProblemDetail',
  components: { Sidebar },
  setup() {
    const store = useStore()
    const route = useRoute()
    const router = useRouter()
    
    const problem = ref({})
    
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
        // 检查路由参数是否有效
        if (!route.params.id) {
          throw new Error('题目ID参数缺失')
        }
        
        // 检查路由参数是否为有效数字
        const problemId = parseInt(route.params.id)
        if (isNaN(problemId) || problemId <= 0) {
          throw new Error(`题目ID参数无效: ${route.params.id}`)
        }
        
        const response = await problemsApi.getProblem(problemId)
        if (!response || !response.problem) {
          throw new Error('无效的响应数据')
        }
        
        problem.value = {
          id: response.problem.ID,
          title: response.problem.title,
          description: response.problem.description,
          difficulty: response.problem.difficulty,
          time_limit: response.problem.time_limit,
          memory_limit: response.problem.memory_limit,
          tags: response.problem.tags || ''
        }
      } catch (error) {
        console.error('获取题目详情失败:', error)
        ElMessage.error('获取题目详情失败: ' + (error.response?.data?.error || error.message))
        problem.value = {}
      } finally {
        loading.value = false
      }
    }
    
    // 监听路由变化
    watch(() => route.params.id, (newId, oldId) => {
      if (newId !== oldId) {
        fetchProblem()
      }
    })
    
    // 获取难度类型
    const getDifficultyType = (difficulty) => {
      const typeMap = {
        'Easy': 'success',
        'Medium': 'warning',
        'Hard': 'danger'
      }
      return typeMap[difficulty] || 'info'
    }
    
    // 获取难度文本
    const getDifficultyText = (difficulty) => {
      const textMap = {
        'Easy': '简单',
        'Medium': '中等',
        'Hard': '困难'
      }
      return textMap[difficulty] || difficulty
    }
    
    // 提交代码
    const submitCode = async () => {
      if (!code.value.trim()) {
        ElMessage.warning('请输入代码')
        return
      }
      
      try {
        const submissionData = {
          problem_id: problem.value.id,
          language: selectedLanguage.value,
          code: code.value
        }

        // 调用提交代码的API接口
        const response = await problemsApi.submitCode(submissionData)

        if (response.status === 'success') {
          ElMessage.success('代码提交成功，正在评测中...')
          await router.push(`/submission/${response.id}`)
        } else {
          ElMessage.error(response.message || '代码提交失败')
        }
      } catch (error) {
        console.error('提交代码时出错:', error)
        ElMessage.error('提交代码时发生错误，请稍后重试')
      }
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
      getDifficultyText,
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
  gap: 5px;
}

.tag-item {
  margin-right: 5px;
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