<template>
  <div class="submission-detail-container">
    <el-container>
      <Sidebar />
      
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
                <el-descriptions-item label="提交时间">{{ submission.created_at }}</el-descriptions-item>
                <el-descriptions-item label="题目">
                  <router-link :to="`/problem/${submission.problem_id}`" class="problem-link">
                    {{ submission.problem ? submission.problem.title : `题目 #${submission.problem_id}` }}
                  </router-link>
                </el-descriptions-item>
                <el-descriptions-item label="状态">
                  <el-tag :type="getStatusType(submission.status)">
                    {{ getStatusText(submission.status) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="语言">{{ submission.language }}</el-descriptions-item>
                <el-descriptions-item label="用户">{{ submission.username }}</el-descriptions-item>
                <el-descriptions-item label="耗时">{{ submission.time_cost }} ms</el-descriptions-item>
                <el-descriptions-item label="内存">{{ submission.memory_cost }} KB</el-descriptions-item>
              </el-descriptions>
            </el-card>
            
            <el-card class="detail-card">
              <div class="card-header">
                <h3>代码</h3>
              </div>
              <pre class="code-block"><code>{{ submission.code }}</code></pre>
            </el-card>
            
            <el-card v-if="submission.result" class="detail-card">
              <div class="card-header">
                <h3>测试结果</h3>
              </div>
              <el-table :data="submission.result.test_cases || []" border stripe style="width: 100%">
                <el-table-column prop="test_case" label="测试点" width="100" />
                <el-table-column prop="status" label="状态" width="120">
                  <template #default="scope">
                    <el-tag :type="getStatusType(scope.row.status)">
                      {{ getStatusText(scope.row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="time_cost" label="耗时" width="100">
                  <template #default="scope">
                    {{ scope.row.time_cost }} ms
                  </template>
                </el-table-column>
                <el-table-column prop="memory_cost" label="内存" width="100">
                  <template #default="scope">
                    {{ scope.row.memory_cost }} KB
                  </template>
                </el-table-column>
                <el-table-column prop="error_message" label="错误信息" min-width="200" />
              </el-table>
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
    const loading = ref(true)
    
    const username = computed(() => {
      const user = store.getters.currentUser
      return user ? user.username : '用户'
    })
    
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
    
    // 获取提交详情
    const fetchSubmissionDetail = async () => {
      const submissionId = route.params.id
      loading.value = true
      
      try {
        const response = await submissionApi.getSubmission(submissionId)
        submission.value = response.data || {}
        
        // 获取提交结果
        try {
          const resultResponse = await submissionApi.getSubmissionResult(submissionId)
          submission.value.result = resultResponse.data || {}
        } catch (error) {
          console.error('获取提交结果失败:', error)
        }
      } catch (error) {
        console.error('获取提交详情失败:', error)
        ElMessage.error('获取提交详情失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    // 获取状态类型
    const getStatusType = (status) => {
      const statusMap = {
        'Accepted': 'success',
        'Wrong Answer': 'danger',
        'Time Limit Exceeded': 'warning',
        'Memory Limit Exceeded': 'warning',
        'Runtime Error': 'danger',
        'Compilation Error': 'info',
        'Pending': 'info',
        'Judging': 'info'
      }
      return statusMap[status] || 'info'
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      const statusTextMap = {
        'Accepted': '通过',
        'Wrong Answer': '答案错误',
        'Time Limit Exceeded': '超时',
        'Memory Limit Exceeded': '内存超限',
        'Runtime Error': '运行错误',
        'Compilation Error': '编译错误',
        'Pending': '等待中',
        'Judging': '评测中'
      }
      return statusTextMap[status] || status
    }
    
    onMounted(() => {
      fetchSubmissionDetail()
    })
    
    return {
      submission,
      loading,
      username,
      getStatusType,
      getStatusText,
      handleCommand
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

.code-block {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
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