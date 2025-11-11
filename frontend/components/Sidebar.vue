<template>
  <el-aside width="200px" class="sidebar">
    <el-menu
      :default-active="activeIndex"
      class="el-menu-vertical"
      @select="handleSelect"
      :unique-opened="true"
    >
      <el-menu-item index="/dashboard">
        <i class="el-icon-house"></i>
        <span>首页</span>
      </el-menu-item>
      
      <el-menu-item index="/problems">
        <i class="el-icon-document"></i>
        <span>题目列表</span>
      </el-menu-item>
      
      <el-menu-item index="/submissions">
        <i class="el-icon-tickets"></i>
        <span>提交记录</span>
      </el-menu-item>
      
      <el-menu-item index="/admin" v-if="isAdmin">
        <i class="el-icon-setting"></i>
        <span>管理系统</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'

export default {
  name: 'Sidebar',
  emits: ['menu-selected'],
  setup(props, { emit }) {
    const store = useStore()
    const route = useRoute()
    
    const activeIndex = ref(route.path)
    
    const isAdmin = computed(() => {
      const user = store.getters.currentUser
      return user && user.role === 'admin'
    })
    
    const handleSelect = (index) => {
      activeIndex.value = index
      emit('menu-selected', index)
    }
    
    // 监听路由变化以更新激活状态
    watch(() => route.path, (newPath) => {
      activeIndex.value = newPath
    })
    
    onMounted(() => {
      activeIndex.value = route.path
    })
    
    return {
      activeIndex,
      isAdmin,
      handleSelect
    }
  }
}
</script>

<style scoped>
.sidebar {
  background-color: #ffffff;
  border-right: 1px solid #e6e6e6;
  box-shadow: 2px 0 2px rgba(0, 0, 0, 0.1);
}

.el-menu-vertical {
  height: 100%;
  border-right: none;
}

.el-menu-item i {
  margin-right: 10px;
}
</style>