<template>
  <div class="submission-detail-container">
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
          <div class="page-header">
            <h2>提交详情</h2>
            <router-link to="/submissions" class="back-link">
              <el-button type="primary" plain>返回提交列表</el-button>
            </router-link>
          </div>
          
          <el-card v-if="loading" class="detail-card">
            <el-skeleton :rows="10" animated />
          </el-card>
          
          <template v-else>
            <el-card class="detail-card">
              <el-descriptions title="基本信息" :column="2" border>
                <el-descriptions-item label="提交ID">{{ submission.id }}</el-descriptions-item>
                <el-descriptions-item label="提交时间">{{ formatDate(submission.created_at) }}</el-descriptions-item>
                <el-descriptions-item label="题目">
                  <router-link :to="`/problem/${submission.problem_id}`" class="problem-link">
                    {{ `题目 #${submission.problem_id}` }}
                  </router-link>
                </el-descriptions-item>
                <el-descriptions-item label="状态">
                  <el-tag :type="getStatusType(submission.status)">
                    {{ getStatusText(submission.status) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="语言">{{ submission.language }}</el-descriptions-item>
                <el-descriptions-item label="耗时">{{ submission.run_time }} ms</el-descriptions-item>
                <el-descriptions-item label="内存">{{ submission.memory }} KB</el-descriptions-item>
                <el-descriptions-item label="评测时间">{{ formatDate(submission.judged_at) }}</el-descriptions-item>
              </el-descriptions>
            </el-card>
            
            <el-card class="detail-card">
              <div class="card-header">
                <h3>代码</h3>
              </div>
              <pre class="code-block"><code>{{ submission.code }}</code></pre>
            </el-card>
            
            <el-card v-if="submission.status !== 'pending' && submission.status !== 'judging'" class="detail-card">
              <div class="card-header">
                <h3>测试结果</h3>
              </div>
              <el-table :data="testCases" border stripe style="width: 100%">
                <el-table-column prop="test_case_id" label="测试点" width="80" />
                <el-table-column label="状态" width="100">
                  <template #default="scope">
                    <el-tag :type="getTestCaseStatusType(scope.row.status)" size="small">
                      {{ getTestCaseStatusText(scope.row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="耗时" width="100">
                  <template #default="scope">
                    {{ scope.row.run_time }} ms
                  </template>
                </el-table-column>
                <el-table-column label="内存" width="100">
                  <template #default="scope">
                    {{ scope.row.memory }} KB
                  </template>
                </el-table-column>
                <el-table-column label="输入" min-width="120">
                  <template #default="scope">
                    <div class="text-content">{{ scope.row.input }}</div>
                  </template>
                </el-table-column>
                <el-table-column label="期望输出" min-width="120">
                  <template #default="scope">
                    <div class="text-content">{{ scope.row.expected_output }}</div>
                  </template>
                </el-table-column>
                <el-table-column label="用户输出" min-width="120">
                  <template #default="scope">
                    <div class="text-content">{{ scope.row.user_output || '-' }}</div>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
            
            <el-card v-if="submission.error_message" class="detail-card">
              <div class="card-header">
                <h3>错误信息</h3>
              </div>
              <pre class="error-block"><code>{{ submission.error_message }}</code></pre>
            </el-card>
          </template>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute, useRouter } from 'vue-router'
import submissionApi from '../api/submission'
import { ElMessage, ElMessageBox } from 'element-plus'
import Sidebar from "../../components/Sidebar.vue"

export default {
  name: 'SubmissionDetail',
  components: { Sidebar },
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    const submission = ref({})
    const testCases = ref([])
    const loading = ref(true)
    
    const username = computed(() => {
      const user = store.getters.currentUser
      return user ? user.username : '用户'
    })
    
    // 处理用户命令
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
        // 跳转到个人信息页面
        ElMessage('个人信息功能开发中')
      } else if (command === 'settings') {
        // 跳转到设置页面
        ElMessage('设置功能开发中')
      }
    }
    
    const onMenuSelect = (index) => {
      console.log(`Selected menu item with index: ${index}`)
      router.push(index)
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }
    
    // 获取提交状态类型
    const getStatusType = (status) => {
      const statusMap = {
        'Accepted': 'success',
        'WrongAnswer': 'danger',
        'TimeLimitExceeded': 'danger',
        'MemoryLimitExceeded': 'danger',
        'RuntimeError': 'danger',
        'CompileError': 'danger',
        'Pending': 'info',
        'Judging': 'info'
      }
      return statusMap[status] || 'info'
    }
    
    // 获取测试用例状态类型
    const getTestCaseStatusType = (status) => {
      const statusMap = {
        'passed': 'success',
        'failed': 'danger',
        'time_limit_exceeded': 'warning',
        'memory_limit_exceeded': 'warning',
        'runtime_error': 'danger'
      }
      return statusMap[status] || 'info'
    }
    
    // 获取提交状态文本
    const getStatusText = (status) => {
      const statusTextMap = {
        'accepted': '通过',
        'wrong_answer': '答案错误',
        'time_limit_exceeded': '超时',
        'memory_limit_exceeded': '内存超限',
        'runtime_error': '运行错误',
        'compilation_error': '编译错误',
        'pending': '等待中',
        'judging': '评测中'
      }
      return statusTextMap[status] || status
    }
    
    // 获取测试用例状态文本
    const getTestCaseStatusText = (status) => {
      const statusTextMap = {
        'passed': '通过',
        'failed': '失败',
        'time_limit_exceeded': '超时',
        'memory_limit_exceeded': '内存超限',
        'runtime_error': '运行错误'
      }
      return statusTextMap[status] || status
    }
    
    // 获取提交详情
    const fetchSubmissionDetail = async () => {
      const submissionId = route.params.id
      loading.value = true
      
      try {
        // 获取提交详情
        const response = await submissionApi.getSubmission(submissionId)
        submission.value = response.submission || {}
        
        // 获取提交结果（包括测试用例）
        try {
          const resultResponse = await submissionApi.getSubmissionResult(submissionId)
          const resultData = resultResponse.data || {}
          testCases.value = resultData.test_cases || []
        } catch (error) {
          console.error('获取测试用例结果失败:', error)
        }
      } catch (error) {
        console.error('获取提交详情失败:', error)
        ElMessage.error('获取提交详情失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    onMounted(() => {
      fetchSubmissionDetail()
    })
    
    return {
      submission,
      testCases,
      loading,
      username,
      formatDate,
      getStatusType,
      getTestCaseStatusType,
      getStatusText,
      getTestCaseStatusText,
      handleCommand,
      onMenuSelect
    }
  }
}
</script>

<style scoped>
.submission-detail-container {
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

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.detail-card {
  margin-bottom: 20px;
}

.card-header {
  margin-bottom: 15px;
}

.code-block, .error-block {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.4;
}

.text-content {
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 100px;
  overflow-y: auto;
}

.problem-link {
  color: #409eff;
  text-decoration: none;
}

.problem-link:hover {
  text-decoration: underline;
}

.back-link {
  text-decoration: none;
}
</style>