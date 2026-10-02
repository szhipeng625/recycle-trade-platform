import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  {
    path: '/admin',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'AdminDashboard', component: () => import('../views/admin/Dashboard.vue') },
      { path: 'goods', name: 'AdminGoods', component: () => import('../views/admin/Goods.vue') },
      { path: 'pricing', name: 'AdminPricing', component: () => import('../views/admin/Pricing.vue') },
      { path: 'transactions', name: 'AdminTransactions', component: () => import('../views/admin/Transactions.vue') },
      { path: 'sellers', name: 'AdminSellers', component: () => import('../views/admin/Sellers.vue') },
      { path: 'inventory', name: 'AdminInventory', component: () => import('../views/admin/Inventory.vue') }
    ]
  },
  {
    path: '/user',
    component: () => import('../layouts/UserLayout.vue'),
    meta: { requiresAuth: true, role: 'user' },
    children: [
      { path: '', redirect: '/user/market' },
      { path: 'market', name: 'UserMarket', component: () => import('../views/user/Market.vue') },
      { path: 'sell', name: 'UserSell', component: () => import('../views/user/Sell.vue') },
      { path: 'records', name: 'UserRecords', component: () => import('../views/user/MyRecords.vue') }
    ]
  },
  { path: '/', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  let user = {}
  try {
    user = JSON.parse(localStorage.getItem('user') || '{}')
  } catch (e) {
    user = {}
  }
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.meta.role && user.role !== to.meta.role) {
    next(user.role === 'admin' ? '/admin' : '/user')
  } else {
    next()
  }
})

export default router
