<template>
  <div class="submission-container">
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
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useStore } from 'vuex'
import submissionApi from '../api/submission'
import { ElMessage } from 'element-plus'

export default {
  name: 'Submission',
  setup() {
    const store = useStore()
    const submissions = ref([])
    const loading = ref(true)
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)
    
    // 获取提交记录
    const fetchSubmissions = async (page = 1) => {
      loading.value = true
      try {
        const response = await submissionApi.getSubmissions({
          page,
          id: localStorage.getItem('user_id'),
          limit: pageSize.value
        })
        submissions.value = response.data.submissions || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('获取提交记录失败:', error)
        ElMessage.error('获取提交记录失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    // 处理页码变化
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchSubmissions(page)
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
        'Pending': '',
        'Judging': ''
      }
      return statusMap[status] || ''
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
      handlePageChange,
      getStatusType,
      getStatusText
    }
  }
}
</script>

<style scoped>
.submission-container {
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