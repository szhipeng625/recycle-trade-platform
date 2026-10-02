<template>
  <div>
    <el-card shadow="hover">
      <template #header><span>我的交易记录</span></template>
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="单号" width="80" />
        <el-table-column label="货物" width="120">
          <template #default="{ row }">{{ row.goods?.name }}</template>
        </el-table-column>
        <el-table-column label="重量" width="160">
          <template #default="{ row }">{{ row.weight }} {{ row.unit }}</template>
        </el-table-column>
        <el-table-column label="单价" width="140">
          <template #default="{ row }">¥{{ row.unit_price.toFixed(2) }}/公斤</template>
        </el-table-column>
        <el-table-column label="金额" width="140">
          <template #default="{ row }"><span style="color:#f56c6c;font-weight:bold">¥{{ row.total_amount.toFixed(2) }}</span></template>
        </el-table-column>
        <el-table-column label="时间">
          <template #default="{ row }">{{ formatTime(row.transaction_time) }}</template>
        </el-table-column>
      </el-table>
      <el-empty v-if="list.length === 0" description="暂无交易记录" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../api'

const list = ref([])

const loadData = async () => {
  // 用户端暂时展示所有交易（实际应按用户关联的卖货人筛选）
  const res = await api.get('/transactions', { params: { size: 100 } })
  list.value = res.data.list
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(loadData)
</script>
