<template>
  <el-container class="layout">
    <el-aside width="200px" class="aside">
      <div class="logo">
        <el-icon :size="22"><Box /></el-icon>
        <span>回收交易平台</span>
      </div>
      <el-menu :default-active="activeMenu" router class="menu" background-color="#001529" text-color="#fff" active-text-color="#67c23a">
        <el-menu-item index="/user/market">
          <el-icon><DataLine /></el-icon><span>今日行情</span>
        </el-menu-item>
        <el-menu-item index="/user/sell">
          <el-icon><ShoppingCart /></el-icon><span>我要卖货</span>
        </el-menu-item>
        <el-menu-item index="/user/records">
          <el-icon><Document /></el-icon><span>我的记录</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">{{ pageTitle }}</div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><UserFilled /></el-icon>
              {{ user.name || user.username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { Box, DataLine, ShoppingCart, Document, UserFilled, ArrowDown } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const user = JSON.parse(localStorage.getItem('user') || '{}')
const activeMenu = computed(() => route.path)
const pageTitle = computed(() => route.name?.replace('User', '') || '')

const handleCommand = (cmd) => {
  if (cmd === 'logout') {
    ElMessageBox.confirm('确定退出登录？', '提示', { type: 'warning' }).then(() => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }).catch(() => {})
  }
}
</script>

<style scoped>
.layout { height: 100vh; }
.aside { background: #001529; }
.logo { height: 60px; display: flex; align-items: center; justify-content: center; gap: 8px; color: #fff; font-size: 15px; font-weight: bold; border-bottom: 1px solid #1f3a5f; }
.menu { border-right: none; }
.header { background: #fff; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #e8e8e8; padding: 0 24px; }
.header-left { font-size: 18px; font-weight: 600; color: #333; }
.user-info { cursor: pointer; display: flex; align-items: center; gap: 6px; color: #666; }
.main { background: #f0f2f5; padding: 20px; overflow-y: auto; }
</style>
