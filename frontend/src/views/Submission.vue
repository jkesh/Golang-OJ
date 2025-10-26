<template>
  <div class="submission-container">
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
            <h2>我的提交记录</h2>
          </div>
          
          <el-card class="submission-card">
            <div v-if="loading" class="loading-container">
              <el-skeleton :rows="10" animated />
            </div>
            
            <div v-else-if="submissions.length === 0" class="empty-container">
              <el-empty description="暂无提交记录" />
            </div>
            
            <el-table v-else :data="submissions" border stripe style="width: 100%">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="problem.title" label="题目" min-width="200">
                <template #default="scope">
                  <router-link :to="`/problem/${scope.row.problem_id}`" class="problem-link">
                    {{ scope.row.problem ? scope.row.problem.title : `题目 #${scope.row.problem_id}` }}
                  </router-link>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="120">
                <template #default="scope">
                  <el-tag :type="getStatusType(scope.row.status)">
                    {{ getStatusText(scope.row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="language" label="语言" width="100" />
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
              <el-table-column prop="created_at" label="提交时间" width="180" />
              <el-table-column label="操作" width="120">
                <template #default="scope">
                  <router-link :to="`/submission/${scope.row.id}`">
                    <el-button type="primary" size="small">查看详情</el-button>
                  </router-link>
                </template>
              </el-table-column>
            </el-table>
            
            <div class="pagination-container">
              <el-pagination
                v-if="total > 0"
                background
                layout="prev, pager, next"
                :total="total"
                :page-size="pageSize"
                :current-page="currentPage"
                @current-change="handlePageChange"
              />
            </div>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import submissionApi from '../api/submission'
import {ElMessage, ElMessageBox} from 'element-plus'
import Sidebar from "../../components/Sidebar.vue";

export default {
  name: 'Submission',
  components: {Sidebar},
  setup() {
    const store = useStore()
    const router = useRouter()
    const submissions = ref([])
    const loading = ref(true)
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)

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
        // 跳转到个人信息页面
        ElMessage('个人信息功能开发中')
      } else if (command === 'settings') {
        // 跳转到设置页面
        ElMessage('设置功能开发中')
      }
    }
    const onMenuSelect = (index) => {
      console.log(`Selected menu item with index: ${index}`)
    }
    const fetchSubmissions = async (page = 1) => {
      loading.value = true
      try {
        const response = await submissionApi.getSubmissions(page, pageSize.value)
        if (response && response.submissions && Array.isArray(response.submissions)) {
          submissions.value = response.submissions.map(item => ({
            id: item.ID,
            problem_id: item.problem_id,
            language: item.language,
            status: item.status,
            time_cost: item.run_time,
            memory_cost: item.memory,
            created_at: item.submitted_at,
            problem: {
              title: `题目 #${item.problem_id}`
            }
          }))
          total.value = response.total
        } else {
          submissions.value = []
          total.value = 0
        }
      } catch (error) {
        console.error('获取提交记录失败:', error)
        submissions.value = []
        total.value = 0
        ElMessage.warning('获取提交记录失败或暂无提交记录，请尝试提交后重试！')
      } finally {
        loading.value = false
      }
    }


    // 处理页码变化
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchSubmissions(page)
    }
    const viewProblem = (problemId) => {
      ElMessage.info('题目详情功能开发中')
      // 未来可以在这里实现跳转逻辑
      // router.push(`/problem/${problemId}`)
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
      fetchSubmissions()
    })
    
    return {
      submissions,
      loading,
      total,
      currentPage,
      pageSize,
      username,
      handlePageChange,
      getStatusType,
      getStatusText,
      onMenuSelect,
      handleCommand,
      viewProblem
    }
  }
}
</script>

<style scoped>
.submission-container {
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
  margin-bottom: 20px;
}

.submission-card {
  margin-bottom: 20px;
}

.loading-container,
.empty-container {
  padding: 40px 0;
  text-align: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.problem-link {
  color: #409eff;
  text-decoration: none;
}

.problem-link:hover {
  text-decoration: underline;
}
</style>