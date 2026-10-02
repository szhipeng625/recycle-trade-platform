<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <el-icon :size="24"><Box /></el-icon>
        <span>回收交易管理</span>
      </div>
      <el-menu :default-active="activeMenu" router class="menu" background-color="#001529" text-color="#fff" active-text-color="#409eff">
        <el-menu-item index="/admin/dashboard">
          <el-icon><DataAnalysis /></el-icon><span>数据概览</span>
        </el-menu-item>
        <el-menu-item index="/admin/pricing">
          <el-icon><PriceTag /></el-icon><span>每日定价</span>
        </el-menu-item>
        <el-menu-item index="/admin/transactions">
          <el-icon><Document /></el-icon><span>交易开单</span>
        </el-menu-item>
        <el-menu-item index="/admin/goods">
          <el-icon><Goods /></el-icon><span>货物管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/sellers">
          <el-icon><User /></el-icon><span>卖货人管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/inventory">
          <el-icon><Warehouse /></el-icon><span>仓库库存</span>
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
import { Box, DataAnalysis, PriceTag, Document, Goods, User, Warehouse, UserFilled, ArrowDown } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const user = JSON.parse(localStorage.getItem('user') || '{}')
const activeMenu = computed(() => route.path)
const pageTitle = computed(() => route.name?.replace('Admin', '') || '')

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
.logo { height: 60px; display: flex; align-items: center; justify-content: center; gap: 8px; color: #fff; font-size: 16px; font-weight: bold; border-bottom: 1px solid #1f3a5f; }
.menu { border-right: none; }
.header { background: #fff; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #e8e8e8; padding: 0 24px; }
.header-left { font-size: 18px; font-weight: 600; color: #333; }
.user-info { cursor: pointer; display: flex; align-items: center; gap: 6px; color: #666; }
.main { background: #f0f2f5; padding: 20px; overflow-y: auto; }
</style>
