<template>
  <div>
    <el-row :gutter="16" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-item">
            <div class="stat-icon blue"><el-icon :size="28"><Money /></el-icon></div>
            <div>
              <div class="stat-label">今日交易额</div>
              <div class="stat-value">¥{{ stats.today_amount?.toFixed(2) || '0.00' }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-item">
            <div class="stat-icon green"><el-icon :size="28"><Scale /></el-icon></div>
            <div>
              <div class="stat-label">今日收货量</div>
              <div class="stat-value">{{ stats.today_weight?.toFixed(2) || '0' }} 公斤</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-item">
            <div class="stat-icon orange"><el-icon :size="28"><Document /></el-icon></div>
            <div>
              <div class="stat-label">今日交易笔数</div>
              <div class="stat-value">{{ stats.today_count || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-item">
            <div class="stat-icon purple"><el-icon :size="28"><Warehouse /></el-icon></div>
            <div>
              <div class="stat-label">库存种类</div>
              <div class="stat-value">{{ inventory.length }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>价格走势（近30天）</span>
              <el-select v-model="selectedGoods" placeholder="选择货物" style="width:160px" @change="loadTrend">
                <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
              </el-select>
            </div>
          </template>
          <div ref="chartRef" style="height:320px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>库存概览</span></template>
          <el-table :data="inventory" size="small" stripe>
            <el-table-column prop="goods.name" label="货物" />
            <el-table-column label="库存(公斤)" width="120">
              <template #default="{ row }">{{ row.quantity.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="hover" style="margin-top:16px">
      <template #header><span>今日定价</span></template>
      <el-table :data="todayPricing" size="small" stripe>
        <el-table-column prop="goods.name" label="货物" width="150" />
        <el-table-column prop="goods.category" label="类别" width="150" />
        <el-table-column label="今日单价(元/公斤)" width="180">
          <template #default="{ row }"><span style="color:#f56c6c;font-weight:bold">¥{{ row.price.toFixed(2) }}</span></template>
        </el-table-column>
        <el-table-column prop="date" label="日期" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { Money, Scale, Document, Warehouse } from '@element-plus/icons-vue'
import api from '../../api'

const stats = ref({})
const inventory = ref([])
const goodsList = ref([])
const todayPricing = ref([])
const selectedGoods = ref('')
const chartRef = ref()
let chart = null

const loadData = async () => {
  const [s, inv, goods, pricing] = await Promise.all([
    api.get('/transactions/stats'),
    api.get('/inventory'),
    api.get('/goods'),
    api.get('/pricing/today')
  ])
  stats.value = s.data
  inventory.value = inv.data
  goodsList.value = goods.data
  todayPricing.value = pricing.data
  if (goodsList.value.length > 0 && !selectedGoods.value) {
    selectedGoods.value = goodsList.value[0].id
    loadTrend()
  }
}

const loadTrend = async () => {
  if (!selectedGoods.value) return
  const res = await api.get('/pricing/trend', { params: { goods_id: selectedGoods.value, days: 30 } })
  const data = res.data
  const dates = [...new Set(data.map(d => d.date))].sort()
  const prices = dates.map(date => {
    const item = data.find(d => d.date === date)
    return item ? item.price : null
  })
  await nextTick()
  if (!chart) {
    chart = echarts.init(chartRef.value)
  }
  const goodsName = goodsList.value.find(g => g.id === selectedGoods.value)?.name || ''
  chart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 50, right: 20, top: 30, bottom: 40 },
    xAxis: { type: 'category', data: dates, axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: { type: 'value', name: '元/公斤' },
    series: [{
      name: goodsName,
      type: 'line',
      data: prices,
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      lineStyle: { width: 2, color: '#409eff' },
      itemStyle: { color: '#409eff' },
      areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: 'rgba(64,158,255,0.3)' },
        { offset: 1, color: 'rgba(64,158,255,0.05)' }
      ])}
    }]
  })
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', () => chart?.resize())
})
</script>

<style scoped>
.stat-cards { margin-bottom: 0; }
.stat-item { display: flex; align-items: center; gap: 16px; }
.stat-icon { width: 56px; height: 56px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: #fff; }
.stat-icon.blue { background: linear-gradient(135deg, #409eff, #66b1ff); }
.stat-icon.green { background: linear-gradient(135deg, #67c23a, #85ce61); }
.stat-icon.orange { background: linear-gradient(135deg, #e6a23c, #ebb563); }
.stat-icon.purple { background: linear-gradient(135deg, #909399, #a6a9ad); }
.stat-label { color: #999; font-size: 13px; margin-bottom: 4px; }
.stat-value { font-size: 22px; font-weight: bold; color: #333; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
