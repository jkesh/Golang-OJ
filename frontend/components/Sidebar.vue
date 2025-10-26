<!-- src/components/Sidebar.vue -->
<template>
  <el-aside class="el-aside">
    <el-menu
        :default-active="activeIndex"
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
        @select="handleSelect">
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
</template>

<script>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

export default {
  name: 'Sidebar',
  emits: ['menu-select'],
  setup(props, { emit }) {
    const route = useRoute()

    // 根据当前路由设置激活菜单项
    const activeIndex = ref('')
    const updateActiveIndex = () => {
      switch(route.path) {
        case '/dashboard': activeIndex.value = '1'; break
        case '/problems': activeIndex.value = '2'; break
        case '/submissions': activeIndex.value = '3'; break
        case '/ranking': activeIndex.value = '4'; break
        default: 
          // 检查是否是题目详情页
          if (route.path.startsWith('/problem/')) {
            activeIndex.value = '2'
          } else {
            activeIndex.value = '1'
          }
      }
    }

    // 监听路由变化更新激活状态
    updateActiveIndex()

    const handleSelect = (index) => {
      emit('menu-select', index)
    }

    return {
      activeIndex,
      handleSelect
    }
  }
}
</script>

<style scoped>

.el-aside {
  background-color: #304156;
  color: #bfcbd9;
  width: 200px;
  height: 100%;
  min-height: 100vh;
}
.sidebar-menu {
  height: 100%;
  border-right: none;
}
</style>