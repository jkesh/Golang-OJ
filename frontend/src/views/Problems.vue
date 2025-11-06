<template>
  <div class="problems-container">
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
            <h2>题目列表</h2>
          </div>
          
          <el-card class="problems-card">
            <div v-if="loading" class="loading-container">
              <el-skeleton :rows="10" animated />
            </div>

            <div v-else-if="problems.length === 0" class="empty-container">
              <el-empty description="暂无题目" />
            </div>
            
            <el-table v-else :data="problems" border stripe style="width: 100%">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="title" label="题目" min-width="200">
                <template #default="scope">
                  <router-link :to="`/problem/${scope.row.id}`" class="problem-link">
                    {{ scope.row.title }}
                  </router-link>
                </template>
              </el-table-column>
              <el-table-column prop="difficulty" label="难度" width="120">
                <template #default="scope">
                  <el-tag :type="getDifficultyType(scope.row.difficulty)">
                    {{ getDifficultyText(scope.row.difficulty) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="tags" label="标签" min-width="200">
                <template #default="scope">
                  <div class="tags">
                    <el-tag 
                      v-for="(tag, index) in scope.row.tags.split(',')" 
                      :key="index" 
                      class="tag-item"
                      size="small"
                    >
                      {{ tag.trim() }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="solved" label="通过率" width="120">
                <template #default="scope">
                  {{ scope.row.solvedRate }}
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
import { ElMessage, ElMessageBox } from 'element-plus'
import Sidebar from "../../components/Sidebar.vue";
import problemsApi from "../../src/api/problems.js";

export default {
  name: 'Problems',
  components: {Sidebar},
  setup() {
    const store = useStore()
    const router = useRouter()
    const problems = ref([])
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
    
    // 获取题目列表
    const fetchProblems = async (page = 1) => {
      loading.value = true
      try {
        // 调用真实的API接口获取题目列表
        const response = await problemsApi.getProblems(page, pageSize.value)
        problems.value = response.problems.map(problem => ({
          id: problem.ID,  // 改为大写ID以匹配后端返回的数据
          title: problem.title,
          difficulty: problem.difficulty,
          tags: problem.tags || '',
        }))
        total.value = response.total || response.problems.length
      } catch (error) {
        console.error('获取题目列表失败:', error)
        problems.value = []
        total.value = 0
        ElMessage.warning('获取题目列表失败！')
      } finally {
        loading.value = false
      }
    }

    // 处理页码变化
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchProblems(page)
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
    
    // 获取难度文本
    const getDifficultyText = (difficulty) => {
      const textMap = {
        'Easy': '简单',
        'Medium': '中等',
        'Hard': '困难'
      }
      return textMap[difficulty] || difficulty
    }
    
    onMounted(() => {
      fetchProblems()
    })
    
    return {
      problems,
      loading,
      total,
      currentPage,
      pageSize,
      username,
      handlePageChange,
      getDifficultyType,
      getDifficultyText,
      onMenuSelect,
      handleCommand
    }
  }
}
</script>

<style scoped>
.problems-container {
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

.problems-card {
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

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.tag-item {
  margin-right: 5px;
}
</style>