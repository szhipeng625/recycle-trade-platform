<template>
  <div>
    <el-card shadow="hover" style="margin-bottom:16px">
      <template #header>
        <div class="card-header">
          <span>今日回收价格</span>
          <el-tag type="success">{{ today }}</el-tag>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :span="8" v-for="p in pricing" :key="p.id">
          <div class="price-card">
            <div class="price-name">{{ p.goods?.name }}</div>
            <div class="price-cat">{{ p.goods?.category }}</div>
            <div class="price-value">¥{{ p.price.toFixed(2) }}<span class="price-unit">/公斤</span></div>
          </div>
        </el-col>
      </el-row>
      <el-empty v-if="pricing.length === 0" description="今日暂无定价" />
    </el-card>

    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>价格走势</span>
          <el-select v-model="selectedGoods" placeholder="选择货物" style="width:160px" @change="loadTrend">
            <el-option v-for="g in goodsList" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </div>
      </template>
      <div ref="chartRef" style="height:300px"></div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import api from '../../api'

const pricing = ref([])
const goodsList = ref([])
const selectedGoods = ref('')
const chartRef = ref()
let chart = null
const today = new Date().toISOString().slice(0, 10)

const loadData = async () => {
  const [p, g] = await Promise.all([api.get('/pricing/today'), api.get('/goods')])
  pricing.value = p.data
  goodsList.value = g.data
  if (g.data.length > 0) {
    selectedGoods.value = g.data[0].id
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
  if (!chart) chart = echarts.init(chartRef.value)
  const goodsName = goodsList.value.find(g => g.id === selectedGoods.value)?.name || ''
  chart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 50, right: 20, top: 30, bottom: 40 },
    xAxis: { type: 'category', data: dates, axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: { type: 'value', name: '元/公斤' },
    series: [{ name: goodsName, type: 'line', data: prices, smooth: true, symbol: 'circle', symbolSize: 6,
      lineStyle: { width: 2, color: '#67c23a' }, itemStyle: { color: '#67c23a' },
      areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: 'rgba(103,194,58,0.3)' }, { offset: 1, color: 'rgba(103,194,58,0.05)' }
      ])}
    }]
  })
}

onMounted(() => { loadData(); window.addEventListener('resize', () => chart?.resize()) })
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.price-card { background: linear-gradient(135deg, #f0f9eb, #e1f3d8); border-radius: 10px; padding: 20px; margin-bottom: 16px; text-align: center; }
.price-name { font-size: 18px; font-weight: bold; color: #333; }
.price-cat { font-size: 12px; color: #999; margin: 4px 0; }
.price-value { font-size: 28px; font-weight: bold; color: #67c23a; }
.price-unit { font-size: 14px; color: #999; font-weight: normal; }
</style>
