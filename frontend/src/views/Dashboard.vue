<template>
  <div class="dashboard-container">
    <el-container>
      <el-aside width="200px">
        <el-menu
          default-active="1"
          class="sidebar-menu"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          router>
          <el-menu-item index="1" route="/dashboard">
            <i class="el-icon-menu"></i>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="2" route="/problems">
            <i class="el-icon-document"></i>
            <span>题目列表</span>
          </el-menu-item>
          <el-menu-item index="3" route="/submissions">
            <i class="el-icon-s-order"></i>
            <span>提交记录</span>
          </el-menu-item>
          <el-menu-item index="4" route="/ranking">
            <i class="el-icon-trophy"></i>
            <span>排行榜</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-container>
        <el-header>
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
        
        <el-main>
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
                    <div class="dashboard-number">100</div>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="hover" class="dashboard-card">
                    <h3>已解决题目</h3>
                    <div class="dashboard-number">25</div>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="hover" class="dashboard-card">
                    <h3>提交次数</h3>
                    <div class="dashboard-number">150</div>
                  </el-card>
                </el-col>
              </el-row>
              
              <el-divider></el-divider>
              
              <h3>最近活动</h3>
              <el-timeline>
                <el-timeline-item
                  v-for="(activity, index) in activities"
                  :key="index"
                  :timestamp="activity.timestamp"
                  :type="activity.type">
                  {{ activity.content }}
                </el-timeline-item>
              </el-timeline>
            </div>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'Dashboard',
  setup() {
    const store = useStore()
    const router = useRouter()
    
    const username = computed(() => {
      const user = store.getters.currentUser
      return user ? user.username : '用户'
    })
    
    const activities = ref([
      {
        content: '成功解决题目 #42 "两数之和"',
        timestamp: '2023-10-21 14:30',
        type: 'success'
      },
      {
        content: '提交了题目 #56 "合并区间"',
        timestamp: '2023-10-21 11:20',
        type: 'warning'
      },
      {
        content: '首次登录系统',
        timestamp: '2023-10-20 09:15',
        type: 'primary'
      }
    ])
    
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
    
    return {
      username,
      activities,
      handleCommand
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  height: 100%;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #e6e6e6;
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.el-aside {
  background-color: #304156;
  color: #bfcbd9;
}

.sidebar-menu {
  height: 100%;
  border-right: none;
}

.el-main {
  background-color: #f0f2f5;
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