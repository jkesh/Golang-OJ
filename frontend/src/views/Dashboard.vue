<template>
  <div class="dashboard-container">
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
          <el-card>
            <template #header>
              <div class="card-header">
                <span>欢迎使用 OJPlus 在线评测系统</span>
              </div>
            </template>
            <div class="dashboard-content">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-card shadow="hover" class="dashboard-card">
                    <h3>题目总数</h3>
                    <div class="dashboard-number">{{ stats.total_problems }}</div>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="hover" class="dashboard-card">
                    <h3>已解决题目</h3>
                    <div class="dashboard-number">{{ stats.solved_problems }}</div>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="hover" class="dashboard-card">
                    <h3>提交次数</h3>
                    <div class="dashboard-number">{{ stats.total_submissions }}</div>
                  </el-card>
                </el-col>
              </el-row>
              
              <el-divider></el-divider>
              
              <h3>最近活动</h3>
              <el-timeline v-if="activities.length > 0">
                <el-timeline-item
                  v-for="(activity, index) in activities"
                  :key="index"
                  :timestamp="formatDate(activity.submitted_at)"
                  :type="getActivityType(activity.status)">
                  {{ getActivityText(activity) }}
                </el-timeline-item>
              </el-timeline>
              <el-empty v-else description="暂无活动记录" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import Sidebar from "../../components/Sidebar.vue"
import dashboardApi from "../api/dashboard.js"

export default {
  name: 'Dashboard',
  components: {Sidebar},
  setup() {
    const store = useStore()
    const router = useRouter()
    
    const username = computed(() => {
      const user = store.getters.currentUser
      return user ? user.username : '用户'
    })
    
    const stats = ref({
      total_problems: 0,
      solved_problems: 0,
      total_submissions: 0
    })
    
    const activities = ref([])
    
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

    // 处理菜单选择事件
    const onMenuSelect = (index) => {
      console.log(`Selected menu item with index: ${index}`)
      router.push(index)
    }
    
    // 获取统计数据
    const fetchStats = async () => {
      try {
        const response = await dashboardApi.getDashboardStats()
        if (response && response.stats) {
          stats.value = response.stats
        }
      } catch (error) {
        console.error('获取统计数据失败:', error)
        ElMessage.error('获取统计数据失败: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // 获取最近活动
    const fetchActivities = async () => {
      try {
        const response = await dashboardApi.getRecentActivities()
        if (response && response.activities) {
          activities.value = response.activities
        }
      } catch (error) {
        console.error('获取最近活动失败:', error)
        ElMessage.error('获取最近活动失败: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }
    
    // 获取活动类型
    const getActivityType = (status) => {
      const typeMap = {
        'accepted': 'success',
        'wrong_answer': 'danger',
        'time_limit_exceeded': 'warning',
        'memory_limit_exceeded': 'warning',
        'runtime_error': 'danger',
        'compile_error': 'info',
        'pending': 'info',
        'judging': 'info'
      }
      return typeMap[status] || 'info'
    }
    
    // 获取活动文本
    const getActivityText = (activity) => {
      const statusTextMap = {
        'accepted': '通过',
        'wrong_answer': '答案错误',
        'time_limit_exceeded': '超时',
        'memory_limit_exceeded': '内存超限',
        'runtime_error': '运行错误',
        'compile_error': '编译错误',
        'pending': '等待中',
        'judging': '评测中'
      }
      
      const statusText = statusTextMap[activity.status] || activity.status
      return `提交了题目 #${activity.problem_id} "${activity.problem_title}" - ${statusText}`
    }
    
    onMounted(() => {
      fetchStats()
      fetchActivities()
    })
    
    return {
      username,
      stats,
      activities,
      handleCommand,
      onMenuSelect,
      formatDate,
      getActivityType,
      getActivityText
    }
  }
}
</script>

<style scoped>

.dashboard-container {
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
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.1); /* 添加这行 */
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

.dashboard-content {
  margin-top: 20px;
}

.dashboard-card {
  text-align: center;
  margin-bottom: 20px;
}

.dashboard-number {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
  margin-top: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>